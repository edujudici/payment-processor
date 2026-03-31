package ports

type AuthInterface interface {
	GetAccessToken() (*string, error)
}
