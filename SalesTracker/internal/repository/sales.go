package repository

import (
	"context"
	"database/sql"
	"errors"
	"l3/SalesTracker/internal/models"

	"github.com/wb-go/wbf/dbpg"
)

type SalesRepo struct {
	db *dbpg.DB
}

func NewSalesRepo(db *dbpg.DB) *SalesRepo {
	return &SalesRepo{
		db: db,
	}
}

func (r *SalesRepo) Create(ctx context.Context, sale *models.Sale) error {
	query := `INSERT INTO sales (title, category, price, quantity, sale_date)
			  VALUES ($1,$2,$3,$4,$5)
			  RETURNING id`
	row := r.db.QueryRowContext(ctx, query,
		sale.Title,
		sale.Category,
		sale.Price,
		sale.Quantity,
		sale.SaleDate)
	return row.Scan(&sale.ID)
}

func (r *SalesRepo) GetByID(ctx context.Context, id string) (*models.Sale, error) {
	query := `SELECT id, title, category, price, quantity, sale_date FROM sales
              WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var sale models.Sale
	err := row.Scan(
		&sale.ID,
		&sale.Title,
		&sale.Category,
		&sale.Price,
		&sale.Quantity,
		&sale.SaleDate,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			//TODO может сделать ошибку на отсутствие чтобы передавать nil, newerr
			return nil, nil
		}
		return nil, err
	}
	return &sale, nil
}

func (r *SalesRepo) GetAll(ctx context.Context) ([]models.Sale, error) {
	query := `SELECT id, title, category, price, quantity, sale_date FROM sales`

	var allSales []models.Sale
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var sale models.Sale

		err := rows.Scan(
			&sale.ID,
			&sale.Title,
			&sale.Category,
			&sale.Price,
			&sale.Quantity,
			&sale.SaleDate,
		)
		if err != nil {
			return nil, err
		}
		allSales = append(allSales, sale)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return allSales, nil
}

func (r *SalesRepo) Update(ctx context.Context, newSale *models.Sale) error {
	query := `UPDATE sales
			  SET title = $1,
				  category = $2,
				  price = $3,
				  quantity = $4,
				  sale_date = $5
			  WHERE id = $6 
    `
	res, err := r.db.ExecContext(ctx, query,
		newSale.Title,
		newSale.Category,
		newSale.Price,
		newSale.Quantity,
		newSale.SaleDate,
		newSale.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		//TODO customErrs
		return errors.New("sales id  not found")
	}
	return nil
}

func (r *SalesRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM sales
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
		return errors.New("sales id  not found")
	}
	return nil
}
