package main

import (
    "fmt"
    "log"
    "os"
    "os/signal"
    "sync"
    "syscall"
    "time"
)


var termChan = make(chan os.Signal)
var goShutChan = make(chan bool)
var wg sync.WaitGroup


func createPeriodicTask(interval time.Duration, since time.Time, f func()) {
    wg.Add(1)
    // AfterFunc func() already is a goroutine
    time.AfterFunc(time.Until(since), func(){
        ticker := time.NewTicker(interval)
        for {
            select {
            case <-ticker.C:
                fmt.Println("Executing task...")
                f()
                fmt.Println("Executing task finished")
            case _, more := <-goShutChan:
                if !more {
                    ticker.Stop()
                    wg.Done()
                    log.Println(interval, " task closed")
                    return
                }
            }
        }
    })
}


func main() {
    signal.Notify(termChan, syscall.SIGTERM, syscall.SIGINT)
    go func() {
       sig := <-termChan
       log.Println(sig)
       close(goShutChan)
       log.Println("goShutChan closed")
    }()

    createPeriodicTask(
        3 * time.Second,
        time.Now(),
        func(){time.Sleep(1 * time.Second);fmt.Println("test 1")})
    createPeriodicTask(
        1 * time.Second,
        time.Now(),
        func(){time.Sleep(2 * time.Second);fmt.Println("test 2")})
    createPeriodicTask(
        2 * time.Second,
        time.Now(),
        func(){time.Sleep(4 * time.Second);fmt.Println("test 3")})

    wg.Wait()
}
