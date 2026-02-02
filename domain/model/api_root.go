package model

type ApiRoot struct {
	Version string             `json:"version"`
	Paths   map[string]ApiLink `json:"paths"`
}

type ApiLink struct {
	Method  string `json:"method"`
	Path    string `json:"path"`
	Example string `json:"example"`
}
