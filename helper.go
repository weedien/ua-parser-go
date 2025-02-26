package uaparser

import (
	"slices"
	"strings"
)

func isAIBot(ua string) bool {
	aibots := []string{

		// AI2
		"ai2bot",

		// Amazon
		"amazonbot",

		// Anthropic
		"anthropic-ai",
		"claude-web",
		"claudebot",

		// Apple
		"applebot",
		"applebot-extended",

		// ByteDance
		"bytespider",

		// Common Crawl
		"ccbot",

		// DataForSeo
		"dataforseobot",

		// Diffbot
		"diffbot",

		// Google
		"googleother",
		"googleother-image",
		"googleother-video",
		"google-extended",

		// Hive AI
		"imagesiftbot",

		// Huawei
		"petalbot",

		// Meta
		"facebookbot",
		"meta-externalagent",

		// OpenAI
		"gptbot",
		"oai-searchbot",

		// Perplexity
		"perplexitybot",

		// Semrush
		"semrushbot-ocob",

		// Timpi
		"timpibot",

		// Velen.io
		"velenpublicwebcrawler",

		// Webz.io
		"omgili",
		"omgilibot",
		"webzio-extended",

		// You.com
		"youbot",

		// Zyte
		"scrapy",
	}

	for _, botName := range aibots {
		if strings.Contains(ua, botName) {
			return true
		}
	}
	return false
}

func isBot(ua string) bool {
	t := NewUAParser(ua).WithExtensions(Bots).Browser().Type
	return slices.Contains([]string{CLI, Crawler, Fetcher, Library}, t)
}
