package main

import (
	"fmt"
	"strings"
)

type Notifier interface {
	Send(to, subject, body string) error
	Name() string
}

type EmailNotifier struct {
	SMTPServer string
}

func (e *EmailNotifier) Name() string { return "Email" }

func (e *EmailNotifier) Send(to, subject, body string) error {
	if strings.Contains(to, "invalid") {
		return fmt.Errorf("failed — invalid recipient address")
	}
	return nil
}

type SMSNotifier struct {
	APIKey string
}

func (s *SMSNotifier) Name() string { return "SMS" }
func (s *SMSNotifier) Send(to, subject, body string) error {
	if !strings.HasPrefix(to, "03") {
		return fmt.Errorf("invalid phone number format")
	}
	return nil
}

type SlackNotifier struct {
	Channel string
}

func (s *SlackNotifier) Name() string { return "Slack" }

func (s *SlackNotifier) Send(to, subject, body string) error {
	return nil
}

type Recipient struct {
	Email string
	Phone string
}

type NotificationService struct {
	notifiers []Notifier
}

func NewNotificationService(notifiers ...Notifier) *NotificationService {
	return &NotificationService{
		notifiers: notifiers,
	}
}

func (ns *NotificationService) Broadcast(r Recipient, subject, body string) map[string]error {
	results := make(map[string]error)
	for _, notifier := range ns.notifiers {
		var to string
		switch notifier.(type) {
		case *EmailNotifier:
			to = r.Email
		case *SMSNotifier, *SlackNotifier:
			to = r.Phone
		}
		results[notifier.Name()] = notifier.Send(to, subject, body)
	}
	return results
}

func (ns *NotificationService) PrintReport(results map[string]error) {
	successCount := 0
	totalCount := len(results)
	order := []string{"Email", "SMS", "Slack"}

	for _, name := range order {
		err, exists := results[name]
		if !exists {
			continue
		}
		if err == nil {
			fmt.Printf("%-6s ✓ sent\n", name+":")
			successCount++
		} else {
			fmt.Printf("%-6s X %v\n", name+":", err)
		}
	}
	fmt.Printf("%d of %d channels successful\n", successCount, totalCount)
}

func main() {
	ns := NewNotificationService(&EmailNotifier{}, &SMSNotifier{}, &SlackNotifier{})
	r1 := Recipient{Email: "ali@example.com", Phone: "0300-1234567"}
	fmt.Printf("Sending alert to %s / %s...\n", r1.Email, r1.Phone)
	report1 := ns.Broadcast(r1, "Alert", "Message body")
	ns.PrintReport(report1)
	fmt.Println()
	r2 := Recipient{Email: "invalid@example.com", Phone: "0300-1234567"}
	fmt.Printf("Sending alert to %s / %s...\n", r2.Email, r2.Phone)
	report2 := ns.Broadcast(r2, "Alert", "Message body")
	ns.PrintReport(report2)
}
