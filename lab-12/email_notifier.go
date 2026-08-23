package main

import (
	"fmt"
	"strings"
)

type EmailNotifier struct {
	SMTPServer string
}

func (e *EmailNotifier) Name() string {
	return "Email"
}

func (e *EmailNotifier) Send(to, subject, body string) error {
	if strings.Contains(to, "invalid") {
		return fmt.Errorf("failed — invalid recipient address")
	}
	return nil
}
