package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/pkg/errors"
)

func UnmarshalLineRequest(data []byte) (LineRequest, error) {
	var r LineRequest
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *LineRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type LineRequest struct {
	Events []linebot.Event `json:"events"`
}

func UnmarshalSummary(data []byte) (SummaryResponse, error) {
	var r SummaryResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *SummaryResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type SummaryResponse struct {
	Message string   `json:"message"`
	Status  int64    `json:"status"`
	Summary []string `json:"summary"`
}

var LineChannelSecret string
var LineAccessToken string
var SummaryApiKey string

func init() {
	LineChannelSecret = os.Getenv("CHANNELSECRET")
	LineAccessToken = os.Getenv("ACCESSTOKEN")
	SummaryApiKey = os.Getenv("SUMMARY_API_KEY")
}
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("*** start")
	defer log.Println("*** end")

	log.Print("*** body")
	log.Println(request.Body)

	myLineRequest, err := UnmarshalLineRequest([]byte(request.Body))
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
	}

	bot, err := linebot.New(LineChannelSecret, LineAccessToken)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
	}

	for _, event := range myLineRequest.Events {
		switch event.Type {
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				bot.ReplyMessage(myLineRequest.Events[0].ReplyToken, linebot.NewTemplateMessage(message.Text)).Do()
				if _, err = bot.ReplyMessage(myLineRequest.Events[0].ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
					return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, nil
				}
				return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: http.StatusOK}, nil

			case *linebot.StickerMessage:
				log.Print("This is a StickerMessage.")

			case *linebot.TemplateMessage:
				log.Print("This is a TemplateMessage.")

			case *linebot.VideoMessage:
				log.Print("This is a VideoMessage.")

			case *linebot.AudioMessage:
				log.Print("This is a AudioMessage.")

			case *linebot.FileMessage:
				log.Print("This is a FileMessage.")

			case *linebot.LocationMessage:
				log.Print("This is a LocationMessage.")
			}

		case linebot.EventTypePostback:
			data := event.Postback.Data
			if err := postbackHandler(data); err != nil {
				log.Println(err)
			}
		}
	}
	return events.APIGatewayProxyResponse{Body: nil, StatusCode: http.StatusOK}, nil
}

func requestSummary(text string) (*http.Response, error) {
	apiUrl := "https://api.a3rt.recruit-tech.co.jp/text_summarization/v1"
	data := url.Values{}
	data.Set("apikey", SummaryApiKey)
	data.Set("sentences", text)

	client := &http.Client{}
	r, err := http.NewRequest("POST", apiUrl, strings.NewReader(data.Encode())) // URL-encoded payload
	if err != nil {
		return nil, errors.New("fail to create NewRequest.")
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	resp, err := client.Do(r)
	if err != nil {
		return nil, errors.New("fail to request.")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("response status code is not http.StatusOK. status code is: %d", resp.StatusCode))
	}
	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
