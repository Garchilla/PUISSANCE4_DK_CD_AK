package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"sync"
)

const NBLignes = 6
const NBColonnes = 7

var (
	mu        sync.Mutex
	gameState = GameState{
		Board:         makeBoard(),
		CurrentPlayer: "red",
	}
)

type GameState struct {
	Board         [][]string
	CurrentPlayer string
}

func makeBoard() [][]string {
	board := make([][]string, NBLignes)
	for i := range board {
		board[i] = make([]string, NBColonnes)
	}
	return board
}

func deposerJeton(board [][]string, colonne int, joueur string) (int, bool) {
	if colonne < 0 || colonne >= NBColonnes {
		return -1, false
	}
	for ligne := NBLignes - 1; ligne >= 0; ligne-- {
		if board[ligne][colonne] == "" {
			board[ligne][colonne] = joueur
			return ligne, true
		}
	}
	return -1, false
}

func switchPlayer(current string) string {
	if current == "red" {
		return "yellow"
	}
	return "red"
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
	tmpl.Execute(w, gameState)
}

func handleAction(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Erreur lors de l'analyse du formulaire", http.StatusBadRequest)
		return
	}

	colStr := r.FormValue("column")
	col, err := strconv.Atoi(colStr)
	if err != nil || col < 0 || col >= NBColonnes {
		http.Error(w, "Colonne invalide", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if _, ok := deposerJeton(gameState.Board, col, gameState.CurrentPlayer); ok {
		gameState.CurrentPlayer = switchPlayer(gameState.CurrentPlayer)
	}

	tmpl, err := template.ParseFiles("about.html")
	if err != nil {
		http.Error(w, "Template introuvable", http.StatusInternalServerError)
		fmt.Println("Parse error:", err)
		return
	}
	err = tmpl.Execute(w, gameState)
	if err != nil {
		http.Error(w, "Erreur lors du rendu du template", http.StatusInternalServerError)
		fmt.Println("Execution error:", err)
	}
}
