package main

import (
    "fmt"
    "time"
    "sync"
)


func periodic_task(wg *sync.WaitGroup, interval int) {
    defer wg.Done()
    for {
        fmt.Println("Starting task with interval: ", interval)
        time.Sleep(time.Duration(interval) * time.Second)
        fmt.Println("Finishing task with interval: ", interval)
    }
}


func task(wg *sync.WaitGroup) {
    defer wg.Done()

    fmt.Println("Starting task")
    time.Sleep(2 * time.Second)
    fmt.Println("Finishing task")
}


func main() {
    var wg sync.WaitGroup

    for i := 0; i <= 1000; i++ {
        wg.Add(1)
        go periodic_task(&wg, 5)
    }

    // wg.Add(1)
    // go periodic_task(&wg, 2)
    // wg.Add(1)
    // go task(&wg)
    wg.Wait()
}
