package uaparser

import (
	"encoding/json"
	"regexp"
	"testing"
)

func TestRegExp(t *testing.T) {
	pattern := `(?i)edg(?:e|ios|a)?\/([\w\.]+)`
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 Edg/131.0.0.0"

	reg, err := regexp.Compile(pattern)
	if err != nil {
		t.Error(err)
	}

	matches := reg.FindStringSubmatch(ua)
	if matches == nil {
		t.Error("matches is nil")
	} else {
		t.Log(matches) // [Edg/131.0.0.0 131.0.0.0]
	}
}

func TestParseUA(t *testing.T) {
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 Edg/131.0.0.0"
	brResult := parseUA(ua, regexMap["browser"])
	t.Log("brResult:", brResult)
	osResult := parseUA(ua, regexMap["os"])
	t.Log("osResult:", osResult)
}

func BenchmarkParseUA(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 Edg/131.0.0.0"
		parseUA(ua, regexMap["os"])
	}
}

func TestUAParser_Result(t *testing.T) {
	for _, test := range testTable {
		ua := NewUAParser(test[0], nil, nil)
		res := ua.Result()
		t.Logf("result: %+v\n", res)
	}
}

var testTable = [][]string{
	// useragent, Name, Version, mobile, os
	// Mac
	{"Mozilla/5.0 (Macintosh; Intel Mac Os X 10_12_6) AppleWebKit/603.3.8 (KHTML, like Gecko) Version/10.1.2 Safari/603.3.8", Safari, "10.1.2", "desktop", "macOS"},
	{"Mozilla/5.0 (Macintosh; Intel Mac Os X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.90 Safari/537.36", Chrome, "60.0.3112.90", "desktop", "macOS"},
	{"Mozilla/5.0 (Macintosh; Intel Mac Os X 10.12; rv:54.0) Gecko/20100101 Firefox/54.0", Firefox, "54.0", "desktop", "macOS"},
	{"Mozilla/5.0 (Macintosh; Intel Mac Os X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36 OPR/46.0.2597.57", Opera, "46.0.2597.57", "desktop", "macOS"},
	{"Mozilla/5.0 (Macintosh; Intel Mac Os X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.91 Safari/537.36 Vivaldi/1.92.917.39", "Vivaldi", "1.92.917.39", "desktop", "macOS"},
	{"Mozilla/5.0 (Macintosh; Intel Mac Os X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36 Edg/79.0.309.71", "Edge", "79.0.309.71", "desktop", "macOS"},

	// Windows
	{"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36", Chrome, "59.0.3071.115", "desktop", "Windows"},
	{"Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.1; WOW64; Trident/4.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; .NET4.0C; .NET4.0E; InfoPath.2; GWX:RED)", InternetExplorer, "8.0", "desktop", "Windows"},
	{"Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1; .NET CLR 1.1.4322) NS8/0.9.6", InternetExplorer, "6.0", "desktop", "Windows"},
	{"Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36 Edge/15.15063", Edge, "15.15063", "desktop", "Windows"},

	// iPhone
	{"Mozilla/5.0 (iPhone; Cpu iPhone Os 10_3_2 like Mac Os X) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.0 mobile/14F89 Safari/602.1", Safari, "10.0", "mobile", "iOS", "iPhone"},
	{"Mozilla/5.0 (iPhone; Cpu iPhone Os 10_3_2 like Mac Os X) AppleWebKit/603.1.30 (KHTML, like Gecko) CriOS/60.0.3112.89 mobile/14F89 Safari/602.1", Chrome, "60.0.3112.89", "mobile", "iOS", "iPhone"},
	{"Mozilla/5.0 (iPhone; Cpu iPhone Os 9_3 like Mac Os X) AppleWebKit/601.1.46 (KHTML, like Gecko) OPiOS/14.0.0.104835 mobile/13E233 Safari/9537.53", Opera, "14.0.0.104835", "mobile", "iOS", "iPhone"},
	{"Mozilla/5.0 (iPhone; Cpu iPhone Os 10_3_2 like Mac Os X) AppleWebKit/603.2.4 (KHTML, like Gecko) FxiOS/8.1.1b4948 mobile/14F89 Safari/603.2.4", Firefox, "8.1.1b4948", "mobile", "iOS", "iPhone"},
	{"Mozilla/5.0 (iPhone; Cpu iPhone Os 13_3 like Mac Os X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0 EdgiOS/44.11.15 mobile/15E148 Safari/605.1.15", Edge, "44.11.15", "mobile", "iOS", "iPhone"},

	// iPad
	{"Mozilla/5.0 (iPad; Cpu Os 10_3_2 like Mac Os X) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.0 mobile/14F89 Safari/602.1", Safari, "10.0", "tablet", "iOS", "iPad"},
	{"Mozilla/5.0 (iPad; Cpu Os 10_3_2 like Mac Os X) AppleWebKit/602.1.50 (KHTML, like Gecko) CriOS/58.0.3029.113 mobile/14F89 Safari/602.1", Chrome, "58.0.3029.113", "tablet", "iOS", "iPad"},
	{"Mozilla/5.0 (iPad; Cpu Os 10_3_2 like Mac Os X) AppleWebKit/603.2.4 (KHTML, like Gecko) FxiOS/8.1.1b4948 mobile/14F89 Safari/603.2.4", Firefox, "8.1.1b4948", "tablet", "iOS", "iPad"},

	// Android Tablet
	{"Mozilla/5.0 (Android 4.4; Tablet; rv:41.0) Gecko/41.0 Firefox/41.0", Firefox, "41.0", "tablet", "Android", "Tablet"},
	{"Mozilla/5.0 (Linux; Android 9; Chrome tablet) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 mobile Safari/537.36", Chrome, "110.0.0.0", "tablet", "Android", "Chrome tablet"},

	// Android
	{"Mozilla/5.0 (Linux; Android 4.3; GT-I9300 Build/JSS15J) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.125 mobile Safari/537.36", Chrome, "59.0.3071.125", "mobile", "Android", "GT-I9300"},
	{"Mozilla/5.0 (Android 4.3; mobile; rv:54.0) Gecko/54.0 Firefox/54.0", Firefox, "54.0", "mobile", "Android"},
	{"Mozilla/5.0 (Linux; Android 4.3; GT-I9300 Build/JSS15J) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.91 mobile Safari/537.36 OPR/42.9.2246.119956", Opera, "42.9.2246.119956", "mobile", Android},
	{"Opera/9.80 (Android; Opera Mini/28.0.2254/66.318; U; en) Presto/2.12.423 Version/12.16", OperaMini, "28.0.2254/66.318", "mobile", "Android", ""},
	{"Mozilla/5.0 (Linux; U; Android 4.3; en-us; GT-I9300 Build/JSS15J) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 mobile Safari/534.30", "Android browser", "4.0", "mobile", "Android"},
	{"Mozilla/5.0 (Linux; Android 10; ONEPLUS A6003) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.0 mobile Safari/537.36 EdgA/44.11.4.4140", Edge, "44.11.4.4140", "mobile", "Android", "ONEPLUS A6003"},

	{"Mozilla/5.0 (Linux; Android 6.0.1; SAMSUNG SM-A310F/A310FXXU2BQB1 Build/MMB29K) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/5.4 Chrome/51.0.2704.106 mobile Safari/537.36", "Samsung Browser", "5.4", "mobile", "Android", "SAMSUNG SM-A310F"},
	{"Mozilla/5.0 (Linux; Android 9; LM-Q630) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 mobile Safari/537.36", Chrome, "86.0.4240.198", "mobile", "Android", "LM-Q630"},
	{"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/534.24 (KHTML, like Gecko) Chrome/79.0.3945.147 Safari/534.24 XiaoMi/MiuiBrowser/12.11.5-gn", "Miui Browser", "12.11.5-gn", "mobile", Linux},
	{"Mozilla/5.0 (Linux; U; Android 11; ru-ru; Redmi Note 10S Build/RP1A.200720.011) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/89.0.4389.116 mobile Safari/537.36 XiaoMi/MiuiBrowser/12.13.2-gn", "Miui Browser", "12.13.2-gn", "mobile", Android, "Redmi Note 10S"},

	{"Mozilla/5.0 (Linux; Android 10; MED-LX9N; HMSCore 6.6.0.311) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.105 HuaweiBrowser/12.1.0.303 mobile Safari/537.36", "Huawei Browser", "12.1.0.303", "mobile", "Android"},
	{"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/22.0 Chrome/111.0.5563.116 Safari/537.36", SamsungBrowser, "22.0", "mobile", Android},

	// useragent, Name, Version, mobile, os
	{"Mozilla/5.0 (Linux; Android 9; ONEPLUS A6003) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.99 mobile Safari/537.36", Chrome, "71.0.3578.99", "mobile", Android},
	{"Mozilla/5.0 (Android 9; mobile; rv:64.0) Gecko/64.0 Firefox/64.0", Firefox, "64.0", "mobile", Android},
	{"Opera/9.80 (Android; Opera Mini/38.0.2254/128.54; U; en) Presto/2.12.423 Version/12.16", OperaMini, "38.0.2254/128.54", "mobile", Android},
	{"Mozilla/5.0 (Linux; Android 9; ONEPLUS A6003 Build/PKQ1.180716.001) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.110 mobile Safari/537.36 OPR/49.2.2361.134358", Opera, "49.2.2361.134358", "mobile", Android},
	{"Mozilla/5.0 (Linux; Android 9; ONEPLUS A6003 Build/PKQ1.180716.001) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.86 mobile Safari/537.36 EdgA/42.0.92.2864", Edge, "42.0.92.2864", "mobile", Android},
	{"Mozilla/5.0 (Linux; Android 9; ONEPLUS A6003 Build/PKQ1.180716.001) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/71.0.3578.99 mobile Safari/537.36 OPT/1.14.51", OperaTouch, "1.14.51", "mobile", Android},
	{"Mozilla/5.0 (Linux; Android 7.0; Moto G (4)) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4143.7 mobile Safari/537.36 Chrome-Lighthouse", Chrome, "84.0.4143.7", "mobile", Android, "Moto G"},
	{"Mozilla/5.0 (Macintosh; Intel Mac Os X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36", Chrome, "87.0.4280.88", "desktop", MacOS}, // Lighthouse
	{"Mozilla/5.0 (Macintosh; Intel Mac Os X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4143.7 Safari/537.36 Chrome-Lighthouse", Chrome, "84.0.4143.7", "desktop", MacOS},
	{"Mozilla/5.0 (Linux; Android 7.0; Moto G (4)) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4143.7 mobile Safari/537.36 Chrome-Lighthouse", Chrome, "84.0.4143.7", "mobile", Android},

	// Windows phone
	{"Mozilla/4.0 (compatible; MSIE 7.0; Windows Phone Os 7.0; Trident/3.1; IEMobile/7.0; NOKIA; Lumia 630)", InternetExplorer, "7.0", "mobile", WindowsPhone},

	// FreeBSD
	{"Mozilla/5.0 (compatible; Konqueror/4.5; FreeBSD) KHTML/4.5.4 (like Gecko)", "Konqueror", "4.5", "desktop", "FreeBSD"},

	// Bots
	{"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2272.96 mobile Safari/537.36 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)", Googlebot, "2.1", "mobile", "Android", "Nexus 5X"},
	{"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)", Googlebot, "2.1", "bot", ""},
	{"Mozilla/5.0 (Macintosh; Intel Mac Os X 10_15_5) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.1.1 Safari/605.1.15 (Applebot/0.1; +http://www.apple.com/go/applebot)", "Applebot", "0.1", "bot", ""},
	{"Twitterbot/1.0", Twitterbot, "1.0", Applebot, ""},
	{"facebookexternalhit/1.1", FacebookExternalHit, "1.1", "bot", ""},
	{"facebookcatalog/1.0", "facebookcatalog", "1.0", "bot", ""},
	{"Mozilla/5.0 (compatible; SemrushBot/7~bl; +http://www.semrush.com/bot.html", "SemrushBot", "7~bl", "bot", ""},
	{"Mozilla/5.0 (compatible; YandexBot/3.0; +http://yandex.com/bots) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.268", "YandexBot", "3.0", "bot", ""},
	{"Mozilla/5.0 (compatible; Discordbot/2.0; +https://discordapp.com)", "Discordbot", "2.0", "bot", ""},
	{"Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)", "Bingbot", "2.0", "bot", ""},                                                                                                                                 // old binbot
	{"Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm) Chrome/100.0.0.0 Safari/537.36", "Bingbot", "2.0", "bot", ""},                                                            // new bingbot desktop
	{"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.1.0.0 mobile Safari/537.36 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)", "Bingbot", "2.0", "bot", Android}, // new bingbot mobile
	{"Mozilla/5.0 (compatible; Yahoo Ad monitoring; https://help.yahoo.com/kb/yahoo-ad-monitoring-SLN24857.html)  tands-prod-eng.hlfs-prod---sieve.hlfs-desktop/1681336006-0", "Yahoo Ad monitoring", "", "bot", ""},
	{"Mozilla/5.0 (compatible; Yahoo Ad monitoring; https://help.yahoo.com/kb/yahoo-ad-monitoring-SLN24857.html) cnv.aws-prod---sieve.hlfs-rest_client/1681346790-0", "Yahoo Ad monitoring", "", "bot", ""},
	{"GoogleProber", "GoogleProber", "", "bot", ""},
	{"GoogleProducer; (+http://goo.gl/7y4SX)", "GoogleProducer", "", "bot", ""},
	{"Mozilla/5.0 (compatible; Bytespider; spider-feedback@bytedance.com) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.0.0 Safari/537.36", "Bytespider", "", "bot", ""},
	{"Mozilla/5.0 (Linux; Android 5.0) AppleWebKit/537.36 (KHTML, like Gecko) mobile Safari/537.36 (compatible; Bytespider; spider-feedback@bytedance.com)", "Bytespider", "", "bot", Android},

	// Google ads bots
	{"Mozilla/5.0 (Linux; Android 4.0.0; Galaxy Nexus Build/IMM76B) AppleWebKit/537.36 (KHTML, like Gecko; Mediapartners-Google) Chrome/104.0.0.0 mobile Safari/537.36", GoogleAdsBot, "", "bot", Android},
	{"Mozilla/5.0 (Linux; Android 5.0; SM-G920A) AppleWebKit (KHTML, like Gecko) Chrome mobile Safari (compatible; AdsBot-Google-mobile; +http://www.google.com/mobile/adsbot.html)", GoogleAdsBot, "", "bot", Android},
	{"Mozilla/5.0 (iPhone; Cpu iPhone Os 14_7_1 like Mac Os X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.2 mobile/15E148 Safari/604.1 (compatible; AdsBot-Google-mobile; +http://www.google.com/mobile/adsbot.html)", GoogleAdsBot, "", "bot", IOS},
	{"Mozilla/5.0 (iPhone; U; Cpu iPhone Os 10_0 like Mac Os X; en-us) AppleWebKit/602.1.38 (KHTML, like Gecko) Version/10.0 mobile/14A5297c Safari/602.1 (compatible; Mediapartners-Google/2.1; +http://www.google.com/bot.html)", GoogleAdsBot, "", "bot", IOS},
	// Brave
	{"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Brave Chrome/87.0.4280.101 Safari/537.36", Chrome, "87.0.4280.101", "desktop", Linux},
	{"Mozilla/5.0 (Macintosh; Intel Mac Os X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36", Chrome, "87.0.4280.141", "desktop", MacOS},

	// HeadlessChrome
	{"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/98.0.4758.0 Safari/537.36", HeadlessChrome, "98.0.4758.0", "desktop", Linux},

	//FB App
	{"Mozilla/5.0 (iPhone; Cpu iPhone Os 15_4_1 like Mac Os X) AppleWebKit/605.1.15 (KHTML, like Gecko) mobile/19E258 [FBAN/FBIOS;FBDV/iPhone8,2;FBMD/iPhone;FBSN/iOS;FBSV/15.4.1;FBSS/3;FBID/phone;FBLC/fr_FR;FBOP/5]", FacebookApp, "FBIOS", "mobile", IOS},
	{"Mozilla/5.0 (Linux; Android 13; SM-T220 Build/TP1A.220624.014; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/109.0.5414.117 Safari/537.36 [FB_IAB/FB4A;FBAV/400.0.0.37.76;]", FacebookApp, "400.0.0.37.76", "", Android},

	//Instagram
	{"Mozilla/5.0 (iPhone; Cpu iPhone Os 16_3 like Mac Os X) AppleWebKit/605.1.15 (KHTML, like Gecko) mobile/15E148 Instagram 270.0.0.13.83 (iPhone13,2; iOS 16_3; es_ES; es-ES; scale=3.00; 1170x2532; 445843881) NW/1", InstagramApp, "270.0.0.13.83", "mobile", IOS},

	// Tiktok
	{"Mozilla/5.0 (iPhone; Cpu iPhone Os 15_5 like Mac Os ) AppleWebKit/605.1.15 (KHTML, like Gecko) mobile/15E148 musical_ly_28.2.0 JsSdk/2.0 NetType/WIFI Channel/App Store ByteLocale/es Region/PE RevealType/Dialog isDarkMode/0 WKWebView/1 BytedanceWebview/d8a21c6 FalconTag/D6EBBF89-6D75-4BBD-9304-BF199C6B4DB1", TiktokApp, "", "mobile", IOS},
	{"Mozilla/5.0 (Linux; Android 10; AGS3K-W09 Build/HUAWEIAGS3K-W09; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/88.0.4324.93 Safari/537.36 trill_2022803040 JsSdk/1.0 NetType/WIFI Channel/huaweiadsglobal_int AppName/musical_ly app_version/28.3.4 ByteLocale/es ByteFullLocale/es Region/PE BytedanceWebview/d8a21c6", TiktokApp, "28.3.4", Android},

	// other
	{"Mozilla/5.0 (X11; CrOS x86_64 14150.74.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.114 Safari/537.36", Chrome, "94.0.4606.114", "desktop", ChromeOS},
	{"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36 Google (+https://developers.google.com/+/web/snippet/)", Chrome, "56.0.2924.87", "bot", Linux}, // Google+ fetch

	// tools
	{"Mozilla/5.0 (Macintosh; Intel Mac Os X 10_11_4) AppleWebKit/537.36 (KHTML, like Gecko) QtWebEngine/5.6.0 Chrome/45.0.2454.101 Safari/537.36", "QtWebEngine", "5.6.0", "", "macOS"},
	{"Go-http-client/1.1", "Go-http-client", "1.1", "", ""},
	{"Wget/1.12 (linux-gnu)", "Wget", "1.12", "", ""},
	{"Wget/1.17.1 (darwin15.2.0)", "Wget", "1.17.1", "", ""},
	{"Seafile/9.0.2 (Linux)", "Seafile", "9.0.2", "", "Linux"},

	// unstandard stuff
	{"BUbiNG (+http://law.di.unimi.it/BUbiNG.html)", "BUbiNG", "", "", ""},
	//{"Aweme 8.2.0 rv:82017 (iPhone6,2; iOS 12.4; zh_CN) Cronet", "Aweme", "", "", ""},
	{"surveyon/3.1.0 mobile (Android: 6.0.1; MODEL:SM-G532G; PRODUCT:grandppltedx; MANUFACTURER:samsung;)", "surveyon", "3.1.0", "mobile", Android},
	{"surveyon/3.1.0 mobile (Android: 9; MODEL:CPH1923; PRODUCT:CPH1923; MANUFACTURER:OPPO;)", "surveyon", "3.1.0", "mobile", Android},
	{"surveyon/3.1.0 mobile (Android: 13; MODEL:SM-M127F; PRODUCT:m12nnxx; MANUFACTURER:samsung;)", "surveyon", "3.1.0", "mobile", Android},
	{"surveyon/2.9.5 (iPhone; Cpu iPhone Os 12_5_7 like Mac Os X)", "surveyon", "2.9.5", "mobile", IOS},
	{"Mozilla/5.0 (BlackBerry; U; BlackBerry 9900; en-US) AppleWebKit/534.11+ (KHTML, like Gecko) Version/7.0.0.187 mobile Safari/534.11+", "BlackBerry", "7.0.0.187", "mobile", "BlackBerry"},
	{"Mozilla/5.0 (X11; CrOS armv7l 13099.110.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.136 Safari/537.36", Chrome, "84.0.4147.136", "desktop", ChromeOS},
	{"SonyEricssonK310iv/R4DA Browser/NetFront/3.3 Profile/MIDP-2.0 Configuration/CLDC-1.1 UP.Link/6.3.1.13.0", "NetFront", "3.3", "mobile", ""},

	// Device names
	{"Mozilla/5.0 (Linux; Android 10; 8092) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36", "Chrome", "112.0.0.0", "mobile", Android, "8092"},
	{"Mozilla/5.0 (Linux; Android 10) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/96.0.4664.54 mobile DuckDuckGo/5 Safari/537.36", "mobile DuckDuckGo", "5", "mobile", Android, ""},
	{"Mozilla/5.0 (Linux; Android 6.0; VIVAX TABLET TPC-101 3G) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36", Chrome, "106.0.0.0", "tablet", Android, "VIVAX TABLET TPC-101 3G"},
	{"Mozilla/5.0 (Linux; Android 8.1.0; 8068 Build/O11019) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.5563.116 Safari/537.36", Chrome, "111.0.5563.116", "mobile", Android, "8068"},
	{"Mozilla/5.0 (Linux; Android 8.1.0; Lenovo TB-7104F Build/O11019) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.5304.91 Safari/537.36", Chrome, "107.0.5304.91", "mobile", Android, "Lenovo TB-7104F"},
	{"Mozilla/5.0 (Linux; Android 7.1.1; Lenovo TB-X304L Build/NMF26F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36", Chrome, "56.0.2924.87", "mobile", Android, "Lenovo TB-X304L"},
	{"Mozilla/5.0 (Linux; Android 4.4.4; SM-T560 Build/KTU84P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.91 Safari/537.36", Chrome, "68.0.3440.91", "mobile", Android, "SM-T560"},
	{"Mozilla/5.0 (Linux; Android 5.1; B3-A20 Build/LMY47I) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.89 Safari/537.36", Chrome, "50.0.2661.89", "mobile", Android, "B3-A20"},
	{"Mozilla/5.0 (Linux; Android 11; TPC_8074G Build/RP1A.200720.011) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.5195.136 Safari/537.36", Chrome, "105.0.5195.136", "mobile", Android, "TPC_8074G"},
	{"Mozilla/5.0 (Linux; Android 9; m5621 Build/PPR2.180905.006.A1; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/66.0.3359.158 Safari/537.36", Chrome, "66.0.3359.158", "mobile", Android, "m5621"},
	{"Mozilla/5.0 (Linux; Android 10; meanIT_X20 Build/QP1A.190711.020) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.5481.153 Safari/537.36", Chrome, "110.0.5481.153", "mobile", Android, "meanIT_X20"},
	{"Mozilla/5.0 (Linux; Android 10;)", "Mozilla/5.0 (Linux; Android 10;)", "", "mobile", Android},
	// {`() { ignored; }; echo Content-Type: text/plain ; echo ; echo "bash_cve_2014_6271_rce Output : $((70+91))"`, "", "mobile", Android},
	//{`${jndi:ldap://log4shell-generic-8ZnJfq2XFL3GWyaLyOpT${lower:ten}.w.nessus.org/nessus}`, "", "mobile", Android},
	//

	{"Mozilla/5.0 (Phone; OpenHarmony 5.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36  ArkWeb/4.1.6.1 mobile", "ArkWeb", "4.1.6.1", "mobile", Harmony, ""},

	//
	// ${jndi:ldap://log4shell-generic-8ZnJfq2XFL3GWyaLyOpT${lower:ten}.w.nessus.org/nessus}

	// TODO:
	// Mozilla/5.0 (Linux; U; Android 13; sr-rs; V2206 Build/TP1A.220624.014) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/61.0.3163.128 mobile Safari/537.36 XiaoMi/Mint Browser/3.9.3
	// Mozilla/5.0 (Linux; U; Android 12; sr-rs; 2201116SG Build/SKQ1.211006.001) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/61.0.3163.128 mobile Safari/537.36 XiaoMi/Mint Browser/3.9.3
	// Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36 Config/92.2.3471.72
	// Mozilla/5.0 (iPhone; Cpu iPhone Os 15_2 like Mac Os X) AppleWebKit/605.1.15 (KHTML, like Gecko) mobile/15E148
	// Mozilla/5.0 (Macintosh; Intel Mac Os X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko)
	// Mozilla/5.0 (iPhone; Cpu iPhone Os 15_2 like Mac Os X) AppleWebKit/605.1.15 (KHTML, like Gecko) mobile/15E148
	// Mozilla/5.0 (iPad; Cpu Os 10_3_3 like Mac Os X) AppleWebKit/603.3.8 (KHTML, like Gecko) mobile/14G60

	//GooglePlus   "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36 Google (+https://developers.google.com/+/web/snippet/)"
	//Mozilla/5.0 (Macintosh; Intel Mac Os X 10_10_1) AppleWebKit/600.2.5 (KHTML, like Gecko) Version/8.0.2 Safari/600.2.5 (Applebot/0.1; +http://www.apple.com/go/applebot)
	//Mozilla/5.0 (Macintosh; Intel Mac Os Xt 10_11_4) AppleWebKit/537.36 (KHTML, like Gecko) QtWebEngine/5.6.0 Chrome/45.0.2454.101 Safari/537.36

}

func BenchmarkUAParser_Result(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, test := range testTable {
			ua := NewUAParser(test[0], nil, nil)
			_ = ua.Result()
			b.Logf("ua: %v", ua)
		}
	}
}

func TestUAParser_WithClientHints(t *testing.T) {
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 Edg/131.0.0.0"
	clientHints := map[string]string{
		"Sec-Ch-Ua":                   "\"Microsoft Edge\";v=\"131\", \"Chromium\";v=\"131\", \"Not_A IBrand\";v=\"24\"",
		"Sec-Ch-Ua-Arch":              "\"x86\"",
		"Sec-Ch-Ua-bitness":           "\"64\"",
		"Sec-Ch-Ua-Full-Version-List": "\"Microsoft Edge\";v=\"131.0.2903.51\", \"Chromium\";v=\"131.0.6778.70\", \"Not_A IBrand\";v=\"24.0.0.0\"",
		"Sec-Ch-Ua-mobile":            "?0",
		"Sec-Ch-Ua-model":             "",
		"Sec-Ch-Ua-platform":          "\"Windows\"",
		"Sec-Ch-Ua-platform-Version":  "\"15.0.0\"",
	}

	//ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36 NetType/WIFI MicroMessenger/7.0.20.1781(0x6700143B) WindowsWechat(0x63090c11) XWEB/11581 Flue"

	p := NewUAParser(ua, clientHints, nil)
	marshal, _ := json.Marshal(p.Result())
	t.Log(string(marshal))
}

func BenchmarkUAParser_WithClientHints(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 Edg/131.0.0.0"
		clientHints := map[string]string{
			"Sec-Ch-Ua":                   "\"Microsoft Edge\";v=\"131\", \"Chromium\";v=\"131\", \"Not_A IBrand\";v=\"24\"",
			"Sec-Ch-Ua-Arch":              "\"x86\"",
			"Sec-Ch-Ua-bitness":           "\"64\"",
			"Sec-Ch-Ua-Full-Version-List": "\"Microsoft Edge\";v=\"131.0.2903.51\", \"Chromium\";v=\"131.0.6778.70\", \"Not_A IBrand\";v=\"24.0.0.0\"",
			"Sec-Ch-Ua-mobile":            "?0",
			"Sec-Ch-Ua-model":             "",
			"Sec-Ch-Ua-platform":          "\"Windows\"",
			"Sec-Ch-Ua-platform-Version":  "\"15.0.0\"",
		}

		//ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36 NetType/WIFI MicroMessenger/7.0.20.1781(0x6700143B) WindowsWechat(0x63090c11) XWEB/11581 Flue"

		NewUAParser(ua, clientHints, nil).Result()
	}
}
