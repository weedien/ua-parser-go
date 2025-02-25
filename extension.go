package uaparser

import "slices"

var (
	CLIs = map[string][]regexItem{
		"browser": {
			{
				// wget / curl / Lynx / ELinks / HTTPie
				patterns: []string{`(?i)(wget|curl|lynx|elinks|httpie)[\\/ ]\\(?([\\w\\.-]+)`},
				output: map[string]string{
					Name:    "$1",
					Version: "$2",
					Type:    CLI,
				},
			},
		},
	}

	Crawlers = map[string][]regexItem{
		"browser": {
			{
				// AhrefsBot, Amazonbot, Bingbot, CCBot, Dotbot, DuckDuckBot, FacebookBot, GPTBot, MJ12bot, MojeekBot, OpenAI's SearchGPT, PerplexityBot, SeznamBot
				patterns: []string{`(?i)((?:ahrefs|amazon|bing|cc|dot|duckduck|exa|facebook|gpt|mj12|mojeek|oai-search|perplexity|semrush|seznam)bot)\/([\w\.-]+)`},
				output: map[string]string{
					Name:    "$1",
					Version: "$2",
					Type:    Crawler,
				},
			},
			{
				// Applebot
				patterns: []string{`(?i)(applebot(?:-extended)?)\/([\w\.]+)`},
				output: map[string]string{
					Name:    "$1",
					Version: "$2",
					Type:    Crawler,
				},
			},
			{
				// Baiduspider
				patterns: []string{`(?i)(baiduspider)[-imagevdonsfcpr]{0,6}\/([\w\.]+)`},
				output: map[string]string{
					Name:    "$1",
					Version: "$2",
					Type:    Crawler,
				},
			},
			{
				// ClaudeBot (Anthropic)
				patterns: []string{`(?i)(claude(?:bot|-web)|anthropic-ai)\/?([\w\.]*)`},
				output: map[string]string{
					Name:    "$1",
					Version: "$2",
					Type:    Crawler,
				},
			},
			{
				// Coc Coc Bot
				patterns: []string{`(?i)(coccocbot-(?:image|web))\/([\w\.]+)`},
				output: map[string]string{
					Name:    "$1",
					Version: "$2",
					Type:    Crawler,
				},
			},
			{
				// Facebook / Meta
				patterns: []string{`(?i)(facebook(?:externalhit|catalog)|meta-externalagent)\/([\w\.]+)`},
				output: map[string]string{
					Name:    "$1",
					Version: "$2",
					Type:    Crawler,
				},
			},
			{
				// Googlebot
				patterns: []string{`(?i)(google(?:bot|other|-inspectiontool)(?:-image|-video|-news)?|storebot-google)\/?([\w\.]*)`},
				output: map[string]string{
					Name:    "$1",
					Version: "$2",
					Type:    Crawler,
				},
			},
			{
				// Internet Archive
				patterns: []string{`(?i)(ia_archiver|archive\.org_bot)\/?([\w\.]*)`},
				output: map[string]string{
					Name:    "$1",
					Version: "$2",
					Type:    Crawler,
				},
			},
			{
				// SemrushBot
				patterns: []string{`(?i)((?:semrush|splitsignal)bot[-abcfimostw]*)\/([\w\.-]+)`},
				output: map[string]string{
					Name:    "$1",
					Version: "$2",
					Type:    Crawler,
				},
			},
			{
				// Sogou Spider
				patterns: []string{`(?i)(sogou (?:pic|head|web|orion|news) spider)\/([\w\.]+)`},
				output: map[string]string{
					Name:    "$1",
					Version: "$2",
					Type:    Crawler,
				},
			},
			{
				// Yahoo! Japan
				patterns: []string{`(?i)(y!?j-(?:asr|br[uw]|dscv|mmp|vsidx|wsc))\/([\w\.]+)`},
				output: map[string]string{
					Name:    "$1",
					Version: "$2",
					Type:    Crawler,
				},
			},
			{
				// Yandex Bots
				patterns: []string{`(?i)(yandex(?:(?:mobile)?(?:accessibility|additional|renderresources|screenshot|sprav)?bot|image(?:s|resizer)|video(?:parser)?|blogs|adnet|favicons|fordomain|market|media|metrika|news|ontodb(?:api)?|pagechecker|partner|rca|tracker|turbo|vertis|webmaster|antivirus))\/([\w\.]+)`},
				output: map[string]string{
					Name:    "$1",
					Version: "$2",
					Type:    Crawler,
				},
			},
			{
				// Yeti (Naver)
				patterns: []string{`(?i)(yeti)\/([\w\.]+)`},
				output: map[string]string{
					Name:    "$1",
					Version: "$2",
					Type:    Crawler,
				},
			},
			{
				// aiHitBot, Diffbot, Magpie-Crawler, Omgilibot, Webzio-Extended, Screaming Frog SEO Spider, Timpibot, VelenPublicWebCrawler, YisouSpider, YouBot
				patterns: []string{`(?i)((?:aihit|diff|timpi|you)bot|omgili(?:bot)?|(?:magpie-|velenpublicweb)crawler|webzio-extended|(?:screaming frog seo |yisou)spider)\/?([\w\.]*)`},
				output: map[string]string{
					Name:    "$1",
					Version: "$2",
					Type:    Crawler,
				},
			},
			{
				// Google Bots
				patterns: []string{`(?i)((?:adsbot|apis|mediapartners)-google(?:-mobile)?|google-?(?:other|cloudvertexbot|extended|safety))`},
				output: map[string]string{
					Name: "$1",
					Type: Crawler,
				},
			},
			{
				// AI2Bot, Bytespider, DataForSeoBot, Huawei AspiegelBot / PetalBot, ImagesiftBot, Qihoo 360Spider, TurnitinBot, Yahoo! Slurp
				patterns: []string{`(?i)\b(360spider-?(?:image|video)?|bytespider|(?:ai2|aspiegel|dataforseo|imagesift|petal|turnitin)bot|teoma|(?=yahoo! )slurp)`},
				output: map[string]string{
					Name: "$1",
					Type: Crawler,
				},
			},
		},
	}

	ExtraDevices = map[string][]regexItem{
		"device": {
			{
				// Nook
				patterns: []string{`(?i)(nook)[\w ]+build\/(\w+)`},
				output: map[string]string{
					Vendor: "Nook",
					Model:  "$1",
					Type:   Tablet,
				},
			},
			{
				// Dell Streak
				patterns: []string{`(?i)(dell) (strea[kpr\d ]*[\dko])`},
				output: map[string]string{
					Vendor: "Dell",
					Model:  "$1",
					Type:   Tablet,
				},
			},
			{
				// Le Pan Tablets
				patterns: []string{`(?i)(le[- ]+pan)[- ]+(\w{1,9}) bui`},
				output: map[string]string{
					Vendor: "Le Pan",
					Model:  "$1",
					Type:   Tablet,
				},
			},
			{
				// Trinity Tablets
				patterns: []string{`(?i)(trinity)[- ]*(t\d{3}) bui`},
				output: map[string]string{
					Vendor: "Trinity",
					Model:  "$1",
					Type:   Tablet,
				},
			},
			{
				// Gigaset Tablets
				patterns: []string{`(?i)(gigaset)[- ]+(q\w{1,9}) bui`},
				output: map[string]string{
					Vendor: "Gigaset",
					Model:  "$1",
					Type:   Tablet,
				},
			},
			{
				// Vodafone
				patterns: []string{`(?i)(vodafone) ([\w ]+)(?:\)| bui)`},
				output: map[string]string{
					Vendor: "Vodafone",
					Model:  "$1",
					Type:   Tablet,
				},
			},
			{
				// AT&T
				patterns: []string{`(?i)(u304aa)`},
				output: map[string]string{
					Model:  "$1",
					Vendor: "AT&T",
					Type:   Mobile,
				},
			},
			{
				// Siemens
				patterns: []string{`(?i)\bsie-(\w*)`},
				output: map[string]string{
					Model:  "$1",
					Vendor: "Siemens",
					Type:   Mobile,
				},
			},
			{
				// RCA Tablets
				patterns: []string{`(?i)\b(rct\w+) b`},
				output: map[string]string{
					Model:  "$1",
					Vendor: "RCA",
					Type:   Tablet,
				},
			},
			{
				// Dell Venue Tablets
				patterns: []string{`(?i)\b(venue[\d ]{2,7}) b`},
				output: map[string]string{
					Model:  "$1",
					Vendor: "Dell",
					Type:   Tablet,
				},
			},
			{
				// Verizon Tablet
				patterns: []string{`(?i)\b(q(?:mv|ta)\w+) b`},
				output: map[string]string{
					Model:  "$1",
					Vendor: "Verizon",
					Type:   Tablet,
				},
			},
			{
				// Barnes & Noble Tablet
				patterns: []string{`(?i)\b(?:barnes[& ]+noble |bn[rt])([\w\+ ]*) b`},
				output: map[string]string{
					Model:  "$1",
					Vendor: "Barnes & Noble",
					Type:   Tablet,
				},
			},
			{
				// NuVision Tablets
				patterns: []string{`(?i)\b(tm\d{3}\w+) b`},
				output: map[string]string{
					Model:  "$1",
					Vendor: "NuVision",
					Type:   Tablet,
				},
			},
			{
				// ZTE K Series Tablet
				patterns: []string{`(?i)\b(k88) b`},
				output: map[string]string{
					Model:  "$1",
					Vendor: "ZTE",
					Type:   Tablet,
				},
			},
			{
				// ZTE Nubia
				patterns: []string{`(?i)\b(nx\d{3}j) b`},
				output: map[string]string{
					Model:  "$1",
					Vendor: "ZTE",
					Type:   Mobile,
				},
			},
			{
				// Swiss GEN Mobile
				patterns: []string{`(?i)\b(gen\d{3}) b.+49h`},
				output: map[string]string{
					Model:  "$1",
					Vendor: "Swiss",
					Type:   Mobile,
				},
			},
			{
				// Swiss ZUR Tablet
				patterns: []string{`(?i)\b(zur\d{3}) b`},
				output: map[string]string{
					Model:  "$1",
					Vendor: "Swiss",
					Type:   Tablet,
				},
			},
			{
				// Zeki Tablets
				patterns: []string{`(?i)^((zeki)?tb.*\b) b`},
				output: map[string]string{
					Model:  "$1",
					Vendor: "Zeki",
					Type:   Tablet,
				},
			},
			{
				// Dragon Touch Tablet
				patterns: []string{`(?i)\b([yr]\d{2}) b`, `(?i)\b(?:dragon[- ]+touch |dt)(\w{5}) b`},
				output: map[string]string{
					Model:  "$1",
					Vendor: "Dragon Touch",
					Type:   Tablet,
				},
			},
			{
				// Insignia Tablets
				patterns: []string{`(?i)\b(ns-?\w{0,9}) b`},
				output: map[string]string{
					Model:  "$1",
					Vendor: "Insignia",
					Type:   Tablet,
				},
			},
			{
				// NextBook Tablets
				patterns: []string{`(?i)\b((nxa|next)-?\w{0,9}) b`},
				output: map[string]string{
					Model:  "$1",
					Vendor: "NextBook",
					Type:   Tablet,
				},
			},
			{
				// Voice Xtreme Phones
				patterns: []string{`(?i)\b(xtreme\_)?(v(1[045]|2[015]|[3469]0|7[05])) b`},
				output: map[string]string{
					Vendor: "Voice",
					Model:  "$1",
					Type:   Mobile,
				},
			},
			{
				// LvTel Phones
				patterns: []string{`(?i)\b(lvtel\-)?(v1[12]) b`},
				output: map[string]string{
					Vendor: "LvTel",
					Model:  "$1",
					Type:   Mobile,
				},
			},
			{
				// Essential PH-1
				patterns: []string{`(?i)\b(ph-1) `},
				output: map[string]string{
					Model:  "$1",
					Vendor: "Essential",
					Type:   Mobile,
				},
			},
			{
				// Envizen Tablets
				patterns: []string{`(?i)\b(v(100md|700na|7011|917g).*\b) b`},
				output: map[string]string{
					Model:  "$1",
					Vendor: "Envizen",
					Type:   Tablet,
				},
			},
			{
				// MachSpeed Tablets
				patterns: []string{`(?i)\b(trio[-\w\. ]+) b`},
				output: map[string]string{
					Model:  "$1",
					Vendor: "MachSpeed",
					Type:   Tablet,
				},
			},
			{
				// Rotor Tablets
				patterns: []string{`(?i)\btu_(1491) b`},
				output: map[string]string{
					Model:  "$1",
					Vendor: "Rotor",
					Type:   Tablet,
				},
			},
		},
	}

	Emails = map[string][]regexItem{
		"browser": {
			{
				// Evolution, Kontact/KMail, [Microsoft/Mac] Outlook, Thunderbird
				patterns: []string{`(?i)(airmail|bluemail|emclient|evolution|foxmail|kmail2?|kontact|(?:microsoft |mac)?outlook(?:-express)?|navermailapp|(?!chrom.+)sparrow|thunderbird|yahoo)(?:m.+ail; |[\\/ ])([\\w\\.]+)`},
				output: map[string]string{
					Name:    "$1",
					Version: "$2",
					Type:    Email,
				},
			},
		},
	}

	Fetchers = map[string][]regexItem{
		"browser": {
			{
				// AhrefsSiteAudit - https://ahrefs.com/robot/site-audit
				// ChatGPT-User - https://platform.openai.com/docs/plugins/bot
				// DuckAssistBot - https://duckduckgo.com/duckassistbot/
				// BingPreview / Mastodon / Pinterestbot / Redditbot / Rogerbot / SiteAuditBot / Telegrambot / Twitterbot / UptimeRobot
				// Google Site Verifier / Meta / Yahoo! Japan
				// Yandex Bots - https://yandex.com/bots
				patterns: []string{`(?i)(ahrefssiteaudit|bingpreview|chatgpt-user|mastodon|(?:discord|duckassist|linkedin|pinterest|reddit|roger|siteaudit|telegram|twitter|uptimero)bot|google-site-verification|meta-externalfetcher|y!?j-dlc|yandex(?:calendar|direct(?:dyn)?|searchshop)|yadirectfetcher)\/([\w\.]+)`},
				output: map[string]string{
					Name:    "$1",
					Version: "$2",
					Type:    Fetcher,
				},
			},
			{
				// Bluesky
				patterns: []string{`(?i)(bluesky) cardyb\/([\w\.]+)`},
				output: map[string]string{
					Name:    "$1",
					Version: "$2",
					Type:    Fetcher,
				},
			},
			{
				// Slackbot
				patterns: []string{`(?i)(slack(?:bot)?(?:-imgproxy|-linkexpanding)?) ([\w\.]+)`},
				output: map[string]string{
					Name:    "$1",
					Version: "$2",
					Type:    Fetcher,
				},
			},
			{
				// WhatsApp
				patterns: []string{`(?i)(whatsapp)\/([\w\.]+)[\/ ][ianw]`},
				output: map[string]string{
					Name:    "$1",
					Version: "$2",
					Type:    Fetcher,
				},
			},
			{
				// Google Bots, Cohere, Snapchat, Vercelbot, Yandex Bots
				patterns: []string{`(?i)(cohere-ai|vercelbot|feedfetcher-google|google(?:-read-aloud|producer)|(?=bot; )snapchat|yandex(?:sitelinks|userproxy))`},
				output: map[string]string{
					Name: "$1",
					Type: Fetcher,
				},
			},
		},
	}

	InApps = map[string][]regexItem{
		"browser": {
			{
				// Slack
				patterns: []string{`(?i)chatlyio\/([\d\.]+)`},
				output: map[string]string{
					Version: "$1",
					Name:    "Slack",
					Type:    InApp,
				},
			},
			{
				// Yahoo! Japan
				patterns: []string{`(?i)jp\.co\.yahoo\.android\.yjtop\/([\d\.]+)`},
				output: map[string]string{
					Version: "$1",
					Name:    "Yahoo! Japan",
					Type:    InApp,
				},
			},
		},
	}

	MediaPlayers = map[string][]regexItem{
		"browser": {
			{
				// Generic Apple CoreMedia
				patterns: []string{`(?i)(apple(?:coremedia|tv))\/([\w\._]+)`},
				output: map[string]string{
					Name:    "$1",
					Version: "$2",
					Type:    MediaPlayer,
				},
			},
			{
				// Ares/Nexplayer/OSSProxy
				patterns: []string{`(?i)(ares|clementine|music player daemon|nexplayer|ossproxy) ([\w\.-]+)`},
				output: map[string]string{
					Name:    "$1",
					Version: "$2",
					Type:    MediaPlayer,
				},
			},
			{
				// Aqualung/Lyssna/BSPlayer/Clementine/MPD, Audacious/AudiMusicStream/Amarok/BASS/OpenCORE/GnomeMplayer/MoC, NSPlayer/PSP-InternetRadioPlayer/Videos, Nero Home/Nero Scout/Nokia, QuickTime/RealMedia/RadioApp/RadioClientApplication, SoundTap/Totem/Stagefright/Streamium, XBMC/gvfs/Xine/XMMS/irapp
				patterns: []string{`(?i)^(aqualung|audacious|audimusicstream|amarok|bass|bsplayer|core|gnomemplayer|gvfs|irapp|lyssna|music on console|nero (?:home|scout)|nokia\d+|nsplayer|psp-internetradioplayer|quicktime|rma|radioapp|radioclientapplication|soundtap|stagefright|streamium|totem|videos|xbmc|xine|xmms)\/([\w\.-]+)`},
				output: map[string]string{
					Name:    "$1",
					Version: "$2",
					Type:    MediaPlayer,
				},
			},
			{
				// NexPlayer/LG Player
				patterns: []string{`(?i)(lg player|nexplayer) ([\d\.]+)`},
				output: map[string]string{
					Name:    "$1",
					Version: "$2",
					Type:    MediaPlayer,
				},
			},
			{
				// Gstreamer
				patterns: []string{`(?i)(gstreamer) souphttpsrc.+libsoup\/([\w\.-]+)`},
				output: map[string]string{
					Name:    "$1",
					Version: "$2",
					Type:    MediaPlayer,
				},
			},
			{
				// HTC Streaming Player
				patterns: []string{`(?i)(htc streaming player) [\w_]+ \/ ([\d\.]+)`},
				output: map[string]string{
					Name:    "$1",
					Version: "$2",
					Type:    MediaPlayer,
				},
			},
			{
				// Lavf (FFMPEG)
				patterns: []string{`(?i)(lavf)([\d\.]+)`},
				output: map[string]string{
					Name:    "$1",
					Version: "$2",
					Type:    MediaPlayer,
				},
			},
			{
				// MPlayer SVN
				patterns: []string{`(?i)(mplayer)(?: |\/)(?:(?:sherpya-){0,1}svn)(?:-| )(r\d+(?:-\d+[\w\.-]+))`},
				output: map[string]string{
					Name:    "$1",
					Version: "$2",
					Type:    MediaPlayer,
				},
			},
			{
				// Songbird/Philips-Songbird
				patterns: []string{`(?i) (songbird)\/([\w\.-]+)`},
				output: map[string]string{
					Name:    "$1",
					Version: "$2",
					Type:    MediaPlayer,
				},
			},
			{
				// Winamp
				patterns: []string{`(?i)(winamp)(?:3 version|mpeg| ) ([\w\.-]+)`},
				output: map[string]string{
					Name:    "$1",
					Version: "$2",
					Type:    MediaPlayer,
				},
			},
			{
				// VLC Videolan
				patterns: []string{`(?i)(vlc)(?:\/| media player - version )([\w\.-]+)`},
				output: map[string]string{
					Name:    "$1",
					Version: "$2",
					Type:    MediaPlayer,
				},
			},
			{
				// Foobar2000/iTunes/SMP
				patterns: []string{`(?i)^(foobar2000|itunes|smp)\/([\d\.]+)`},
				output: map[string]string{
					Name:    "$1",
					Version: "$2",
					Type:    MediaPlayer,
				},
			},
			{
				// RiseUP Radio Alarm
				patterns: []string{`(?i)com\.(riseupradioalarm)\/([\d\.]*)`},
				output: map[string]string{
					Name:    "$1",
					Version: "$2",
					Type:    MediaPlayer,
				},
			},
			{
				// MPlayer
				patterns: []string{`(?i)(mplayer)(?:\s|\/| unknown-)([\w\.\-]+)`},
				output: map[string]string{
					Name:    "$1",
					Version: "$2",
					Type:    MediaPlayer,
				},
			},
			{
				// Windows Media Server
				patterns: []string{`(?i)(windows)\/([\w\.-]+) upnp\/[\d\.]+ dlnadoc\/[\d\.]+ home media server`},
				output: map[string]string{
					Name:    "$1",
					Version: "$2",
					Type:    MediaPlayer,
				},
			},
			{
				// Flip Player
				patterns: []string{`(?i)(flrp)\/([\w\.-]+)`},
				output: map[string]string{
					Name:    "Flip Player",
					Version: "$2",
					Type:    MediaPlayer,
				},
			},
			{
				// FStream/NativeHost/QuerySeekSpider, MPlayer (no other info)/Media Player Classic/Nero ShowTime, OCMS-bot/tap in radio/tunein/unknown/winamp (no other info), inlight radio / YourMuze
				patterns: []string{`(?i)(fstream|media player classic|inlight radio|mplayer|nativehost|nero showtime|ocms-bot|queryseekspider|tapinradio|tunein radio|winamp|yourmuze)`},
				output: map[string]string{
					Name: "$1",
					Type: MediaPlayer,
				},
			},
			{
				// HTC One S / Windows Media Player
				patterns: []string{`(?i)(htc_one_s|windows-media-player|wmplayer)\/([\w\.-]+)`},
				output: map[string]string{
					Name:    "$1",
					Version: "$2",
					Type:    MediaPlayer,
				},
			},
			{
				// Rad.io
				patterns: []string{`(?i)(rad.io|radio.(?:de|at|fr)) ([\d\.]+)`},
				output: map[string]string{
					Name:    "rad.io",
					Version: "$2",
					Type:    MediaPlayer,
				},
			},
		},
	}

	Libraries = map[string][]regexItem{
		"browser": {
			{
				// Apache-HttpClient/Axios/go-http-client/got/GuzzleHttp/Java[-HttpClient]/jsdom/libwww-perl/lua-resty-http/Needle/node-fetch/OkHttp/PHP-SOAP/PostmanRuntime/python-urllib/python-requests/Scrapy/superagent
				patterns: []string{`(?i)^(apache-httpclient|axios|(?:go|java)-http-client|got|guzzlehttp|java|libwww-perl|lua-resty-http|needle|node-(?:fetch|superagent)|okhttp|php-soap|postmanruntime|python-(?:urllib|requests)|scrapy)\/([\w\.]+)`, `(?i)(jsdom|(?<=\()java)\/([\w\.]+)`},
				output: map[string]string{
					Name:    "$1",
					Version: "$2",
					Type:    Library,
				},
			},
		},
	}

	Vehicles = map[string][]regexItem{
		"device": {
			{
				// BYD
				patterns: []string{`(?i)dilink.+(byd) auto`},
				output: map[string]string{
					Vendor: "$1",
				},
			},
			{
				// Rivian
				patterns: []string{`(?i)(rivian) (r1t)`},
				output: map[string]string{
					Vendor: "$1",
					Model:  "$2",
				},
			},
			{
				// Volvo
				patterns: []string{`(?i)vcc.+netfront`},
				output: map[string]string{
					Vendor: "Volvo",
				},
			},
		},
	}

	Bots = map[string][]regexItem{
		"browser": slices.Concat(
			CLIs["browser"],
			Crawlers["browser"],
			Fetchers["browser"],
			Libraries["browser"],
		),
	}
)
