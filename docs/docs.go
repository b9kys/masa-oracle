// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Masa API Support",
            "url": "https://masa.ai",
            "email": "support@masa.ai"
        },
        "license": {
            "name": "MIT",
            "url": "https://opensource.org/license/mit"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
			"/peers": {
				"get": {
					"description": "Retrieves a list of peers connected to the node",
					"consumes": [
						"application/json"
					],
					"produces": [
						"application/json"
					],
					"tags": [
						"Peers"
					],
					"summary": "Get list of peers",
					"responses": {
						"200": {
							"description": "List of peer IDs",
							"schema": {
								"type": "array",
								"items": {
									"type": "string"
								}
							}
						}
					}
				}
			},
			"/peer/addresses": {
				"get": {
					"description": "Retrieves a list of peer addresses connected to the node",
					"consumes": [
						"application/json"
					],
					"produces": [
						"application/json"
					],
					"tags": [
						"Peers"
					],
					"summary": "Get peer addresses",
					"responses": {
						"200": {
							"description": "List of peer addresses",
							"schema": {
								"type": "array",
								"items": {
									"type": "string"
								}
							}
						}
					}
				}
			},
			"/ads": {
				"get": {
					"description": "Retrieves a list of ads from the network",
					"consumes": [
						"application/json"
					],
					"produces": [
						"application/json"
					],
					"tags": [
						"Ads"
					],
					"summary": "Get ads",
					"responses": {
						"200": {
							"description": "List of ads",
							"schema": {
								"type": "array",
								"items": {
									"$ref": "#/definitions/Ad"
								}
							}
						}
					}
				},
				"post": {
					"description": "Adds a new ad to the network",
					"consumes": [
						"application/json"
					],
					"produces": [
						"application/json"
					],
					"tags": [
						"Ads"
					],
					"summary": "Post an ad",
					"parameters": [
						{
							"description": "Ad Content",
							"name": "ad",
							"in": "body",
							"required": true,
							"schema": {
								"$ref": "#/definitions/Ad"
							}
						}
					],
					"responses": {
						"200": {
							"description": "Ad successfully posted",
							"schema": {
								"$ref": "#/definitions/AdResponse"
							}
						},
						"400": {
							"description": "Invalid ad data",
							"schema": {
								"$ref": "#/definitions/ErrorResponse"
							}
						}
					}
				}
			},
			"/ads/subscribe": {
				"post": {
					"description": "Subscribes the user to receive ad notifications",
					"consumes": [
						"application/json"
					],
					"produces": [
						"application/json"
					],
					"tags": [
						"Ads"
					],
					"summary": "Subscribe to ads",
					"parameters": [
						{
							"description": "Subscription details",
							"name": "subscription",
							"in": "body",
							"required": true,
							"schema": {
								"$ref": "#/definitions/Subscription"
							}
						}
					],
					"responses": {
						"200": {
							"description": "Successfully subscribed to ads",
							"schema": {
								"$ref": "#/definitions/SubscriptionResponse"
							}
						},
						"400": {
							"description": "Invalid subscription data",
							"schema": {
								"$ref": "#/definitions/ErrorResponse"
							}
						}
					}
				}
			},
			"/data/twitter/profile/{username}": {
				"get": {
					"description": "Retrieves tweets from a specific Twitter profile",
					"consumes": [
						"application/json"
					],
					"produces": [
						"application/json"
					],
					"tags": [
						"Twitter"
					],
					"summary": "Search Twitter Profile",
					"parameters": [
						{
							"type": "string",
							"description": "Twitter Username",
							"name": "username",
							"in": "path",
							"required": true
						}
					],
					"responses": {
						"200": {
							"description": "List of tweets from the profile",
							"schema": {
								"type": "array",
								"items": {
									"$ref": "#/definitions/Tweet"
								}
							}
						},
						"400": {
							"description": "Invalid username or error fetching tweets",
							"schema": {
								"$ref": "#/definitions/ErrorResponse"
							}
						}
					}
				}
			},
			"/data/twitter/tweets/recent": {
				"post": {
					"description": "Retrieves recent tweets based on query parameters",
					"consumes": [
						"application/json"
					],
					"produces": [
						"application/json"
					],
					"tags": [
						"Twitter"
					],
					"summary": "Search recent tweets",
					"parameters": [
						{
							"in": "body",
							"name": "body",
							"description": "Search parameters",
							"required": true,
							"schema": {
								"type": "object",
								"properties": {
									"query": {
										"type": "string",
										"description": "Search Query"
									},
									"count": {
										"type": "integer",
										"description": "Number of tweets to return"
									}
								}
							}
						}
					],
					"responses": {
						"200": {
							"description": "List of recent tweets",
							"schema": {
								"type": "array",
								"items": {
									"$ref": "#/definitions/Tweet"
								}
							}
						},
						"400": {
							"description": "Invalid query or error fetching tweets",
							"schema": {
								"$ref": "#/definitions/ErrorResponse"
							}
						}
					}
				}
			},
			"/data/twitter/tweets/trends": {
				"get": {
					"description": "Retrieves the latest Twitter trending topics",
					"consumes": [
						"application/json"
					],
					"produces": [
						"application/json"
					],
					"tags": [
						"Twitter"
					],
					"summary": "Twitter Trends",
					"responses": {
						"200": {
							"description": "List of trending topics",
							"schema": {
								"type": "array",
								"items": {
									"$ref": "#/definitions/Trend"
								}
							}
						},
						"400": {
							"description": "Error fetching Twitter trends",
							"schema": {
								"$ref": "#/definitions/ErrorResponse"
							}
						}
					}
				}
			},
			"/data/web": {
				"post": {
					"description": "Retrieves data from the web",
					"consumes": [
						"application/json"
					],
					"produces": [
						"application/json"
					],
					"tags": [
						"Web"
					],
					"summary": "Web Data",
					"parameters": [
						{
							"in": "body",
							"name": "body",
							"description": "Search parameters",
							"required": true,
							"schema": {
								"type": "object",
								"properties": {
									"url": {
										"type": "string",
										"description": "Url"
									},
									"depth": {
										"type": "integer",
										"description": "Number of pages to scrape"
									}
								}
							}
						}
					],
					"responses": {
						"200": {
							"description": "Successfully retrieved web data",
							"schema": {
								"$ref": "#/definitions/WebDataResponse"
							}
						},
						"400": {
							"description": "Invalid query or error fetching web data",
							"schema": {
								"$ref": "#/definitions/ErrorResponse"
							}
						}
					}
				}
			},
			"/dht": {
				"get": {
					"description": "Retrieves data from the DHT (Distributed Hash Table)",
					"consumes": [
						"application/json"
					],
					"produces": [
						"application/json"
					],
					"tags": [
						"DHT"
					],
					"summary": "Get DHT Data",
					"parameters": [
						{
							"in": "query",
							"name": "key",
							"description": "Key to retrieve data for",
							"required": true,
							"type": "string"
						}
					],
					"responses": {
						"200": {
							"description": "Successfully retrieved data from DHT",
							"schema": {
								"$ref": "#/definitions/DHTResponse"
							}
						},
						"400": {
							"description": "Error retrieving data from DHT",
							"schema": {
								"$ref": "#/definitions/ErrorResponse"
							}
						}
					}
				},
				"post": {
					"description": "Adds data to the DHT (Distributed Hash Table)",
					"consumes": [
						"application/json"
					],
					"produces": [
						"application/json"
					],
					"tags": [
						"DHT"
					],
					"summary": "Post to DHT",
					"parameters": [
						{
							"description": "Data to store in DHT",
							"name": "data",
							"in": "body",
							"required": true,
							"schema": {
								"type": "object",
								"properties": {
									"key": {
										"type": "string"
									},
									"value": {
										"type": "string"
									}
								}
							}
						}
					],
					"responses": {
						"200": {
							"description": "Successfully added data to DHT",
							"schema": {
								"$ref": "#/definitions/SuccessResponse"
							}
						},
						"400": {
							"description": "Error adding data to DHT",
							"schema": {
								"$ref": "#/definitions/ErrorResponse"
							}
						}
					}
				}
			},
			"/llm/models": {
				"get": {
					"description": "Retrieves the available LLM models",
					"consumes": [
						"application/json"
					],
					"produces": [
						"application/json"
					],
					"tags": [
						"LLM"
					],
					"summary": "Get LLM Models",
					"responses": {
						"200": {
							"description": "Successfully retrieved LLM models",
							"schema": {
								"$ref": "#/definitions/LLMModelsResponse"
							}
						},
						"400": {
							"description": "Error retrieving LLM models",
							"schema": {
								"$ref": "#/definitions/ErrorResponse"
							}
						}
					}
				}
			},
			"/chat": {
				"post": {
					"description": "Initiates a chat session with an AI model that accepts common ollama formatted requests",
					"consumes": [
						"application/json"
					],
					"produces": [
						"application/json"
					],
					"tags": [
						"Chat"
					],
					"summary": "Chat with AI",
					"parameters": [
						{
							"description": "Message to send to AI",
							"name": "message",
							"in": "body",
							"required": true,
							"schema": {
								"type": "string"
							}
						}
					],
					"responses": {
						"200": {
							"description": "Successfully received response from AI",
							"schema": {
								"$ref": "#/definitions/ChatResponse"
							}
						},
						"400": {
							"description": "Error communicating with AI",
							"schema": {
								"$ref": "#/definitions/ErrorResponse"
							}
						}
					}
				}
			},
			"/node/data": {
				"get": {
					"description": "Retrieves data from the node",
					"consumes": [
						"application/json"
					],
					"produces": [
						"application/json"
					],
					"tags": [
						"Node"
					],
					"summary": "Node Data",
					"responses": {
						"200": {
							"description": "Successfully retrieved node data",
							"schema": {
								"$ref": "#/definitions/NodeDataResponse"
							}
						},
						"400": {
							"description": "Error retrieving node data",
							"schema": {
								"$ref": "#/definitions/ErrorResponse"
							}
						}
					}
				}
			},
			"/node/data/{peerid}": {
				"get": {
					"description": "Retrieves data for a specific node identified by peer ID",
					"consumes": [
						"application/json"
					],
					"produces": [
						"application/json"
					],
					"tags": [
						"Node"
					],
					"summary": "Get Node Data by Peer ID",
					"parameters": [
						{
							"type": "string",
							"description": "Peer ID",
							"name": "peerid",
							"in": "path",
							"required": true
						}
					],
					"responses": {
						"200": {
							"description": "Successfully retrieved node data by peer ID",
							"schema": {
								"$ref": "#/definitions/NodeDataResponse"
							}
						},
						"400": {
							"description": "Error retrieving node data by peer ID",
							"schema": {
								"$ref": "#/definitions/ErrorResponse"
							}
						}
					}
				}
			},
			"/sentiment/tweets": {
				"post": {
					"description": "Searches for tweets and analyzes their sentiment",
					"consumes": [
						"application/json"
					],
					"produces": [
						"application/json"
					],
					"tags": [
						"Sentiment"
					],
					"summary": "Analyze Sentiment of Tweets",
					"parameters": [
						{
							"in": "body",
							"name": "body",
							"description": "Sentiment analysis request body",
							"required": true,
							"schema": {
								"type": "object",
								"properties": {
									"query": {
										"type": "string",
										"description": "Search Query"
									},
									"count": {
										"type": "integer",
										"description": "Number of tweets to analyze"
									},
									"model": {
										"type": "string",
										"description": "Sentiment analysis model to use"
									}
								}
							}
						}
					],
					"responses": {
						"200": {
							"description": "Successfully analyzed sentiment of tweets",
							"schema": {
								"$ref": "#/definitions/SentimentAnalysisResponse"
							}
						},
						"400": {
							"description": "Error analyzing sentiment of tweets",
							"schema": {
								"$ref": "#/definitions/ErrorResponse"
							}
						}
					}
				}
			},
			"/sentiment/web": {
				"post": {
					"description": "Searches for web content and analyzes its sentiment",
					"consumes": [
						"application/json"
					],
					"produces": [
						"application/json"
					],
					"tags": [
						"Sentiment"
					],
					"summary": "Analyze Sentiment of Web Content",
					"parameters": [
						{
							"in": "body",
							"name": "body",
							"description": "Sentiment analysis request body",
							"required": true,
							"schema": {
								"type": "object",
								"properties": {
									"url": {
										"type": "string",
										"description": "URL of the web content"
									},
									"depth": {
										"type": "integer",
										"description": "Depth of web crawling"
									},
									"model": {
										"type": "string",
										"description": "Sentiment analysis model to use"
									}
								}
							}
						}
					],
					"responses": {
						"200": {
							"description": "Successfully analyzed sentiment of web content",
							"schema": {
								"$ref": "#/definitions/SentimentAnalysisResponse"
							}
						},
						"400": {
							"description": "Error analyzing sentiment of web content",
							"schema": {
								"$ref": "#/definitions/ErrorResponse"
							}
						}
					}
				}
			},
		},
		"ChatResponse": {
			"type": "object",
			"properties": {
				"message": {
					"type": "string"
				}
			}
		},
		"DHTResponse": {
			"type": "object",
			"properties": {
				"key": {
					"type": "string"
				},
				"value": {
					"type": "string"
				}
			}
		},
		"SuccessResponse": {
			"type": "object",
			"properties": {
				"message": {
					"type": "string"
				}
			}
		},
		"WebDataRequest": {
			"type": "object",
			"properties": {
				"query": {
					"type": "string"
				},
				"url": {
					"type": "string"
				},
				"depth": {
					"type": "integer"
				}
			}
		},
		"WebDataResponse": {
			"type": "object",
			"properties": {
				"data": {
					"type": "string"
				}
			}
		},
		"SentimentAnalysisResponse": {
			"type": "object",
			"properties": {
				"sentiment": {
					"type": "string"
				},
				"data": {
					"type": "string"
				}
			}
		},
		"definitions": {
			"ErrorResponse": {
				"type": "object",
				"properties": {
					"error": {
						"type": "string"
					}
				}
			},				
			"Tweet": {
				"type": "object",
				"properties": {
					"id": {
						"type": "string"
					},
					"text": {
						"type": "string"
					},
					"created_at": {
						"type": "string"
					},
					"user": {
						"type": "object",
						"properties": {
							"id": {
								"type": "string"
							},
							"name": {
								"type": "string"
							},
							"screen_name": {
								"type": "string"
							}
						}
					}
				}
			},
			"Trend": {
				"type": "object",
				"properties": {
					"name": {
						"type": "string"
					},
					"url": {
						"type": "string"
					},
					"tweet_volume": {
						"type": "integer"
					}
				}
			},
			"SentimentAnalysisResponse": {
				"type": "object",
				"properties": {
					"sentiment": {
						"type": "string"
					},
					"data": {
						"type": "string"
					}
				}
			},
			"WebDataResponse": {
				"type": "object",
				"properties": {
					"data": {
						"type": "string"
					}
				}
			},
			"LLMModelsResponse": {
				"type": "object",
				"properties": {
					"models": {
						"type": "array",
						"items": {
							"type": "string"
						}
					}
				}
			},
			"NodeDataResponse": {
				"type": "object",
				"properties": {
					"peer_id": {
						"type": "string"
					},
					"data": {
						"type": "string"
					}
				}
			}
		}
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.0.2-beta",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Masa Oracle API",
	Description:      "The Worlds Personal Data Network Masa Oracle Node API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
