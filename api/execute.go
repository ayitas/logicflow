package api

import (
	"net/http"

	// Register all algorithms
	_ "logicflow/internal/algorithm"
	"logicflow/internal/handler"
)

// Handler is the Vercel serverless function entrypoint for POST /execute
func Handler(w http.ResponseWriter, r *http.Request) {
	// Let the logicflow handler process the request
	handler.ExecuteHandler(w, r)
}
