package tool

import (
	"runtime"
	"path"
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

func Sprintf(format string, a ...interface{}) string {
	callStartIndex := 1
	for {
		_, f, _, _ := runtime.Caller(callStartIndex)
		if !strings.HasSuffix(f, "error.go") {
			break
		} else {
			callStartIndex += 1
		}

	}
	_, file, line, _ := runtime.Caller(callStartIndex)
	_, file2, line2, _ := runtime.Caller(callStartIndex+1)
	fileName := path.Base(file)
	fileName2 := path.Base(file2)
	callStack := fmt.Sprintf("[%s:%d=>%s:%d]", fileName2, line2, fileName, line)
	return fmt.Sprintf(callStack+format, a...)
}


func Error(err error) error {
	return errors.New(Sprintf("%v", err))
}

func ErrorWithMSG(msg string, err error) error {
	return errors.New(Sprintf(msg+":%v", err))
}

func ErrorMSG(msg string) error {
	return errors.New(Sprintf(msg))
}

func Println(a ...interface{})  {
	var msg string
	for _, v := range a {
		msg = fmt.Sprintf("%s %v", msg, v)
	}
	fmt.Println(ErrorMSG(msg))
}

func Printf(format string, a ...interface{}) {
	fmt.Print(Sprintf(format, a...))
}

func Printfln(format string, a ...interface{}) {
	fmt.Println(Sprintf(format+"", a...))
}
