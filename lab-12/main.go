package main

import "fmt"

type Notifier interface {
	Send(to, subject, body string) error
	Name() string
}

func main() {
	ns := NewNotificationService(&EmailNotifier{}, &SMSNotifier{}, &SlackNotifier{})

	email1, phone1 := "ali@example.com", "0300-1234567"
	to1 := email1 + "/" + phone1
	fmt.Printf("Sending alert to %s / %s...\n", email1, phone1)
	report1 := ns.Broadcast(to1, "Alert", "Message body")
	ns.PrintReport(report1)
	fmt.Println()

	email2, phone2 := "invalid@example.com", "0300-1234567"
	to2 := email2 + "/" + phone2
	fmt.Printf("Sending alert to %s / %s...\n", email2, phone2)
	report2 := ns.Broadcast(to2, "Alert", "Message body")
	ns.PrintReport(report2)
}
