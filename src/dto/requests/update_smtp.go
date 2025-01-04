package requests

import "time"

// UpdateSmtp struct for updating a SMTP record.
type UpdateSmtp struct {
	App                  string    `json:"app" validate:"required"`
	Mail                 string    `json:"mail" validate:"required,email"`
	Username             string    `json:"username" validate:"required"`
	Password             string    `json:"password"`
	Host                 string    `json:"host" validate:"required"`
	Port                 int       `json:"port" validate:"required"`
	DkimPrivateKey       string    `json:"dkimPrivateKey"`
	DkimDomain           string    `json:"dkimDomain"`
	DkimCanonicalization string    `json:"dkimCanonicalization"`
	Primary              bool      `json:"primary"`
	UpdatedAt            time.Time `json:"updatedAt"`
}
