package uaparser

import (
	"regexp"
	"strings"
)

var ie11Regexp = regexp.MustCompile("^rv:(.+)$")

func (u *UserAgent) detectBrowser(sections []section) {
	if u.Browser.Name != "" && u.Browser.Version != "" {
		return
	}

	slen := len(sections)

	if sections[0].name == "Opera" {
		u.Browser.Name = "Opera"
		u.Browser.Version = sections[0].version
		//u.Browser.Engine = "Presto"
		//if slen > 1 {
		//	u.Browser.EngineVersion = sections[1].version
		//}
	} else if sections[0].name == "Dalvik" {
		// When Dalvik VM is in use, there is no Browser info attached to ua.
		// Although Browser is still a Mozilla/5.0 compatible.
		//u.mozilla = "5.0"
	} else if slen > 1 {
		engine := sections[1]
		//u.Browser.Engine = engine.name
		//u.Browser.EngineVersion = engine.version
		if slen > 2 {
			sectionIndex := 2
			// The version after the engine comments is empty on e.g. Ubuntu
			// platforms so if this is the case, let's use the next in line.
			if sections[2].version == "" && slen > 3 {
				sectionIndex = 3
			}
			u.Browser.Version = sections[sectionIndex].version
			if engine.name == "AppleWebKit" || engine.name == "Blink" {
				for _, comment := range engine.comments {
					if len(comment) > 5 &&
						(strings.HasPrefix(comment, "Googlebot") || strings.HasPrefix(comment, "bingbot")) {
						u.undecided = true
						break
					}
				}
				switch sections[slen-1].name {
				case "Edge":
					u.Browser.Name = "Edge"
					u.Browser.Version = sections[slen-1].version
				case "EdgA":
					u.Browser.Name = "Edge"
					u.Browser.Version = sections[slen-1].version
					//u.Browser.Engine = "EdgeHTML"
					//u.Browser.EngineVersion = ""
				case "Edg":
					if !u.undecided {
						u.Browser.Name = "Edge"
						u.Browser.Version = sections[slen-1].version
						//u.Browser.Engine = "AppleWebKit"
						//u.Browser.EngineVersion = sections[slen-2].version
					}
				case "OPR":
					u.Browser.Name = "Opera"
					u.Browser.Version = sections[slen-1].version
				case "mobile":
					u.Browser.Name = "mobile App"
					u.Browser.Version = ""
				default:
					switch sections[slen-3].name {
					case "YaBrowser":
						u.Browser.Name = "Yandex"
						u.Browser.Version = sections[slen-3].version
					case "coc_coc_Browser":
						u.Browser.Name = "Coc Coc"
						u.Browser.Version = sections[slen-3].version
					default:
						switch sections[slen-2].name {
						case "Electron":
							u.Browser.Name = "Electron"
							u.Browser.Version = sections[slen-2].version
						case "DuckDuckGo":
							u.Browser.Name = "DuckDuckGo"
							u.Browser.Version = sections[slen-2].version
						case "PhantomJS":
							u.Browser.Name = "PhantomJS"
							u.Browser.Version = sections[slen-2].version
						default:
							switch sections[sectionIndex].name {
							case "Chrome", "CriOS":
								u.Browser.Name = "Chrome"
							case "HeadlessChrome":
								u.Browser.Name = "Headless Chrome"
							case "Chromium":
								u.Browser.Name = "Chromium"
							case "GSA":
								u.Browser.Name = "Google App"
							case "FxiOS":
								u.Browser.Name = "Firefox"
							default:
								u.Browser.Name = "Safari"
							}
						}
					}
					// It's possible the google-bot emulates these now
					for _, comment := range engine.comments {
						if len(comment) > 5 &&
							(strings.HasPrefix(comment, "Googlebot") || strings.HasPrefix(comment, "bingbot")) {
							u.undecided = true
							break
						}
					}
				}
			} else if engine.name == "Gecko" {
				name := sections[2].name
				if name == "MRA" && slen > 4 {
					name = sections[4].name
					u.Browser.Version = sections[4].version
				}
				u.Browser.Name = name
			} else if engine.name == "like" && sections[2].name == "Gecko" {
				// This is the new user agent from Internet Explorer 11.
				//u.Browser.Engine = "Trident"
				u.Browser.Name = "Internet Explorer"
				for _, c := range sections[0].comments {
					version := ie11Regexp.FindStringSubmatch(c)
					if len(version) > 0 {
						u.Browser.Version = version[1]
						return
					}
				}
				u.Browser.Version = ""
			}
		}
	} else if slen == 1 && len(sections[0].comments) > 1 {
		comment := sections[0].comments
		if comment[0] == "compatible" && strings.HasPrefix(comment[1], "MSIE") {
			//u.Browser.Engine = "Trident"
			u.Browser.Name = "Internet Explorer"
			// The MSIE version may be reported as the compatibility version.
			// For IE 8 through 10, the Trident token is more accurate.
			// http://msdn.microsoft.com/en-us/library/ie/ms537503(v=vs.85).aspx#VerToken
			for _, v := range comment {
				if strings.HasPrefix(v, "Trident/") {
					switch v[8:] {
					case "4.0":
						u.Browser.Version = "8.0"
					case "5.0":
						u.Browser.Version = "9.0"
					case "6.0":
						u.Browser.Version = "10.0"
					}
					break
				}
			}
			// If the Trident token is not provided, fall back to MSIE token.
			if u.Browser.Version == "" {
				u.Browser.Version = strings.TrimSpace(comment[1][4:])
			}
		}
	}
}
