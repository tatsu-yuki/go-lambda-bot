package main

import (
	"log"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/tatsu-yuki/question_bot/template"
)

func textHandler(message *linebot.TextMessage, replyToken string) error {
	var msg linebot.Message
	switch message.Text {
	case "buttons":
		btn := template.NewButtons()
		if err := btn.AddButtons(
			linebot.NewPostbackTemplateAction("Say hello1", "hello こんにちは", "", "hello こんにちは"),
			linebot.NewPostbackTemplateAction("言 hello2", "hello こんにちは", "hello こんにちは", ""),
			linebot.NewPostbackTemplateAction("言 hello2", "hello こんにちは", "hello こんにちは", ""),
			linebot.NewPostbackTemplateAction("言 hello2", "hello こんにちは", "hello こんにちは", ""),
		); err != nil {
			return err
		}
		msg = btn.ButtonsTemplate()

	case "confirm":
		confirm := template.NewConfirms()
		msg = confirm.ConfirmsTemplate()

	case "carousel":
		carousel := template.NewCarousel()
		btn := template.NewButtons()
		if err := btn.AddButtons(
			linebot.NewPostbackTemplateAction("Say hello1", "hello こんにちは", "", "hello こんにちは"),
			linebot.NewPostbackTemplateAction("言 hello2", "hello こんにちは", "hello こんにちは", ""),
			linebot.NewPostbackTemplateAction("言 hello2", "hello こんにちは", "hello こんにちは", ""),
		); err != nil {
			return err
		}

		if err := carousel.SetColumns(btn, btn); err != nil {
			return err
		}
		log.Println(carousel.Columns)
		msg = carousel.CarouselTemplate()

	case "image carousel":
		col := template.NewImageColumns()
		col.SetImageAction(linebot.NewURITemplateAction("Go to LINE", "https://line.me"))

		c := template.NewImageCarousel()
		if err := c.SetImageCarousel(col, col, col); err != nil {
			return err
		}
		msg = c.CarouselTemplate()

	case "flex json":
		//TODO: adding flex_json template func

		contents, err := linebot.UnmarshalFlexMessageJSON([]byte(template.jsonString))
		if err != nil {
			return err
		}
		if _, err := bot.ReplyMessage(
			replyToken,
			linebot.NewFlexMessage("Flex message alt text", contents),
		).Do(); err != nil {
			return err
		}

	default:
		log.Println(message.Text)
	}
	if msg != nil {
		if _, err := bot.ReplyMessage(replyToken, msg).Do(); err != nil {
			return err
		}
	}
	return nil
}