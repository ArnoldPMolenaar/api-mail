package requests

type SendMail struct {
	App      string   `json:"app" validate:"required"`
	Mail     string   `json:"mail"`
	Type     string   `json:"type"`
	FromName string   `json:"fromName"`
	FromMail string   `json:"fromMail" validate:"email"`
	To       string   `json:"to" validate:"required,email"`
	Subject  string   `json:"subject" validate:"required"`
	Body     string   `json:"body" validate:"required"`
	MimeType string   `json:"mimeType"`
	Ccs      []string `json:"ccs"`
	Bccs     []string `json:"bccs"`
}
