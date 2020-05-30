package app

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/dennis/hello_go/handlers"
	"github.com/dennis/hello_go/context"
	"github.com/dennis/hello_go/models"
)

type App struct {
	Router *mux.Router
	Context context.Context
}

func (a *App) Initialize() {
	a.Router = mux.NewRouter()

	a.Router.HandleFunc("/api/messages", a.handleRequest(handlers.GetMessages)).Methods("GET")
	a.Router.HandleFunc("/api/messages/{id}", a.handleRequest(handlers.GetMessage)).Methods("GET")
	a.Router.HandleFunc("/api/messages", a.handleRequest(handlers.CreateMessage)).Methods("POST")
	a.Router.HandleFunc("/api/messages/{id}", a.handleRequest(handlers.UpdateMessage)).Methods("PUT")
	a.Router.HandleFunc("/api/messages/{id}", a.handleRequest(handlers.DeleteMessage)).Methods("DELETE")

	a.Context = context.Context {}
	a.Context.Messages = append(a.Context.Messages, models.Message{ID: "1", Topic: "Hello World", Body: "Lorem lipsum"})
	a.Context.Messages = append(a.Context.Messages, models.Message{ID: "2", Topic: "re: Hello World", Body: "Really?"})
}

func (a *App) Run() {
	log.Fatal(http.ListenAndServe(":8080", a.Router))
}

func (a *App) handleRequest(handler func(ctx *context.Context, w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	// TODO: Authentication/Authorization
	return func(w http.ResponseWriter, r *http.Request) {
		handler(&a.Context, w, r)
	}
}
