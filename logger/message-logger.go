package logger

import (
	"github.com/dgmann/ma-shared"
	"fmt"
)

type MessageLogger struct {
	input <-chan shared.Message
}

func NewMessageLogger(input <-chan shared.Message) MessageLogger {
	return MessageLogger{input:input}
}

func(logger *MessageLogger) Print() {
	for msg := range logger.input {
		fmt.Printf("%+v\n", msg)
	}
}

func(logger *MessageLogger) PrintResults() {
	for msg := range logger.input {
		fmt.Printf("%+v\n", msg.Result)
	}
}

func(logger *MessageLogger) PrintFoundPlates() {
	for msg := range logger.input {
		fmt.Printf("%+v\n", msg.Result.OpenALPR.Results)
	}
}

func(logger *MessageLogger) PrintWanted() {
	for msg := range logger.input {
		fmt.Printf("%+v\n", msg.Result.WantedNumbers)
	}
}
