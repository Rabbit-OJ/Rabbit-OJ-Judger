package models

type TestCaseType struct {
	Id         int    `json:"id"`
	StdinPath  string `json:"stdin_path"`
	StdoutPath string `json:"stdout_path"`
	Stdin      []byte `json:"stdin"`
	Stdout     []byte `json:"stdout"`
}
