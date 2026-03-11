package repository

import (
	"context"
	"database/sql"
	"errors"
	"l3/EventBooker/internal/customErrs"
	"l3/EventBooker/internal/models"
	"time"

	"github.com/wb-go/wbf/dbpg"
)

type BookingRepository struct {
	db *dbpg.DB
}

func NewBookingRepo(db *dbpg.DB) *BookingRepository {
	return &BookingRepository{
		db: db,
	}
}
func (r *BookingRepository) Create(ctx context.Context, booking models.Booking) error {
	query := `INSERT INTO bookings (event_id, username, status, created_at, expired_at)
			VALUES ($1, $2, $3, $4,$5)
			RETURNING id`
	//created_at = time.Now, expired_at = time.Now().Add(time.Minute*15)
	//status = pending or confirmed(if without price)
	row := r.db.QueryRowContext(ctx, query,
		booking.EventID,
		booking.Username,
		booking.Status,
		booking.CreatedAt,
		booking.ExpiredAt)
	return row.Scan(&booking.ID)
}

func (r *BookingRepository) GetByID(ctx context.Context, id string) (*models.Booking, error) {
	query := `SELECT id, event_id, username, status, created_at, expired_at, confirmed_at 
				FROM bookings WHERE id = $1`

	var booking models.Booking
	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&booking.ID,
		&booking.EventID,
		&booking.Username,
		&booking.Status,
		&booking.CreatedAt,
		&booking.ExpiredAt,
		&booking.ConfirmedAt,
	)
	if err != nil {
		//TODO err?
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &booking, nil
}

func (r *BookingRepository) UpdateStatus(ctx context.Context, id string, status string, confirmedAt *time.Time) error {
	query := `UPDATE bookings
			SET status = $1, confirmed_at = $2
			WHERE id = $3`

	result, err := r.db.ExecContext(ctx, query, status, confirmedAt, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		//TODO может кастомную реализовать
		return sql.ErrNoRows
	}

	return err
}

func (r *BookingRepository) GetExpiredBookings(ctx context.Context) ([]models.Booking, error) {
	query := `SELECT id, event_id, username, status, created_at, expired_at, confirmed_at 
			FROM bookings 
			WHERE status = 'pending' AND expired_at < NOW()`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expiredBooking []models.Booking

	for rows.Next() {
		var booking models.Booking
		err = rows.Scan(
			&booking.ID,
			&booking.EventID,
			&booking.Username,
			&booking.Status,
			&booking.CreatedAt,
			&booking.ExpiredAt,
			&booking.ConfirmedAt,
		)
		if err != nil {
			return expiredBooking, err
		}
		expiredBooking = append(expiredBooking, booking)
	}
	if err = rows.Err(); err != nil {
		return expiredBooking, err
	}
	return expiredBooking, nil
}

func (r *BookingRepository) CreateTx(ctx context.Context, tx *sql.Tx, booking *models.Booking) error {
	query := `INSERT INTO bookings (event_id, username, status, created_at, expired_at )
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
`
	row := tx.QueryRowContext(ctx, query,
		booking.EventID,
		booking.Username,
		booking.Status,
		booking.CreatedAt,
		booking.ExpiredAt,
	)
	return row.Scan(&booking.ID)
}

func (r *BookingRepository) GetByIDForUpdateTx(ctx context.Context, tx *sql.Tx, id string) (*models.Booking, error) {
	query := `SELECT id, event_id, username, created_at, expired_at,confirmed_at FROM bookings
			WHERE id = $1
			FOR UPDATE`
	var booking models.Booking
	row := tx.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&booking.ID,
		&booking.EventID,
		&booking.Username,
		&booking.CreatedAt,
		&booking.ExpiredAt,
		&booking.ConfirmedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &booking, nil
}

func (r *BookingRepository) UpdateStatusTx(ctx context.Context, tx *sql.Tx, id string, status string, confirmedAt *time.Time) error {
	query := `UPDATE bookings
			SET status = $1, confirmed_at = $2, expired_at = NULL
			WHERE id = $3`

	result, err := tx.ExecContext(ctx, query, status, confirmedAt, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		//TODO может кастомную реализовать
		return customErrs.ErrBookingNotFound
	}

	return nil
}
