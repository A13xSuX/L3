package handler

import (
	"context"
	"time"

	"github.com/wb-go/wbf/ginext"
)

func Health(deps Deps) ginext.HandlerFunc {
	return func(c *ginext.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var res int
		err := deps.DB.QueryRowContext(ctx, "select 1").Scan(&res)
		if err != nil {
			c.JSON(500, ginext.H{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}
		c.JSON(200, ginext.H{"status": "ok"})
	}
}
