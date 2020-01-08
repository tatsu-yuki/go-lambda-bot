package template

import "github.com/line/line-bot-sdk-go/linebot"

type Confirms struct {
	Title string
	Left  linebot.TemplateAction
	Right linebot.TemplateAction
}

// new default confirms template
func NewConfirms() Confirms {
	return Confirms{
		Title: "confirm",
		Left: &linebot.MessageAction{
			Label: "1",
			Text:  "Yes",
		},
		Right: &linebot.MessageAction{
			Label: "2",
			Text:  "No",
		},
	}
}

// set left template action
func (c *Confirms) SetLeft(label, text string) {
	c.Left = &linebot.MessageAction{
		Label: label,
		Text:  text,
	}
}

// set right template action
func (c *Confirms) SetRight(label, text string) {
	c.Right = &linebot.MessageAction{
		Label: label,
		Text:  text,
	}
}

func (c *Confirms) ConfirmsTemplate() *linebot.TemplateMessage {
	return linebot.NewTemplateMessage("confirms template",
		&linebot.ConfirmTemplate{
			Text:    c.Title,
			Actions: []linebot.TemplateAction{c.Left, c.Right},
		})
}
