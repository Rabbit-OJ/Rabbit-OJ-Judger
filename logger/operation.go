package logger

import "fmt"

type PrintlnType = func(a ...interface{})
type PrintfType = func(format string, a ...interface{})

var (
	Println = func(a ...interface{}) {
		fmt.Println(a...)
	}
	Printf  = func(format string, a ...interface{}) {
		fmt.Printf(format, a...)
	}
)
