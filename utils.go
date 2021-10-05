package main

import (
	"errors"
	"fmt"

	resty "github.com/go-resty/resty/v2"
)

// SendSlackNotification will post to an 'Incoming Webook' url setup in Slack Apps. It accepts
// some text and the slack channel is saved within Slack.
func sendSlackNotification(msg string, webhookURL string) error {
	if webhookURL == "" {
		return nil
	}

	client := resty.New()
	response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(fmt.Sprintf(`{"text":"%v"}`, msg)).
		Post(webhookURL)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	if string(response.Body()) != "ok" {
		fmt.Printf("Error: %v\n", string(response.Body()))
		return errors.New("Non-ok response returned from Slack")
	}
	return nil
}
