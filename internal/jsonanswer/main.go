package jsonanswer

import (
	"encoding/json"
	"log"
	"net/http"
)

func Error(w http.ResponseWriter, code string, message string) {
	w.WriteHeader(http.StatusOK)
	resp := make(map[string]interface{})
	resp["error"] = true
	resp["code"] = code
	resp["message"] = message
	returnJson(w, resp)
}

func Response(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusCreated)
	resp := make(map[string]interface{})
	resp["error"] = false
	resp["message"] = message
	returnJson(w, resp)
}

func returnJson(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}
