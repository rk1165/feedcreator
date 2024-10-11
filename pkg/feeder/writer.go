package scraper

import (
	"encoding/xml"
	"github.com/rk1165/feedcreator/internal/models"
	"log"
	"os"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
	IsPermalink bool   `xml:"isPermalink,attr"`
}

type Channel struct {
	Title       string  `xml:"title"`
	Link        string  `xml:"link"`
	Description string  `xml:"description"`
	Items       *[]Item `xml:"item"`
}

func CreateFeedFile(feed *models.Feed) {
	channel := CreateChannel(feed)
	rssXml, err := createRSSXML(channel)
	if err != nil {
		log.Printf("Error creating RSS XML: %v", err)
		return
	}
	xmlFile, err := os.Create("./ui/static/rss/" + feed.Name)
	if err != nil {
		log.Println("Error creating RSS file", err)
		return
	}
	defer xmlFile.Close()
	xmlFile.WriteString(xml.Header)
	_, err = xmlFile.Write(rssXml)
	if err != nil {
		log.Println("Error creating RSS file", err)
		return
	}
	log.Println("RSS file created successfully")
}

func createRSSXML(channel *Channel) ([]byte, error) {
	rss := RSS{
		Version: "2.0",
		Channel: *channel,
	}

	output, err := xml.MarshalIndent(rss, "", "  ")
	if err != nil {
		log.Println("Error: ", err)
		return nil, err
	}
	//os.Stdout.Write(output)
	return output, nil
}
