package handler

import (
	"context"
	"strconv"
	"time"

	"github.com/wb-go/wbf/ginext"
)

func DeleteComment(deps Deps) ginext.HandlerFunc {
	return func(c *ginext.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id <= 0 {
			c.JSON(400, ginext.H{"status": "error", "error": "id must be positive int"})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		const q = `DELETE FROM comments WHERE id = $1`
		res, err := deps.DB.ExecContext(ctx, q, id)
		if err != nil {
			c.JSON(500, ginext.H{"status": "error", "error": err.Error()})
			return
		}

		affected, err := res.RowsAffected()
		if err != nil {
			c.JSON(500, ginext.H{"status": "error", "error": err.Error()})
			return
		}
		if affected == 0 {
			c.JSON(404, ginext.H{"status": "error", "error": "not found"})
			return
		}

		c.JSON(200, ginext.H{"status": "deleted"})
	}
}
