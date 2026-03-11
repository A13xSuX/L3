package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"l3/EventBooker/internal/customErrs"
	"l3/EventBooker/internal/models"

	"github.com/wb-go/wbf/dbpg"
)

type EventRepository struct {
	db *dbpg.DB
}

func NewEventRepository(db *dbpg.DB) *EventRepository {
	return &EventRepository{
		db: db,
	}
}

func (r *EventRepository) Create(ctx context.Context, event *models.Event) error {
	query := `INSERT INTO events (title, description, date, total_seats, available_seats,
    price, payment_required) 
	VALUES ($1, $2, $3, $4, $5, $6, $7) 
	RETURNING id, created_at`

	row := r.db.QueryRowContext(ctx, query, event.Title, event.Description, event.Date,
		event.TotalSeats, event.AvailableSeats, event.Price, event.PaymentRequired)

	err := row.Scan(&event.ID, &event.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *EventRepository) GetByID(ctx context.Context, id string) (*models.Event, error) {
	query := `SELECT id, title, description, date, total_seats,
       available_seats, price, payment_required, created_at 
		FROM events WHERE id = $1`

	var event models.Event
	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&event.ID,
		&event.Title,
		&event.Description,
		&event.Date,
		&event.TotalSeats,
		&event.AvailableSeats,
		&event.Price,
		&event.PaymentRequired,
		&event.CreatedAt,
	)
	if err != nil {
		//TODO check exactly err
		if errors.Is(err, errors.New("sql: no rows in result set")) {
			return nil, nil
		}
		return nil, err
	}
	return &event, nil
}

func (r *EventRepository) GetAll(ctx context.Context) ([]models.Event, error) {
	query := `SELECT id, title, description, date, total_seats,
       available_seats, price, payment_required, created_at 
		FROM events
		ORDER BY date`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		err = rows.Scan(&event.ID,
			&event.Title,
			&event.Description,
			&event.Date,
			&event.TotalSeats,
			&event.AvailableSeats,
			&event.Price,
			&event.PaymentRequired,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func (r *EventRepository) UpdateSeats(ctx context.Context, eventID string, newSeats int) error {
	query := `UPDATE events
			  SET available_seats = available_seats + $1
			  WHERE  id = $2
			  AND available_seats + $1 >= 0
			  AND available_seats + $1 <= total_seats`
	result, err := r.db.ExecContext(ctx, query, newSeats, eventID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no seats updated")
	}
	return nil
}

// for tx
func (r *EventRepository) GetByIDForUpdateTx(ctx context.Context, tx *sql.Tx, id string) (*models.Event, error) {
	query := `SELECT  id, title, description, date, total_seats,
       available_seats, price, payment_required, created_at  FROM events
	   WHERE id = $1
       FOR UPDATE
`
	row := tx.QueryRowContext(ctx, query, id)

	var event models.Event
	err := row.Scan(&event.ID,
		&event.Title,
		&event.Description,
		&event.Date,
		&event.TotalSeats,
		&event.AvailableSeats,
		&event.Price,
		&event.PaymentRequired,
		&event.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			//TODO в остальных поменять
			return nil, customErrs.ErrEventNotFound
		}
		return nil, err
	}
	return &event, nil
}

func (r *EventRepository) UpdateSeatsTx(ctx context.Context, tx *sql.Tx, eventID string, newSeats int) error {
	query := `UPDATE events
			  SET available_seats = available_seats + $1
			  WHERE  id = $2
			  AND available_seats + $1 >= 0
			  AND available_seats + $1 <= total_seats`
	result, err := tx.ExecContext(ctx, query, newSeats, eventID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return customErrs.ErrNoAvailableSeats
	}
	return nil
}
