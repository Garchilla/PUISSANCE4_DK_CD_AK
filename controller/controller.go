package controller

import (
	"html/template"
	"net/http"
)

func renderTemplate(w http.ResponseWriter, filename string, data map[string]string) {
	tmpl := template.Must(template.ParseFiles("template/" + filename))
	tmpl.Execute(w, data)
}

func Home(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Title":   "PUISSANCE 4",
		"Message": "Bienvenue au jeu de Puissance 4 en ligne!",
	}
	renderTemplate(w, "index.html", data)
}

func About(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Title":   "A propos",
		"Message": "Ceci est la page à propos",
	}
	renderTemplate(w, "about.html", data)
}

func Contact(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		msg := r.FormValue("msg")

		data := map[string]string{
			"Title":   "Contact",
			"Message": "Merci " + name + " pour votre message : " + msg,
		}
		renderTemplate(w, "contact.html", data)
		return
	}

	data := map[string]string{
		"Title":   "Contact",
		"Message": "Envoie-nous un message!",
	}
	renderTemplate(w, "contact.html", data)
}
