package main

import (
	"embed"
	"flag"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/crypto/acme/autocert"
)

var domain string
var dev bool
var log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Timestamp().Logger()

//go:embed static/*
var static embed.FS
var fs http.FileSystem

func redirectHTTPtoHTTPS(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+domain+":443"+r.RequestURI, http.StatusMovedPermanently)
}

func main() {
	flag.StringVar(&domain, "domain", os.Getenv("DOMAIN"), "domain name to process HTTP/s server")
	flag.BoolVar(&dev, "dev", false, "enable dev mode")
	flag.Parse()

	devStr := os.Getenv("DEV")
	if devStr == "1" {
		dev = true
	}

	if len(domain) == 0 && !dev {
		log.Fatal().Msg("The domain parameter is required")
	}

	// Redirect HTTP to HTTPS
	if !dev {
		go func() {
			if err := http.ListenAndServe(":80", http.HandlerFunc(redirectHTTPtoHTTPS)); err != nil {
				log.Fatal().Err(err).Msg("ListenAndServe error")
			}
		}()
	}

	// HTTPS server
	mux := http.NewServeMux()
	if dev {
		fs = http.Dir("./static")
		fsHandler := http.FileServer(fs)
		mux.Handle("/static/", http.StripPrefix("/static/", fsHandler))
	} else {
		fs = http.FS(static)
		fsHandler := http.FileServer(fs)
		mux.Handle("/static/", fsHandler)
	}
	mux.HandleFunc("/", serveTemplate)

	if dev {
		err := http.ListenAndServe(":8080", mux)
		log.Fatal().Err(err).Msg("Server error")
	} else {
		err := http.Serve(autocert.NewListener(domain), mux)
		log.Fatal().Err(err).Msg("Server error")
	}
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		r.URL.Path = "/index.html"
	}

	lp := filepath.Join("layouts", "main.html")
	fp := filepath.Join("templates", filepath.Clean(r.URL.Path))

	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			httpNotFound(w)
			return
		} else {
			log.Error().Err(err).Msg("Cannot get file stat")
		}
	}

	if info.IsDir() {
		httpNotFound(w)
		return
	}

	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		log.Error().Err(err).Msg("Cannot parse template files")
		http.Error(w, http.StatusText(500), 500)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		log.Error().Err(err).Msg("Cannot execute teamplate")
		http.Error(w, http.StatusText(500), 500)
	}
}

func httpNotFound(w http.ResponseWriter) {
	var f io.ReadCloser
	var err error

	if dev {
		f, err = fs.Open("/pages/404.html")
	} else {
		f, err = static.Open("static/pages/404.html")
	}
	if err != nil {
		log.Error().Err(err).Msg("Cannot open 404.html file")
		return
	}
	defer f.Close()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusNotFound)
	
	if _, err := io.Copy(w, f); err != nil {
		log.Error().Err(err).Msg("Cannot read 404.html file")
		return
	}
}