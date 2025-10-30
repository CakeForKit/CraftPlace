package productrep

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/CakeForKit/CraftPlace.git/internal/cnfg"
	"github.com/CakeForKit/CraftPlace.git/internal/models/models"
	dberrors "github.com/CakeForKit/CraftPlace.git/internal/repository/db_errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PgProductRep struct {
	db         *sql.DB
	isReadOnly bool
}

func NewPgProductRep(
	ctx context.Context,
	pgCreds *cnfg.PostgresCredentials,
	dbConf *cnfg.DatebaseConnConfig,
) (ProductRep, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		pgCreds.Username, pgCreds.Password, pgCreds.Host, pgCreds.Port, pgCreds.DbName)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("NewPgProductRep: %w: %w", dberrors.ErrOpenConnect, err)
	}

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("NewPgProductRep: %w: %w", dberrors.ErrPing, err)
	}
	// Настраиваем пул соединений
	db.SetMaxOpenConns(dbConf.MaxOpenConns)
	db.SetMaxIdleConns(dbConf.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(dbConf.ConnMaxLifetime.Hours()))

	return &PgProductRep{db: db, isReadOnly: pgCreds.ReadOnly}, nil
}

func (pg *PgProductRep) parseProductsRows(rows *sql.Rows) ([]*models.Product, error) {
	baseErr := errors.New("parseProductsRows")
	var resProducts []*models.Product
	for rows.Next() {
		var id, shopID uuid.UUID
		var title, description string
		var cost uint64
		var updateTime time.Time
		if err := rows.Scan(&id, &title, &description, &cost, &shopID, &updateTime); err != nil {
			return nil, fmt.Errorf("%w scan error: %w", baseErr, err)
		}
		post, err := models.NewProduct(id, title, description, cost, shopID, uuid.UUIDs{}, updateTime)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", baseErr, err)
		}
		resProducts = append(resProducts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w rows iteration error: %w", baseErr, err)
	}
	return resProducts, nil
}

func (pg *PgProductRep) addFilterParams(query sq.SelectBuilder, filterOps *models.ProductFilter) sq.SelectBuilder {
	if filterOps.Title != "" {
		query = query.Where(sq.ILike{"products.title": "%" + filterOps.Title + "%"})
	}
	if filterOps.ShopID != uuid.Nil {
		query = query.Where(sq.Eq{"products.shop_id": filterOps.ShopID})
	}
	if filterOps.MinCost > 0 || filterOps.MaxCost < math.MaxUint64 {
		if filterOps.MinCost > 0 && filterOps.MaxCost < math.MaxUint64 {
			query = query.Where(sq.And{
				sq.GtOrEq{"products.cost": filterOps.MinCost},
				sq.LtOrEq{"products.cost": filterOps.MaxCost},
			})
		} else if filterOps.MinCost > 0 {
			query = query.Where(sq.GtOrEq{"products.cost": filterOps.MinCost})
		} else if filterOps.MaxCost < math.MaxUint64 {
			query = query.Where(sq.LtOrEq{"products.cost": filterOps.MaxCost})
		}
	}
	query = query.Offset(filterOps.Offset)
	if filterOps.Limit != 0 {
		query = query.Limit(filterOps.Limit)
	}
	return query
}

func (pg *PgProductRep) addSortParams(query sq.SelectBuilder) sq.SelectBuilder {
	query = query.OrderBy("products.update_time DESC ")
	return query
}

func (pg *PgProductRep) execSelectQuery(ctx context.Context, query sq.SelectBuilder) ([]*models.Product, error) {
	querySQL, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", dberrors.ErrQueryBuilds, err)
	}

	rows, err := pg.db.QueryContext(ctx, querySQL, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", dberrors.ErrQueryExec, err)
	}
	defer rows.Close()

	arts, err := pg.parseProductsRows(rows)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return arts, nil
}

func (pg *PgProductRep) GetByID(ctx context.Context, productID uuid.UUID) (*models.Product, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select(
		"products.id", "products.title", "products.description",
		"products.cost", "products.shop_id", "products.update_time").
		From("products").
		Where(sq.Eq{"products.id": productID})

	products, err := pg.execSelectQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("PgProductRep.GetByID: %w", err)
	}

	if len(products) == 0 {
		return nil, ErrProductNotFound
	} else if len(products) > 1 {
		return nil, fmt.Errorf("PgProductRep.GetByID: %w", dberrors.ErrExpectedOne)
	}
	resProduct := products[0]
	categoryIDs, err := pg.getCategoryIDs(ctx, resProduct.GetID())
	if err != nil {
		return nil, fmt.Errorf("PgProductRep.GetByID: %w", err)
	}
	err = resProduct.AddCategoryIDs(categoryIDs)
	if err != nil {
		return nil, fmt.Errorf("PgProductRep.GetByID: %w", err)
	}
	return resProduct, nil
}

func (pg *PgProductRep) getCategoryIDs(ctx context.Context, productID uuid.UUID) (uuid.UUIDs, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select("category_id").
		From("product_category").
		Where(sq.Eq{"product_id": productID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", dberrors.ErrQueryBuilds, err)
	}
	rows, err := pg.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", dberrors.ErrQueryExec, err)
	}
	defer rows.Close()

	var categoryIDs uuid.UUIDs
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("PgProductRep.getCategoryIDs: %w", err)
		}
		categoryIDs = append(categoryIDs, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("PgProductRep.getCategoryIDs rows iteration error: %w", err)
	}
	return categoryIDs, nil
}

func (pg *PgProductRep) joinCategoryIDsToProducts(ctx context.Context, products []*models.Product, filterOps *models.ProductFilter) ([]*models.Product, error) {
	resProducts := make([]*models.Product, 0)
	for _, p := range products {
		categoryIDs, err := pg.getCategoryIDs(ctx, p.GetID())
		if err != nil {
			return nil, fmt.Errorf("join CategoryIDs %w", err)
		}
		if filterOps.CategoryID != uuid.Nil {
			in := false
			for _, v := range categoryIDs {
				if filterOps.CategoryID == v {
					in = true
					break
				}
			}
			if in {
				if err := p.AddCategoryIDs(categoryIDs); err != nil {
					return nil, fmt.Errorf("join CategoryIDs %w", err)
				}
				resProducts = append(resProducts, p)
			}
		} else {
			if err := p.AddCategoryIDs(categoryIDs); err != nil {
				return nil, fmt.Errorf("join CategoryIDs %w", err)
			}
			resProducts = append(resProducts, p)
		}
	}
	return resProducts, nil
}

func (pg *PgProductRep) GetByFilter(ctx context.Context, filterOps *models.ProductFilter) ([]*models.Product, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select(
		"products.id", "products.title", "products.description",
		"products.cost", "products.shop_id", "products.update_time").
		From("products")

	query = pg.addFilterParams(query, filterOps)
	query = pg.addSortParams(query)
	products, err := pg.execSelectQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("PgProductRep.GetByFilter: %w", err)
	}
	products, err = pg.joinCategoryIDsToProducts(ctx, products, filterOps)
	if err != nil {
		return nil, fmt.Errorf("PgProductRep.GetByFilter: %w", err)
	}
	return products, nil
}

func (pg *PgProductRep) execChangeQuery(ctx context.Context, query sq.Sqlizer) error {
	if pg.isReadOnly {
		return ErrReadOnly
	}
	querySQL, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("%w: %w", dberrors.ErrQueryBuilds, err)
	}
	result, err := pg.db.ExecContext(ctx, querySQL, args...)
	if err != nil {
		return fmt.Errorf("%w: %w", dberrors.ErrQueryExec, err)
	}
	// проверка количества затронутых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %w", dberrors.ErrRowsAffected, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%w: no products changed", dberrors.ErrRowsAffected)
	}
	return nil
}

func (pg *PgProductRep) Add(ctx context.Context, e *models.Product) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Insert("products").
		Columns("id", "title", "description", "cost", "shop_id", "update_time").
		Values(e.GetID(), e.GetTitle(), e.GetDescription(), e.GetCost(), e.GetShopID(), e.GetUpdateTime())

	err := pg.execChangeQuery(ctx, query)
	if err != nil {
		return fmt.Errorf("PgProductRep.Add: %w", err)
	}
	return nil
}

func (pg *PgProductRep) Delete(ctx context.Context, id uuid.UUID) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Delete("products").
		Where(sq.Eq{"id": id})
	err := pg.execChangeQuery(ctx, query)
	if err != nil {
		return fmt.Errorf("PgProductRep.Delete: %w", err)
	}
	return nil
}

func (pg *PgProductRep) Update(ctx context.Context,
	productID uuid.UUID,
	funcUpdate func(*models.Product) (*models.Product, error),
) (*models.Product, error) {
	baseErr := fmt.Errorf("PgProductRep Update")
	p, err := pg.GetByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	updatedProduct, err := funcUpdate(p)
	if err != nil {
		return nil, fmt.Errorf("%w: %w (%w)", baseErr, ErrUpdateProduct, err)
	}
	query := psql.Update("products").
		Set("title", updatedProduct.GetTitle()).
		Set("description", updatedProduct.GetDescription()).
		Set("cost", updatedProduct.GetCost()).
		Set("update_time", updatedProduct.GetUpdateTime()).
		Where(sq.Eq{"id": productID})
	err = pg.execChangeQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%w %w", baseErr, err)
	}
	return updatedProduct, nil
}

func (pg *PgProductRep) Ping(ctx context.Context) error {
	return pg.db.PingContext(ctx)
}

func (pg *PgProductRep) Close() {
	pg.db.Close()
}
