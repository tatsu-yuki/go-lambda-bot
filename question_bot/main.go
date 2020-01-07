package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"net/http"
	"os"
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
	Events      []Event `json:"events"`
	Destination string  `json:"destination"`
}

type Event struct {
	Type       string  `json:"type"`
	ReplyToken string  `json:"replyToken"`
	Source     Source  `json:"source"`
	Timestamp  int64   `json:"timestamp"`
	Message    Message `json:"message"`
}

type Message struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Text string `json:"text"`
}

type Source struct {
	UserID string `json:"userId"`
	Type   string `json:"type"`
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

	leftBtn := linebot.NewMessageAction("left", "left clicked")
	rightBtn := linebot.NewMessageAction("right", "right clicked")
	template := linebot.NewConfirmTemplate("Hello World", leftBtn, rightBtn)
	message := linebot.NewTemplateMessage("Sorry :(, please update your app.", template)

	myLineRequest, err := UnmarshalLineRequest([]byte(request.Body))
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
	}

	bot, err := linebot.New(LineChannelSecret, LineAccessToken)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
	}

	fmt.Println("*** reply")
	var messages []linebot.SendingMessage

	messages = append(messages, message)

	// append some message to messages
	_, err = bot.ReplyMessage(myLineRequest.Events[0].ReplyToken, messages...).Do()
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
	}
	return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: http.StatusOK}, nil
}

func main() {
	lambda.Start(Handler)
}
