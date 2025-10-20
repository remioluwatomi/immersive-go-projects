package models

type FailedJson struct {
	Ok       bool   `json:"ok"`
	Message  string `json:"message"`
	Code     int    `json:"code"`
}


func CreateFailedJson(message string, code int) FailedJson {
	return FailedJson{Ok: false, Message: message, Code: code }
}
