package urlmanagerutil

import (
	"github.com/nguyendhoan9/coderschool.go/assignment.1/urlmanager"
)

func GetURLShortenerBy(urlManager urlmanager.URLManager, shortURL string) urlmanager.URLShortener {
	return urlManager.GetURLShortenerBy(shortURL)
}

func IncreaseTimeOfUsage(urlManager urlmanager.URLManager, shortURL string) {
	urlManager.IncreaseTimeOfUsage(shortURL)
}

func GetURLShorteners(urlManager urlmanager.URLManager, URLShorteners *[]urlmanager.URLShortener) error {
	return urlManager.GetURLShorteners(URLShorteners)
}

func AddShortURL(urlmanager urlmanager.URLManager, shortURL, redirectURL string) {
	urlmanager.AddShortURL(shortURL, redirectURL)
}

func RemoveShortURL(urlmanager urlmanager.URLManager, removeURL string) {
	urlmanager.RemoveShortURL(removeURL)
}
