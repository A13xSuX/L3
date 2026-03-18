package repository

import (
	"context"
	"database/sql"
	"errors"
	"l3/WarehouseControl/internal/models"

	"github.com/wb-go/wbf/dbpg"
)

type UserRepo struct {
	db *dbpg.DB
}

func NewUserRepo(db *dbpg.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `SELECT id, username, password_hash, role, created_at FROM users
				WHERE username = $1`
	var user models.User
	row := r.db.QueryRowContext(ctx, query, username)
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
