package urlmanagerutil

import (
	"github.com/nguyendhoan9/coderschool.go/assignment.1/urlmanager"
)

func GetURLShorteners(urlManager urlmanager.URLManager, URLShorteners *[]urlmanager.URLShortener) error {
	return urlManager.GetURLShorteners(URLShorteners)
}

func AddShortURL(urlmanager urlmanager.URLManager, shortURL, redirectURL string) {
	urlmanager.AddShortURL(shortURL, redirectURL)
}

func RemoveShortURL(urlmanager urlmanager.URLManager, removeURL string) {
	urlmanager.RemoveShortURL(removeURL)
}
