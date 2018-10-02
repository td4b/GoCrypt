package main

import (
    "log"
    "net/http"
)

func apicall(w http.ResponseWriter, r *http.Request) {
    //...
}

func main() {
    http.HandleFunc("/api/", apicall)
    log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
