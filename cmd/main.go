package main

import (
	"portfolyo/internal/infrastructure/app"
	"portfolyo/internal/router"
)

func main() {

	r := router.NewRouter()
	a := app.New(r)
	a.Start()
}
