package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
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
	yamlurlmanager := urlmanager.Yamlurlmanager{}
	if commad.isHelpCommand() {
		flag.PrintDefaults()
	} else if commad.isListRedirectionCommand() {
		listRedirectionURL(yamlurlmanager)
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

func listRedirectionURL(urlManager urlmanager.Yamlurlmanager) {
	shortendURLs := make([]urlmanager.URLShortener, 0)
	urlmanagerutil.GetURLShorteners(urlManager, &shortendURLs)
	if len(shortendURLs) > 0 {
		fmt.Println("ShortURL\t\tRedirectURL\t\t\t\tUsed")
		for _, v := range shortendURLs {
			fmt.Printf("%s\t\t\t%s\t\t\t%d\n", v.ShortURL, v.RedirectURL, v.Used)
		}
	} else {
		fmt.Println("There is no short URL.")
	}
}

func startServer(cmd *command) {
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	logger.Println("Server is starting...")

	router := http.NewServeMux()
	router.Handle("/onetwo", handlerRequest())

	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	server := &http.Server{
		Addr:         ":" + cmd.port,
		Handler:      tracing(nextRequestID)(logging(logger)(router)),
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		logger.Println("Server is shutting down...")
		atomic.StoreInt32(&healthy, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	logger.Println("Server is ready to handle requests at", server.Addr)
	atomic.StoreInt32(&healthy, 1)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", server.Addr, err)
	}

	<-done
	logger.Println("Server stopped")

	// http.HandleFunc("/", handlerRequest)
	// http.ListenAndServe(":"+cmd.port, logRequest(http.DefaultServeMux))
}

func handlerRequest() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// if req.URL.Path != "/" {
		// 	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		// 	return
		// }
		http.Redirect(w, req, "https://zing.vn", http.StatusMovedPermanently)
	})
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

const requestIDKey int = 0

var (
	healthy int32
)

func logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				requestID, ok := r.Context().Value(requestIDKey).(string)
				if !ok {
					requestID = "unknown"
				}
				logger.Println(requestID, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func tracing(nextRequestID func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = nextRequestID()
			}
			ctx := context.WithValue(r.Context(), requestIDKey, requestID)
			w.Header().Set("X-Request-Id", requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
