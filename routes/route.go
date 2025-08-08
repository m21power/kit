package routes

import (
	"log"
	"net/http"

	handler "kit/cmd/handlers"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Router struct {
	route *mux.Router
}

func NewRouter(r *mux.Router) *Router {
	return &Router{route: r}
}

// Register the new routes for testing
func (r *Router) RegisterRoute() {
	path := r.route.PathPrefix("/api/v1").Subrouter()
	path.HandleFunc("/init", handler.InitHandler).Methods("POST")
	path.HandleFunc("/add", handler.AddHandler).Methods("POST")
	path.HandleFunc("/commit", handler.CommitHandler).Methods("POST")
	path.HandleFunc("/log", handler.LogHandler).Methods("POST")
	path.HandleFunc("/status", handler.StatusHandler).Methods("POST")
	path.HandleFunc("/restore", handler.RestoreHandler).Methods("POST")
	path.HandleFunc("/branch", handler.CreateBranch).Methods("POST")
	path.HandleFunc("/checkout", handler.CheckoutBranch).Methods("POST")
	path.HandleFunc("/branches", handler.ListBranches).Methods("GET")
	path.HandleFunc("/dir", handler.GetAllDir).Methods("GET")
	path.HandleFunc("/check", handler.CheckDir).Methods("GET")
	path.HandleFunc("/reset", handler.ResetHandler).Methods("POST")

}

func (r *Router) Run(addr string) error {
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},                                       // Allow all origins
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Allowed methods
		AllowedHeaders: []string{"Content-Type", "Authorization"},           // Allowed headers
	})

	handler := corsHandler.Handler(r.route)

	log.Println("Server running on port: ", addr)
	return http.ListenAndServe(addr, handler)
}
