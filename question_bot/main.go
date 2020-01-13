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

var LineChannelSecret string
var LineAccessToken string

var questionMap = map[int]string{
	1:  "①「私なんてどうせ無理だ」と思うことがよくある",
	2:  "②自分の行動に対して、○○すべき、××しなければならないとよく思う",
	3:  "③いままで自分だけが頑張り、他の人は頑張らず体調を崩した事がある",
	4:  "④褒められると本当は良いと思っていても「そんな事ないよ」と謙遜してしまう",
	5:  "⑤すごく気を使っているのに人間関係が上手くいかない",
	6:  "⑥プレゼントの金額や感謝の言葉が重すぎて引かれてしまった事がある",
	7:  "⑦誰も私の心をわかってくれないと思う時がよくある",
	8:  "⑧自分を好きになってくれる異性などいない",
	9:  "⑨こんな自分がどうやって生きられるのかと不安になる",
	10: "⑩周囲の人は自分の揚げ足ばかりとるといつも思ってしまう",
}

var answerMap = map[int]string{
	1: "たまに自己の事を嫌いになってしまう事があるようです。\n自分の好きなものを見つけたり、不必要な謙遜をしなくてもいいように、自己分析をしてみませんか？",
	2: "人間関係にとても悩みがあるようです。\n自分の事をわかってくれない、まわりは全然頑張ってくれないと思う事はありませんか？\n人間関係、アサーティブなコミュニケーションを学んでみませんか？",
	3: "人間関係にとても疲れているようです。\n自分自身を傷つけてしまったり、目の前が真っ暗になるような事はありませんか？\nまずはご自分の事を認めてあげれるような良いところ探しが必要です、面談を通じてご自分を好きになれるよう行動パターンを変えてみませんか？",
}

func init() {
	LineChannelSecret = os.Getenv("CHANNELSECRET")
	LineAccessToken = os.Getenv("ACCESSTOKEN")
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
				log.Print("message: " + message.Text)
				templateMessage := linebot.NewTemplateMessage("confirms template",
					&linebot.ConfirmTemplate{
						Text: questionMap[1],
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

			if currentQuestionNo == len(questionMap) {
				if currentCount <= 2 {
					if _, err = bot.ReplyMessage(myLineRequest.Events[0].ReplyToken, linebot.NewTextMessage(answerMap[1])).Do(); err != nil {
						return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, nil
					}
				}
				if 3 <= currentCount && currentCount <= 4 {
					if _, err = bot.ReplyMessage(myLineRequest.Events[0].ReplyToken, linebot.NewTextMessage(answerMap[2])).Do(); err != nil {
						return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, nil
					}
				}
				if 5 <= currentCount {
					if _, err = bot.ReplyMessage(myLineRequest.Events[0].ReplyToken, linebot.NewTextMessage(answerMap[3])).Do(); err != nil {
						return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, nil
					}
				}

				pushText := "ご興味のある方はこちらをご覧ください。\nhttps://youtu.be/IOyI5H8sioc"
				if _, err = bot.PushMessage(myLineRequest.Events[0].Source.UserID, linebot.NewTextMessage(pushText)).Do(); err != nil {
					return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, nil
				}
				return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: http.StatusOK}, nil
			}

			templateMessage := linebot.NewTemplateMessage("confirm template",
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
		}
	}
	return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: http.StatusOK}, nil
}

func main() {
	lambda.Start(Handler)
}
