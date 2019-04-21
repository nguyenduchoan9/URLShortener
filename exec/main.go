package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"syscall"
	"time"

	"github.com/nguyendhoan9/coderschool.go/assignment.1/urlmanager"
	"github.com/nguyendhoan9/coderschool.go/assignment.1/urlmanagerutil"
)

type command struct {
	addURL, redirectionURLOfAddURL, removeURL, port string
	listRedirections, help                          bool
	args                                            []string
}

func (cmd *command) isHelpCommand() bool {
	return cmd.help
}

func (cmd *command) isListRedirectionCommand() bool {
	return cmd.listRedirections
}

func (cmd *command) isAddURLCommand() bool {
	return len(cmd.addURL) > 0 && len(cmd.redirectionURLOfAddURL) > 0 && len(cmd.args) > 0 && cmd.args[0] == "configure"
}

func (cmd *command) isRemoveURLCommand() bool {
	return len(cmd.removeURL) > 0
}

func (cmd *command) isStartServerInPort() bool {
	return len(cmd.port) > 0 && len(cmd.args) > 0 && cmd.args[0] == "run"
}

var yamlurlmanager = urlmanager.Yamlurlmanager{}

func main() {
	if command, valid := parseCommand(); valid {
		proceedCommand(&command)
	} else {
		printUsage()
		os.Exit(1)
	}
}

func parseCommand() (command, bool) {
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
	return command{
		addURL:                 *addURL,
		redirectionURLOfAddURL: *redirectURLOfAddURL,
		removeURL:              *removedURL,
		listRedirections:       *listRedirections,
		port:                   *port,
		help:                   *help,
		args:                   args,
	}, valid
}

func printUsage() {
	fmt.Println("Usage: urlshorten [arg] [flag]")
	flag.PrintDefaults()
}

func proceedCommand(commad *command) {
	if commad.isHelpCommand() {
		printUsage()
	} else if commad.isListRedirectionCommand() {
		listRedirectionURL()
	} else if commad.isRemoveURLCommand() {
		urlmanagerutil.RemoveShortURL(yamlurlmanager, commad.removeURL)
	} else if commad.isAddURLCommand() {
		urlmanagerutil.AddShortURL(yamlurlmanager, commad.addURL, commad.redirectionURLOfAddURL)
	} else if commad.isStartServerInPort() {
		startServer(commad)
	} else {
		printUsage()
	}
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

func startServer(cmd *command) {
	log.Println("Starting serving at http://localhost:" + cmd.port)
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

	http.ListenAndServe(":"+cmd.port, nil)
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
