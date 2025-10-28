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
		CurrentPlayer: "Pizza",
		GameOver:      false,
		Tour:          1,
	}
)

type GameState struct {
	Board         [][]string
	CurrentPlayer string
	GameOver      bool
	Winner        string
	Tour          int
}

func Creation() [][]string {
	board := make([][]string, NBLignes)
	for i := range board {
		board[i] = make([]string, NBColonnes)
	}
	return board
}

func deposerJeton(board [][]string, colonne int, joueur string) (int, bool) {
	joueur = gameState.CurrentPlayer
	if colonne < 0 || colonne >= NBColonnes {
		return -1, false
	}
	for ligne := NBLignes - 1; ligne >= 0; ligne-- {
		if board[ligne][colonne] == "" {
			if gameState.CurrentPlayer == "Pizza" {
				board[ligne][colonne] = joueur
				gameState.Tour++
			} else if gameState.CurrentPlayer == "Burger" {
				board[ligne][colonne] = joueur
				gameState.Tour++
			}
			return ligne, true
		}
	}
	return -1, false
}

func switchPlayer(current string) string {
	if current == "Pizza" {
		return "Burger"
	}
	return "Pizza"
}

func checkWin(board [][]string, ligne, colonne int, joueur string) bool {
	directions := [][2]int{
		{0, 1}, {1, 0}, {1, 1}, {1, -1},
	}
	joueur = gameState.CurrentPlayer
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

	http.HandleFunc("/", handleMenu)
	http.HandleFunc("/about", handleAbout)
	http.HandleFunc("/action", handleAction)

	fmt.Println("Serveur démarré sur http://localhost:8084/")
	http.ListenAndServe(":8084", nil)

}

func handleMenu(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("victoire.html"))
	tmpl.Execute(w, nil)
}

func handleAbout(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("About.html"))
	tmpl.Execute(w, gameState)
}

func handleAction(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Erreur lors de l'analyse du formulaire", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if r.FormValue("Reset") == "true" {
		gameState.Board = Creation()
		gameState.CurrentPlayer = "Pizza"
		gameState.Tour = 1
		gameState.GameOver = false
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if gameState.GameOver {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodPost {
		colStr := r.FormValue("column")
		col, err := strconv.Atoi(colStr)
		if err != nil || col < 0 || col >= NBColonnes {
			http.Error(w, "Colonne invalide", http.StatusBadRequest)
			return
		}
		if ligne, ok := deposerJeton(gameState.Board, col, gameState.CurrentPlayer); ok {
			if gameState.Tour == NBLignes*NBColonnes+1 && !checkWin(gameState.Board, ligne, col, gameState.CurrentPlayer) {
				fmt.Println("Match nul!")
				gameState.CurrentPlayer = "MATCH NUL"
				gameState.GameOver = true
				gameState.Winner = gameState.CurrentPlayer

				tmpl, err := template.ParseFiles("victory.html")
				if err != nil {
					http.Error(w, "Template introuvable", http.StatusInternalServerError)
				}
				tmpl.Execute(w, gameState)
				return
			}
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
