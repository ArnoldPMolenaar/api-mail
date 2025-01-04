package enums

// DkimCanonicalization is an enum that contains Simple and Relaxed for the Dkim settings.
type DkimCanonicalization string

const (
	Simple  DkimCanonicalization = "Simple"
	Relaxed DkimCanonicalization = "Relaxed"
)

// ToDkimCanonicalization converts a string to a DkimCanonicalization enum.
// Default is Relaxed on failure.
func ToDkimCanonicalization(s string) DkimCanonicalization {
	switch s {
	case string(Simple):
		return Simple
	case string(Relaxed):
		return Relaxed
	default:
		return Relaxed
	}
}
