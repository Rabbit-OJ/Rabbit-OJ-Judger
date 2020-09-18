package models

type TestCaseType struct {
	Id     int `json:"id"`
	Path   string `json:"path"`
	Stdout []byte `json:"stdout"`
}
