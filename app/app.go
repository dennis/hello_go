package app

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/dennis/hello_go/context"
	"github.com/dennis/hello_go/handlers"
	"github.com/dennis/hello_go/models"
)

type App struct {
	Router  *mux.Router
	Context context.Context
}

func (a *App) Initialize() {
	a.setupRoutes()
	a.populateData()
}

func (a *App) setupRoutes() {
	a.Router = mux.NewRouter()

	a.Router.HandleFunc("/api/messages", a.handleRequest(handlers.GetMessages)).Methods("GET")
	a.Router.HandleFunc("/api/messages/{id}", a.handleRequest(handlers.GetMessage)).Methods("GET")
	a.Router.HandleFunc("/api/messages", a.handleRequest(handlers.CreateMessage)).Methods("POST")
	a.Router.HandleFunc("/api/messages/{id}", a.handleRequest(handlers.UpdateMessage)).Methods("PUT")
	a.Router.HandleFunc("/api/messages/{id}", a.handleRequest(handlers.DeleteMessage)).Methods("DELETE")
}

func (a *App) populateData() {
	a.Context = context.Context{}

	a.Context.MessageRepository.Insert(models.Message{ID: "1", Author: "Dennis", Topic: "Hello World", Body: "Lorem lipsum"})
	a.Context.MessageRepository.Insert(models.Message{ID: "2", Author: "Marianne", Topic: "re: Hello World", Body: "Really?"})

	a.Context.UserRepository.Insert(models.User{Username: "Dennis", AuthToken: "authtokendennis"})
	a.Context.UserRepository.Insert(models.User{Username: "Marianne", AuthToken: "authtokenmarianne"})
}

func (a *App) Run() {
	log.Fatal(http.ListenAndServe(":8080", a.Router))
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingresponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (l *loggingResponseWriter) WriteHeader(code int) {
	l.statusCode = code
	l.ResponseWriter.WriteHeader(code)
}

func (a *App) handleRequest(handler func(ctx *context.Context, w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(original_w http.ResponseWriter, r *http.Request) {
		w := newLoggingresponseWriter(original_w)

		start := time.Now()

		if user := handlers.Authenticate(&a.Context, r); user != nil {
			a.Context.CurrentUser = *user

			handler(&a.Context, w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}

		log.Printf(
			"%3d %-10s%-6s\t%s\t%s",
			w.statusCode,
			a.Context.CurrentUser.Username,
			r.Method,
			r.RequestURI,
			time.Since(start))
	}
}
