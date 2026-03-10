package dto

type ApiRootDTO struct {
	Version string                `json:"version"`
	Paths   map[string]ApiLinkDTO `json:"paths"`
}

type ApiLinkDTO struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}
