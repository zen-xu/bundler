package utils

import (
	"fmt"
	color "github.com/mgutz/ansi"
	"os"
	"runtime"
)

var (
	Purple = color.ColorFunc("magenta+h")
	Red    = color.ColorFunc("red+h")
	Yellow = color.ColorFunc("yellow+h")
	Green  = color.ColorFunc("green+h")
	Blue   = color.ColorFunc("blue+h")
	Gray   = color.ColorFunc("black+h")
	Bold   = color.ColorFunc("default+b")
)

// ExitWithErrorMessage will exit with return code 1 and output an error message
func ExitWithErrorMessage(msg string) {
	_, _ = fmt.Fprintln(os.Stderr, Red(msg))
	os.Exit(1)
}

// CheckError will exit upon the presence of an error, showing a message upon error
func CheckError(err error, message string) {
	if err != nil {
		fmt.Println(Red("Error:"))
		_, file, line, _ := runtime.Caller(1)
		fmt.Println("Line:", line, "\tFile:", file, "\n", err)
		ExitWithErrorMessage(message)
	}
}
