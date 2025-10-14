package main

import (
	"fmt"
	"html/template"
	"net/http"
	"sync"
)

const (
	rows    = 6
	columns = 7
)

var (
	board  = make([][]string, rows)
	player = "R"
	mu     sync.Mutex
)

func init() {
	for i := range board {
		board[i] = make([]string, columns)
	}
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("about.html"))
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/action", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && player == "R" {
			fmt.Println("Faire apparaitre en rouge")
			w.Write([]byte("Action effectuée côté serveur ! "))
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Serveur démarré sur http://localhost:8081/")
	http.ListenAndServe(":8081", nil)
}
