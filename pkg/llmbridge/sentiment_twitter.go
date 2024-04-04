package llmbridge

import (
	"strings"

	twitterscraper "github.com/n0madic/twitter-scraper"
	"github.com/sirupsen/logrus"
)

// AnalyzeSentimentTweets analyzes the sentiment of the provided tweets by sending them to the Claude API.
// It concatenates the tweets, creates a payload, sends a request to Claude, parses the response,
// and returns the concatenated tweets content, a sentiment summary, and any error.
func AnalyzeSentimentTweets(tweets []*twitterscraper.Tweet, model string, prompt string) (string, string, error) {
	// check if we are using claude or gpt, can add others easily
	if strings.Contains(model, "claude-") {
		client := NewClaudeClient() // Adjusted to call without arguments
		tweetsContent := ConcatenateTweets(tweets)
		payloadBytes, err := CreatePayload(tweetsContent, model, prompt)
		if err != nil {
			logrus.Errorf("Error creating payload: %v", err)
			return "", "", err
		}
		resp, err := client.SendRequest(payloadBytes)
		if err != nil {
			logrus.Errorf("Error sending request to Claude API: %v", err)
			return "", "", err
		}
		defer resp.Body.Close()
		sentimentSummary, err := ParseResponse(resp)
		if err != nil {
			logrus.Errorf("Error parsing response from Claude: %v", err)
			return "", "", err
		}
		return tweetsContent, sentimentSummary, nil

	} else if strings.Contains(model, "gpt-") {
		client := NewGPTClient()
		tweetsContent := ConcatenateTweets(tweets)
		sentimentSummary, err := client.SendRequest(tweetsContent, model, prompt)
		if err != nil {
			logrus.Errorf("Error sending request to GPT: %v", err)
			return "", "", err
		}
		return tweetsContent, sentimentSummary, nil
	} else {
		return "", "", nil
	}
}

// ConcatenateTweets concatenates the text of the provided tweets into a single string,
// with each tweet separated by a newline character.
func ConcatenateTweets(tweets []*twitterscraper.Tweet) string {
	var tweetsTexts []string
	for _, tweet := range tweets {
		tweetsTexts = append(tweetsTexts, tweet.Text)
	}
	return strings.Join(tweetsTexts, "\n")
}

// AnalyzeSentimentWeb analyzes the sentiment of the provided web page text data by sending them to the Claude API.
// It concatenates the text, creates a payload, sends a request to Claude, parses the response,
// and returns the concatenated content, a sentiment summary, and any error.
func AnalyzeSentimentWeb(data string, model string, prompt string) (string, string, error) {
	// check if we are using claude or gpt, can add others easily
	if strings.Contains(model, "claude-") {
		client := NewClaudeClient() // Adjusted to call without arguments
		payloadBytes, err := CreatePayload(data, model, prompt)
		if err != nil {
			logrus.Errorf("Error creating payload: %v", err)
			return "", "", err
		}
		resp, err := client.SendRequest(payloadBytes)
		if err != nil {
			logrus.Errorf("Error sending request to Claude API: %v", err)
			return "", "", err
		}
		defer resp.Body.Close()
		sentimentSummary, err := ParseResponse(resp)
		if err != nil {
			logrus.Errorf("Error parsing response from Claude: %v", err)
			return "", "", err
		}
		return data, sentimentSummary, nil

	} else if strings.Contains(model, "gpt-") {
		client := NewGPTClient()
		sentimentSummary, err := client.SendRequest(data, model, prompt)
		if err != nil {
			logrus.Errorf("Error sending request to GPT: %v", err)
			return "", "", err
		}
		return data, sentimentSummary, nil
	} else {
		return "", "", nil
	}
}
