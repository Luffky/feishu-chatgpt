package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"start-feishubot/handlers"
	"start-feishubot/initialization"
	"start-feishubot/services/openai"

	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"

	sdkginext "github.com/larksuite/oapi-sdk-gin"
)

var (
	cfg = pflag.StringP("config", "c", "./config.yaml", "apiserver config file path.")
)


func main() {
	pflag.Parse()
	config := initialization.LoadConfig(*cfg)
	initialization.LoadLarkClient(*config)
	gpt := openai.NewChatGPT(*config)
	handlers.InitHandlers(gpt, *config)

	//eventHandler := dispatcher.NewEventDispatcher(
	//	config.FeishuAppVerificationToken, config.FeishuAppEncryptKey).
	//	OnP2MessageReceiveV1(handlers.Handler).
	//	OnP2MessageReadV1(func(ctx context.Context, event *larkim.P2MessageReadV1) error {
	//		return handlers.ReadHandler(ctx, event)
	//	})

	cardHandler := larkcard.NewCardActionHandler(
		config.FeishuAppVerificationToken, config.FeishuAppEncryptKey,
		handlers.CardHandler())

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/webhook/event",
		echo)
	r.POST("/webhook/card",
		sdkginext.NewCardActionHandlerFunc(
			cardHandler))

	err := initialization.StartServer(*config, r)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}

func echo(c *gin.Context) {
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	log.Printf(string(bodyBytes))
	c.JSON(200, gin.H{
		"message": string(bodyBytes),
	})
}
