package helper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ResponseStatus string

const (
	Success ResponseStatus = "success"
	Error   ResponseStatus = "error"
)

type Response_struct struct {
	Timestamp  string         `json:"timestamp,omitempty"`
	ApiVersion string         `json:"api_version,omitempty"`
	DurationMs int64          `json:"duration_ms,omitempty"`
	Error      string         `json:"error,omitempty"`
	Message    string         `json:"message,omitempty"`
	Data       any            `json:"data,omitempty"`
	Status     ResponseStatus `json:"status,omitempty"`
	HttpInfo   any            `json:"http,omitempty"`
}

func Response(res Response_struct, w http.ResponseWriter, status int, msg string, httpInfo any, DurationMs int64) {
	w.Header().Set("Content-Type", "application/json")

	res.Timestamp = time.Now().UTC().Format(time.RFC3339)
	res.ApiVersion = "v1"
	res.Message = msg
	res.DurationMs = DurationMs
	res.HttpInfo = httpInfo

	if status >= 200 && status < 300 {
		res.Status = Success
	} else {
		res.Status = Error
	}

	data, err := json.Marshal(res)
	if err != nil {
		fmt.Println("Erro ao fazer marshal no json ", err)
		fallback := Response_struct{Timestamp: time.Now().UTC().Format(time.RFC3339), ApiVersion: "v1", Error: "something went wrong", Status: Error}
		fallbackData, _ := json.Marshal(fallback)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(fallbackData)
		return
	}

	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		fmt.Println("erro ao enviar resposta: ", err)
		return
	}
}
