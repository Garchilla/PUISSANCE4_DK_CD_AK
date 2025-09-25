package router

import (
	"net/http"
	"power4/controller"
)

func New() *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("/", controller.Home)
	mux.HandleFunc("/about", controller.About)
	mux.HandleFunc("/contact", controller.Contact)

	return mux
}
