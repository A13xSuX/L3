package handler

import (
	"context"
	"time"

	"github.com/wb-go/wbf/ginext"
)

type createCommentReq struct {
	ParentID *int64 `json:"parent_id"`
	Text     string `json:"text"`
}

func CreateComment(deps Deps) ginext.HandlerFunc {
	return func(c *ginext.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var req createCommentReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, ginext.H{"status": "error", "error": err.Error()})
			return
		}
		if req.Text == "" {
			c.JSON(400, ginext.H{"status": "error", "error": "text is required"})
			return
		}

		const q = `
				INSERT INTO comments(parent_id, text)
				VALUES ($1,$2)
				RETURNING id, created_at`

		var id int64
		var createdAt time.Time
		if err := deps.DB.Master.QueryRowContext(ctx, q, req.ParentID, req.Text).Scan(&id, &createdAt); err != nil {
			c.JSON(500, ginext.H{"status": "error", "error": err.Error()})
			return
		}
		c.JSON(201, ginext.H{"status": "created", "id": id, "created_at": createdAt})
	}
}
