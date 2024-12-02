package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gay/plugintools/internal/core"
)

// Server represents the HTTP server for the tools platform
type Server struct {
	registry core.ToolRegistry
}

// NewServer creates a new server instance
func NewServer(registry core.ToolRegistry) *Server {
	return &Server{
		registry: registry,
	}
}

// Start starts the HTTP server
func (s *Server) Start(addr string) error {
	// Register routes with middleware
	http.HandleFunc("/api/v1/tools", Chain(s.handleTools, Logger, Auth))
	http.HandleFunc("/api/v1/tools/", Chain(s.handleToolOperation, Logger, Auth))

	fmt.Printf("Server starting on %s\n", addr)
	return http.ListenAndServe(addr, nil)
}

// handleTools handles GET /api/v1/tools
func (s *Server) handleTools(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tools := s.registry.List()
	toolInfos := make([]core.ToolInfo, 0, len(tools))
	for _, tool := range tools {
		toolInfos = append(toolInfos, tool.GetInfo())
	}

	s.writeJSON(w, toolInfos)
}

// handleToolOperation handles operations on specific tools
func (s *Server) handleToolOperation(w http.ResponseWriter, r *http.Request) {
	toolID := r.URL.Path[len("/api/v1/tools/"):]
	if toolID == "" {
		http.Error(w, "Tool ID required", http.StatusBadRequest)
		return
	}

	tool, err := s.registry.Get(toolID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		if r.URL.Query().Get("params") == "true" {
			s.writeJSON(w, tool.GetParams())
		} else {
			s.writeJSON(w, tool.GetInfo())
		}
	case http.MethodPost:
		s.handleToolExecution(w, r, tool)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleToolExecution handles tool execution
func (s *Server) handleToolExecution(w http.ResponseWriter, r *http.Request, tool core.Tool) {
	var params map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required parameters
	for _, param := range tool.GetParams() {
		if param.Required {
			if _, exists := params[param.Name]; !exists {
				http.Error(w, fmt.Sprintf("Missing required parameter: %s", param.Name), http.StatusBadRequest)
				return
			}
		}
	}

	result, err := tool.Execute(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, result)
}

// writeJSON writes JSON response
func (s *Server) writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
