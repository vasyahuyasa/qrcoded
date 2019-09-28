package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/skip2/go-qrcode"
)

var levels = map[string]qrcode.RecoveryLevel{
	"l": qrcode.Low,
	"m": qrcode.Medium,
	"q": qrcode.High,
	"h": qrcode.Highest,
}

type handler struct {
	debug bool
}

func run(addr string, debug bool) {
	h := handler{
		debug: debug,
	}

	srv := http.Server{
		Addr:    addr,
		Handler: http.HandlerFunc(h.handle),
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)

		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)

		<-sigint
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	log.Printf("HTTP server run on http://%s/\n", addr)
	if debug {
		log.Println("Debug is enabled")
	}
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}

func (h handler) handle(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/favicon.ico" {
		h.error(w, r, "no favicon", http.StatusNotFound)
		return
	}

	start := time.Now()

	if r.Method != http.MethodGet {
		h.error(w, r, "only GET method is supported", http.StatusMethodNotAllowed)
		return
	}

	// text for encoding
	text := r.FormValue("q")
	if text == "" {
		h.error(w, r, "q param is required", http.StatusBadRequest)
		return
	}

	// error correction level
	level := qrcode.Medium
	strLevel := r.FormValue("r")
	if strLevel != "" {
		var ok bool
		level, ok = levels[strLevel]
		if !ok {
			h.error(w, r, "r param must be one of l, m, q, h", http.StatusBadRequest)
			return
		}
	}

	// size of qr code in pixels
	size := 256
	strSize := r.FormValue("s")
	if strSize != "" {
		var err error
		size, err = strconv.Atoi(strSize)
		if err != nil {
			h.error(w, r, "s param must number", http.StatusBadRequest)
			return
		}
	}

	data, err := qrcode.Encode(text, level, size)
	if err != nil {
		h.error(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-type", "image/png")
	_, err = w.Write(data)
	if err != nil {
		log.Printf("from=%s error=can not send image to client: %v\n", r.RemoteAddr, err)
		return
	}

	if h.debug {
		log.Printf("request=%q from=%s text=%s size=%d level=%d time=%v", r.RequestURI, r.RemoteAddr, text, size, level, time.Since(start))
	}
}

func (h handler) error(w http.ResponseWriter, r *http.Request, error string, code int) {
	http.Error(w, error, code)

	if h.debug {
		log.Printf("request=%q from=%s error=%s\n", r.RequestURI, r.RemoteAddr, error)
	}
}
