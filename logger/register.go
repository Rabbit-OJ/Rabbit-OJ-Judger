package logger

func RegisterLogger(println PrintlnType, printf PrintfType) {
	Println = println
	Printf = printf
}
