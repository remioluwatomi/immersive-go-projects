package types

type Image struct {
    Title   string `json:"title"`
		AltText string `json:"alt_text"`
		URL     string `json:"url"`
}

type FailureJson struct {
	  Ok         bool `json:"ok"`
		Message    string `json:"message"`
		Code       int `json:"code"`
	}
