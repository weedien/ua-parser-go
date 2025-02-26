package uaparser

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestMergeExtension(t *testing.T) {
	wget := "Wget/1.21.1"
	facebookBot := "Mozilla/5.0 (compatible; FacebookBot/1.0; +https://developers.facebook.com/docs/sharing/webmasters/facebookbot/)"

	crawlersAndCLIs := map[string][]regexItem{
		"browser": append(Crawlers["browser"], CLIs["browser"]...),
	}
	crawlersAndCLIsParser := NewUAParser("", nil, crawlersAndCLIs)
	assert.Equal(t, IBrowser{Name: "Wget", Version: "1.21.1", Major: "1", Type: "cli"}, crawlersAndCLIsParser.SetUA(wget).Browser())
	assert.Equal(t, IBrowser{Name: "FacebookBot", Version: "1.0", Major: "1", Type: "crawler"}, crawlersAndCLIsParser.SetUA(facebookBot).Browser())
}

func TestExtension(t *testing.T) {
	alltestcases := []struct {
		label     string
		extension map[string][]regexItem
		list      []TestCase
	}{
		{
			label:     "CLIs",
			extension: CLIs,
			list:      loadJson("./data/ua/extension/cli.json"),
		},
		{
			label:     "Crawlers",
			extension: Crawlers,
			list:      loadJson("./data/ua/extension/crawler.json"),
		},
		{
			label:     "Emails",
			extension: Emails,
			list:      loadJson("./data/ua/extension/email.json"),
		},
		{
			label:     "Fetchers",
			extension: Fetchers,
			list:      loadJson("./data/ua/extension/fetcher.json"),
		},
		{
			label:     "Libraries",
			extension: Libraries,
			list:      loadJson("./data/ua/extension/library.json"),
		},
	}

	for _, singleCase := range alltestcases {
		for _, tc := range singleCase.list {
			browser := NewUAParser(tc.Ua, nil, singleCase.extension).Browser()
			for key, val := range tc.Expect {
				if isUndefined(val) {
					continue
				}

				if err := compareField(reflect.ValueOf(browser), key, val, singleCase.label, tc.Ua); err != nil {
					t.Error(err.Error())
				}
			}
		}
	}

	outlook := "Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 10.0; WOW64; Trident/7.0; .NET4.0C; .NET4.0E; .NET CLR 2.0.50727; .NET CLR 3.0.30729; .NET CLR 3.5.30729; Microsoft Outlook 16.0.9126; Microsoft Outlook 16.0.9126; ms-office; MSOffice 16)"
	thunderbird := "Mozilla/5.0 (X11; Linux x86_64; rv:78.0) Gecko/20100101 Thunderbird/78.13.0"
	axios := "axios/1.3.5"
	jsdom := "Mozilla/5.0 (darwin) AppleWebKit/537.36 (KHTML, like Gecko) jsdom/20.0.3"
	scrapy := "Scrapy/1.5.0 (+https://scrapy.org)"
	bluesky := "Mozilla/5.0 (compatible; Bluesky Cardyb/1.1; +mailto:support@bsky.app)"

	// Test for Scrapy
	parser := NewUAParser(scrapy, nil, Bots)
	assert.Equal(t, "Scrapy", parser.Browser().Name)

	// Test for Emails
	emailParser := NewUAParser("", nil, Emails)
	assert.Equal(t, IBrowser{Name: "Microsoft Outlook", Version: "16.0.9126", Major: "16", Type: "email"}, emailParser.SetUA(outlook).Browser())
	assert.Equal(t, IBrowser{Name: "Thunderbird", Version: "78.13.0", Major: "78", Type: "email"}, emailParser.SetUA(thunderbird).Browser())

	// Test for Libraries
	libraryParser := NewUAParser("", nil, Libraries)
	assert.Equal(t, IBrowser{Name: "axios", Version: "1.3.5", Major: "1", Type: "library"}, libraryParser.SetUA(axios).Browser())
	assert.Equal(t, IBrowser{Name: "jsdom", Version: "20.0.3", Major: "20", Type: "library"}, libraryParser.SetUA(jsdom).Browser())
	assert.Equal(t, IBrowser{Name: "Scrapy", Version: "1.5.0", Major: "1", Type: "library"}, libraryParser.SetUA(scrapy).Browser())

	// Test for Bluesky
	assert.Equal(t, IBrowser{Name: "Bluesky", Version: "1.1", Major: "1", Type: "fetcher"}, NewUAParser(bluesky, nil, Bots).Browser())
}
