package main

import (
	"fmt"
	"html/template"
	"net/http"
	"sync"
)

const NBLignes = 6
const NBColonnes = 7

var (
	joueur  = "R"
	mu      sync.Mutex
	board   [NBLignes][NBColonnes]string
	colonne int
	ligne   int
	counter int
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
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))

	http.HandleFunc("/", handleAbout)

	http.HandleFunc("/action", handleAction)

	fmt.Println("Serveur démarré sur http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}

func handleAbout(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("about.html"))
	tmpl.Execute(w, map[string]interface{}{
		"Count": counter,
	})
}

func handleAction(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && joueur == "R" {
		mu.Lock()
		counter++
		joueur = "J"
		mu.Unlock()
		http.Redirect(w, r, "/", http.StatusSeeOther)
		fmt.Println("Jeton Rouge ajouté")
	} else if r.Method == http.MethodPost && joueur == "J" {
		mu.Lock()
		counter++
		joueur = "R"
		mu.Unlock()
		http.Redirect(w, r, "/", http.StatusSeeOther)
		fmt.Println("Jeton Jaune ajouté")
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
