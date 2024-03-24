package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

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
