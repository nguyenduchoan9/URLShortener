package urlmanager

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const mURLPersistentName = "URLShorter.yml"

type Yamlurlmanager struct{}

func (manager Yamlurlmanager) GetURLShortenerBy(shortURL string) URLShortener {
	shortenerurls := make([]URLShortener, 0)
	if err := manager.GetURLShorteners(&shortenerurls); err == nil {
		for _, v := range shortenerurls {
			if v.ShortURL == shortURL {
				return v
			}
		}
	}
	return URLShortener{}
}

func (manager Yamlurlmanager) IncreaseTimeOfUsage(shortURL string) {
	shortenerurls := make([]URLShortener, 0)
	if err := manager.GetURLShorteners(&shortenerurls); err == nil {
		success := false
		index := 0
		for i, v := range shortenerurls {
			if v.ShortURL == shortURL {
				index = i
				success = true
			}
		}
		if success {
			shortenerurls[index].Used++
			saveURLShortener(&shortenerurls)
		}
	}
}

func (manager Yamlurlmanager) GetURLShorteners(URLShorteners *[]URLShortener) error {
	if _, err := os.Stat(getYAMLFilePath()); err == nil {
		yamlFileBytes, readError := ioutil.ReadFile((getYAMLFilePath()))
		if readError == nil {
			yaml.Unmarshal(yamlFileBytes, URLShorteners)
			return nil
		}
		return readError
	} else if os.IsNotExist(err) {
		os.Create(mURLPersistentName)
	}
	return nil
}

func (manager Yamlurlmanager) AddShortURL(shortURL, redirectURL string) {
	shortenerurls := make([]URLShortener, 0)
	if err := manager.GetURLShorteners(&shortenerurls); err == nil {
		exist := false
		for _, v := range shortenerurls {
			if v.ShortURL == shortURL {
				exist = true
				break
			}
		}
		if !exist {
			shortenerurls = append(shortenerurls, URLShortener{
				ShortURL:    shortURL,
				RedirectURL: redirectURL,
				Used:        0,
			})
			saveURLShortener(&shortenerurls)
		}
	}
}

func (manager Yamlurlmanager) RemoveShortURL(removeURL string) {
	shortenerurls := make([]URLShortener, 0)
	if err := manager.GetURLShorteners(&shortenerurls); err == nil {
		index := -1
		for i, v := range shortenerurls {
			if v.ShortURL == removeURL {
				index = i
				break
			}
		}
		if index != -1 {
			shortenerurls = append(shortenerurls[:index], shortenerurls[index+1:]...)
			saveURLShortener(&shortenerurls)
		}
	}
}

func saveURLShortener(shortenerurls *[]URLShortener) {
	if data, err := yaml.Marshal(*shortenerurls); err == nil {
		ioutil.WriteFile(getYAMLFilePath(), data, os.ModeDevice)
	} else {
		fmt.Println("Error occured while saving state", err)
	}
}
func getYAMLFilePath() string {
	pwd, _ := os.Getwd()
	return filepath.Join(pwd, mURLPersistentName)
}
