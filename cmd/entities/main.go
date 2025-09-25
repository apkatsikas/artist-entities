package main

import (
    "net/http"

    "github.com/apkatsikas/artist-entities"
)

func main() {
    router := entities.ServiceContainer().Setup()
    http.ListenAndServe(":8080", router)
}
