package userrep

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

type PgUserRep struct {
	db *sql.DB
}

func NewPgUserRep(
	ctx context.Context,
	pgCreds *cnfg.PostgresCredentials,
	dbConf *cnfg.DatebaseConnConfig,
) (UserRep, error) {
	// connStr := "postgres://puser:ppassword@postgres_artworks:5432/artworks"
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		pgCreds.Username, pgCreds.Password, pgCreds.Host, pgCreds.Port, pgCreds.DbName)
	fmt.Printf("connStr: %s\n", connStr)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("NewPgUserRep: %w: %w", dberrors.ErrOpenConnect, err)
	}

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("NewPgUserRep: %w: %w", dberrors.ErrPing, err)
	}
	// Настраиваем пул соединений
	db.SetMaxOpenConns(dbConf.MaxOpenConns)
	db.SetMaxIdleConns(dbConf.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(dbConf.ConnMaxLifetime.Hours()))

	return &PgUserRep{db: db}, nil
}

func (pg *PgUserRep) parseUsersRows(rows *sql.Rows) ([]*models.User, error) {
	baseErr := errors.New("parseUsersRows")
	var resUsers []*models.User
	for rows.Next() {
		var id uuid.UUID
		var login, hashedPassword string
		if err := rows.Scan(&id, &login, &hashedPassword); err != nil {
			return nil, fmt.Errorf("%w scan error: %w", baseErr, err)
		}
		user, err := models.NewUser(id, login, hashedPassword)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", baseErr, err)
		}
		resUsers = append(resUsers, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w rows iteration error: %w", baseErr, err)
	}
	return resUsers, nil
}

func (pg *PgUserRep) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	baseErr := errors.New("PgUserRep.GetByID")
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select("id", "login", "hashed_password").
		From("users").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %w: %w", baseErr, dberrors.ErrQueryBuilds, err)
	}

	rows, err := pg.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %w: %w", baseErr, dberrors.ErrQueryExec, err)
	}
	defer rows.Close()
	users, err := pg.parseUsersRows(rows)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	if len(users) == 0 {
		return nil, ErrUserNotFound
	} else if len(users) > 1 {
		return nil, fmt.Errorf("%w: %w: %w", baseErr, dberrors.ErrExpectedOne, err)
	}
	return users[0], nil
}

func (pg *PgUserRep) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	baseErr := errors.New("PgUserRep.GetByLogin")
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select("id", "login", "hashed_password").
		From("users").
		Where(sq.Eq{"login": login}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %w: %w", baseErr, dberrors.ErrQueryBuilds, err)
	}

	rows, err := pg.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %w: %w", baseErr, dberrors.ErrQueryExec, err)
	}
	defer rows.Close()
	users, err := pg.parseUsersRows(rows)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	if len(users) == 0 {
		return nil, ErrUserNotFound
	} else if len(users) > 1 {
		return nil, fmt.Errorf("%w: %w: %w", baseErr, dberrors.ErrExpectedOne, err)
	}
	return users[0], nil
}

func (pg *PgUserRep) Add(ctx context.Context, e *models.User) error {
	baseErr := errors.New("PgUserRep.Add")
	_, err := pg.GetByLogin(ctx, e.GetLogin())
	if err == nil {
		return ErrDuplicateLoginUser
	} else if err != ErrUserNotFound {
		return fmt.Errorf("%w %w", baseErr, err)
	}
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Insert("users").
		Columns("id", "login", "hashed_password").
		Values(e.GetID(), e.GetLogin(), e.GetHashedPassword()).
		ToSql()
	if err != nil {
		return fmt.Errorf("%w %w: %w", baseErr, dberrors.ErrQueryBuilds, err)
	}
	result, err := pg.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%w %w: %w", baseErr, dberrors.ErrQueryExec, err)
	}
	// проверка количества затронутых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w %w: %w", baseErr, dberrors.ErrRowsAffected, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%w %w: no user added", baseErr, dberrors.ErrRowsAffected)
	}
	return nil
}

func (pg *PgUserRep) Delete(ctx context.Context, id uuid.UUID) error {
	baseErr := errors.New("PgUserRep.Delete")
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Delete("Users").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%w %w: %w", baseErr, dberrors.ErrQueryBuilds, err)
	}
	result, err := pg.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%w %w: %w", baseErr, dberrors.ErrQueryExec, err)
	}
	// проверка количества затронутых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w %w: %w", baseErr, dberrors.ErrRowsAffected, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%w %w: no user with id %s", baseErr, dberrors.ErrRowsAffected, id)
	}
	return nil
}

func (pg *PgUserRep) Update(ctx context.Context,
	id uuid.UUID,
	funcUpdate func(*models.User) (*models.User, error),
) (*models.User, error) {
	baseErr := errors.New("PgUserRep.Update")
	user, err := pg.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	updatedUser, err := funcUpdate(user)
	if err != nil {
		return nil, fmt.Errorf("%w funcUpdate: %w", baseErr, err)
	}
	query, args, err := psql.Update("users").
		Set("login", updatedUser.GetLogin()).
		Set("hashedPassword", updatedUser.GetHashedPassword()).
		Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %w: %w", baseErr, dberrors.ErrQueryBuilds, err)
	}
	result, err := pg.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %w: %w", baseErr, dberrors.ErrQueryExec, err)
	}
	// проверка количества затронутых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("%w: %w: %w", baseErr, dberrors.ErrRowsAffected, err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("%w: %w: no user updated", baseErr, dberrors.ErrRowsAffected)
	}
	return updatedUser, nil
}

func (pg *PgUserRep) Ping(ctx context.Context) error {
	return pg.db.PingContext(ctx)
}

func (pg *PgUserRep) Close() {
	if pg.db != nil {
		pg.db.Close()
	}
}
