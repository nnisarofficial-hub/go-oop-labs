package main

import (
	"fmt"
	"strings"
)

type SMSNotifier struct {
	APIKey string
}

func (s *SMSNotifier) Name() string {
	return "SMS"
}
func (s *SMSNotifier) Send(to, subject, body string) error {
	if strings.HasPrefix(to, "03") {
		return fmt.Errorf("invalid phone number format")
	}
	return nil
}
