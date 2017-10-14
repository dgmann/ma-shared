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

func(logger *MessageLogger) Print() chan shared.Message {
	output := make(chan shared.Message, 10000)
	go func() {
		for msg := range logger.input {
			output <- msg
			fmt.Printf("%+v\n", msg)
		}
		close(output)
	}()
	return output
}

func(logger *MessageLogger) PrintResults() chan shared.Message {
	output := make(chan shared.Message, 10000)
	go func() {
		for msg := range logger.input {
			output <- msg
			fmt.Printf("%+v\n", msg.Result)
		}
		close(output)
	}()
	return output
}

func(logger *MessageLogger) PrintFoundPlates() chan shared.Message {
	output := make(chan shared.Message, 10000)
	go func() {
		for msg := range logger.input {
			output <- msg
			fmt.Printf("%+v\n", msg.Result.OpenALPR.Results)
		}
		close(output)
	}()
	return output
}

func(logger *MessageLogger) PrintWanted() chan shared.Message {
	output := make(chan shared.Message, 10000)
	go func() {
		lastCount := 0
		for msg := range logger.input {
			output <- msg
			if len(msg.Result.WantedNumbers) > lastCount {
				fmt.Printf("Wanted Numbers:\n%+v\n", msg.Result.WantedNumbers)
				lastCount = len(msg.Result.WantedNumbers)
			}
		}
		close(output)
	}()
	return output
}
