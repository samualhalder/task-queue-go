package model

type SendMailTemplate struct {
	Template string         `json:"template"`
	To       []string       `json:"to"`
	Subject  string         `json:"subject"`
	Data     map[string]any `json:"data"`
}

