package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Router struct {
	route *mux.Router
}

func NewRouter(r *mux.Router) *Router {
	return &Router{route: r}
}

// Handler to create nested directory and add content
func CreateNestedDirHandler(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		DirPath  string `json:"dir_path"`
		FileName string `json:"file_name"`
		Content  string `json:"content"`
	}
	var req Req
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request"))
		return
	}
	if err := os.MkdirAll(req.DirPath, 0755); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to create directory"))
		return
	}
	filePath := filepath.Join(req.DirPath, req.FileName)
	if err := os.WriteFile(filePath, []byte(req.Content), 0644); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to write file"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Directory and file created successfully"))
}

// Handler to list directory structure
func ListDirStructureHandler(w http.ResponseWriter, r *http.Request) {
	type Node struct {
		Name     string `json:"name"`
		IsDir    bool   `json:"is_dir"`
		Children []Node `json:"children,omitempty"`
	}
	var walk func(string) Node
	walk = func(path string) Node {
		info, err := os.Stat(path)
		if err != nil {
			return Node{Name: path, IsDir: false}
		}
		node := Node{Name: info.Name(), IsDir: info.IsDir()}
		if info.IsDir() {
			entries, _ := os.ReadDir(path)
			for _, entry := range entries {
				childPath := filepath.Join(path, entry.Name())
				node.Children = append(node.Children, walk(childPath))
			}
		}
		return node
	}
	root := r.URL.Query().Get("root")
	if root == "" {
		root = "." // default to current dir
	}
	structure := walk(root)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(structure)
}

// Register the new routes for testing
func (r *Router) RegisterRoute() {
	path := r.route.PathPrefix("/api/v1").Subrouter()
	path.HandleFunc("/create", CreateNestedDirHandler).Methods("POST")
	path.HandleFunc("/list", ListDirStructureHandler).Methods("GET")
}

func (r *Router) Run(addr string) error {
	// CORS configuration to allow all origins
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},                                       // Allow all origins
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Allowed methods
		AllowedHeaders: []string{"Content-Type", "Authorization"},           // Allowed headers
	})

	// Wrap the mux router with CORS middleware
	handler := corsHandler.Handler(r.route)

	// Run the server with CORS enabled
	log.Println("Server running on port: ", addr)
	return http.ListenAndServe(addr, handler)
}
