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
}
