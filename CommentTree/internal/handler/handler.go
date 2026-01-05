package handler

import (
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/ginext"
)

type Deps struct {
	DB *dbpg.DB
}

func Register(r *ginext.Engine, deps Deps) {
	r.GET("/health", Health(deps))
	r.POST("/comments", CreateComment(deps))
	r.GET("/comments", GetCommentsTree(deps))
	r.DELETE("/comments/:id", DeleteComment(deps))
	r.GET("/comments/search", SearchComments(deps))
	r.GET("/comments/tree", CommentsTree(deps))
}
