package handlers

import (
	"net/http"

	"github.com/wb-go/wbf/ginext"
)

func Template() ginext.HandlerFunc {
	return func(c *ginext.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	}
}
