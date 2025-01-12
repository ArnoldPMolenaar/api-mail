package enums

// AppMailPrimaryType is an enum that contains Azure, Gmail or SMTP.
type AppMailPrimaryType string

const (
	Azure AppMailPrimaryType = "Azure"
	Gmail AppMailPrimaryType = "Gmail"
	SMTP  AppMailPrimaryType = "SMTP"
)

var (
	azure = Azure
	gmail = Gmail
	smtp  = SMTP
)

// ToAppMailPrimaryType converts a string to a MailType enum.
func ToAppMailPrimaryType(s *string) *AppMailPrimaryType {
	switch *s {
	case string(Azure):
		return &azure
	case string(Gmail):
		return &gmail
	case string(SMTP):
		return &smtp
	default:
		return nil
	}
}

// ToString converts a MailType enum to a string.
func (primaryType *AppMailPrimaryType) ToString() *string {
	var primary string

	switch *primaryType {
	case Azure:
		primary = "Azure"
	case Gmail:
		primary = "Gmail"
	case SMTP:
		primary = "SMTP"
	default:
		return nil
	}

	return &primary
}
