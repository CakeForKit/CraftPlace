package postrep

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/CakeForKit/CraftPlace.git/internal/cnfg"
	"github.com/CakeForKit/CraftPlace.git/internal/models/models"
	reqresp "github.com/CakeForKit/CraftPlace.git/internal/models/req_resp"
	dberrors "github.com/CakeForKit/CraftPlace.git/internal/repository/db_errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type PgPostRep struct {
	db *sql.DB
}

func NewPgPostRep(
	ctx context.Context,
	pgCreds *cnfg.PostgresCredentials,
	dbConf *cnfg.DatebaseConnConfig,
) (PostRep, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		pgCreds.Username, pgCreds.Password, pgCreds.Host, pgCreds.Port, pgCreds.DbName)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("NewPgPostRep: %w: %w", dberrors.ErrOpenConnect, err)
	}

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("NewPgPostRep: %w: %w", dberrors.ErrPing, err)
	}
	// Настраиваем пул соединений
	db.SetMaxOpenConns(dbConf.MaxOpenConns)
	db.SetMaxIdleConns(dbConf.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(dbConf.ConnMaxLifetime.Hours()))

	return &PgPostRep{db: db}, nil
}

func (pg *PgPostRep) parsePostsRows(rows *sql.Rows) ([]*models.Post, error) {
	baseErr := errors.New("parsePostsRows")
	var resPosts []*models.Post
	for rows.Next() {
		var id, shopID uuid.UUID
		var description string
		var timePublication time.Time
		if err := rows.Scan(&id, &description, &timePublication, &shopID); err != nil {
			return nil, fmt.Errorf("%w scan error: %w", baseErr, err)
		}
		post, err := models.NewPost(id, description, timePublication, shopID)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", baseErr, err)
		}
		resPosts = append(resPosts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w rows iteration error: %w", baseErr, err)
	}
	return resPosts, nil
}

func (pg *PgPostRep) addFilterParams(query sq.SelectBuilder, filterOps *reqresp.PostFilter) sq.SelectBuilder {
	if filterOps.ShopID != uuid.Nil {
		query = query.Where(sq.Eq{"posts.shop_id": filterOps.ShopID})
	}
	return query
}

func (pg *PgPostRep) addSortParams(query sq.SelectBuilder) sq.SelectBuilder {
	query = query.OrderBy("artworks.publication_time DESC ")
	return query
}

func (pg *PgPostRep) execSelectQuery(ctx context.Context, query sq.SelectBuilder) ([]*models.Post, error) {
	querySQL, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", dberrors.ErrQueryBuilds, err)
	}

	rows, err := pg.db.QueryContext(ctx, querySQL, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", dberrors.ErrQueryExec, err)
	}
	defer rows.Close()

	arts, err := pg.parsePostsRows(rows)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return arts, nil
}

func (pg *PgPostRep) GetByFilter(ctx context.Context, filterOps *reqresp.PostFilter) ([]*models.Post, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select(
		"posts.id", "posts.description",
		"posts.publication_time", "posts.shop_id").
		From("posts")

	query = pg.addFilterParams(query, filterOps)
	query = pg.addSortParams(query)
	arts, err := pg.execSelectQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("PgPostRep.GetByFilter: %w", err)
	}
	return arts, nil
}

func (pg *PgPostRep) execChangeQuery(ctx context.Context, query sq.Sqlizer) error {
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
		return fmt.Errorf("%w: no posts changed", dberrors.ErrRowsAffected)
	}
	return nil
}

func (pg *PgPostRep) Add(ctx context.Context, e *models.Post) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Insert("posts").
		Columns("id", "description", "publication_time", "shop_id").
		Values(e.GetID(), e.GetDescription(), e.GetTimePublication(), e.GetShopID())

	err := pg.execChangeQuery(ctx, query)
	if err != nil {
		return fmt.Errorf("PgPostRep.Add: %w", err)
	}
	return nil
}

func (pg *PgPostRep) Delete(ctx context.Context, id uuid.UUID) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Delete("posts").
		Where(sq.Eq{"id": id})
	err := pg.execChangeQuery(ctx, query)
	if err != nil {
		return fmt.Errorf("PgPostRep.Delete: %w", err)
	}
	return nil
}

func (pg *PgPostRep) Ping(ctx context.Context) error {
	return pg.db.PingContext(ctx)
}

func (pg *PgPostRep) Close() {
	pg.db.Close()
}
