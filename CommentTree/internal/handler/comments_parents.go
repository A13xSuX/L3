package handler

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/wb-go/wbf/ginext"
)

type commentRow struct {
	ID        int64     `json:"id"`
	ParentID  *int64    `json:"parent_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

func GetCommentsTree(deps Deps) ginext.HandlerFunc {
	return func(c *ginext.Context) {
		parentStr := c.Query("parent")
		if parentStr == "" {
			c.JSON(400, ginext.H{"status": "error", "error": "parent is required"})
			return
		}
		parentID, err := strconv.ParseInt(parentStr, 10, 64)
		if err != nil || parentID <= 0 {
			c.JSON(400, ginext.H{"status": "error", "error": "parent must be positive int"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		const q = `
				WITH RECURSIVE tree AS (
				  SELECT id, parent_id, text, created_at
				  FROM comments
				  WHERE id = $1
				
				  UNION ALL
				
				  SELECT c.id, c.parent_id, c.text, c.created_at
				  FROM comments c
				  JOIN tree t ON c.parent_id = t.id
				)
				SELECT id, parent_id, text, created_at
				FROM tree
				ORDER BY created_at, id;`
		rows, err := deps.DB.QueryContext(ctx, q, parentID)
		if err != nil {
			c.JSON(500, ginext.H{"status": "error", "error": err.Error()})
			return
		}
		defer rows.Close()

		out := make([]commentRow, 0, 16)
		for rows.Next() {
			var r commentRow
			var parent sql.NullInt64
			if err := rows.Scan(&r.ID, &parent, &r.Text, &r.CreatedAt); err != nil {
				c.JSON(500, ginext.H{"status": "error", "error": err.Error()})
				return
			}
			if parent.Valid {
				v := parent.Int64
				r.ParentID = &v
			}
			out = append(out, r)
		}
		if err := rows.Err(); err != nil {
			c.JSON(500, ginext.H{"status": "error", "error": err.Error()})
			return
		}

		// если id не найден, вернётся пустой список
		if len(out) == 0 {
			c.JSON(404, ginext.H{"status": "error", "error": "not found"})
			return
		}

		c.JSON(200, ginext.H{"status": "ok", "items": out})
	}
}
