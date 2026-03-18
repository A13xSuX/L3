package repository

import (
	"context"
	"l3/WarehouseControl/internal/models"

	"github.com/wb-go/wbf/dbpg"
)

type AuditRepository struct {
	db *dbpg.DB
}

func NewAuditRepository(db *dbpg.DB) *AuditRepository {
	return &AuditRepository{
		db: db,
	}
}

func (r *AuditRepository) GetByItemID(ctx context.Context, itemID int64) ([]models.Audit, error) {
	query := `SELECT 
       id,
       item_id,
       action,
       COALESCE(old_data, 'null'::jsonb),
       COALESCE(new_data, 'null'::jsonb),
       changed_by_user_id,
       changed_by_username,
       changed_by_role,
       changed_at 
			FROM audit
			WHERE item_id = $1
			ORDER BY changed_at DESC`
	var audit []models.Audit
	rows, err := r.db.QueryContext(ctx, query, itemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var row models.Audit
		err = rows.Scan(
			&row.ID,
			&row.ItemID,
			&row.Action,
			&row.OldData,
			&row.NewData,
			&row.ChangedByUserID,
			&row.ChangedByUsername,
			&row.ChangedByRole,
			&row.ChangedAt,
		)
		if err != nil {
			return nil, err
		}
		audit = append(audit, row)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return audit, nil
}
