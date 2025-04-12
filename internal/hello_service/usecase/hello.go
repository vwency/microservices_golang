package usecase_hello

import "strings"

type HelloUsecase interface {
	ProcessGreeting(text string) string
}

type helloUsecase struct{}

func NewHelloUsecase() HelloUsecase {
	return &helloUsecase{}
}

func (u *helloUsecase) ProcessGreeting(text string) string {
	if strings.Contains(strings.ToLower(text), "hello") {
		return "hello"
	}
	return "None"
}
