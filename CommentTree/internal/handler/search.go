package handler

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/wb-go/wbf/ginext"
)

func SearchComments(deps Deps) ginext.HandlerFunc {
	return func(c *ginext.Context) {
		q := c.Query("q")
		if q == "" {
			c.JSON(400, ginext.H{"status": "error", "error": "q is required"})
			return
		}
		limit := int64(20)
		if s := c.Query("limit"); s != "" {
			v, err := strconv.ParseInt(s, 10, 64)
			if err != nil || v <= 0 {
				c.JSON(400, ginext.H{"status": "error", "error": "limit must be positive int"})
				return
			}
			if v > 100 { //max limit
				v = 100
			}
			limit = v
		}
		offset := int64(0)
		if s := c.Query("offset"); s != "" {
			v, err := strconv.ParseInt(s, 10, 64)
			if err != nil || v < 0 {
				c.JSON(400, ginext.H{"status": "error", "error": "offset must be >= 0"})
				return
			}
			offset = v
		}
		sort := c.Query("sort")
		orderBy := "rank DESC, created_at DESC, id DESC"
		switch sort {
		case "", "rank":
			orderBy = "rank DESC, created_at DESC, id DESC"
		case "created_at":
			orderBy = "created_at DESC, id DESC"
		default:
			c.JSON(400, ginext.H{"status": "error", "error": "invalid sort"})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		query := fmt.Sprintf(`SELECT id, parent_id, text, created_at,
ts_rank_cd(to_tsvector('simple', text), websearch_to_tsquery('simple', $1)) AS rank
      FROM comments
      WHERE to_tsvector('simple', text) @@ websearch_to_tsquery('simple', $1)
      ORDER BY %s
      LIMIT $2 OFFSET $3
    `, orderBy)
		rows, err := deps.DB.QueryContext(ctx, query, q, limit, offset)
		if err != nil {
			c.JSON(500, ginext.H{"status": "error", "error": err.Error()})
			return
		}
		defer rows.Close()
		type item struct {
			commentRow
			Rank float64 `json:"rank"`
		}
		out := make([]item, 0, limit)
		for rows.Next() {
			var r item
			var parent sql.NullInt64
			if err := rows.Scan(&r.ID, &parent, &r.Text, &r.CreatedAt, &r.Rank); err != nil {
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

		c.JSON(200, ginext.H{"status": "ok", "items": out, "limit": limit, "offset": offset})
	}
}
