package main

import (
	"log"

	"learning-growth-platform/internal/http/router"
)

func main() {
	r := router.NewRouter(nil)
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
