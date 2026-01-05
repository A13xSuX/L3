package handler

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/wb-go/wbf/ginext"
)

func CommentsTree(deps Deps) ginext.HandlerFunc {
	return func(c *ginext.Context) {
		idStr := c.Query("id")
		if idStr == "" {
			c.JSON(400, ginext.H{"status": "error", "error": "id is required"})
			return
		}
		rootID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || rootID <= 0 {
			c.JSON(400, ginext.H{"status": "error", "error": "id must be positive int"})
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
		rows, err := deps.DB.QueryContext(ctx, q, rootID)
		if err != nil {
			c.JSON(500, ginext.H{"status": "error", "error": err.Error()})
			return
		}
		defer rows.Close()
		type node struct {
			ID        int64     `json:"id"`
			ParentID  *int64    `json:"parent_id"`
			Text      string    `json:"text"`
			CreatedAt time.Time `json:"created_at"`
			Children  []*node   `json:"children,omitempty"`
		}
		nodes := make(map[int64]*node)
		order := make([]int64, 0)

		for rows.Next() {
			var id int64
			var parent sql.NullInt64
			var text string
			var created time.Time
			if err := rows.Scan(&id, &parent, &text, &created); err != nil {
				c.JSON(500, ginext.H{"status": "error", "error": err.Error()})
				return
			}
			n := &node{
				ID:        id,
				Text:      text,
				CreatedAt: created,
			}
			if parent.Valid {
				v := parent.Int64
				n.ParentID = &v
			}
			nodes[id] = n
			order = append(order, id)
		}
		if err := rows.Err(); err != nil {
			c.JSON(500, ginext.H{"status": "error", "error": err.Error()})
			return
		}
		if len(order) == 0 {
			c.JSON(404, ginext.H{"status": "error", "error": "not found"})
			return
		}
		for _, id := range order {
			n := nodes[id]
			if n.ParentID == nil {
				continue
			}
			p := nodes[*n.ParentID]
			if p != nil {
				p.Children = append(p.Children, n)
			}
		}
		c.JSON(200, ginext.H{"status": "ok", "item": nodes[rootID]})
	}
}
