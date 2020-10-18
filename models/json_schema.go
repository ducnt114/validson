package models

type JsonSchema struct {
	SchemaURL   string                   `json:"schema"`
	Title       string                   `json:"title"`
	Description string                   `json:"description"`
	Type        string                   `json:"type"`
	Properties  map[string]*JsonProperty `json:"properties"`
	Required    []string                 `json:"required"`
}

type JsonProperty struct {
	Type             string  `json:"type"`
	Description      string  `json:"description"`
	Minimum          float64 `json:"minimum"`
	Maximum          float64 `json:"maximum"`
	ExclusiveMinimum bool    `json:"exclusive_minimum"`
	ExclusiveMaximum bool    `json:"exclusive_maximum"`
}
