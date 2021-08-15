package main

import (
    "encoding/json"
    "log"
    "os"
    "os/signal"
    "sync"
    "syscall"
    "time"
    "net/http"
    "context"
    
)


var termChan = make(chan os.Signal)
var goShutChan = make(chan bool)
var wg sync.WaitGroup


func signalRegister(wg *sync.WaitGroup) {
    wg.Add(1)
    signal.Notify(termChan, syscall.SIGTERM, syscall.SIGINT)
    go func() {
        defer wg.Done()
        sig := <-termChan
        log.Println(sig)
        close(goShutChan)
        log.Println("goShutChan closed")
     }()
}


func createPeriodicTask(interval time.Duration, since time.Time, f func()) {
    wg.Add(1)
    // AfterFunc func() already is a goroutine
    time.AfterFunc(time.Until(since), func(){
        ticker := time.NewTicker(interval)
        for {
            select {
            case <-ticker.C:
                log.Println("Executing task: ", interval)
                f()
                log.Println("Finishing task: ", interval)
            case _, more := <-goShutChan:
                if !more {
                    ticker.Stop()
                    wg.Done()
                    log.Println("Closing task: ", interval)
                    return
                }
            }
        }
    })
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


func main() {
    server := &http.Server{
        Addr:           "0.0.0.0:8080",
        ReadTimeout:    5 * time.Second,
        WriteTimeout:   10 * time.Second,
        IdleTimeout:    15 * time.Second,
    }
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

    signalRegister(&wg)
    startServer(server)
    createPeriodicTask(
        3 * time.Second,
        time.Now(),
        func(){time.Sleep(1 * time.Second);log.Println("test 1")})
    createPeriodicTask(
        1 * time.Second,
        time.Now(),
        func(){time.Sleep(2 * time.Second);log.Println("test 2")})
    createPeriodicTask(
        2 * time.Second,
        time.Now(),
        func(){time.Sleep(4 * time.Second);log.Println("test 3")})
    wg.Wait()
    stopServer(server)
    log.Println("Graceful shutdown.")    
}
