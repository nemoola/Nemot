package commands

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type Wiki struct{}

type WikiClient struct {
	urlBuilder url.URL
}

func NewWikiClient() WikiClient {
	return WikiClient{
		urlBuilder: url.URL{
			Scheme: "https",
			Host:   "ja.wikipedia.org",
			Path:   "/w/api.php",
		},
	}
}

type WikiSearch struct {
	Query struct {
		Search []struct {
			Title  string `json:"title"`
			PageID int    `json:"pageid"`
		} `json:"search"`
		Pages map[string]struct {
			PageID  int    `json:"pageid"`
			Title   string `json:"title"`
			Extract string `json:"extract"`
		} `json:"pages"`
	} `json:"query"`
}

type WikiResult struct {
	Title   string
	Content string
}

func (wc *WikiClient) GetPageContent(pageID int) WikiResult {
	query := wc.urlBuilder.Query()
	query.Add("action", "query")
	query.Add("prop", "extracts")
	query.Add("exintro", "1")
	query.Add("explaintext", "1")
	query.Add("exsectionformat", "plain")
	query.Add("pageids", strconv.Itoa(pageID))
	query.Add("format", "json")
	wc.urlBuilder.RawQuery = query.Encode()

	res, _ := http.Get(wc.urlBuilder.String())
	body, _ := io.ReadAll(res.Body)
	var search WikiSearch
	_ = json.Unmarshal(body, &search)
	result := WikiResult{
		Title:   search.Query.Pages[strconv.Itoa(pageID)].Title,
		Content: search.Query.Pages[strconv.Itoa(pageID)].Extract,
	}
	return result
}

func (wc *WikiClient) Search(title string) WikiResult {
	query := wc.urlBuilder.Query()
	query.Add("action", "query")
	query.Add("list", "search")
	query.Add("srsearch", title)
	query.Add("srlimit", "1")
	query.Add("srprop", "size")
	query.Add("srenablerewrites", "1")
	query.Add("format", "json")
	wc.urlBuilder.RawQuery = query.Encode()

	res, _ := http.Get(wc.urlBuilder.String())
	body, _ := io.ReadAll(res.Body)
	var search WikiSearch
	_ = json.Unmarshal(body, &search)
	wc.urlBuilder.RawQuery = url.Values{}.Encode()
	return wc.GetPageContent(search.Query.Search[0].PageID)
}

func (Wiki) Help() string {
	return "Wikipediaの検索結果を返します"
}

func (Wiki) Execute(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	wc := NewWikiClient()
	result := wc.Search(args[0])
	_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		URL:         fmt.Sprintf("https://ja.wikipedia.org/wiki/%s", result.Title),
		Title:       result.Title,
		Description: result.Content,
		Color:       0x00ff00,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "Wikipedia",
			IconURL: "https://upload.wikimedia.org/wikipedia/commons/thumb/8/80/Wikipedia-logo-v2.svg/103px-Wikipedia-logo-v2.svg.png",
		},
	})
}
