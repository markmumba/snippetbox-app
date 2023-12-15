package main

import "net/http"

func (app *application) routes() http.Handler {

	router := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	router.Handle("/static/", http.StripPrefix("/static", fileServer))

	router.HandleFunc("/", app.home)
	router.HandleFunc("/snippet/view", app.snippetView)
	router.HandleFunc("/snippet/create", app.snippetCreate)

	return secureHeaders(router)
}
