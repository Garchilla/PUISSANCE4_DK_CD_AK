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
		Board:         Creation(),
		CurrentPlayer: "Rouge",
		GameOver:      false,
	}
)

type GameState struct {
	Board         [][]string
	CurrentPlayer string
	GameOver      bool
	Winner        string
}

func Creation() [][]string {
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
	if current == "Rouge" {
		return "Jaune"
	}
	return "Rouge"
}

func checkWin(board [][]string, ligne, colonne int, joueur string) bool {
	directions := [][2]int{
		{0, 1}, {1, 0}, {1, 1}, {1, -1},
	}
	for _, dir := range directions {
		count := 1
		for _, step := range []int{1, -1} {
			dx, dy := dir[0]*step, dir[1]*step
			x, y := ligne+dx, colonne+dy
			for x >= 0 && x < NBLignes && y >= 0 && y < NBColonnes && board[x][y] == joueur {
				count++
				if count >= 4 {
					return true
				}
				x += dx
				y += dy
			}
		}
	}
	return false
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))

	http.HandleFunc("/", handleAbout)
	http.HandleFunc("/action", handleAction)

	fmt.Println("Serveur démarré sur http://localhost:8081/")
	http.ListenAndServe(":8081", nil)

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

	mu.Lock()
	defer mu.Unlock()

	if gameState.GameOver {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.FormValue("Reset") == "true" {
		gameState.Board = Creation()
		gameState.CurrentPlayer = "Rouge"
	} else {
		colStr := r.FormValue("column")
		col, err := strconv.Atoi(colStr)
		if err != nil || col < 0 || col >= NBColonnes {
			http.Error(w, "Colonne invalide", http.StatusBadRequest)
			return
		}
		if ligne, ok := deposerJeton(gameState.Board, col, gameState.CurrentPlayer); ok {
			if checkWin(gameState.Board, ligne, col, gameState.CurrentPlayer) {
				fmt.Printf("Le joueur %s a gagné!\n", gameState.CurrentPlayer)
				gameState.CurrentPlayer = "Gagnant: " + gameState.CurrentPlayer
				gameState.GameOver = true
				gameState.Winner = gameState.CurrentPlayer

				tmpl, err := template.ParseFiles("victory.html")
				if err != nil {
					http.Error(w, "Template introuvable", http.StatusInternalServerError)
				}
				tmpl.Execute(w, gameState)
				return
			} else {
				gameState.CurrentPlayer = switchPlayer(gameState.CurrentPlayer)
			}
		}
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

func handleRematch(w http.ResponseWriter, r *http.Request) {
	gameState.Board = Creation()
	gameState.CurrentPlayer = "Rouge"
	gameState.GameOver = false
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
