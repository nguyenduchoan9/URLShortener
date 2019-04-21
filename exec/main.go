package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"sort"
	"strconv"
	"syscall"
	"time"

	"github.com/nguyendhoan9/coderschool.go/assignment.1/models"
	"github.com/nguyendhoan9/coderschool.go/assignment.1/urlmanager"
	"github.com/nguyendhoan9/coderschool.go/assignment.1/urlmanagerutil"
)

var yamlurlmanager = urlmanager.Yamlurlmanager{}

func main() {
	if command, valid := parseCommand(); valid {
		proceedCommand(&command)
	} else {
		printUsage()
		os.Exit(1)
	}
}

func parseCommand() (models.Command, bool) {
	addURL := flag.String("a", "", "Append URL.")
	redirectURLOfAddURL := flag.String("u", "", "Redirection URL of add URL.")
	removedURL := flag.String("d", "", "Remove URL.")
	listRedirections := flag.Bool("l", false, "List redirections.")
	help := flag.Bool("h", false, "Prints usage info.")
	port := flag.String("p", "", "Run HTTP server run on this port")
	flag.Parse()

	args := flag.Args()
	valid := true
	if flag.NArg() == 0 && flag.NFlag() == 0 {
		valid = false
	} else if len(args) > 0 {
		if args[0] == "configure" || args[0] == "run" {
			flag.CommandLine.Parse(args[1:])
		}
	}
	return models.Command{
		AddURL:                 *addURL,
		RedirectionURLOfAddURL: *redirectURLOfAddURL,
		RemoveURL:              *removedURL,
		ListRedirections:       *listRedirections,
		Port:                   *port,
		Help:                   *help,
		Args:                   args,
	}, valid
}

func printUsage() {
	fmt.Println("Usage: urlshorten [arg] [flag]")
	flag.PrintDefaults()
}

func proceedCommand(commad *models.Command) {
	if commad.IsHelpCommand() {
		printUsage()
	} else if commad.IsListRedirectionCommand() {
		listRedirectionURL()
	} else if commad.IsRemoveURLCommand() {
		urlmanagerutil.RemoveShortURL(yamlurlmanager, commad.RemoveURL)
	} else if commad.IsAddURLCommand() {
		if len(commad.AddURL) == 0 {
			urlmanagerutil.AddShortURL(yamlurlmanager, randStringRunes(8, []rune(removeNonChar(commad.RedirectionURLOfAddURL))), commad.RedirectionURLOfAddURL)
		} else {
			urlmanagerutil.AddShortURL(yamlurlmanager, commad.AddURL, commad.RedirectionURLOfAddURL)
		}
	} else if commad.IsStartServerInPort() {
		startServer(commad)
	} else {
		printUsage()
	}
}

func removeNonChar(words string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	return reg.ReplaceAllString(words, "")
}

func randStringRunes(n int, letterRunes []rune) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func listRedirectionURL() {
	shortendURLs := make([]urlmanager.URLShortener, 0)
	urlmanagerutil.GetURLShorteners(yamlurlmanager, &shortendURLs)
	sort.Slice(shortendURLs, func(firstIndex, secondIndex int) bool {
		return shortendURLs[firstIndex].Used > shortendURLs[secondIndex].Used
	})
	if len(shortendURLs) > 0 {
		printRow("ShortURL", "RedirectURL", "Used")
		printRow("--------", "---------------------------", "-------")
		for _, v := range shortendURLs {
			printRow(v.ShortURL, v.RedirectURL, strconv.Itoa(v.Used))
		}
	} else {
		fmt.Println("There is no short URL.")
	}
}

func printRow(col1, col2, col3 string) {
	fmt.Printf("|%-15s|%-40s|%-7s|\n", col1, col2, col3)
}

func startServer(cmd *models.Command) {
	log.Println("Starting serving at http://localhost:" + cmd.Port)
	http.HandleFunc("/", handlerRequest)

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		fmt.Printf("caught sig: %+v", sig)
		fmt.Println("Wait for 2 second to finish processing")
		time.Sleep(2 * time.Second)
		log.Println("The server stopped.")
		os.Exit(0)
	}()

	http.ListenAndServe(":"+cmd.Port, nil)
}

func handlerRequest(w http.ResponseWriter, req *http.Request) {
	log.Printf("Request: %s %s %s\n", req.URL.Hostname(), req.URL.Port(), req.URL.Path)
	defaultRedirectURL := "http://coderschool.vn"
	if len(req.URL.Path) > 0 {
		urlShortener := urlmanagerutil.GetURLShortenerBy(yamlurlmanager, req.URL.Path[1:])
		if len(urlShortener.RedirectURL) > 0 {
			urlmanagerutil.IncreaseTimeOfUsage(yamlurlmanager, urlShortener.ShortURL)
			http.Redirect(w, req, urlShortener.RedirectURL, http.StatusFound)
			return
		}
	}
	http.Redirect(w, req, defaultRedirectURL, http.StatusFound)
}
