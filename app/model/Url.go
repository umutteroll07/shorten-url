package model

type Url struct {
	Base
	Original_url string `json:"original_url"`
	Short_url    string `json:"short_url"`
	Expires_at   string `json:"expires_at"`
	Usage_count  int    `json:"usage_count"`
}
