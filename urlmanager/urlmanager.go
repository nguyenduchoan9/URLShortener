package urlmanager

type URLManager interface {
	GetURLShorteners(URLShorteners *[]URLShortener) error
	AddShortURL(shortURL, redirectURL string)
	RemoveShortURL(removeURL string)
}
