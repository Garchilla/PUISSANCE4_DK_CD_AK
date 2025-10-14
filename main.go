package main

import (
	"fmt"
	"html/template"
	"net/http"
	"sync"
)

func main() {
	r := router.New()

	fmt.Println("Serveur démarré sur http://localhost:8080")
	http.ListenAndServe(":8080", r)
}

const NBLignes = 6
const NBColonnes = 7

var (
	joueur  = "R"
	mu      sync.Mutex
	board   [NBLignes][NBColonnes]string
	colonne int
	ligne   int
)

func deposerJeton(board *[NBLignes][NBColonnes]string, colonne int, joueur string) (int, bool) {
	if colonne < 0 || colonne >= NBColonnes {
		return -1, false
	}

	for ligne := NBLignes - 1; ligne >= 0; ligne-- {
		if board[ligne][colonne] == "." {
			board[ligne][colonne] = joueur
			return ligne, true
		}
	}
	return -1, false
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
			player = "J"
		} else if r.Method == http.MethodPost && player == "J" {
			fmt.Println("Faire apparaitre en jaune")
			w.Write([]byte("Action effectuée côté serveur"))
			player = "R"
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Serveur démarré sur http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
