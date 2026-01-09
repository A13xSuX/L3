package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/wb-go/wbf/dbpg"
)

type ImagesRepo struct {
	db *dbpg.DB
}

type Image struct {
	ID            uuid.UUID
	Status        string
	OriginalPath  string
	ProcessedPath *string
	ThumbPath     *string
	Error         *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

var ErrNotFound = errors.New("image not found")

func NewImagesRepo(db *dbpg.DB) *ImagesRepo {
	return &ImagesRepo{db: db}
}

func (r *ImagesRepo) Create(ctx context.Context, id uuid.UUID, originalPath string) error {
	const q = `
	INSERT INTO images (id, status, original_path)
	VALUES ($1, 'queued', $2)
`
	_, err := r.db.ExecContext(ctx, q, id, originalPath)
	return err
}

func (r *ImagesRepo) Get(ctx context.Context, id uuid.UUID) (*Image, error) {
	const q = `SELECT id, status, original_path, processed_path, thumb_path,error, created_at, updated_at 
			   FROM images
			   WHERE id = $1`
	var img Image
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&img.ID,
		&img.Status,
		&img.OriginalPath,
		&img.ProcessedPath,
		&img.ThumbPath,
		&img.Error,
		&img.CreatedAt,
		&img.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &img, nil
}

func (r *ImagesRepo) MarkProcessing(ctx context.Context, id uuid.UUID) error {
	const q = `UPDATE images
		       SET status = 'processing', updated_at = now()
			   WHERE id = $1`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
func (r *ImagesRepo) MarkReady(ctx context.Context, id uuid.UUID, processedPath string, thumbPath string) error {
	const q = `UPDATE images
		       SET status = 'ready',
		    	processed_path = $2,
		    	thumb_path = $3,
		    	error = NULL,
		    	updated_at = now()
			   WHERE id = $1`
	res, err := r.db.ExecContext(ctx, q, id, processedPath, thumbPath)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *ImagesRepo) MarkFailed(ctx context.Context, id uuid.UUID, errText string) error {
	const q = `UPDATE images
SET status = 'failed',
    error = $2,
    updated_at = now()
    WHERE id = $1	`
	res, err := r.db.ExecContext(ctx, q, id, errText)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()

	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *ImagesRepo) Delete(ctx context.Context, id uuid.UUID) error {
	const q = `DELETE FROM images WHERE id = $1`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
