package app

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/dennis/hello_go/context"
	"github.com/dennis/hello_go/handlers"
	"github.com/dennis/hello_go/repositories"
	"github.com/dennis/hello_go/services"
)

// A wrapper for ResponseWriter, that also captures the StatusCode. Used for logging
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (l *loggingResponseWriter) WriteHeader(code int) {
	l.statusCode = code
	l.ResponseWriter.WriteHeader(code)
}

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
	messageRepository := repositories.MessageRepository{}
	userRepository := repositories.UserRepository{}

	PopulateMessages(&messageRepository)
	PopulateUsers(&userRepository)

	a.Context = context.Context{
		AuthenticationService: services.AuthenticationService{UserRepository: &userRepository},
		MessageService:        services.MessageService{MessageRepository: &messageRepository},
	}
}

func (a *App) Run() {
	log.Fatal(http.ListenAndServe(":8080", a.Router))
}

// This dispatches a request to a Handler as configured in setupRoutes.
// It performs a number of tasks:
// 1) It logs the request, its duration and statuscode
// 2) It performs Authentication and only allows authenticated requests to
//    reach our handlers
// 3) It provides Context, ResponseWriter, Request and our URL vars to the handler
func (a *App) handleRequest(handler func(ctx *context.Context, session *context.Session, w http.ResponseWriter, r *http.Request, vars map[string]string)) http.HandlerFunc {
	return func(original_w http.ResponseWriter, r *http.Request) {
		w := newLoggingResponseWriter(original_w)

		// To avoid that mux leaks into the handlers, we capture any
		// variables it as, and provide them as a plain map[string]string
		// /api/message/{id} in the route will result in a vars[id] = <value>
		vars := mux.Vars(r)

		start := time.Now()
		username := "unknown"

		if user := handlers.Authenticate(&a.Context, r); user != nil {
			session := context.Session{CurrentUser: *user}

			username = user.Username

			handler(&a.Context, &session, w, r, vars)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}

		log.Printf(
			"%3d %-10s%-6s\t%s\t%s",
			w.statusCode,
			username,
			r.Method,
			r.RequestURI,
			time.Since(start))
	}
}
