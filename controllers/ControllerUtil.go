package controllers

import (
    "encoding/json"
    "net/http"
)

func encodeRes(res http.ResponseWriter, v any) {
    json.NewEncoder(res).Encode(v)
}

func handleRes(res http.ResponseWriter, v any, status int) {
    res.WriteHeader(status)
    encodeRes(res, v)
}
