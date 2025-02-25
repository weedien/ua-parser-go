package uaparser

const (
	Model               = "model"
	Mobile              = "mobile"
	Brands              = "brands"
	FormFactors         = "formFactors"
	FullVerList         = "fullVersionList"
	Platform            = "platform"
	PlatformVer         = "platformVersion"
	Bitness             = "bitness"
	CHHeader            = "sec-ch-ua"
	CHHeaderFullVerList = CHHeader + "-full-Version-list"
	CHHeaderArch        = CHHeader + "-arch"
	CHHeaderBitness     = CHHeader + "-" + Bitness
	CHHeaderFormFactors = CHHeader + "-form-factors"
	CHHeaderMobile      = CHHeader + "-" + Mobile
	CHHeaderModel       = CHHeader + "-" + Model
	CHHeaderPlatform    = CHHeader + "-" + Platform
	CHHeaderPlatformVer = CHHeaderPlatform + "-Version"
)

// for regex-base impl
const (
	Name         = "Name"
	Version      = "Version"
	Type         = "type"
	Architecture = "architecture"
	Vendor       = "vendor"
	Console      = "console"
	Desktop      = "desktop"
	Bot          = "bot"

	Windows      = "Windows"
	WindowsPhone = "Windows Phone"
	Android      = "Android"
	MacOS        = "macOS"
	IOS          = "iOS"
	Linux        = "Linux"
	ChromeOS     = "ChromeOS"
	Harmony      = "Harmony"

	Opera            = "Opera"
	OperaMini        = "Opera Mini"
	OperaTouch       = "Opera Touch"
	Chrome           = "Chrome"
	HeadlessChrome   = "Headless Chrome"
	Firefox          = "Firefox"
	InternetExplorer = "Internet Explorer"
	Safari           = "Safari"
	Edge             = "Edge"

	SamsungBrowser = "Samsung Browser"

	GoogleAdsBot        = "Google Ads Bot"
	Googlebot           = "Googlebot"
	Twitterbot          = "Twitterbot"
	FacebookExternalHit = "facebookexternalhit"
	Applebot            = "Applebot"

	FacebookApp  = "Facebook App"
	InstagramApp = "Instagram App"
	TiktokApp    = "TikTok App"

	Tablet = "Tablet"

	Major      = "major"
	SmartTV    = "smarttv"
	Wearable   = "wearable"
	XR         = "xr"
	Embedded   = "embedded"
	UABrowser  = "browser"
	UACpu      = "cpu"
	UADevice   = "device"
	UAEngine   = "engine"
	UAOS       = "os"
	UAResult   = "result"
	Amazon     = "Amazon"
	Apple      = "Apple"
	ASUS       = "ASUS"
	Blackberry = "BlackBerry"
	Google     = "Google"
	Huawei     = "Huawei"
	Lenovo     = "Lenovo"
	Honor      = "Honor"
	LG         = "LG"
	Microsoft  = "Microsoft"
	Motorola   = "Motorola"
	OPPO       = "OPPO"
	Samsung    = "Samsung"
	Sharp      = "Sharp"
	Sony       = "Sony"
	Xiaomi     = "Xiaomi"
	Zebra      = "Zebra"
	Chromecast = "Chromecast"

	Facebook = "Facebook"
)
