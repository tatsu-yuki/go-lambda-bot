package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/line/line-bot-sdk-go/linebot"
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

func init() {
	LineChannelSecret = os.Getenv("CHANNELSECRET")
	LineAccessToken = os.Getenv("ACCESSTOKEN")
}
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	questionMap := map[int]string{
		1:  "「私なんてどうせ無理だ」と思うことがよくある",
		2:  "自分の行動に対して、○○すべき、××しなければならないとよく思う",
		3:  "いままで自分だけが頑張り、他の人は頑張らず体調を崩した事がある",
		4:  "褒められると本当は良いと思っていても「そんな事ないよ」と謙遜してしまう",
		5:  "すごく気を使っているのに人間関係が上手くいかない",
		6:  "プレゼントの金額や感謝の言葉が重すぎて引かれてしまった事がある",
		7:  "誰も私の心をわかってくれないと思う時がよくある",
		8:  "自分を好きになってくれる異性などいない",
		9:  "こんな自分がどうやって生きられるのかと不安になる",
		10: "周囲の人は自分の揚げ足ばかりとるといつも思ってしまう",
	}

	answerMap := map[int]string{
		1: "たまに自己の事を嫌いになってしまう事があるようです。自分の好きなものを見つけたり、不必要な謙遜をしなくてもいいように、自己分析をしてみませんか？",
		2: "人間関係にとても悩みがあるようです。自分の事をわかってくれない、まわりは全然頑張ってくれないと思う事はありませんか？人間関係、アサーティブなコミュニケーションを学んでみませんか？",
		3: "人間関係にとても疲れているようです。自分自身を傷つけてしまったり、目の前が真っ暗になるような事はありませんか？まずはご自分の事を認めてあげれるような良いところ探しが必要です、面談を通じてご自分を好きになれるよう行動パターンを変えてみませんか？",
	}

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
				log.Print("message: " + message.Text)
				templateMessage := linebot.NewTemplateMessage("confirms template",
					&linebot.ConfirmTemplate{
						Text: "「私なんてどうせ無理だ」と思うことがよくある",
						Actions: []linebot.TemplateAction{
							&linebot.PostbackAction{
								Label: "はい",
								Data:  "1,1",
							},
							&linebot.PostbackAction{
								Label: "いいえ",
								Data:  "1,0",
							},
						},
					})

				_, err := bot.ReplyMessage(myLineRequest.Events[0].ReplyToken, templateMessage).Do()
				if err != nil {
					return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, nil
				}
				return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: http.StatusOK}, nil

			case *linebot.ImageMessage:
				log.Print(message)
			case *linebot.VideoMessage:
				log.Print(message)
			case *linebot.AudioMessage:
				log.Print(message)
			case *linebot.FileMessage:
				log.Print(message)
			case *linebot.LocationMessage:
				log.Print(message)
			case *linebot.StickerMessage:
				log.Print(message)
			default:
				log.Printf("Unknown message: %v", message)
			}

		case linebot.EventTypePostback:
			log.Print("This is a EventTypePostback.")

			data := event.Postback.Data
			log.Print("data: " + data)
			s := strings.Split(data, ",")

			currentQuestionNo, err := strconv.Atoi(s[0])
			if err != nil {
				return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, nil
			}

			currentCount, err := strconv.Atoi(s[1])
			if err != nil {
				return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, nil
			}

			if currentQuestionNo == 10 {
				if currentCount < 3 {
					if _, err = bot.ReplyMessage(myLineRequest.Events[0].ReplyToken, linebot.NewTextMessage(answerMap[1])).Do(); err != nil {
						return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, nil
					}
					return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: http.StatusOK}, nil
				}
				if currentCount < 5 {
					if _, err = bot.ReplyMessage(myLineRequest.Events[0].ReplyToken, linebot.NewTextMessage(answerMap[2])).Do(); err != nil {
						return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, nil
					}
					return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: http.StatusOK}, nil
				}
				if _, err = bot.ReplyMessage(myLineRequest.Events[0].ReplyToken, linebot.NewTextMessage(answerMap[3])).Do(); err != nil {
					return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, nil
				}
				return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: http.StatusOK}, nil
			}

			templateMessage := linebot.NewTemplateMessage("confirms template",
				&linebot.ConfirmTemplate{
					Text: questionMap[currentQuestionNo+1],
					Actions: []linebot.TemplateAction{
						&linebot.PostbackAction{
							Label: "はい",
							Data:  strconv.Itoa(currentQuestionNo+1) + "," + strconv.Itoa(currentCount+1),
						},
						&linebot.PostbackAction{
							Label: "いいえ",
							Data:  strconv.Itoa(currentQuestionNo+1) + "," + strconv.Itoa(currentCount),
						},
					},
				})

			_, err = bot.ReplyMessage(myLineRequest.Events[0].ReplyToken, templateMessage).Do()
			if err != nil {
				return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, nil
			}
			log.Println(data)
		}
	}
	return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: http.StatusOK}, nil
}

func main() {
	lambda.Start(Handler)
}
