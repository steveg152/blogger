package feedFetch

import (
	"encoding/xml"
	"io"
	"net/http"
)

type RSSFeed struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	XMLName xml.Name  `xml:"channel"`
	Items   []RSSItem `xml:"item"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
}

// FetchXMLFeed fetches an XML feed from the given URL and returns it as a string.
func FetchXMLFeed(url string) (RSSFeed, error) {
	res, err := http.Get(url)

	if err != nil {
		return RSSFeed{}, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return RSSFeed{}, err
	}

	var feed RSSFeed

	err = xml.Unmarshal(body, &feed)

	if err != nil {
		return RSSFeed{}, err
	}

	// for _, item := range feed.Channel.Items {
	// 	println(item.Title)
	// }

	return feed, nil

}
