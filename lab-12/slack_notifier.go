package main

type SlackNotifier struct {
	Channel string
}

func (s *SlackNotifier) Name() string {
	return "Slack"
}

func (s *SlackNotifier) Send(to, subject, body string) error {
	return nil
}
