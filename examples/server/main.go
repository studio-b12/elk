package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/studio-b12/elk"
)

func main() {
	db := NewDatabase("db.json")
	ctl := NewController(db)

	mux := http.NewServeMux()

	mux.HandleFunc("/count", handleCount(ctl))

	http.ListenAndServe(":8080", mux)
}

func handleCount(ctl *Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch strings.ToUpper(r.Method) {
		case "GET":
			handleGetCount(ctl, w, r)
		case "POST":
			handlePostCount(ctl, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func handleGetCount(ctl *Controller, w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := ctl.GetCount(id)
	if err != nil {
		switch elk.Cast(err).Code() {
		case ErrorCountNotFound:
			w.WriteHeader(http.StatusNotFound)
		default:
			log.Printf("error: %+.5v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write(elk.MustJson(err))
		return
	}

	d, _ := json.MarshalIndent(res, "", "  ")
	w.Write(d)
}

func handlePostCount(ctl *Controller, w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := ctl.IncrementCount(id)
	if err != nil {
		switch elk.Cast(err).Code() {
		default:
			log.Printf("error: %#.5v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write(elk.MustJson(err))
		return
	}

	d, _ := json.MarshalIndent(res, "", "  ")
	w.Write(d)
}
