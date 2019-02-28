package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"log"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type CatalogItem struct {
    Type	string			`json:"type"`
    Id		string			`json:"id"`
}

type Manifest struct {
    Id		string		`json:"id"`
    Version	string		`json:"version"`
    Name	string		`json:"name"`
    Description	string		`json:"description"`
    Types	[]string	`json:"types"`
    Catalogs	[]CatalogItem	`json:"catalogs"`
    Resources	[]string	`json:"resources"`
}

type MetaItem struct {
    Name	string			`json:"name"`
    Genres	[]string		`json:"genres,omitempty"`
}

var CATALOG_ID = "Jupiter Broadcasting Shows"

var MANIFEST = Manifest{
	Id:		"org.stremio.video.jupiterbroadcasting",
	Version:	"0.0.1",
	Name:		"Jupiter Broadcasting",
	Description:	"Watch shows from the Jupiter Broadcasting Network including Linux Action News, TechSNAP, Ask Noah, Coder Radio, and more.",
	Types:		[]string{ "series" },
	Catalogs:	[]CatalogItem{},
	Resources:	[]string{ "stream", "catalog", "meta" },
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

	fmt.Printf("func StreamHandler(")

	if params["type"] != "series" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fmt.Printf("StreamHandler: " + params["type"] + ":     id is " + params["id"] + "\n")
	
	var showID int
	var episodeId int
	var Show *JupiterShow = nil

	fmt.Sscanf(params["id"], "%d:%d", &showID, &episodeId)
	fmt.Printf("StreamHandler: after scan show id " + string(showID) + "  :    episode id  " + string(episodeId))
	for _, show := range jupiterShows {
		fmt.Printf("No show " + show.Id + " " + string(showID) )
		if show.Id == string(showID) {
			Show = show
			fmt.Printf("Found show " + show.Id + string(showID))
			break
		}
	}
	if Show == nil {
		fmt.Printf("No show ")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	stream := StreamItem{}
	for _, episode := range Show.EpisodeList {
		if episodeId == episode.Episode {
			stream = episode.Streams[0]
			break
		}
	}
	if stream == (StreamItem{}) {
		fmt.Printf("No steam/episode ")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([] byte(`{"streams": [`))
	streamJson, _ := json.Marshal(stream)
	w.Write(streamJson)
	w.Write([] byte(`]}`))
}

func CatalogHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	fmt.Printf("func CatalogHandler(")

	if params["type"] != "series" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([] byte(`{"metas": `))

	metas := []JupiterShow{}
	for _, show := range jupiterShows {
		item := JupiterShow{
			Id: show.Id,
			Type: params["type"],
			Name: show.Name,
			Genres: show.Genres,
			Poster: show.Logo,
		}
		metas = append(metas, item)
	}

	catalogJson, _ := json.Marshal(metas)
	w.Write(catalogJson)
	w.Write([] byte(`}`))
}

func MetaHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	fmt.Printf("func MetaHandler(" + params["id"])
	//XXX: update episodes
	if params["type"] != "series" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	for _, show := range jupiterShows {
		fmt.Printf("MetaHandler() " +  show.Id)
		if show.Id == params["id"] {
			fmt.Printf(show.Name + ": episodes " + string(len(show.EpisodeList)))
			w.Header().Set("Content-Type", "application/json")
			w.Write([] byte(`{"meta": `))
			streamJson, _ := json.Marshal(show)
			w.Write(streamJson)
			w.Write([] byte(`}`))
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}
