package main

import (
	"fmt"
	"l3/EventBooker/internal/appcfg"
)

func main() {
	cfg, err := appcfg.NewAppConfig()
	if err != nil {
		fmt.Printf("%w", err)
		return
	}
	fmt.Println(cfg)
}
