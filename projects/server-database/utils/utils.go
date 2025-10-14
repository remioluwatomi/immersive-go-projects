package utils 

import (
	"net/url"
	"strconv"
)

func GetConvIndent(u *url.URL) int {
	paramsValues := u.Query()
	indent := paramsValues.Get("indent")
	if indent == "" {indent = "0"}
  convIndent, err := strconv.Atoi(indent)
  if err != nil {return 0}
	return convIndent
}

