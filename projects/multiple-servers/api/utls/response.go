package utls

import (
	"servers/api/models"
	"net/http"
	"net/url"
	"strconv"
	"encoding/json"
	"strings"
)

func JSONError(w http.ResponseWriter, data models.FailedJson, status, indent int ) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	
	jsonRes, _ := json.MarshalIndent(data, "", strings.Repeat(" ", indent))
  w.Write(jsonRes)
}

func JSON(w http.ResponseWriter, data []models.Image, status, indent int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	jsonRes, _ := json.MarshalIndent(data, "", strings.Repeat(" ", indent))
	w.Write(jsonRes)
}


func ResponseIndent(u *url.URL) int {
  params := u.Query()
	indent := params.Get("indent")
	if indent == "" {indent = "0"}

	indentVal, err := strconv.Atoi(indent)
	if err != nil {return 0}
	return indentVal
}
