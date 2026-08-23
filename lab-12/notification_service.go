package main

import "fmt"

type NotificationService struct {
	notifiers []Notifier
}

func NewNotificationService(notifiers ...Notifier) *NotificationService {
	return &NotificationService{
		notifiers: notifiers,
	}
}

func (ns *NotificationService) Broadcast(to, subject, body string) map[string]error {
	results := make(map[string]error)
	for _, notifier := range ns.notifiers {
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
