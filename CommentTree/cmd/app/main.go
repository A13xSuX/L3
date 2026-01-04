package main

import (
	"fmt"
	"l3/CommentTree/internal/appcfg"
	"log"
)

func main() {

	cfg, err := appcfg.NewAppConfig()
	if err != nil {
		log.Fatal("Config doesnt load")
	}
	fmt.Println(cfg)
}
