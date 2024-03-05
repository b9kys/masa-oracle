package api

import (
	"github.com/gin-gonic/gin"

	masa "github.com/masa-finance/masa-oracle/pkg"
)

func SetupRoutes(node *masa.OracleNode) *gin.Engine {
	router := gin.Default()

	API := NewAPI(node)

	router.GET("/peers", API.GetPeersHandler())
	router.GET("/peerAddresses", API.GetPeerAddresses())

	router.POST("/ads", API.PostAd())
	router.GET("/ads", API.GetAds())
	router.POST("/subscribeToAds", API.SubscribeToAds())

	router.GET("/nodeData", API.GetNodeDataHandler())
	router.GET("/nodeData/:peerID", API.GetNodeHandler())

	router.GET("/publicKeys", API.GetPublicKeysHandler())
	router.POST("/publishPublicKey", API.PublishPublicKeyHandler())

	router.POST("/createTopic", API.CreateNewTopicHandler())
	router.POST("/postToTopic", API.PostToTopicHandler())

	return router
}
