package repository

import (
	"context"
	"database/sql"
	"l3/WarehouseControl/internal/customErrs"
	"l3/WarehouseControl/internal/models"
	"strconv"

	"github.com/wb-go/wbf/dbpg"
)

type ItemRepository struct {
	db *dbpg.DB
}

func NewItemRepository(db *dbpg.DB) *ItemRepository {
	return &ItemRepository{
		db: db,
	}
}

func (r *ItemRepository) Create(ctx context.Context, user *models.CurrentUser, item *models.Item) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = r.setAuditContext(ctx, tx, user)
	if err != nil {
		return err
	}
	query := `INSERT INTO items (title, sku, quantity) 
				VALUES ($1, $2, $3)
				RETURNING id, created_at, updated_at`
	row := tx.QueryRowContext(ctx, query, item.Title, item.Sku, item.Quantity)
	err = row.Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (r *ItemRepository) GetAll(ctx context.Context) ([]models.Item, error) {
	query := `SELECT id, title, sku, quantity, created_at, updated_at FROM items`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Sku,
			&item.Quantity,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return items, nil
}

func (r *ItemRepository) Update(ctx context.Context, user *models.CurrentUser, id int64, newItem *models.Item) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	err = r.setAuditContext(ctx, tx, user)
	if err != nil {
		return err
	}
	query := `UPDATE items
			SET title = $1,
				sku = $2,
				quantity = $3,
				updated_at = NOW()
			WHERE id = $4`
	res, err := tx.ExecContext(ctx, query, newItem.Title, newItem.Sku, newItem.Quantity, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return customErrs.ErrNotFoundID
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (r *ItemRepository) Delete(ctx context.Context, user *models.CurrentUser, id int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = r.setAuditContext(ctx, tx, user)
	if err != nil {
		return err
	}
	query := `DELETE FROM items
				WHERE id = $1`
	res, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return customErrs.ErrNotFoundID
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (r *ItemRepository) setAuditContext(ctx context.Context, tx *sql.Tx, user *models.CurrentUser) error {
	if user == nil {
		return customErrs.ErrUnauthorized
	}
	_, err := tx.ExecContext(ctx,
		`SELECT
				set_config('app.current_user_id', $1, true),
				set_config('app.current_username', $2, true),
				set_config('app.current_role', $3, true)`,
		strconv.FormatInt(user.UserID, 10),
		user.Username,
		user.Role,
	)
	return err
}
