package main

import (
	"log"
	"net/http"
	"time"

	"github.com/a-know/moshi-moshi/handlers"
	"github.com/go-chi/chi"
)

const location = "Asia/Tokyo"

func main() {
	r := chi.NewRouter()

	// timezone
	loc, err := time.LoadLocation(location)
	if err != nil {
		loc = time.FixedZone(location, 9*60*60)
	}
	time.Local = loc

	r.Get("/heartbeat", handlers.HandleHeartbeat)
	r.Get("/moshimoshi/{site}/{path}", handlers.HandleMoshimoshi) // expect format : /moshimoshi/a_know_blog/2018/06/24/224424?title=xxxx

	log.Printf("moshi-moshi started.")
	http.ListenAndServe(":8080", r)
}
