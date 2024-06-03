package workers

import (
	"encoding/json"
	"errors"

	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/scrapers/discord"
	"github.com/masa-finance/masa-oracle/pkg/scrapers/twitter"
	"github.com/masa-finance/masa-oracle/pkg/scrapers/web"
)

type LLMChatHandler struct{}

// HandleWork implements the WorkHandler interface for LLMChatHandler.
// It contains the logic for processing LLMChat work.
func (h *LLMChatHandler) HandleWork(data map[string]interface{}) (interface{}, error) {
	logrus.Infof("[+] LLM Chat %s", data)
	uri := config.GetInstance().LLMChatUrl
	if uri == "" {
		return nil, errors.New("missing env variable LLM_CHAT_URL")
	}
	jsnBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return Post(uri, jsnBytes, nil)
}

// DiscordHandler is a struct that implements the WorkHandler interface for Discord work.
type DiscordHandler struct{}

// HandleWork implements the WorkHandler interface for DiscordHandler.
func (h *DiscordHandler) HandleWork(data map[string]interface{}) (interface{}, error) {
	logrus.Infof("[+] Discord %s", data)
	userID := data["userID"].(string)
	botToken := data["botToken"].(string)
	return discord.GetUserProfile(userID, botToken)
}

// All of the twitter handlers implement the WorkHandler interface.

type TwitterQueryHandler struct{}
type TwitterFollowersHandler struct{}
type TwitterProfileHandler struct{}
type TwitterSentimentHandler struct{}
type TwitterTrendsHandler struct{}

func (h *TwitterQueryHandler) HandleWork(data map[string]interface{}) (interface{}, error) {
	logrus.Infof("[+] TwitterQueryHandler %s", data)
	count := int(data["count"].(float64))
	query := data["query"].(string)
	return twitter.ScrapeTweetsByQuery(query, count)
}

func (h *TwitterFollowersHandler) HandleWork(data map[string]interface{}) (interface{}, error) {
	logrus.Infof("[+] TwitterFollowersHandler %s", data)
	username := data["username"].(string)
	count := int(data["count"].(float64))
	return twitter.ScrapeFollowersForProfile(username, count)
}

func (h *TwitterProfileHandler) HandleWork(data map[string]interface{}) (interface{}, error) {
	logrus.Infof("[+] TwitterProfileHandler %s", data)
	username := data["username"].(string)
	return twitter.ScrapeTweetsProfile(username)
}

func (h *TwitterSentimentHandler) HandleWork(data map[string]interface{}) (interface{}, error) {
	logrus.Infof("[+] TwitterSentimentHandler %s", data)
	count := int(data["count"].(float64))
	query := data["query"].(string)
	model := data["model"].(string)
	_, resp, err := twitter.ScrapeTweetsForSentiment(query, count, model)
	return resp, err
}

func (h *TwitterTrendsHandler) HandleWork(data map[string]interface{}) (interface{}, error) {
	logrus.Infof("[+] TwitterTrendsHandler %s", data)
	return twitter.ScrapeTweetsByTrends()
}

// All of the web handlers implement the WorkHandler interface.
type WebHandler struct{}
type WebSentimentHandler struct{}

func (h *WebHandler) HandleWork(data map[string]interface{}) (interface{}, error) {
	logrus.Infof("[+] WebHandler %s", data)
	depth := int(data["depth"].(float64))
	urls := []string{data["url"].(string)}
	return web.ScrapeWebData(urls, depth)
}

func (h *WebSentimentHandler) HandleWork(data map[string]interface{}) (interface{}, error) {
	logrus.Infof("[+] WebSentimentHandler %s", data)
	depth := int(data["depth"].(float64))
	urls := []string{data["url"].(string)}
	model := data["model"].(string)

	_, resp, err := web.ScrapeWebDataForSentiment(urls, depth, model)
	return resp, err
}
