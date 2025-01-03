package requests

// CreateSmtp struct for creating a new SMTP.
type CreateSmtp struct {
	App                  string `json:"app" validate:"required"`
	Mail                 string `json:"mail" validate:"required,email"`
	Username             string `json:"username" validate:"required"`
	Password             string `json:"password" validate:"required"`
	Host                 string `json:"host" validate:"required"`
	Port                 int    `json:"port" validate:"required"`
	DkimPrivateKey       string `json:"dkimPrivateKey"`
	DkimDomain           string `json:"dkimDomain"`
	DkimCanonicalization string `json:"dkimCanonicalization"`
	Primary              bool   `json:"primary"`
}
