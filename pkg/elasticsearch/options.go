package elasticsearch

// Options is a set of flags to configure a Client.
type Options struct {
	Host     string
	User     string
	Password string
	Debug    bool
}
