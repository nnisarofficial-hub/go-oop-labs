package main

import (
	"fmt"
	"strings"
)

type Notification struct {
	id      int
	Title   string
	Message string
	Type    string
	read    bool
}

func NewNotification(id int, title, message, notifType string) *Notification {
	return &Notification{
		id:      id,
		Title:   title,
		Message: message,
		Type:    notifType,
		read:    false,
	}
}

func (n *Notification) MarkAsRead() {
	n.read = true
}

func (n Notification) IsRead() bool {
	return n.read
}

func (n Notification) Summary() string {
	return fmt.Sprintf("[%s] %s: %s", strings.ToUpper(n.Type), n.Title, n.Message)
}
func (n Notification) ID() int {
	return n.id
}
func main() {
	notif := NewNotification(1, "Low disk space", "Your disk is 90% full", "warning")
	fmt.Printf("ID: %d\n", notif.ID())
	fmt.Printf("Summary: %s\n", notif.Summary())
	fmt.Printf("Is read: %t\n\n", notif.IsRead())
	notif.MarkAsRead()
	fmt.Println("After marking as read:")
	fmt.Printf("Is read: %t\n", notif.IsRead())
}
