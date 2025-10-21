package categoryrep

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

type PgCategoryRep struct {
	db *sql.DB
}

func NewPgCategoryRep(
	ctx context.Context,
	pgCreds *cnfg.PostgresCredentials,
	dbConf *cnfg.DatebaseConnConfig,
) (CategoryRep, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		pgCreds.Username, pgCreds.Password, pgCreds.Host, pgCreds.Port, pgCreds.DbName)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("NewPgCategoryRep: %w: %w", dberrors.ErrOpenConnect, err)
	}

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("NewPgCategoryRep: %w: %w", dberrors.ErrPing, err)
	}
	// Настраиваем пул соединений
	db.SetMaxOpenConns(dbConf.MaxOpenConns)
	db.SetMaxIdleConns(dbConf.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(dbConf.ConnMaxLifetime.Hours()))

	return &PgCategoryRep{db: db}, nil
}

func (pg *PgCategoryRep) parseCategorysRows(rows *sql.Rows) ([]*models.Category, error) {
	baseErr := errors.New("parseCategorysRows")
	var resCategorys []*models.Category
	for rows.Next() {
		var id uuid.UUID
		var title, description string
		if err := rows.Scan(&id, &title, &description); err != nil {
			return nil, fmt.Errorf("%w scan error: %w", baseErr, err)
		}
		category, err := models.NewCategory(id, title, description)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", baseErr, err)
		}
		resCategorys = append(resCategorys, category)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w rows iteration error: %w", baseErr, err)
	}
	return resCategorys, nil
}

func (pg *PgCategoryRep) addFilterParams(query sq.SelectBuilder, filterOps *models.CategoryFilter) sq.SelectBuilder {
	if filterOps.Title != "" {
		query = query.Where(sq.ILike{"categories.title": "%" + filterOps.Title + "%"})
	}
	query = query.Offset(filterOps.Offset)
	if filterOps.Limit != 0 {
		query = query.Limit(filterOps.Limit)
	}
	return query
}

func (pg *PgCategoryRep) execSelectQuery(ctx context.Context, query sq.SelectBuilder) ([]*models.Category, error) {
	querySQL, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", dberrors.ErrQueryBuilds, err)
	}

	rows, err := pg.db.QueryContext(ctx, querySQL, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", dberrors.ErrQueryExec, err)
	}
	defer rows.Close()

	arts, err := pg.parseCategorysRows(rows)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return arts, nil
}

func (pg *PgCategoryRep) GetByID(ctx context.Context, categoryID uuid.UUID) (*models.Category, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select(
		"categories.id", "categories.title", "categories.description").
		From("categories").
		Where(sq.Eq{"categories.id": categoryID})

	arts, err := pg.execSelectQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("PgCategoryRep.GetByID: %w", err)
	}

	if len(arts) == 0 {
		return nil, ErrCategoryNotFound
	} else if len(arts) > 1 {
		return nil, fmt.Errorf("PgCategoryRep.GetByID: %w", dberrors.ErrExpectedOne)
	}
	return arts[0], nil
}

func (pg *PgCategoryRep) GetByFilter(ctx context.Context, filterOps *models.CategoryFilter) ([]*models.Category, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select(
		"categories.id", "categories.title", "categories.description").
		From("categories")

	query = pg.addFilterParams(query, filterOps)
	arts, err := pg.execSelectQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("PgCategoryRep.GetByFilter: %w", err)
	}
	return arts, nil
}

func (pg *PgCategoryRep) execChangeQuery(ctx context.Context, query sq.Sqlizer) error {
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
		return fmt.Errorf("%w: no categories changed", dberrors.ErrRowsAffected)
	}
	return nil
}

func (pg *PgCategoryRep) Add(ctx context.Context, e *models.Category) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Insert("categories").
		Columns("id", "title", "description").
		Values(e.GetID(), e.GetTitle(), e.GetDescription())

	err := pg.execChangeQuery(ctx, query)
	if err != nil {
		return fmt.Errorf("PgCategoryRep.Add: %w", err)
	}
	return nil
}

func (pg *PgCategoryRep) Delete(ctx context.Context, id uuid.UUID) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Delete("categories").
		Where(sq.Eq{"id": id})
	err := pg.execChangeQuery(ctx, query)
	if err != nil {
		return fmt.Errorf("PgCategoryRep.Delete: %w", err)
	}
	return nil
}

func (pg *PgCategoryRep) Ping(ctx context.Context) error {
	return pg.db.PingContext(ctx)
}

func (pg *PgCategoryRep) Close() {
	pg.db.Close()
}
