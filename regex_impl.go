package uaparser

import (
	"fmt"
	"github.com/dlclark/regexp2"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

const (
	UAMinLength = 5
	UAMaxLength = 500
)

var (
	extractNumberReg      = regexp.MustCompile(`\d+`)
	dollarReplaceReg      = regexp.MustCompile(`\$\d+`)
	notBrandReg           = regexp.MustCompile(`(?i)not.a.brand`)
	NonNumericOrDotReg    = regexp.MustCompile(`[^\d.]`)
	NonNumericSequenceReg = regexp.MustCompile(`[^\d.]+.`)
	quoteReplaceReg       = regexp.MustCompile(`\\?"`)
)

type mapperItem struct {
	field string
	fn    func(string) string
}

type regexItem struct {
	patterns    []string
	output      map[string]string
	mapperItems []mapperItem
}

var regexMap = map[string][]regexItem{
	"browser": {
		// Most common regardless engine
		{
			patterns: []string{`(?i)\b(?:crmo|crios)\/([\w\.]+)`},
			output: map[string]string{
				Name:    "Mobile Chrome",
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)edg(?:e|ios|a)?\/([\w\.]+)`},
			output: map[string]string{
				Name:    Edge,
				Version: "$1",
			},
		},
		// Presto based
		{
			patterns: []string{`(?i)(opera mini)\/([-\w\.]+)`,
				`(?i)(opera [mobiletab]{3,6})\b.+Version\/([-\w\.]+)`,
				`(?i)(opera)(?:.+Version\/|[\/ ]+)([\w\.]+)`,
			},
			output: map[string]string{
				Name:    "$1",
				Version: "$2",
			},
		},
		{
			patterns: []string{`(?i)opios[\/ ]+([\w\.]+)`},
			output: map[string]string{
				Name:    Opera + " Mini",
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)\bop(?:rg)?x\/([\w\.]+)`},
			output: map[string]string{
				Name:    Opera + " GX",
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)\bopr\/([\w\.]+)`},
			output: map[string]string{
				Name:    Opera,
				Version: "$1",
			},
		},
		// Mixed
		{
			patterns: []string{`(?i)\bb[ai]*d(?:uhd|[ub]*[aekoprswx]{5,6})[\/ ]?([\w\.]+)`},
			output: map[string]string{
				Name:    Baidu,
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)\b(?:mxbrowser|mxios|myie2)\/?([-\w\.]*)\b`},
			output: map[string]string{
				Name:    "Maxthon",
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)(kindle)\/([\w\.]+)`},
			output: map[string]string{
				Name:    "Kindle",
				Version: "$2",
			},
		},
		{
			patterns: []string{`(?i)(lunascape|maxthon|netfront|jasmine|blazer|sleipnir)[\/ ]?([\w\.]*)`},
			output: map[string]string{
				Name:    "$1",
				Version: "$2",
			},
		},
		{
			patterns: []string{`(?i)(avant|iemobile|slim(?:browser|boat|jet))[\/ ]?([\d\.]*)`,
				`(?i)(?:ms|\()(ie) ([\w\.]+)`,
			},
			output: map[string]string{
				Name:    "$1",
				Version: "$2",
			},
		},
		// Blink/Webkit/KHTML based
		{
			patterns: []string{`(?i)(flock|rockmelt|midori|epiphany|silk|skyfire|ovibrowser|bolt|iron|vivaldi|iridium|phantomjs|bowser|qupzilla|falkon|rekonq|puffin|brave|whale(?!.+naver)|qqbrowserlite|duckduckgo|klar|helio|(?=comodo_)?dragon)\/([-\w\.]+)`,
				`(?i)(heytap|ovi|115)browser\/([\d\.]+)`,
				`(?i)(weibo)__([\d\.]+)`},
			output: map[string]string{
				Name:    "$1",
				Version: "$2",
			},
		},
		{
			patterns: []string{`(?i)quark(?:pc)?\/([-\w\.]+)`},
			output: map[string]string{
				Name:    "Quark",
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)\bddg\/([\w\.]+)`},
			output: map[string]string{
				Name:    "DuckDuckGo",
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)(?:\buc? ?browser|(?:juc.+)ucweb)[\/ ]?([\w\.]+)`},
			output: map[string]string{
				Name:    "UCBrowser",
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)microm.+\bqbcore\/([\w\.]+)`,
				`(?i)\bqbcore\/([\w\.]+).+microm`,
				`(?i)micromessenger\/([\w\.]+)`,
			},
			output: map[string]string{
				Name:    "WeChat",
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)konqueror\/([\w\.]+)`},
			output: map[string]string{
				Name:    "Konqueror",
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)trident.+rv[: ]([\w\.]{1,9})\b.+like gecko`},
			output: map[string]string{
				Name:    "IE",
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)ya(?:search)?browser\/([\w\.]+)`},
			output: map[string]string{
				Name:    "Yandex",
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)slbrowser\/([\w\.]+)`},
			output: map[string]string{
				Name:    "Smart Lenovo Browser",
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)(avast|avg)\/([\w\.]+)`},
			output: map[string]string{
				Name:    "$1 Secure Browser",
				Version: "$2",
			},
		},
		{
			patterns: []string{`(?i)\bfocus\/([\w\.]+)`},
			output: map[string]string{
				Name:    "Firefox Focus",
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)\bopt\/([\w\.]+)`},
			output: map[string]string{
				Name:    "Opera Touch",
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)coc_coc\w+\/([\w\.]+)`},
			output: map[string]string{
				Name:    "Coc Coc",
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)dolfin\/([\w\.]+)`},
			output: map[string]string{
				Name:    "Dolphin",
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)coast\/([\w\.]+)`},
			output: map[string]string{
				Name:    "Opera Coast",
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)miuibrowser\/([\w\.]+)`},
			output: map[string]string{
				Name:    "MIUI Browser",
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)fxios\/([\w\.-]+)`},
			output: map[string]string{
				Name:    "mobile Firefox",
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)\bqihoobrowser\/?([\w\.]*)`},
			output: map[string]string{
				Name:    "360",
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)\b(qq)\/([\w\.]+)`},
			output: map[string]string{
				Name:    "$1Browser",
				Version: "$2",
			},
		},
		{
			patterns: []string{`(?i)(oculus|sailfish|huawei|vivo|pico)browser\/([\w\.]+)`},
			output: map[string]string{
				Name:    "$1 Browser",
				Version: "$2",
			},
		},
		{
			patterns: []string{`(?i)samsungbrowser\/([\w\.]+)`},
			output: map[string]string{
				Name:    "Samsung Internet",
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)metasr[\/ ]?([\d\.]+)`},
			output: map[string]string{
				Name:    "Sogou Explorer",
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)(sogou)mo\w+\/([\d\.]+)`},
			output: map[string]string{
				Name:    "Sogou mobile",
				Version: "$2",
			},
		},
		{
			patterns: []string{`(?i)(electron)\/([\w\.]+) safari`,
				`(?i)(tesla)(?: qtcarbrowser|\/(20\d\d\.[-\w\.]+))`,
				`(?i)m?(qqbrowser|2345(?=browser|chrome|explorer))\w*[\/ ]?v?([\w\.]+)`,
			},
			output: map[string]string{
				Name:    "$1",
				Version: "$2",
			},
		},
		{
			patterns: []string{`(?i)(lbbrowser|rekonq)`},
			output: map[string]string{
				Name: "$1",
			},
		},
		{
			patterns: []string{`(?i)ome\/([\w\.]+) \w* ?(iron) saf`,
				`(?i)ome\/([\w\.]+).+qihu (360)[es]e`,
			},
			output: map[string]string{
				Name:    "$2",
				Version: "$1",
			},
		},
		// WebView
		{
			patterns: []string{`(?i)((?:fban\/fbios|fb_iab\/fb4a)(?!.+fbav)|;fbav\/([\w\.]+);)`},
			output: map[string]string{
				Name:    Facebook,
				Version: "$2",
				Type:    "InApp",
			},
		},
		{
			patterns: []string{`(?i)(Klarna)\/([\w\.]+)`,
				`(?i)(kakao(?:talk|story))[\/ ]([\w\.]+)`,
				`(?i)(naver)\(.*?(\d+\.[\w\.]+).*\)`,
				`(?i)(daum)apps[\/ ]([\w\.]+)`,
				`(?i)safari (line)\/([\w\.]+)`,
				`(?i)\b(line)\/([\w\.]+)\/iab`,
				`(?i)(alipay)client\/([\w\.]+)`,
				`(?i)(twitter)(?:and| f.+e\/([\w\.]+))`,
				`(?i)(instagram|snapchat)[\/ ]([-\w\.]+)`,
			},
			output: map[string]string{
				Name:    "$1",
				Version: "$2",
				Type:    InApp,
			},
		},
		{
			patterns: []string{`(?i)\bgsa\/([\w\.]+) .*safari\/`},
			output: map[string]string{
				Name:    "GSA",
				Version: "$1",
				Type:    InApp,
			},
		},
		{
			patterns: []string{`(?i)musical_ly(?:.+app_?Version\/|_)([\w\.]+)`},
			output: map[string]string{
				Name:    "TikTok",
				Version: "$1",
				Type:    InApp,
			},
		},
		{
			patterns: []string{`(?i)\[(linkedin)app\]`},
			output: map[string]string{
				Name: "LinkedIn",
				Type: InApp,
			},
		},
		{
			patterns: []string{`(?i)(chromium)[\/ ]([-\w\.]+)`},
			output: map[string]string{
				Name:    "Chromium",
				Version: "$2",
			},
		},
		{
			patterns: []string{`(?i)headlesschrome(?:\/([\w\.]+)| )`},
			output: map[string]string{
				Name:    "Chrome Headless",
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i) wv\).+(chrome)\/([\w\.]+)`},
			output: map[string]string{
				Name:    "Chrome WebView",
				Version: "$2",
			},
		},
		{
			patterns: []string{`(?i)droid.+ Version\/([\w\.]+)\b.+(?:mobile safari|safari)`},
			output: map[string]string{
				Name:    "Android Browser",
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)chrome\/([\w\.]+) mobile`},
			output: map[string]string{
				Name:    "Mobile Chrome",
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)(chrome|omniweb|arora|[tizenoka]{5} ?browser)\/v?([\w\.]+)`},
			output: map[string]string{
				Name:    "$1",
				Version: "$2",
			},
		},
		{
			patterns: []string{`(?i)Version\/([\w\.\,]+) .*mobile(?:\/\w+ | ?)safari`},
			output: map[string]string{
				Name:    "mobile Safari",
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)iphone .*mobile(?:\/\w+ | ?)safari`},
			output: map[string]string{
				Name: "mobile Safari",
			},
		},
		{
			patterns: []string{`(?i)Version\/([\w\.\,]+) .*(safari)`},
			output: map[string]string{
				Name:    "$2",
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)webkit.+?(mobile ?safari|safari)(\/[\w\.]+)`},
			output: map[string]string{
				Name:    "$1",
				Version: "1",
			},
		},
		{
			patterns: []string{`(?i)(webkit|khtml)\/([\w\.]+)`},
			output: map[string]string{
				Name:    "$1",
				Version: "$2",
			},
		},
		// Gecko based
		{
			patterns: []string{`(?i)(?:mobile|tablet);.*(firefox)\/([\w\.-]+)`},
			output: map[string]string{
				Name:    "mobile Firefox",
				Version: "$2",
			},
		},
		{
			patterns: []string{`(?i)(navigator|netscape\d?)\/([-\w\.]+)`},
			output: map[string]string{
				Name:    "Netscape",
				Version: "$2",
			},
		},
		{
			patterns: []string{`(?i)(wolvic|librewolf)\/([\w\.]+)`},
			output: map[string]string{
				Name:    "$1",
				Version: "$2",
			},
		},
		{
			patterns: []string{`(?i)mobile vr; rv:([\w\.]+)\).+firefox`},
			output: map[string]string{
				Name:    "Firefox Reality",
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)ekiohf.+(flow)\/([\w\.]+)`,
				`(?i)(swiftfox)`,
				`(?i)(icedragon|iceweasel|camino|chimera|fennec|maemo browser|minimo|conkeror)[\/ ]?([\w\.\+]+)`,
				`(?i)(seamonkey|k-meleon|icecat|iceape|firebird|phoenix|palemoon|basilisk|waterfox)\/([-\w\.]+)$`,
				`(?i)(firefox)\/([\w\.]+)`,
				`(?i)(mozilla)\/([\w\.]+) .+rv\:.+gecko\/\d+`,
			},
			output: map[string]string{
				Name:    "$1",
				Version: "$2",
			},
		},
		// Other
		{
			patterns: []string{`(?i)(amaya|dillo|doris|icab|ladybird|lynx|mosaic|netsurf|obigo|polaris|w3m|(?:go|ice|up)[\. ]?browser)[-\/ ]?v?([\w\.]+)`,
				`(?i)\b(links) \(([\w\.]+)`,
			},
			output: map[string]string{
				Name:    "$1",
				Version: "$2",
			},
			mapperItems: []mapperItem{
				{
					field: Version,
					fn:    func(s string) string { return strings.ReplaceAll(s, "_", ".") },
				},
			},
		},
		{
			patterns: []string{`(?i)(cobalt)\/([\w\.]+)`},
			output: map[string]string{
				Name:    "$1",
				Version: "$2",
			},
			mapperItems: []mapperItem{
				{
					field: Version,
					fn:    func(s string) string { return NonNumericSequenceReg.ReplaceAllString(s, "") },
				},
			},
		},
	},
	"cpu": {
		{
			patterns: []string{`(?i)\b(?:(amd|x|x86[-_]?|wow|win)64)\b`},
			output: map[string]string{
				Architecture: "amd64",
			},
		},
		{
			patterns: []string{`(?i)(ia32(?=;))`, `(?i)\b((i[346]|x)86)(pc)?\b`},
			output: map[string]string{
				Architecture: "ia32",
			},
		},
		{
			patterns: []string{`(?i)\b(aarch64|arm(v?8e?l?|_?64))\b`},
			output: map[string]string{
				Architecture: "arm64",
			},
		},
		{
			patterns: []string{`(?i)\b(arm(?:v[67])?ht?n?[fl]p?)\b`},
			output: map[string]string{
				Architecture: "armhf",
			},
		},
		{
			patterns: []string{`(?i)( (ce|mobile); ppc;|\/[\w\.]+arm\b)`},
			output: map[string]string{
				Architecture: "arm",
			},
		},
		{
			patterns: []string{`(?i)((?:ppc|powerpc)(?:64)?)(?: mac|;|\))`},
			output: map[string]string{
				Architecture: "$1",
			},
			mapperItems: []mapperItem{
				{
					field: Architecture,
					fn: func(s string) string {
						return strings.ToLower(strings.ReplaceAll(s, "ower", ""))
					},
				},
			},
		},
		{
			patterns: []string{`(?i)(sun4\w)[;\)]`},
			output: map[string]string{
				Architecture: "sparc",
			},
		},
		{
			patterns: []string{`(?i)((?:avr32|ia64(?=;))|68k(?=\))|\barm(?=v(?:[1-7]|[5-7]1)l?|;|eabi)|(?=atmel )avr|(?:irix|mips|sparc)(?:64)?\b|pa-risc)`},
			output: map[string]string{
				Architecture: "$1",
			},
			mapperItems: []mapperItem{
				{
					field: Architecture,
					fn:    strings.ToLower,
				},
			},
		},
	},
	"device": {
		//////////////////////////
		// MobileS & TabletS
		/////////////////////////
		// Samsung
		{
			patterns: []string{`(?i)\b(sch-i[89]0\d|shw-m380s|sm-[ptx]\w{2,4}|gt-[pn]\d{2,4}|sgh-t8[56]9|nexus 10)`},
			output: map[string]string{
				Model:  "$1",
				Vendor: Samsung,
				Type:   Tablet,
			},
		},
		{
			patterns: []string{`(?i)\b((?:s[cgp]h|gt|sm)-(?![lr])\w+|sc[g-]?[\d]+a?|galaxy nexus)`,
				`(?i)samsung[- ]((?!sm-[lr])[-\w]+)`,
				`(?i)sec-(sgh\w+)`},
			output: map[string]string{
				Model:  "$1",
				Vendor: Samsung,
				Type:   Mobile,
			},
		},
		// Apple
		{
			patterns: []string{`(?i)(?:\/|\()(ip(?:hone|od)[\w, ]*)(?:\/|;)`},
			output: map[string]string{
				Model:  "$1",
				Vendor: Apple,
				Type:   Mobile,
			},
		},
		{
			patterns: []string{`(?i)\((ipad);[-\w\),; ]+apple`,
				`(?i)applecoremedia\/[\w\.]+ \((ipad)`,
				`(?i)\b(ipad)\d\d?,\d\d?[;\]].+ios`},
			output: map[string]string{
				Model:  "$1",
				Vendor: Apple,
				Type:   Tablet,
			},
		},
		{
			patterns: []string{`(?i)(macintosh);`},
			output: map[string]string{
				Model:  "$1",
				Vendor: Apple,
			},
		},
		// Sharp
		{
			patterns: []string{`(?i)\b(sh-?[altvz]?\d\d[a-ekm]?)`},
			output: map[string]string{
				Model:  "$1",
				Vendor: Sharp,
				Type:   Mobile,
			},
		},
		// Honor
		{
			patterns: []string{`(?i)\b((?:brt|eln|hey2?|gdi|jdn)-a?[lnw]09|(?:ag[rm]3?|jdn2|kob2)-a?[lw]0[09]hn)(?: bui|\)|;)`},
			output: map[string]string{
				Model:  "$1",
				Vendor: Honor,
				Type:   Tablet,
			},
		},
		{
			patterns: []string{`(?i)honor([-\w ]+)[;\)]`},
			output: map[string]string{
				Model:  "$1",
				Vendor: Honor,
				Type:   Mobile,
			},
		},
		// Huawei
		{
			patterns: []string{`(?i)\b((?:ag[rs][2356]?k?|bah[234]?|bg[2o]|bt[kv]|cmr|cpn|db[ry]2?|jdn2|got|kob2?k?|mon|pce|scm|sht?|[tw]gr|vrd)-[ad]?[lw][0125][09]b?|605hw|bg2-u03|(?:gem|fdr|m2|ple|t1)-[7a]0[1-4][lu]|t1-a2[13][lw]|mediapad[\w\. ]*(?= bui|\)))\b(?!.+d\/s)`},
			output: map[string]string{
				Model:  "$1",
				Vendor: Huawei,
				Type:   Tablet,
			},
		},
		{
			patterns: []string{`(?i)(?:huawei)([-\w ]+)[;\)]`,
				`(?i)\b(nexus 6p|\w{2,4}e?-[atu]?[ln][\dx][012359c][adn]?)\b(?!.+d\/s)`,
			},
			output: map[string]string{
				Model:  "$1",
				Vendor: Huawei,
				Type:   Mobile,
			},
		},
		// Xiaomi
		{
			patterns: []string{`(?i)oid[^\)]+; (2[\dbc]{4}(182|283|rp\w{2})[cgl]|m2105k81a?c)(?: bui|\))`,
				`(?i)\b((?:red)?mi[-_ ]?pad[\w- ]*)(?: bui|\))`},
			output: map[string]string{
				Model:  "$1",
				Vendor: Xiaomi,
				Type:   Tablet,
			},
		},
		{
			patterns: []string{`(?i)\b(poco[\w ]+|m2\d{3}j\d\d[a-z]{2})(?: bui|\))`,
				`(?i)\b; (\w+) build\/hm\1`,
				`(?i)\b(hm[-_ ]?note?[_ ]?(?:\d\w)?) bui`,
				`(?i)\b(redmi[\-_ ]?(?:note|k)?[\w_ ]+)(?: bui|\))`,
				`(?i)oid[^\)]+; (m?[12][0-389][01]\w{3,6}[c-y])( bui|; wv|\))`,
				`(?i)\b(mi[-_ ]?(?:a\d|one|one[_ ]plus|note lte|max|cc)?[_ ]?(?:\d?\w?)[_ ]?(?:plus|se|lite|pro)?)(?: bui|\))`,
				`(?i) ([\w ]+) miui\/v?\d`,
			},
			output: map[string]string{
				Model:  "$1",
				Vendor: Xiaomi,
				Type:   Mobile,
			},
		},
		// OPPO
		{
			patterns: []string{`(?i); (\w+) bui.+ oppo`,
				`(?i)\b(cph[12]\d{3}|p(?:af|c[al]|d\w|e[ar])[mt]\d0|x9007|a101op)\b`,
			},
			output: map[string]string{
				Model:  "$1",
				Vendor: OPPO,
				Type:   Mobile,
			},
		},
		{
			patterns: []string{`(?i)\b(opd2(\d{3}a?))(?: bui|\))`},
			output: map[string]string{
				Model:  "$1",
				Vendor: "$2",
				Type:   Tablet,
			},
			mapperItems: []mapperItem{
				{
					field: Vendor,
					fn: func(s string) string {
						return strMapper(s, map[string][]string{"OnePlus": {"304", "403", "203"}, "*": {OPPO}})
					},
				},
			},
		},
		// Vivo
		{
			patterns: []string{`(?i)vivo (\w+)(?: bui|\))`,
				`(?i)\b(v[12]\d{3}\w?[at])(?: bui|;)`,
			},
			output: map[string]string{
				Model:  "$1",
				Vendor: "Vivo",
				Type:   Mobile,
			},
		},
		// Realme
		{
			patterns: []string{`(?i)\b(rmx[1-3]\d{3})(?: bui|;|\))`},
			output: map[string]string{
				Model:  "$1",
				Vendor: "Realme",
				Type:   Mobile,
			},
		},
		// Motorola
		{
			patterns: []string{`(?i)\b(milestone|droid(?:[2-4x]| (?:bionic|x2|pro|razr))?:?( 4g)?)\b[\w ]+build\/`,
				`(?i)\bmot(?:orola)?[- ](\w*)`,
				`(?i)((?:moto(?! 360)[\w\(\) ]+|xt\d{3,4}|nexus 6)(?= bui|\)))`,
			},
			output: map[string]string{
				Model:  "$1",
				Vendor: Motorola,
				Type:   Mobile,
			},
		},
		{
			patterns: []string{`(?i)\b(mz60\d|xoom[2 ]{0,2}) build\/`},
			output: map[string]string{
				Model:  "$1",
				Vendor: Motorola,
				Type:   Tablet,
			},
		},
		// LG
		{
			patterns: []string{`(?i)((?=lg)?[vl]k\-?\d{3}) bui| 3\.[-\w; ]{10}lg?-([06cv9]{3,4})`},
			output: map[string]string{
				Model:  "$1",
				Vendor: LG,
				Type:   Tablet,
			},
		},
		{
			patterns: []string{`(?i)(lm(?:-?f100[nv]?|-[\w\.]+)(?= bui|\))|nexus [45])`,
				`(?i)\blg[-e;\/ ]+((?!browser|netcast|android tv|watch)\w+)`,
				`(?i)\blg-?([\d\w]+) bui`,
			},
			output: map[string]string{
				Model:  "$1",
				Vendor: LG,
				Type:   Mobile,
			},
		},
		// Lenovo
		{
			patterns: []string{`(?i)(ideatab[-\w ]+|602lv|d-42a|a101lv|a2109a|a3500-hv|s[56]000|pb-6505[my]|tb-?x?\d{3,4}(?:f[cu]|xu|[av])|yt\d?-[jx]?\d+[lfmx])( bui|;|\)|\/)`,
				`(?i)lenovo ?(b[68]0[08]0-?[hf]?|tab(?:[\w- ]+?)|tb[\w-]{6,7})( bui|;|\)|\/)`,
			},
			output: map[string]string{
				Model:  "$1",
				Vendor: Lenovo,
				Type:   Tablet,
			},
		},
		// Nokia
		{
			patterns: []string{`(?i)(nokia) (t[12][01])`},
			output: map[string]string{
				Vendor: "$1",
				Model:  "$2",
				Type:   Tablet,
			},
		},
		{
			patterns: []string{`(?i)(?:maemo|nokia).*(n900|lumia \d+|rm-\d+)`,
				`(?i)nokia[-_ ]?(([-\w\. ]*))`,
			},
			output: map[string]string{
				Model:  "$1",
				Vendor: "Nokia",
				Type:   Mobile,
			},
			mapperItems: []mapperItem{
				{
					field: Model,
					fn:    func(s string) string { return strings.ReplaceAll(s, "_", " ") },
				},
			},
		},
		// Google
		{
			patterns: []string{`(?i)(pixel (c|tablet))\b`},
			output: map[string]string{
				Model:  "$1",
				Vendor: Google,
				Type:   Tablet,
			},
		},
		{
			patterns: []string{`(?i)droid.+; (pixel[\daxl ]{0,6})(?: bui|\))`},
			output: map[string]string{
				Model:  "$1",
				Vendor: Google,
				Type:   Mobile,
			},
		},
		// Sony
		{
			patterns: []string{`(?i)droid.+; (a?\d[0-2]{2}so|[c-g]\d{4}|so[-gl]\w+|xq-a\w[4-7][12])(?= bui|\).+chrome\/(?![1-6]{0,1}\d\.))`},
			output: map[string]string{
				Model:  "$1",
				Vendor: Sony,
				Type:   Mobile,
			},
		},
		{
			patterns: []string{`(?i)sony tablet [ps]`, `(?i)\b(?:sony)?sgp\w+(?: bui|\))`},
			output: map[string]string{
				Model:  "Xperia Tablet",
				Vendor: Sony,
				Type:   Tablet,
			},
		},
		// OnePlus
		{
			patterns: []string{`(?i) (kb2005|in20[12]5|be20[12][59])\b`, `(?i)(?:one)?(?:plus)? (a\d0\d\d)(?: b|\))`},
			output: map[string]string{
				Model:  "$1",
				Vendor: "OnePlus",
				Type:   Mobile,
			},
		},
		// Amazon
		{
			patterns: []string{`(?i)(alexa)webm`, `(?i)(kf[a-z]{2}wi|aeo(?!bc)\w\w)( bui|\))`, `(?i)(kf[a-z]+)( bui|\)).+silk\/`},
			output: map[string]string{
				Model:  "$1",
				Vendor: Amazon,
				Type:   Tablet,
			},
		},
		{
			patterns: []string{`(?i)((?:sd|kf)[0349hijorstuw]+)( bui|\)).+silk\/`},
			output: map[string]string{
				Model:  "Fire Phone $1",
				Vendor: Amazon,
				Type:   Mobile,
			},
		},
		// BlackBerry
		{
			patterns: []string{`(?i)(playbook);[-\w\),; ]+(rim)`},
			output: map[string]string{
				Model:  "$1",
				Vendor: "$2",
				Type:   Tablet,
			},
		},
		{
			patterns: []string{`(?i)\b((?:bb[a-f]|st[hv])100-\d)`, `(?i)\(bb10; (\w+)`},
			output: map[string]string{
				Model:  "$1",
				Vendor: Blackberry,
				Type:   Mobile,
			},
		},
		// Asus
		{
			patterns: []string{`(?i)(?:\b|asus_)(transfo[prime ]{4,10} \w+|eeepc|slider \w+|nexus 7|padfone|p00[cj])`},
			output: map[string]string{
				Model:  "$1",
				Vendor: ASUS,
				Type:   Tablet,
			},
		},
		{
			patterns: []string{`(?i) (z[bes]6[027][012][km][ls]|zenfone \d\w?)\b`},
			output: map[string]string{
				Model:  "$1",
				Vendor: ASUS,
				Type:   Mobile,
			},
		},
		// HTC
		{
			patterns: []string{`(?i)(nexus 9)`},
			output: map[string]string{
				Model:  "$1",
				Vendor: "HTC",
				Type:   Tablet,
			},
		},
		{
			patterns: []string{`(?i)(htc)[-;_ ]{1,2}([\w ]+(?=\)| bui)|\w+)`},
			output: map[string]string{
				Model:  "$2",
				Vendor: "$1",
				Type:   Mobile,
			},
		},
		// ZTE
		{
			patterns: []string{`(?i)(zte)[- ]([\w ]+?)(?: bui|\/|\))`, `(?i)(alcatel|geeksphone|nexian|panasonic(?!(?:;|\.))|sony(?!-bra))[-_ ]?([-\w]*)`},
			output: map[string]string{
				Vendor: "$1",
				Model:  "$2",
				Type:   Mobile,
			},
		},
		// TCL
		{
			patterns: []string{`(?i)tcl (xess p17aa)`, `(?i)droid [\w\.]+; ((?:8[14]9[16]|9(?:0(?:48|60|8[01])|1(?:3[27]|66)|2(?:6[69]|9[56])|466))[gqswx])(_\w(\w|\w\w))?(\)| bui)`},
			output: map[string]string{
				Model:  "$1",
				Vendor: "TCL",
				Type:   Tablet,
			},
		},
		{
			patterns: []string{`(?i)droid [\w\.]+; (418(?:7d|8v)|5087z|5102l|61(?:02[dh]|25[adfh]|27[ai]|56[dh]|59k|65[ah])|a509dl|t(?:43(?:0w|1[adepqu])|50(?:6d|7[adju])|6(?:09dl|10k|12b|71[efho]|76[hjk])|7(?:66[ahju]|67[hw]|7[045][bh]|71[hk]|73o|76[ho]|79w|81[hks]?|82h|90[bhsy]|99b)|810[hs]))(_\w(\w|\w\w))?(\)| bui)`},
			output: map[string]string{
				Model:  "$1",
				Vendor: "TCL",
				Type:   Mobile,
			},
		},
		// itel
		{
			patterns: []string{`(?i)(itel) ((\w+))`},
			output: map[string]string{
				Vendor: "$1",
				Model:  "$2",
				Type:   "$3",
			},
			mapperItems: []mapperItem{
				{
					field: Vendor,
					fn:    strings.ToLower,
				},
				{
					field: Type,
					fn: func(s string) string {
						return strMapper(s, map[string][]string{Tablet: {"p10001l", "w7001"}, "*": {Mobile}})
					},
				},
			},
		},
		// Acer
		{
			patterns: []string{`(?i)droid.+; ([ab][1-7]-?[0178a]\d\d?)`},
			output: map[string]string{
				Model:  "$1",
				Vendor: "Acer",
				Type:   Tablet,
			},
		},
		// Meizu
		{
			patterns: []string{`(?i)droid.+; (m[1-5] note) bui`, `(?i)\bmz-([-\w]{2,})`},
			output: map[string]string{
				Model:  "$1",
				Vendor: "Meizu",
				Type:   Mobile,
			},
		},
		// Ulefone
		{
			patterns: []string{`(?i); ((?:power )?armor(?:[\w ]{0,8}))(?: bui|\))`},
			output: map[string]string{
				Model:  "$1",
				Vendor: "Ulefone",
				Type:   Mobile,
			},
		},
		// Energizer
		{
			patterns: []string{`(?i); (energy ?\w+)(?: bui|\))`, `(?i); energizer ([\w ]+)(?: bui|\))`},
			output: map[string]string{
				Model:  "$1",
				Vendor: "Energizer",
				Type:   Mobile,
			},
		},
		// Cat
		{
			patterns: []string{`(?i); cat (b35);`, `(?i); (b15q?|s22 flip|s48c|s62 pro)(?: bui|\))`},
			output: map[string]string{
				Model:  "$1",
				Vendor: "Cat",
				Type:   Mobile,
			},
		},
		// Smartfren
		{
			patterns: []string{`(?i)((?:new )?andromax[\w- ]+)(?: bui|\))`},
			output: map[string]string{
				Model:  "$1",
				Vendor: "Smartfren",
				Type:   Mobile,
			},
		},
		// Nothing
		{
			patterns: []string{`(?i)droid.+; (a(?:015|06[35]|142p?))`},
			output: map[string]string{
				Model:  "$1",
				Vendor: "Nothing",
				Type:   Mobile,
			},
		},
		// MIXED
		{
			patterns: []string{`(?i)(imo) (tab \w+)`, `(?i)(infinix) (x1101b?)`},
			output: map[string]string{
				Vendor: "$1",
				Model:  "$2",
				Type:   Tablet,
			},
		},
		{
			patterns: []string{
				`(?i)(blackberry|benq|palm(?=\-)|sonyericsson|acer|asus(?! zenw)|dell|jolla|meizu|motorola|polytron|infinix|tecno|micromax|advan)[-_ ]?([-\w]*)`,
				`(?i); (hmd|imo) ([\w ]+?)(?: bui|\))`,
				`(?i)(hp) ([\w ]+\w)`,
				`(?i)(asus)-?(\w+)`,
				`(?i)(microsoft); (lumia[\w ]+)`,
				`(?i)(lenovo)[-_ ]?([-\w ]+?)(?: bui|\)|\/)`,
				`(?i)(oppo) ?([\w ]+) bui`,
			},
			output: map[string]string{
				Vendor: "$1",
				Model:  "$2",
				Type:   Mobile,
			},
		},
		{
			patterns: []string{`(?i)(kobo)\s(ereader|touch)`, `(archos) (gamepad2?)`, `(hp).+(touchpad(?!.+tablet)|tablet)`, `(kindle)\/([\w\.]+)`},
			output: map[string]string{
				Vendor: "$1",
				Model:  "$2",
				Type:   Tablet,
			},
		},
		{
			patterns: []string{`(?i)(surface duo)`},
			output: map[string]string{
				Model:  "$1",
				Vendor: Microsoft,
				Type:   Tablet,
			},
		},
		{
			patterns: []string{`(?i)droid [\d\.]+; (fp\du?)(?: b|\))`},
			output: map[string]string{
				Model:  "$1",
				Vendor: "Fairphone",
				Type:   Mobile,
			},
		},
		{
			patterns: []string{`(?i)((?:tegranote|shield t(?!.+d tv))[\w- ]*?)(?: b|\))`},
			output: map[string]string{
				Model:  "$1",
				Vendor: "Nvidia",
				Type:   Tablet,
			},
		},
		{
			patterns: []string{`(?i)(sprint) (\w+)`},
			output: map[string]string{
				Vendor: "$1",
				Model:  "$2",
				Type:   Mobile,
			},
		},
		{
			patterns: []string{`(?i)(kin\.[onetw]{3})`},
			output: map[string]string{
				Model:  "$1",
				Vendor: Microsoft,
				Type:   Mobile,
			},
		},
		{
			patterns: []string{`(?i)droid.+; ([c6]+|et5[16]|mc[239][23]x?|vc8[03]x?)\)`, `droid.+; (ec30|ps20|tc[2-8]\d[kx])\)`},
			output: map[string]string{
				Model:  "$1",
				Vendor: Zebra,
				Type:   Mobile,
			},
		},

		///////////////////
		// SmartTVs
		///////////////////

		{
			patterns: []string{`(?i)smart-tv.+(samsung)`},
			output: map[string]string{
				Vendor: Samsung,
				Type:   SmartTV,
			},
		},
		{
			patterns: []string{`(?i)hbbtv.+maple;(\d+)`},
			output: map[string]string{
				Model:  "SmartTV$1",
				Vendor: Samsung,
				Type:   SmartTV,
			},
		},
		{
			patterns: []string{`(?i)(nux; netcast.+smarttv|lg (netcast\.tv-201\d|android tv))`},
			output: map[string]string{
				Vendor: LG,
				Type:   SmartTV,
			},
		},
		{
			patterns: []string{`(?i)(apple) ?tv`},
			output: map[string]string{
				Vendor: Apple,
				Model:  "Apple TV",
				Type:   SmartTV,
			},
		},
		{
			patterns: []string{`(?i)crkey.*devicetype\/chromecast`},
			output: map[string]string{
				Model:  Chromecast + " Third Generation",
				Vendor: Google,
				Type:   SmartTV,
			},
		},
		{
			patterns: []string{`(?i)crkey.*devicetype\/([^/]*)`},
			output: map[string]string{
				Model:  Chromecast + " $1",
				Vendor: Google,
				Type:   SmartTV,
			},
		},
		{
			patterns: []string{`(?i)fuchsia.*crkey`},
			output: map[string]string{
				Model:  Chromecast + " Nest Hub",
				Vendor: Google,
				Type:   SmartTV,
			},
		},
		{
			patterns: []string{`(?i)crkey`},
			output: map[string]string{
				Model:  Chromecast,
				Vendor: Google,
				Type:   SmartTV,
			},
		},
		{
			patterns: []string{`(?i)droid.+aft(\w+)( bui|\))`},
			output: map[string]string{
				Model:  "$1",
				Vendor: Amazon,
				Type:   SmartTV,
			},
		},
		{
			patterns: []string{`(?i)(shield \w+ tv)`},
			output: map[string]string{
				Model:  "$1",
				Vendor: "Nvidia",
				Type:   SmartTV,
			},
		},
		{
			patterns: []string{`(?i)\(dtv[\);].+(aquos)`, `(?i)(aquos-tv[\w ]+)\)`},
			output: map[string]string{
				Model:  "$1",
				Vendor: Sharp,
				Type:   SmartTV,
			},
		},
		{
			patterns: []string{`(?i)(bravia[\w ]+)( bui|\))`},
			output: map[string]string{
				Model:  "$1",
				Vendor: Sony,
				Type:   SmartTV,
			},
		},
		{
			patterns: []string{`(?i)(mi(tv|box)-?\w+) bui`},
			output: map[string]string{
				Model:  "$1",
				Vendor: Xiaomi,
				Type:   SmartTV,
			},
		},
		{
			patterns: []string{`(?i)Hbbtv.*(technisat) (.*);`},
			output: map[string]string{
				Vendor: "$1",
				Model:  "$2",
				Type:   SmartTV,
			},
		},
		{
			patterns: []string{`(?i)\b(roku)[\dx]*[\)\/]((?:dvp-)?[\d\.]*)`, `(?i)hbbtv\/\d+\.\d+\.\d+ +\([\w\+ ]*; *([\w\d][^;]*);([^;]*)`},
			output: map[string]string{
				Vendor: "$1",
				Model:  "$2",
				Type:   SmartTV,
			},
			mapperItems: []mapperItem{
				{
					field: Vendor,
					fn:    strings.TrimSpace,
				},
				{
					field: Model,
					fn:    strings.TrimSpace,
				},
			},
		},
		{
			patterns: []string{`(?i)droid.+; ([\w- ]+) (?:android tv|smart[- ]?tv)`},
			output: map[string]string{
				Model: "$1",
				Type:  SmartTV,
			},
		},
		{
			patterns: []string{`(?i)\b(android tv|smart[- ]?tv|opera tv|tv; rv:)\b`},
			output: map[string]string{
				Type: SmartTV,
			},
		},

		///////////////////
		// Consoles
		///////////////////

		{
			patterns: []string{`(?i)(ouya)`},
			output: map[string]string{
				Vendor: "Ouya",
				Type:   Console,
			},
		},
		{
			patterns: []string{`(?i)(nintendo) (\w+)`},
			output: map[string]string{
				Vendor: "Nintendo",
				Model:  "$2",
				Type:   Console,
			},
		},
		{
			patterns: []string{`(?i)droid.+; (shield)( bui|\))`},
			output: map[string]string{
				Model:  "$1",
				Vendor: "Nvidia",
				Type:   Console,
			},
		},
		{
			patterns: []string{`(?i)(playstation \w+)`},
			output: map[string]string{
				Model:  "$1",
				Vendor: Sony,
				Type:   Console,
			},
		},
		{
			patterns: []string{`(?i)\b(xbox(?: one)?(?!; xbox))[\); ]`},
			output: map[string]string{
				Model:  "$1",
				Vendor: Microsoft,
				Type:   Console,
			},
		},

		///////////////////
		// Wearables
		///////////////////

		{
			patterns: []string{`(?i)\b(sm-[lr]\d\d[0156][fnuw]?s?|gear live)\b`},
			output: map[string]string{
				Model:  "$1",
				Vendor: Samsung,
				Type:   Wearable,
			},
		},
		{
			patterns: []string{`(?i)((pebble))app`},
			output: map[string]string{
				Vendor: "Pebble",
				Type:   Wearable,
			},
		},
		{
			patterns: []string{`(?i)(asus|google|lg|oppo) ((pixel |zen)?watch[\w ]*)( bui|\))`},
			output: map[string]string{
				Vendor: "$1",
				Model:  "$2",
				Type:   Wearable,
			},
		},
		{
			patterns: []string{`(?i)(ow(?:19|20)?we?[1-3]{1,3})`},
			output: map[string]string{
				Vendor: OPPO,
				Model:  "$1",
				Type:   Wearable,
			},
		},
		{
			patterns: []string{`(?i)(watch)(?: ?os[,\/]|\d,\d\/)[\d\.]+`},
			output: map[string]string{
				Model:  "$1",
				Vendor: Apple,
				Type:   Wearable,
			},
		},
		{
			patterns: []string{`(?i)(opwwe\d{3})`},
			output: map[string]string{
				Model:  "$1",
				Vendor: "OnePlus",
				Type:   Wearable,
			},
		},
		{
			patterns: []string{`(?i)(moto 360)`},
			output: map[string]string{
				Model:  "$1",
				Vendor: "Motorola",
				Type:   Wearable,
			},
		},
		{
			patterns: []string{`(?i)(smartwatch 3)`},
			output: map[string]string{
				Model:  "$1",
				Vendor: "Sony",
				Type:   Wearable,
			},
		},
		{
			patterns: []string{`(?i)(g watch r)`},
			output: map[string]string{
				Model:  "$1",
				Vendor: "LG",
				Type:   Wearable,
			},
		},
		{
			patterns: []string{`(?i)droid.+; (wt63?0{2,3})\)`},
			output: map[string]string{
				Model:  "$1",
				Vendor: Zebra,
				Type:   Wearable,
			},
		},

		///////////////////
		// XR
		///////////////////

		{
			patterns: []string{`(?i)droid.+; (glass) \d`},
			output: map[string]string{
				Model:  "$1",
				Vendor: Google,
				Type:   XR,
			},
		},
		{
			patterns: []string{`(?i)(pico) (4|neo3(?: link|pro)?)`},
			output: map[string]string{
				Vendor: "Pico",
				Model:  "$2",
				Type:   XR,
			},
		},
		{
			patterns: []string{`(?i); (quest( \d| pro)?)`},
			output: map[string]string{
				Model:  "$1",
				Vendor: Facebook,
				Type:   XR,
			},
		},

		///////////////////
		// Embedded
		///////////////////

		{
			patterns: []string{`(?i)(tesla)(?: qtcarbrowser|\/[-\w\.]+)`},
			output: map[string]string{
				Vendor: "Tesla",
				Type:   Embedded,
			},
		},
		{
			patterns: []string{`(?i)(aeobc)\b`},
			output: map[string]string{
				Model:  "$1",
				Vendor: Amazon,
				Type:   Embedded,
			},
		},
		{
			patterns: []string{`(?i)(homepod).+mac os`},
			output: map[string]string{
				Model:  "$1",
				Vendor: Apple,
				Type:   Embedded,
			},
		},
		{
			patterns: []string{`(?i)windows iot`},
			output: map[string]string{
				Type: Embedded,
			},
		},

		////////////////////
		// Mixed (Generic)
		///////////////////

		{
			patterns: []string{`(?i)droid .+?; ([^;]+?)(?: bui|; wv\)|\) applew).+?(mobile|vr|\d) safari`},
			output: map[string]string{
				Model: "$1",
				Type:  "$2",
			},
			mapperItems: []mapperItem{
				{
					field: Type,
					fn: func(s string) string {
						return strMapper(s, map[string][]string{Mobile: {Mobile}, "xr": {"VR"}, "*": {Tablet}})
					},
				},
			},
		},
		{
			patterns: []string{`(?i)droid .+?; ([^;]+?)(?: bui|\) applew).+?(?! mobile) safari`},
			output: map[string]string{
				Model: "$1",
				Type:  Tablet,
			},
		},
		{
			patterns: []string{`(?i)\b((tablet|tab)[;\/]|focus\/\d(?!.+mobile))`},
			output: map[string]string{
				Type: Tablet,
			},
		},
		{
			patterns: []string{`(?i)(phone|mobile(?:[;\/]| [ \w\/\.]*safari)|pda(?=.+windows ce))`},
			output: map[string]string{
				Type: Mobile,
			},
		},
		{
			patterns: []string{`(?i)droid .+?; ([\w\. -]+)( bui|\))`},
			output: map[string]string{
				Model:  "$1",
				Vendor: "Generic",
			},
		},
	},
	"engine": {
		{
			patterns: []string{`(?i)windows.+ edge\/([\w\.]+)`},
			output: map[string]string{
				Name:    "EdgeHTML",
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)(arkweb)\/([\w\.]+)`},
			output: map[string]string{
				Name:    "$1",
				Version: "$2",
			},
		},
		{
			patterns: []string{`(?i)webkit\/537\.36.+chrome\/(?!27)([\w\.]+)`},
			output: map[string]string{
				Name:    "Blink",
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)(presto)\/([\w\.]+)`,
				`(?i)(webkit|trident|netfront|netsurf|amaya|lynx|w3m|goanna|servo)\/([\w\.]+)`,
				`(?i)ekioh(flow)\/([\w\.]+)`,
				`(?i)(khtml|tasman|links)[\/ ]\(?([\w\.]+)`,
				`(?i)(icab)[\/ ]([23]\.[\d\.]+)`,
				`(?i)\b(libweb)`,
			},
			output: map[string]string{
				Name:    "$1",
				Version: "$2",
			},
		},
		{
			patterns: []string{`(?i)ladybird\/`},
			output: map[string]string{
				Name: "LibWeb",
			},
		},
		{
			patterns: []string{`(?i)rv\:([\w\.]{1,9})\b.+(gecko)`},
			output: map[string]string{
				Name:    "$2",
				Version: "$1",
			},
		},
	},
	"os": {
		// Windows (iTunes)
		{
			patterns: []string{`(?i)microsoft (windows) (vista|xp)`},
			output: map[string]string{
				Name:    "$1",
				Version: "$2",
			},
		},
		// Windows Phone
		{
			patterns: []string{`(?i)(windows (?:phone(?: os)?|mobile|iot))[\/ ]?([\d\.\w ]*)`},
			output: map[string]string{
				Name:    "$1",
				Version: "$2",
			},
			mapperItems: []mapperItem{
				{
					field: Version,
					fn: func(s string) string {
						return strMapper(s, windowsVersionMap)
					},
				},
			},
		},
		// Windows RT
		{
			patterns: []string{`(?i)windows nt 6\.2; (arm)`,
				`(?i)windows[\/ ]([ntce\d\. ]+\w)(?!.+xbox)`,
				`(?i)(?:win(?=3|9|n)|win 9x )([nt\d\.]+)`,
			},
			output: map[string]string{
				Name:    Windows,
				Version: "$1",
			},
			mapperItems: []mapperItem{
				{
					field: Version,
					fn: func(s string) string {
						return strMapper(s, windowsVersionMap)
					},
				},
			},
		},
		// iOS/macOS
		{
			patterns: []string{
				`(?i)[adehimnop]{4,7}\b(?:.*os ([\w]+) like mac|; opera)`,
				`(?i)(?:ios;fbsv\/|iphone.+ios[\/ ])([\d\.]+)`,
				`(?i)cfnetwork\/.+darwin`,
			},
			output: map[string]string{
				Name:    "iOS",
				Version: "$1",
			},
			mapperItems: []mapperItem{
				{
					field: Version,
					fn: func(str string) string {
						return strings.ReplaceAll(str, "_", ".")
					},
				},
			},
		},
		{
			patterns: []string{
				`(?i)(mac os x) ?([\w\. ]*)`,
				`(?i)(macintosh|mac_powerpc\b)(?!.+haiku)`,
			},
			output: map[string]string{
				Name:    "macOS",
				Version: "$2",
			},
			mapperItems: []mapperItem{
				{
					field: Version,
					fn: func(str string) string {
						return strings.ReplaceAll(str, "_", ".")
					},
				},
			},
		},
		// Google Chromecast
		{
			patterns: []string{`(?i)android ([\d\.]+).*crkey`},
			output: map[string]string{
				Name:    Chromecast + " Android",
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)fuchsia.*crkey\/([\d\.]+)`},
			output: map[string]string{
				Name:    Chromecast + " Fuchsia",
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)crkey\/([\d\.]+).*devicetype\/smartspeaker`},
			output: map[string]string{
				Name:    Chromecast + " SmartSpeaker",
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)linux.*crkey\/([\d\.]+)`},
			output: map[string]string{
				Name:    Chromecast + " Linux",
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)crkey\/([\d\.]+)`},
			output: map[string]string{
				Name:    Chromecast,
				Version: "$1",
			},
		},
		// mobile OSes
		{
			patterns: []string{`(?i)droid ([\w\.]+)\b.+(android[- ]x86|harmonyos)`},
			output: map[string]string{
				Name:    "$2",
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)(ubuntu) ([\w\.]+) like android`},
			output: map[string]string{
				Name:    "$1 Touch",
				Version: "$2",
			},
		},
		{
			patterns: []string{
				`(?i)(android|bada|blackberry|kaios|maemo|meego|openharmony|qnx|rim tablet os|sailfish|series40|symbian|tizen|webos)\w*[-\/; ]?([\d\.]*)`,
			},
			output: map[string]string{
				Name:    "$1",
				Version: "$2",
			},
		},
		{
			patterns: []string{`(?i)\(bb(10);`},
			output: map[string]string{
				Name:    Blackberry,
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)(?:symbian ?os|symbos|s60(?=;)|series ?60)[-\/ ]?([\w\.]*)`},
			output: map[string]string{
				Name:    "Symbian",
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)mozilla\/[\d\.]+ \((?:mobile|tablet|tv|mobile; [\w ]+); rv:.+ gecko\/([\w\.]+)`},
			output: map[string]string{
				Name:    "Firefox Os",
				Version: "$1",
			},
		},
		{
			patterns: []string{
				`(?i)web0s;.+rt(tv)`,
				`(?i)\b(?:hp)?wos(?:browser)?\/([\w\.]+)`,
			},
			output: map[string]string{
				Name:    "webOS",
				Version: "$1",
			},
		},
		{
			patterns: []string{`(?i)watch(?: ?os[,\/]|\d,\d\/)([\d\.]+)`},
			output: map[string]string{
				Name:    "watchOS",
				Version: "$1",
			},
		},
		// Google ChromeOS
		{
			patterns: []string{`(?i)(cros) [\w]+(?:\)| ([\w\.]+)\b)`},
			output: map[string]string{
				Name:    "Chrome Os",
				Version: "$2",
			},
		},
		// Smart TVs
		{
			patterns: []string{`(?i)panasonic;(viera)`},
			output: map[string]string{
				Name:    "Viera",
				Version: "",
			},
		},
		{
			patterns: []string{`(?i)(netrange)mmh`},
			output: map[string]string{
				Name:    "Netrange",
				Version: "",
			},
		},
		{
			patterns: []string{`(?i)(nettv)\/(\d+\.[\w\.]+)`},
			output: map[string]string{
				Name:    "NetTV",
				Version: "$2",
			},
		},
		// Console
		{
			patterns: []string{
				`(?i)(nintendo|playstation) (\w+)`,
				`(?i)(xbox); +xbox ([^\);]+)`,
				`(?i)(pico) .+os([\w\.]+)`,
			},
			output: map[string]string{
				Name:    "$1",
				Version: "$2",
			},
		},
		// Other
		{
			patterns: []string{
				`(?i)\b(joli|palm)\b ?(?:os)?\/?([\w\.]*)`,
				`(?i)(mint)[\/\(\) ]?(\w*)`,
				`(?i)(mageia|vectorlinux)[; ]`,
				`(?i)([kxln]?ubuntu|debian|suse|opensuse|gentoo|arch(?= linux)|slackware|fedora|mandriva|centos|pclinuxos|red ?hat|zenwalk|linpus|raspbian|plan 9|minix|risc os|contiki|deepin|manjaro|elementary os|sabayon|linspire)(?: gnu\/linux)?(?: enterprise)?(?:[- ]linux)?(?:-gnu)?[-\/ ]?(?!chrom|package)([-\w\.]*)`,
				`(?i)(hurd|linux)(?: arm\w*| x86\w*| ?)([\w\.]*)`,
				`(?i)(gnu) ?([\w\.]*)`,
				`(?i)\b([-frentopcghs]{0,5}bsd|dragonfly)[\/ ]?(?!amd|[ix346]{1,2}86)([\w\.]*)`,
				`(?i)(haiku) (\w+)`,
			},
			output: map[string]string{
				Name:    "$1",
				Version: "$2",
			},
		},
		{
			patterns: []string{`(?i)(sunos) ?([\w\.\d]*)`},
			output: map[string]string{
				Name:    "Solaris",
				Version: "$2",
			},
		},
		{
			patterns: []string{
				`(?i)((?:open)?solaris)[-\/ ]?([\w\.]*)`,
				`(?i)(aix) ((\d)(?=\.|\)| )[\w\.])*`,
				`(?i)\b(beos|os\/2|amigaos|morphos|openvms|fuchsia|hp-ux|serenityos)`,
				`(?i)(unix) ?([\w\.]*)`,
			},
			output: map[string]string{
				Name:    "$1",
				Version: "$2",
			},
		},
	},
}

var windowsVersionMap = map[string][]string{
	"ME":      {"4.90"},
	"NT 3.11": {"NT3.51"},
	"NT 4.0":  {"NT4.0"},
	"2000":    {"NT 5.0"},
	"XP":      {"NT 5.1", "NT 5.2"},
	"Vista":   {"NT 6.0"},
	"7":       {"NT 6.1"},
	"8":       {"NT 6.2"},
	"8.1":     {"NT 6.3"},
	"10":      {"NT 6.4", "NT 10.0"},
	"RT":      {"ARM"},
}

var formFactorsMap = map[string][]string{
	Embedded: {"Automotive"},
	Mobile:   {Mobile},
	Tablet:   {Tablet, "EInk"},
	SmartTV:  {"TV"},
	Wearable: {"Watch"},
	XR:       {"VR", "XR"},
	"":       {"Desktop", "Unknown"},
	"*":      {},
}

var clientHintsBrowserMap = map[string][]string{
	"Chrome":          {"Google Chrome"},
	"Edge":            {"Microsoft Edge"},
	"Chrome WebView":  {"Android WebView"},
	"Chrome Headless": {"HeadlessChrome"},
	"Huawei Browser":  {"HuaweiBrowser"},
	"MIUI Browser":    {"Miui Browser"},
	"Opera Mobi":      {"OperaMobile"},
	"Yandex":          {"YaBrowser"},
}

func strMapper(str string, m map[string][]string) string {
	for key, valueList := range m {
		for _, value := range valueList {
			if strings.Contains(strings.ToLower(str), strings.ToLower(value)) {
				return key
			}
		}
	}
	if value, ok := m["*"]; ok {
		return value[0]
	}
	return str
}

func majorize(version string) string {
	return strings.Split(NonNumericOrDotReg.ReplaceAllString(version, ""), ".")[0]
}

var (
	reCache  sync.Map
	re2Cache sync.Map
)

func getCachedRegexp(pattern string) (*regexp.Regexp, error) {
	if re, exists := reCache.Load(pattern); exists {
		return re.(*regexp.Regexp), nil
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	reCache.Store(pattern, re)
	return re, nil
}

func getCachedRegexp2(pattern string) (*regexp2.Regexp, error) {
	if re, exists := re2Cache.Load(pattern); exists {
		return re.(*regexp2.Regexp), nil
	}
	re := regexp2.MustCompile(pattern, 0)
	re2Cache.Store(pattern, re)
	return re, nil
}

func deepCopyMap(original map[string]string) map[string]string {
	c := make(map[string]string, len(original))
	for key, value := range original {
		c[key] = value
	}
	return c
}

func extractIndex(value string) (int, error) {
	if idx := strings.Index(value, "$"); idx != -1 {
		if match := extractNumberReg.FindString(value[idx:]); match != "" {
			return strconv.Atoi(match)
		}
		return 0, fmt.Errorf("no number found after $")
	}
	return 0, nil
}

func applyPattern(ua string, pattern string, output map[string]string) (map[string]string, bool) {
	processMatches := func(matches []string, output map[string]string) map[string]string {
		result := deepCopyMap(output)
		for key, value := range output {
			if idx, err := extractIndex(value); err == nil && idx > 0 && idx < len(matches) {
				result[key] = dollarReplaceReg.ReplaceAllStringFunc(result[key], func(s string) string {
					if i, err := strconv.Atoi(s[1:]); err == nil && i < len(matches) {
						return matches[i]
					}
					return s
				})
			}
		}
		return result
	}

	// Attempt to get the regex from cache
	re, err := getCachedRegexp(pattern)
	if err == nil {
		matches := re.FindStringSubmatch(ua)
		if len(matches) > 0 {
			return processMatches(matches, output), true
		}
	}

	// Fallback to regex2 cache and matching
	re2, err := getCachedRegexp2(pattern)
	if err == nil {
		matches, err := re2.FindStringMatch(ua)
		if err == nil && matches != nil {
			groups := make([]string, len(matches.Groups()))
			for i, group := range matches.Groups() {
				groups[i] = group.String()
			}
			return processMatches(groups, output), true
		}
	}

	return nil, false
}

func parseUA(ua string, regexItems []regexItem) map[string]string {
	for _, regItem := range regexItems {
		for _, pattern := range regItem.patterns {
			result, matched := applyPattern(ua, pattern, regItem.output)
			if matched {
				// Apply mapping functions
				for _, mp := range regItem.mapperItems {
					if mp.field != "" {
						result[mp.field] = mp.fn(result[mp.field])
					}
				}
				return result
			}
		}
	}
	return make(map[string]string)
}

// ClientHints User Agent Client Hints data
type ClientHints struct {
	brands       []IBrand
	fullVerList  []IBrand
	mobile       bool
	model        string
	platform     string
	platformVer  string
	formFactors  []string
	architecture string
	bitness      string
}

func NewClientHints(headers map[string]string) ClientHints {

	stripQuotes := func(str string) string {
		return strings.ReplaceAll(str, "\"", "")
	}

	itemListToArray := func(header string) []IBrand {
		if header == "" {
			return nil
		}
		var arr []IBrand

		tokens := strings.Split(quoteReplaceReg.ReplaceAllString(header, ""), ",")
		for _, token := range tokens {
			if strings.Contains(token, ";") {
				parts := strings.Split(strings.Trim(token, " "), ";v=")
				arr = append(arr, IBrand{
					Name:    parts[0],
					Version: parts[1],
				})
			} else {
				arr = append(arr, IBrand{
					Name: strings.TrimSpace(token),
				})
			}
		}
		return arr
	}

	itemListToArray2 := func(header string) []string {
		if header == "" {
			return nil
		}
		var arr []string

		tokens := strings.Split(quoteReplaceReg.ReplaceAllString(header, ""), ",")
		for _, token := range tokens {
			arr = append(arr, strings.TrimSpace(token))
		}
		return arr
	}

	mobile, _ := regexp.MatchString(`\?1`, headers[CHHeaderMobile])

	return ClientHints{
		brands:       itemListToArray(headers[CHHeader]),
		fullVerList:  itemListToArray(headers[CHHeaderFullVerList]),
		mobile:       mobile,
		model:        stripQuotes(headers[CHHeaderModel]),
		platform:     stripQuotes(headers[CHHeaderPlatform]),
		platformVer:  stripQuotes(headers[CHHeaderPlatformVer]),
		formFactors:  itemListToArray2(headers[CHHeaderFormFactors]),
		architecture: stripQuotes(headers[CHHeaderArch]),
		bitness:      stripQuotes(headers[CHHeaderBitness]),
	}
}

type UAItem struct {
	itemType string // browser, engine, os, device, cpu.
	ua       string
	uaCH     ClientHints
	rgxMap   map[string][]regexItem
	data     map[string]string // ua 解析结果
}

func NewUAItem(itemType string, ua string, rgxMap map[string][]regexItem, uaCH ClientHints) *UAItem {
	return &UAItem{
		itemType: itemType,
		ua:       ua,
		rgxMap:   rgxMap,
		uaCH:     uaCH,
		data:     make(map[string]string),
	}
}

func (item *UAItem) parseCH() *UAItem {
	uaCh := item.uaCH
	rgxMap := item.rgxMap

	switch item.itemType {
	case UABrowser, UAEngine:
		brands := uaCh.fullVerList
		if brands == nil {
			brands = uaCh.brands
		}
		var prevName string
		for _, brand := range brands {
			brandName := brand.Name
			brandVersion := brand.Version
			// Mapping brands for more readable names
			mappedName := strMapper(brandName, clientHintsBrowserMap)

			if mappedName != "" {
				brandName = mappedName
			}

			if item.itemType == UABrowser && !notBrandReg.MatchString(brandName) &&
				(prevName == "" || (strings.Contains(strings.ToLower(prevName), "chrom") && brandName != Chromium)) {
				item.data[Name] = brandName
				item.data[Version] = brandVersion
				item.data[Major] = majorize(brandVersion)
				prevName = brandName
			}

			if item.itemType == UAEngine && brandName == Chromium {
				item.data[Version] = brandVersion
			}
		}
	case UACpu:
		archName := uaCh.architecture
		if archName != "" {
			if uaCh.bitness == "64" {
				archName += "64"
			}
			item.data = parseUA(archName+";", rgxMap[item.itemType])
		}

	case UADevice:
		if uaCh.mobile {
			item.data[Type] = Mobile
		}
		if uaCh.model != "" {
			item.data[Model] = uaCh.model
			if item.data[Type] == "" || item.data[Vendor] == "" {
				reParse := map[string]string{}
				reParse = parseUA("droid 9; "+uaCh.model+")", rgxMap[item.itemType])
				if item.data[Type] == "" && reParse[Type] != "" {
					item.data[Type] = reParse[Type]
				}
				if item.data[Vendor] == "" && reParse[Vendor] != "" {
					item.data[Vendor] = reParse[Vendor]
				}
			}
		}
		if len(item.uaCH.formFactors) > 0 {
			var ff string
			for _, formFactor := range item.uaCH.formFactors {
				ff = strMapper(formFactor, formFactorsMap)
				if ff != "" {
					break
				}
			}
			item.data[Type] = ff
		}
	case UAOS:
		osName := uaCh.platform
		if osName != "" {
			osVersion := uaCh.platformVer
			if osName == Windows {
				v, _ := strconv.Atoi(majorize(osVersion))
				if v >= 13 {
					osVersion = "11"
				} else {
					osVersion = "10"
				}
			}
			item.data[Name] = osName
			item.data[Version] = osVersion
		}

		// Xbox-Specific Detection
		if item.data[Name] == Windows && uaCh.model == Xbox {
			item.data[Name] = Xbox
			item.data[Version] = ""
		}
	}
	return item
}

func (item *UAItem) parseUA() *UAItem {
	if item.itemType != UAResult {
		item.data = parseUA(item.ua, item.rgxMap[item.itemType])
	}
	if item.itemType == UABrowser {
		originVersion := item.data[Version]
		item.data[Version] = originVersion
		item.data[Major] = majorize(item.data[Version])
	}
	return item
}

func (item *UAItem) getData() map[string]string {
	return item.data
}

type UAParser struct {
	ua       string
	httpUACH ClientHints
	withCH   bool
	regexMap map[string][]regexItem
}

func NewUAParser(ua string) *UAParser {
	if len(ua) >= UAMaxLength {
		ua = ""
		slog.Warn("User-Agent string is too long, it will be ignored")
	}
	return &UAParser{
		ua:       ua,
		withCH:   false,
		httpUACH: ClientHints{},
		regexMap: regexMap,
	}
}

func (p *UAParser) WithUA(ua string) *UAParser {
	p.ua = ua
	return p
}

func (p *UAParser) WithHeaders(headers map[string]string) *UAParser {
	if len(headers) > 0 {
		p.withCH = true
	}

	for k, v := range headers {
		headers[strings.ToLower(k)] = v
	}

	p.httpUACH = NewClientHints(headers)
	if ua, ok := headers[UserAgent]; ok && p.ua == "" && len(ua) <= UAMaxLength {
		p.ua = ua
	}
	return p
}

func (p *UAParser) WithExtensions(extensions map[string][]regexItem) *UAParser {
	mergedMap := make(map[string][]regexItem)
	if extensions != nil && len(extensions) > 0 {
		for key, value := range regexMap {
			if extValue, exists := extensions[key]; exists {
				extensions[key] = append(extValue, value...)
			} else {
				extensions[key] = value
			}
		}
		mergedMap = extensions
	} else {
		mergedMap = regexMap
	}
	p.regexMap = mergedMap
	return p
}

func (p *UAParser) getData(itemType string) map[string]string {
	if len(p.ua) < UAMinLength && !p.withCH {
		return make(map[string]string)
	}

	uaItem := NewUAItem(itemType, p.ua, p.regexMap, p.httpUACH)
	var data map[string]string
	if p.withCH {
		data = uaItem.parseUA().parseCH().getData()
	} else {
		data = uaItem.parseUA().getData()
	}
	return data
}

func (p *UAParser) Browser() IBrowser {
	data := p.getData(UABrowser)
	return IBrowser{
		Name:    data[Name],
		Version: data[Version],
		Major:   data[Major],
		Type:    data[Type],
	}
}

func (p *UAParser) CPU() ICpu {
	data := p.getData(UACpu)
	return ICpu{
		Architecture: data[Architecture],
	}
}

func (p *UAParser) Device() IDevice {
	data := p.getData(UADevice)
	return IDevice{
		Type:   data[Type],
		Vendor: data[Vendor],
		Model:  data[Model],
	}
}

func (p *UAParser) Engine() IEngine {
	data := p.getData(UAEngine)
	return IEngine{
		Name:    data[Name],
		Version: data[Version],
	}
}

func (p *UAParser) Os() IOs {
	data := p.getData(UAOS)
	return IOs{
		Name:    data[Name],
		Version: data[Version],
	}
}

func (p *UAParser) Result() IResult {
	return IResult{
		UA:      p.ua,
		Browser: p.Browser(),
		Engine:  p.Engine(),
		Os:      p.Os(),
		Device:  p.Device(),
		Cpu:     p.CPU(),
	}
}
