package urlmanager

type URLManager interface {
	GetURLShortenerBy(shortURL string) URLShortener
	IncreaseTimeOfUsage(shortURL string)
	GetURLShorteners(URLShorteners *[]URLShortener) error
	AddShortURL(shortURL, redirectURL string)
	RemoveShortURL(removeURL string)
}
