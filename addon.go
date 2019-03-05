package main

import (
	// 	"fmt"
	"encoding/json"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

type CatalogItem struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

type Manifest struct {
	Id          string        `json:"id"`
	Version     string        `json:"version"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Types       []string      `json:"types"`
	Catalogs    []CatalogItem `json:"catalogs"`
	Resources   []string      `json:"resources"`
}

type MetaItem struct {
	Name   string   `json:"name"`
	Genres []string `json:"genres,omitempty"`
}

var CATALOG_ID = "Jupiter Broadcasting Shows"

var MANIFEST = Manifest{
	Id:          "org.stremio.video.jupiterbroadcasting",
	Version:     "0.0.1",
	Name:        "Jupiter Broadcasting",
	Description: "Watch shows from the Jupiter Broadcasting Network including Linux Action News, TechSNAP, Ask Noah, Coder Radio, and more.",
	Types:       []string{"series"},
	Catalogs:    []CatalogItem{},
	Resources:   []string{"stream", "catalog", "meta"},
}

var jupiterShows []*JupiterShow

func main() {
	jupiterShows = InitShows()

	//MANIFEST.Catalogs = append(MANIFEST.Catalogs, CatalogItem{"channels", CATALOG_ID})
	MANIFEST.Catalogs = append(MANIFEST.Catalogs, CatalogItem{"series", CATALOG_ID})

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/manifest.json", ManifestHandler)
	r.HandleFunc("/stream/{type}/{id}.json", StreamHandler)
	r.HandleFunc("/catalog/{type}/{id}.json", CatalogHandler)
	r.HandleFunc("/meta/{type}/{id}.json", MetaHandler)
	http.Handle("/", r)

	// CORS configuration
	headersOk := handlers.AllowedHeaders([]string{
		"Content-Type",
		"X-Requested-With",
		"Accept",
		"Accept-Language",
		"Accept-Encoding",
		"Content-Language",
		"Origin",
	})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET"})
	// Listen

	err := http.ListenAndServe("0.0.0.0:3592", handlers.CORS(originsOk, headersOk, methodsOk)(r))

	if err != nil {
		log.Fatalf("Listen: %s", err.Error())
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	type jsonObj map[string]interface{}

	jr, _ := json.Marshal(jsonObj{"Path": '/'})
	w.Header().Set("Content-Type", "application/json")
	w.Write(jr)
}

func ManifestHandler(w http.ResponseWriter, r *http.Request) {
	jr, _ := json.Marshal(MANIFEST)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jr)
}

func StreamHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	if params["type"] != "series" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var Show *JupiterShow = nil
	itemIds := strings.Split(params["id"], ":")
	showID := itemIds[0]

	for _, show := range jupiterShows {
		if show.Id == showID {
			Show = show
			break
		}
	}
	if Show == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var stream *StreamItem

	for _, episode := range Show.EpisodeList {
		if params["id"] == episode.Id {
			stream = &episode.Streams[0]
			break
		}
	}
	if stream == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"streams": [`))
	streamJson, _ := json.Marshal(stream)
	w.Write(streamJson)
	w.Write([]byte(`]}`))
}

func CatalogHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	if params["type"] != "series" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"metas": `))

	metas := []JupiterShow{}
	for _, show := range jupiterShows {
		item := JupiterShow{
			Id:     show.Id,
			Type:   params["type"],
			Name:   show.Name,
			Genres: show.Genres,
			Poster: show.Logo,
		}
		metas = append(metas, item)
	}

	catalogJson, _ := json.Marshal(metas)
	w.Write(catalogJson)
	w.Write([]byte(`}`))
}

func MetaHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	if params["type"] != "series" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	for _, show := range jupiterShows {
		if show.Id == params["id"] {
			UpdateEpisodes(*show)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"meta": `))
			streamJson, _ := json.Marshal(show)
			w.Write(streamJson)
			w.Write([]byte(`}`))
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}
