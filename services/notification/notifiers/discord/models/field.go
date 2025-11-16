package models

type Field struct {
	Key      string `json:"name"`
	Value    string `json:"value"`
	IsInline bool   `json:"inline"`
}
