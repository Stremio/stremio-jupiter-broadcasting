package main

import (
	"io/ioutil"
	"encoding/json"
	"net/http"
	"fmt"
	"log"
)

type FeedAttachment struct {
	Url		string		`json:"url"`
	MimeType	string		`json:"mime_type"`
	SizeInBytes	uint64		`json:"size_in_bytes"`
}

type FeedEpisode struct {
	Id		string			`json:"id"`
	Url		string			`json:"url"`
	Title		string			`json:"title"`
	Html		string			`json:"content_html"`
	Published	string			`json:"date_published"`
	Attachments	[]FeedAttachment	`json:"attachments,omitempty"`
}

type Episode struct {
	Id		string		`json:"id"` //XXX
	Season		int		`json:"season"`
	Episode		int		`json:"episode"`
	Title		string		`json:"title"`
	Thumbnail	string		`json:"thumbnail,omitempty"`
	Released	string		`json:"released"`
	Streams		[]StreamItem	`json:"streams,omitempty"`
	Overview	string		`json:"overview,omitempty"`
}

// Map each show to a series item
type JupiterShow struct {
	Id		string		`json:"id"`
	Type		string		`json:"type"`
	Name		string		`json:"name"`
	Description	string		`json:"description"`
	Logo		string		`json:"logo"`
	Poster      	string		`json:"poster"`
	Genres		[]string	`json:"genres"`
	EpisodeList	[]Episode	`json:"videos,omitempty"`
	Feed		string		`json:"-"`
}

type JsonResponse struct {
	Version		string		`json:"version"`
	Title		string		`json:"title"`
	HomePage	string		`json:"home_page_url"`
	FeedUrl		string		`json:"feed_url"`
	Episodes	[]FeedEpisode	`json:"items"`
}

type StreamItem struct {
    Title	string			`json:"title"`
    InfoHash	string			`json:"infoHash,omitempty"`
    FileIdx	uint8			`json:"fileIdx,omitempty"`
    Url		string			`json:"url"`
    YtId	string			`json:"ytId,omitempty"`
    ExternalUrl	string			`json:"externalUrl,omitempty"`
}

//var FEED_BURNER = "http://feeds2.feedburner.com/"
var JUPITER_COM = "http://www.jupiterbroadcasting.com/"
var DEFAULT_GENRES = []string{ "Education", "Technology" }
var FEED_BASE = "http://feed.jupiter.zone/"
var ALT_FEED_BASE = "http://feedpress.me/"
var FEED_FORMAT = "?format=json"

func InitShows() (showList []*JupiterShow) {
	showList = append(showList, &JupiterShow{
		Id: "30020",
		Type: "series",
		Name: "BSD Now",
		Description: "A weekly show covering the latest developments in the world of the BSD" +
			     "family of operating systems. News, Tutorials and Interviews for new " +
			     " users and long time developers alike.",
		Logo: "https://static.feedpress.it/logo/bsdnowvid-5abd82a07a0e6.jpg",
		Poster: "https://static.feedpress.it/logo/bsdnowvid-5abd82a07a0e6.jpg", //XXX
		Genres: DEFAULT_GENRES,
		Feed: FEED_BASE + "bsdvid" + FEED_FORMAT,
	})

	showList = append(showList, &JupiterShow{
		Id: "30017",
		Type: "series",
		Name: "Coder Radio",
		Description: "A weekly talk show taking a pragmatic look at the art and business of" +
				" Software Development and related technologies.",
		Logo: "https://static.feedpress.it/logo/codervideo-5aafe52c954f4.jpg",
		Poster: "https://static.feedpress.it/logo/codervideo-5aafe52c954f4.jpg",
		Genres: DEFAULT_GENRES,
		Feed: FEED_BASE + "codervid" + FEED_FORMAT,
	})

	showList = append(showList, &JupiterShow{
		Id: "30019",
		Type: "series",
		Name: "LINUX Unplugged",
		Description: "The Linux Action Show with no prep, no limits, and tons of opinion." + 
			" An open show powered by community LINUX Unplugged takes the best attributes" + 
			" of open collaboration and focuses them into a weekly lifestyle show about Linux.",
		Logo: "https://static.feedpress.it/logo/lupvid-5ab1c61d12ac2.jpg",
		Poster: "https://static.feedpress.it/logo/lupvid-5ab1c61d12ac2.jpg",
		Genres: DEFAULT_GENRES,
		Feed: FEED_BASE + "lupvid" + FEED_FORMAT,
	})

	showList = append(showList, &JupiterShow{
		Id: "30008",
		Type: "series",
		Name: "TechSNAP",
		Description: "TechSNAP our weekly Systems, Network, and Administration Podcast. Every week" +
			" TechSNAP covers the stories that impact those of us in the tech industry, and all" +
			" of us that follow it. Every episode we dedicate a portion of the show to answer " +
			"audience questions, discuss best practices, and solving your problems.",
		Logo: "https://static.feedpress.it/logo/techsnapvid-5a208ae7c62dc.jpg",
		Poster: "https://static.feedpress.it/logo/techsnapvid-5a208ae7c62dc.jpg",
		Genres: DEFAULT_GENRES,
		Feed: ALT_FEED_BASE + "techsnapvid" + FEED_FORMAT,
	})

	showList = append(showList, &JupiterShow{
		Id: "30024",
		Type: "series",
		Name: "User Error",
		Description: "Life is a series of mistakes, but that's what makes it interesting. A show " +
			"about life, Linux, the universe, and everything in between.",
		Logo: "https://static.feedpress.it/logo/uevideo-57c6160ae2eab.png",
		Poster: "https://static.feedpress.it/logo/uevideo-57c6160ae2eab.png",
		Genres: DEFAULT_GENRES,
		Feed: ALT_FEED_BASE + "uevideo" + FEED_FORMAT,
	})

	for _, show := range showList {
		show.EpisodeList = UpdateEpisodes(*show)
	}

	return showList
}

func UpdateEpisodes(show JupiterShow) ([]Episode) {
	resp, _ := http.Get(show.Feed)
	bytes, _ := ioutil.ReadAll(resp.Body)

	var jsonResponse JsonResponse
	err := json.Unmarshal(bytes, &jsonResponse)
	episodes := []Episode{}

	if err != nil {
		log.Println("Error fetching jupiter feed: %q", err)
		log.Println("Error fetching jupiter feed: %v", err.Error())
		return episodes
	}

	for j, feedItem := range jsonResponse.Episodes {
		episode := Episode{
			Id:		fmt.Sprintf("%s:%d", show.Id, len(jsonResponse.Episodes) - j),
			Season:		2019,
			Episode:	len(jsonResponse.Episodes) - j,
			Title:		feedItem.Title,
			Thumbnail:	show.Logo, //XXX: fetch better?
			Released:	feedItem.Published,
			Overview:	feedItem.Html,
			Streams:	[]StreamItem{},
		}

		for _, attachment := range feedItem.Attachments {
			stream := StreamItem{
				/// XXX? Id:	"
				Title:	"Jupiter Zone",
				Url:	attachment.Url,
			}
			episode.Streams = append(episode.Streams, stream)
		}
		
		episodes = append(episodes, episode)
	}
	
	return episodes
}
