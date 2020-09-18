package models

type TestCaseType struct {
	Id     int    `json:"id"`
	Path   string `json:"path"`
	Stdin  []byte `json:"stdin"`
	Stdout []byte `json:"stdout"`
}
