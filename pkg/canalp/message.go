package canalp

type Message[T any] struct {
	Data     []T    `json:"data"`
	Database string `json:"database"`
	Table    string `json:"table"`
	Type     string `json:"type"`
}
