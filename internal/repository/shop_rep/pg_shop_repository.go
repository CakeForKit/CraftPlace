package shoprep

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/CakeForKit/CraftPlace.git/internal/cnfg"
	"github.com/CakeForKit/CraftPlace.git/internal/models/models"
	dberrors "github.com/CakeForKit/CraftPlace.git/internal/repository/db_errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PgShopRep struct {
	db *sql.DB
}

func NewPgShopRep(
	ctx context.Context,
	pgCreds *cnfg.PostgresCredentials,
	dbConf *cnfg.DatebaseConnConfig,
) (ShopRep, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		pgCreds.Username, pgCreds.Password, pgCreds.Host, pgCreds.Port, pgCreds.DbName)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("NewPgShopRep: %w: %w", dberrors.ErrOpenConnect, err)
	}

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("NewPgShopRep: %w: %w", dberrors.ErrPing, err)
	}
	// Настраиваем пул соединений
	db.SetMaxOpenConns(dbConf.MaxOpenConns)
	db.SetMaxIdleConns(dbConf.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(dbConf.ConnMaxLifetime.Hours()))

	return &PgShopRep{db: db}, nil
}

func (pg *PgShopRep) parseShopsRows(rows *sql.Rows) ([]*models.Shop, error) {
	baseErr := errors.New("parseShopsRows")
	var resShops []*models.Shop
	for rows.Next() {
		var id, userID uuid.UUID
		var title, description string
		var updateTime time.Time
		if err := rows.Scan(&id, &title, &description, &userID, &updateTime); err != nil {
			return nil, fmt.Errorf("%w scan error: %w", baseErr, err)
		}
		shop, err := models.NewShop(id, title, description, userID, updateTime)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", baseErr, err)
		}
		resShops = append(resShops, shop)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w rows iteration error: %w", baseErr, err)
	}
	return resShops, nil
}

func (pg *PgShopRep) addFilterParams(query sq.SelectBuilder, filterOps *models.ShopFilter) sq.SelectBuilder {
	if filterOps.Title != "" {
		query = query.Where(sq.ILike{"shops.title": "%" + filterOps.Title + "%"})
	}
	if filterOps.UserID != uuid.Nil {
		query = query.Where(sq.Eq{"shops.user_id": filterOps.UserID})
	}
	query = query.Offset(filterOps.Offset)
	if filterOps.Limit != 0 {
		query = query.Limit(filterOps.Limit)
	}
	return query
}

func (pg *PgShopRep) execSelectQuery(ctx context.Context, query sq.SelectBuilder) ([]*models.Shop, error) {
	querySQL, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", dberrors.ErrQueryBuilds, err)
	}

	rows, err := pg.db.QueryContext(ctx, querySQL, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", dberrors.ErrQueryExec, err)
	}
	defer rows.Close()

	arts, err := pg.parseShopsRows(rows)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return arts, nil
}

func (pg *PgShopRep) GetByID(ctx context.Context, id uuid.UUID) (*models.Shop, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select(
		"shops.id", "shops.title", "shops.description",
		"shops.user_id", "shops.update_time").
		From("shops").
		Where(sq.Eq{"shops.id": id})

	arts, err := pg.execSelectQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("PgShopRep.GetByID: %w", err)
	}

	if len(arts) == 0 {
		return nil, ErrShopNotFound
	} else if len(arts) > 1 {
		return nil, fmt.Errorf("PgShopRep.GetByID: %w", dberrors.ErrExpectedOne)
	}
	return arts[0], nil
}

func (pg *PgShopRep) GetByFilter(ctx context.Context, filterOps *models.ShopFilter) ([]*models.Shop, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select(
		"shops.id", "shops.title", "shops.description",
		"shops.user_id", "shops.update_time").
		From("shops")

	query = pg.addFilterParams(query, filterOps)
	arts, err := pg.execSelectQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("PgShopRep.GetByFilter: %w", err)
	}
	return arts, nil
}

func (pg *PgShopRep) execChangeQuery(ctx context.Context, query sq.Sqlizer) error {
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
		return fmt.Errorf("%w: no shops changed", dberrors.ErrRowsAffected)
	}
	return nil
}

func (pg *PgShopRep) Add(ctx context.Context, e *models.Shop) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Insert("shops").
		Columns("id", "title", "description", "user_id", "update_time").
		Values(e.GetID(), e.GetTitle(), e.GetDescription(), e.GetUserID(), e.GetUpdateTime())

	err := pg.execChangeQuery(ctx, query)
	if err != nil {
		return fmt.Errorf("PgShopRep.Add: %w", err)
	}
	return nil
}

func (pg *PgShopRep) Delete(ctx context.Context, id uuid.UUID) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Delete("shops").
		Where(sq.Eq{"id": id})
	err := pg.execChangeQuery(ctx, query)
	if err != nil {
		return fmt.Errorf("PgShopRep.Delete: %w", err)
	}
	return nil
}

func (pg *PgShopRep) Update(ctx context.Context,
	shopID uuid.UUID,
	funcUpdate func(*models.Shop) (*models.Shop, error),
) (*models.Shop, error) {
	art, err := pg.GetByID(ctx, shopID)
	if err != nil {
		return nil, fmt.Errorf("PgShopRep.Update: %w", err)
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	updatedArtwork, err := funcUpdate(art)
	if err != nil {
		return nil, fmt.Errorf("PgShopRep.Update: %w (%w)", ErrUpdateShop, err)
	}
	query := psql.Update("shops").
		Set("title", updatedArtwork.GetTitle()).
		Set("description", updatedArtwork.GetDescription()).
		Set("update_time", updatedArtwork.GetUpdateTime()).
		Where(sq.Eq{"id": shopID})
	err = pg.execChangeQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("PgShopRep.Update %w", err)
	}
	return updatedArtwork, nil
}

func (pg *PgShopRep) Ping(ctx context.Context) error {
	return pg.db.PingContext(ctx)
}

func (pg *PgShopRep) Close() {
	pg.db.Close()
}
