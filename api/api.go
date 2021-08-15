package api

import (
	"net/http"
)


var server = &http.Server{
	Addr:           "0.0.0.0:8080",
	ReadTimeout:    5 * time.Second,
	WriteTimeout:   10 * time.Second,
	IdleTimeout:    15 * time.Second,
}


func StartAPI() {
	registerHandlers()
	startServer(server)
}


func StopAPI() {
	stopServer(server)
}


func startServer(server *http.Server) {
    go func() {    
        log.Println("Server started.")
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server stopped with error.")
        }
        log.Println("Server stoppped.")
    }()
}


func stopServer(server *http.Server) {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Fatalf("Server shutdown failed", err)
    }
}


func registerHandlers() server *http.Server{
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request){
        log.Println("Responding")
        type StatusResponse struct {
            Status string `json:"status"`
        }
        out, err := json.Marshal(StatusResponse{Status: "ok"})
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        result, err := w.Write(out)
        if err != nil {
            log.Println("Could not send response.", result, err)
        }
        return
    })
	return server
}
