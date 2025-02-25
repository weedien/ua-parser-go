package uaparser

import (
	"strings"
)

func (u *UserAgent) detectEngine() {
	ua := u.UA

	if strings.Contains(strings.ToLower(ua), "windows") && strings.Contains(strings.ToLower(ua), "edge/") {
		u.Engine.Name = "EdgeHTML"
		u.Engine.Version = extractVersion(ua, "edge/")
		return
	}

	if strings.Contains(strings.ToLower(ua), "arkweb/") {
		u.Engine.Name = "ArkWeb"
		u.Engine.Version = extractVersion(ua, "arkweb/")
		return
	}

	if strings.Contains(strings.ToLower(ua), "webkit/537.36") &&
		strings.Contains(strings.ToLower(ua), "chrome/") &&
		!strings.Contains(strings.ToLower(ua), "chrome/27") {
		u.Engine.Name = "Blink"
		u.Engine.Version = extractVersion(ua, "chrome/")
		return
	}

	commonEngines := []string{"presto", "webkit", "trident", "netfront", "netsurf", "amaya", "lynx", "w3m", "goanna", "servo", "ekioh(flow)", "khtml", "tasman", "links", "icab", "libweb"}
	for _, engine := range commonEngines {
		if strings.Contains(strings.ToLower(ua), engine) {
			u.Engine.Name = engine
			u.Engine.Version = extractVersion(ua, engine)
			return
		}
	}

	if strings.Contains(strings.ToLower(ua), "rv:") && strings.Contains(strings.ToLower(ua), "gecko") {
		u.Engine.Name = "Gecko"
		u.Engine.Version = extractVersion(ua, "rv:")
		return
	}
}

func extractVersion(ua, key string) string {
	start := strings.Index(strings.ToLower(ua), key) + len(key)
	end := start
	for end < len(ua) && (ua[end] == '.' || ua[end] == '_' || ua[end] == '-' || (ua[end] >= '0' && ua[end] <= '9') || (ua[end] >= 'a' && ua[end] <= 'z') || (ua[end] >= 'A' && ua[end] <= 'Z')) {
		end++
	}
	return ua[start:end]
}

//var (
//	edgeHTMLEngineRegex = regexp.MustCompile(`(?i)windows.+ edge/([\w.]+)`)
//	arkWebEngineRegex   = regexp.MustCompile(`(?i)arkweb/([\w.]+)`)
//	blinkEngineRegex    = regexp.MustCompile(`(?i)webkit/537\.36.+chrome/(?!27)([\w.]+)`)
//
//	commonEngineRegex = []*regexp.Regexp{
//		regexp.MustCompile(`(?i)(presto)/([\w.]+)`),
//		regexp.MustCompile(`(?i)(webkit|trident|netfront|netsurf|amaya|lynx|w3m|goanna|servo)/([\w.]+)`),
//		regexp.MustCompile(`(?i)ekioh(flow)/([\w.]+)`),
//		regexp.MustCompile(`(?i)(khtml|tasman|links)[/ ]\(?([\w.]+)`),
//		regexp.MustCompile(`(?i)(icab)[/ ]([23]\.[\d.]+)`),
//		regexp.MustCompile(`(?i)\b(libweb)`),
//	}
//
//	geckoEngineRegex = regexp.MustCompile(`(?i)rv:([\w.]{1,9})\b.+(gecko)`)
//)
//
//func (u *UserAgent) detectEngine() {
//	ua := u.UA
//	if edgeHTMLEngineRegex.MatchString(ua) {
//		m := edgeHTMLEngineRegex.FindStringSubmatch(ua)
//		u.Engine.Name = "EdgeHTML"
//		u.Engine.Version = m[1]
//		return
//	}
//
//	if arkWebEngineRegex.MatchString(ua) {
//		m := arkWebEngineRegex.FindStringSubmatch(ua)
//		u.Engine.Name = "ArkWeb"
//		u.Engine.Version = m[1]
//		return
//	}
//
//	if blinkEngineRegex.MatchString(ua) {
//		m := blinkEngineRegex.FindStringSubmatch(ua)
//		u.Engine.Name = "Blink"
//		u.Engine.Version = m[1]
//		return
//	}
//
//	for _, r := range commonEngineRegex {
//		if r.MatchString(ua) {
//			m := r.FindStringSubmatch(ua)
//			u.Engine.Name = m[1]
//			u.Engine.Version = m[2]
//			return
//		}
//	}
//
//	if geckoEngineRegex.MatchString(ua) {
//		m := geckoEngineRegex.FindStringSubmatch(ua)
//		u.Engine.Name = "Gecko"
//		u.Engine.Version = m[1]
//		return
//	}
//}
