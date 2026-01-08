package jsonresponse

import (
	"encoding/json"
	"net/http"
)

func writeJson(w http.ResponseWriter,status int,data any)error{
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func readJson(w http.ResponseWriter,r *http.Request,data any) error{
	maxBytes:= 1_048_578
	r.Body=http.MaxBytesReader(w,r.Body,int64(maxBytes))
	decoder:= json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(data)
}

func Error(w http.ResponseWriter,status int,message string) error {
	 type envelope struct{
		Error bool `json:"error"`
		Message string `json:"message"`
	 }
	 return writeJson(w,status,envelope{Error:true, Message:message})
}

func Success(w http.ResponseWriter,status int,message string,data any) error {
	 type envelope struct{
		Error bool `json:"error"`
		Message string `json:"message"`
		Data any `json:"data"`
	 }

	 return writeJson(w,status,envelope{Error:false, Message:message,Data:data })
}