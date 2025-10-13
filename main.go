package main

import (
	"fmt"
	"net/http"
	"power4/router"
	"sync"
)

func main() {
	r := router.New()

	fmt.Println("Serveur démarré sur http://localhost:8080")
	http.ListenAndServe(":8080", r)
}

const (
	rows    = 6
	columns = 7
)

var (
	board  = make([][]string, rows)
	player = "R"
	mu     sync.Mutex
)
