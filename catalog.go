package main

import (
	"io/ioutil"
	"encoding/json"
	"net/http"
	"fmt"
//	"log"
)

var FEED_BURNER = "http://feeds2.feedburner.com/"
var JUPITER_COM = "http://www.jupiterbroadcasting.com/"
var DEFAULT_GENRES = []string{ "Education", "Technology" }
var FEED_BASE = "http://feed.jupiter.zone/"
var FEED_FORMAT = "?format=json"


type FeedAttachment struct {
	Url		string		`json:"id"`
	MimeType	string		`json:"mime_type"`
	SizeInBytes	uint64		`json:"size_in_bytes"`
}

type FeedEpisode struct {
	Id		string			`json:"id"`
	Url		string			`json:"url"`
	Html		string			`json:"content_html"`
	Published	string			`json:"date_published"`
	Attachments	[]FeedAttachment	`json:"attachments,omitempty"`
}

type Episode struct {
	Id		string		`json:"id"`
}

// Map each show to a channel item
type JupiterShow struct {
	Id		string		`json:"id"`
	Type		string		`json:"type"`
	Name		string		`json:"name"`
	Description	string		`json:"description"`
	Logo		string		`json:"logo"`
	Genres		[]string	`json:"genres"`
	EpisodeList	[]Episode	`json:"videos,omitempty"`
	Feed		string		`json:"-"`
}

type JsonResponse struct {
	Version		string		`json:"version"`
	Title		string		`json:"title"`
	HomePage	string		`json:"home_page_url"`
	FeedUrl		string		`json:"feed_url"`
	Episodes	[]Episode	`json:"items"`
}

func initShows() (showList []JupiterShow){
	showList = make([]JupiterShow, 0 )

	showList = append(showList, JupiterShow{
		Id: "30020",
		Type: "channel",
		Name: "BSD Now",
		Description: "A weekly show covering the latest developments in the world of the BSD family of operating systems. News, Tutorials and Interviews for new users and long time developers alike.",
		Logo: "https://raw.githubusercontent.com/JupiterBroadcasting/plugin.video.jupiterbroadcasting/krypton/resources/media/bsd-now.jpg",
		Genres: DEFAULT_GENRES,
		Feed: FEED_BASE + "bsdvid" + FEED_FORMAT,
	})

	showList = append(showList, JupiterShow{
		Id: "30017",
		Type: "channel",
		Name: "Coder Radio",
		Description: "A weekly talk show taking a pragmatic look at the art and business of Software Development and related technologies.",
		Logo: "https://raw.githubusercontent.com/JupiterBroadcasting/plugin.video.jupiterbroadcasting/krypton/resources/media/coder-radio.jpg",
		Genres: DEFAULT_GENRES,
		Feed: FEED_BASE + "coderradiovideo" + FEED_FORMAT,
	})

	showList = append(showList, JupiterShow{
		Id: "30019",
		Type: "channel",
		Name: "LINUX Unplugged",
		Description: "The Linux Action Show with no prep, no limits, and tons of opinion. An open show powered by community LINUX Unplugged takes the best attributes of open collaboration and focuses them into a weekly lifestyle show about Linux.",
		Logo: "https://raw.githubusercontent.com/JupiterBroadcasting/plugin.video.jupiterbroadcasting/krypton/resources/media/linux-unplugged.jpg",
		Genres: DEFAULT_GENRES,
		Feed: FEED_BASE + "linuxunvid" + FEED_FORMAT,
	})

	showList = append(showList, JupiterShow{
		Id: "30008",
		Type: "channel",
		Name: "TechSNAP",
		Description: "TechSNAP our weekly Systems, Network, and Administration Podcast. Every week TechSNAP covers the stories that impact those of us in the tech industry, and all of us that follow it. Every episode we dedicate a portion of the show to answer audience questions, discuss best practices, and solving your problems.",
		Logo: "https://raw.githubusercontent.com/JupiterBroadcasting/plugin.video.jupiterbroadcasting/krypton/resources/media/techsnap.jpg",
		Genres: DEFAULT_GENRES,
		Feed: FEED_BASE + "techsnapvid" + FEED_FORMAT,
	})

	showList = append(showList, JupiterShow{
		Id: "30008",
		Type: "channel",
		Name: "User Error",
		Description: "Life is a series of mistakes, but that's what makes it interesting. A show about life, Linux, the universe, and everything in between.",
		Logo: "https://raw.githubusercontent.com/JupiterBroadcasting/plugin.video.jupiterbroadcasting/krypton/resources/media/usererror.png",
		Genres: DEFAULT_GENRES,
		Feed: FEED_BASE + "uevideo" + FEED_FORMAT,
	})

	return showList
}

func fetchEpisodes(show JupiterShow) ([]Episode) {
	resp, _ := http.Get(show.Feed)
	bytes, _ := ioutil.ReadAll(resp.Body)


	fmt.Println("HTML:\n\n", string(bytes))

	var jsonResponse JsonResponse
	err := json.Unmarshal(bytes, &jsonResponse)
	
	fmt.Println(err)

	return jsonResponse.Episodes
}

