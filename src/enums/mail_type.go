package enums

import "errors"

// MailType is an enum that contains Azure, Gmail or SMTP.
type MailType string

const (
	Azure MailType = "Azure"
	Gmail MailType = "Gmail"
	SMTP  MailType = "SMTP"
)

// ToMailType converts a string to a MailType enum.
func ToMailType(s string) (MailType, error) {
	switch s {
	case string(Azure):
		return Azure, nil
	case string(Gmail):
		return Gmail, nil
	case string(SMTP):
		return SMTP, nil
	default:
		return "", errors.New("invalid MailType")
	}
}

// ToString converts a MailType enum to a string.
func (mt *MailType) ToString() string {
	switch *mt {
	case Azure:
		return "Azure"
	case Gmail:
		return "Gmail"
	case SMTP:
		return "SMTP"
	default:
		return ""
	}
}
