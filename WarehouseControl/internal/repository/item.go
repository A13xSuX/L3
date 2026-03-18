package repository

import (
	"context"
	"l3/WarehouseControl/internal/customErrs"
	"l3/WarehouseControl/internal/models"

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

func (r *ItemRepository) Create(ctx context.Context, item *models.Item) error {
	query := `INSERT INTO items (title, sku, quantity) 
				VALUES ($1, $2, $3)
				RETURNING id, created_at, updated_at`
	row := r.db.QueryRowContext(ctx, query, item.Title, item.Sku, item.Quantity)
	return row.Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
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

func (r *ItemRepository) Update(ctx context.Context, id int64, newItem *models.Item) error {
	query := `UPDATE items
			SET title = $1,
				sku = $2,
				quantity = $3,
				updated_at = NOW()
			WHERE id = $4`
	res, err := r.db.ExecContext(ctx, query, newItem.Title, newItem.Sku, newItem.Quantity, id)
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
	return nil
}

func (r *ItemRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM items
				WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
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
	return nil
}
