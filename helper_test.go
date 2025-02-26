package uaparser

import "testing"

func TestIsAIBot(t *testing.T) {
	tests := []struct {
		ua       string
		expected bool
	}{
		{"Mozilla/5.0 (compatible; ai2bot/1.0; +http://www.ai2.com/bot.html)", true},
		{"Mozilla/5.0 (compatible; googlebot/2.1; +http://www.google.com/bot.html)", false},
		{"Mozilla/5.0 (compatible; gptbot/1.0; +http://www.openai.com/bot.html)", true},
		{"Mozilla/5.0 (compatible; someotherbot/1.0; +http://www.someotherbot.com/bot.html)", false},
	}

	for _, test := range tests {
		result := isAIBot(test.ua)
		if result != test.expected {
			t.Errorf("isAIBot(%q) = %v; want %v", test.ua, result, test.expected)
		}
	}
}

func TestIsBot(t *testing.T) {
	tests := []struct {
		ua       string
		expected bool
	}{
		{"Mozilla/5.0 (compatible; googlebot/2.1; +http://www.google.com/bot.html)", true},
		{"Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)", true},
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3", false},
		{"Mozilla/5.0 (compatible; YandexBot/3.0; +http://yandex.com/bots)", true},
	}

	for _, test := range tests {
		result := isBot(test.ua)
		if result != test.expected {
			t.Errorf("isBot(%q) = %v; want %v", test.ua, result, test.expected)
		}
	}
}
