package uaparser

import "strings"

// Normalize the name of the operating system. By now, this just
// affects to Windows NT.
//
// Returns a string containing the normalized name for the Operating System.
func normalizeOS(name string) string {
	sp := strings.SplitN(name, " ", 3)
	if len(sp) != 3 || sp[1] != "NT" {
		return name
	}

	switch sp[2] {
	case "5.0":
		return "Windows 2000"
	case "5.01":
		return "Windows 2000, Service Pack 1 (SP1)"
	case "5.1":
		return "Windows XP"
	case "5.2":
		return "Windows XP x64 Edition"
	case "6.0":
		return "Windows Vista"
	case "6.1":
		return "Windows 7"
	case "6.2":
		return "Windows 8"
	case "6.3":
		return "Windows 8.1"
	case "10.0":
		return "Windows 10"
	}
	return name
}

// Guess the Os, the localization and if this is a mobile device for a
// Webkit-powered browser.
//
// The first argument p is a reference to the current UserAgent and the second
// argument is a slice of strings containing the comments.
func webkit(p *UserAgent, comment []string) {
	if p.OS.Platform == "webOS" {
		p.Browser.Name = p.OS.Platform
		p.OS.Name = "Palm"
		p.Device.Type = Mobile
	} else if p.OS.Platform == "Symbian" {
		p.Device.Type = Mobile
		p.Browser.Name = p.OS.Platform
		p.OS.Name = comment[0]
	} else if p.OS.Platform == "Linux" {
		p.Device.Type = Mobile
		if p.Browser.Name == "Safari" {
			p.Browser.Name = "Android"
		}
		if len(comment) > 1 {
			if comment[1] == "U" || comment[1] == "arm_64" {
				if len(comment) > 2 {
					p.OS.Name = comment[2]
				} else {
					//p.mobile = false
					p.OS.Name = comment[0]
				}
			} else {
				p.OS.Name = comment[1]
			}
		}
		//if len(comments) > 3 {
		//	p.localization = comments[3]
		//} else
		if len(comment) == 3 {
			_ = p.googleOrBingBot()
		}
	} else if len(comment) > 0 {
		//if len(comments) > 3 {
		//	p.localization = comments[3]
		//}
		if strings.HasPrefix(comment[0], "Windows NT") {
			p.OS.Name = normalizeOS(comment[0])
		} else if len(comment) < 2 {
			//p.localization = comments[0]
		} else if len(comment) < 3 {
			if !p.googleOrBingBot() && !p.iMessagePreview() {
				p.OS.Name = normalizeOS(comment[1])
			}
		} else {
			p.OS.Name = normalizeOS(comment[2])
		}
		if p.OS.Platform == "BlackBerry" {
			p.Browser.Name = p.OS.Platform
			if p.OS.Name == "Touch" {
				p.OS.Name = p.OS.Platform
			}
		}
	}

	// Special case for Firefox on iPad, where the platform is advertised as Macintosh instead of iPad
	if p.OS.Platform == "Macintosh" && p.Engine.Name == "AppleWebKit" && p.Browser.Name == "Firefox" {
		p.OS.Platform = "iPad"
		p.Device.Type = Mobile
	}
}

// Guess the Os, the localization and if this is a mobile device
// for a Gecko-powered browser.
//
// The first argument p is a reference to the current UserAgent and the second
// argument is a slice of strings containing the comments.
func gecko(p *UserAgent, comment []string) {
	if len(comment) > 1 {
		if comment[1] == "U" || comment[1] == "arm_64" {
			if len(comment) > 2 {
				p.OS.Name = normalizeOS(comment[2])
			} else {
				p.OS.Name = normalizeOS(comment[1])
			}
		} else {
			if strings.Contains(p.OS.Platform, "Android") {
				p.Device.Type = Mobile
				p.OS.Platform, p.OS.Name = normalizeOS(comment[1]), p.OS.Platform
			} else if comment[0] == "mobile" || comment[0] == "Tablet" {
				p.Device.Type = Mobile
				p.OS.Name = "FirefoxOS"
			} else {
				if p.OS.Name == "" {
					p.OS.Name = normalizeOS(comment[1])
				}
			}
		}
		// Only parse 4th comments as localization if it doesn't start with rv:.
		// For example Firefox on Ubuntu contains "rv:XX.X" in this field.
		//if len(comments) > 3 && !strings.HasPrefix(comments[3], "rv:") {
		//	p.localization = comments[3]
		//}
	}
}

// Guess the Os, the localization and if this is a mobile device
// for Internet Explorer.
//
// The first argument p is a reference to the current UserAgent and the second
// argument is a slice of strings containing the comments.
func trident(p *UserAgent, comment []string) {
	// Internet Explorer only runs on Windows.
	p.OS.Platform = "Windows"

	// The Os can be set before to handle a new case in IE11.
	if p.OS.Name == "" {
		if len(comment) > 2 {
			p.OS.Name = normalizeOS(comment[2])
		} else {
			p.OS.Name = "Windows NT 4.0"
		}
	}

	// Last but not least, let's detect if it comes from a mobile device.
	for _, v := range comment {
		if strings.HasPrefix(v, "IEMobile") {
			p.Device.Type = Mobile
			return
		}
	}
}

// Guess the Os, the localization and if this is a mobile device
// for Opera.
//
// The first argument p is a reference to the current UserAgent and the second
// argument is a slice of strings containing the comments.
func opera(p *UserAgent, comment []string) {
	slen := len(comment)

	if strings.HasPrefix(comment[0], "Windows") {
		p.OS.Platform = "Windows"
		p.OS.Name = normalizeOS(comment[0])
	} else {
		if strings.HasPrefix(comment[0], "Android") {
			p.Device.Type = Mobile
		}
		p.OS.Platform = comment[0]
		if slen > 1 {
			p.OS.Name = comment[1]
		} else {
			p.OS.Name = comment[0]
		}
	}
}

// Guess the Os. Android browsers send Dalvik as the user agent in the
// request header.
//
// The first argument p is a reference to the current UserAgent and the second
// argument is a slice of strings containing the comments.
func dalvik(p *UserAgent, comment []string) {
	slen := len(comment)

	if strings.HasPrefix(comment[0], "Linux") {
		p.OS.Platform = comment[0]
		if slen > 2 {
			p.OS.Name = comment[2]
		}
		p.Device.Type = Mobile
	}
}

// Given the comments of the first section of the UserAgent string,
// get the platform.
func getPlatform(comment []string) string {
	if len(comment) > 0 {
		if comment[0] != "compatible" {
			if strings.HasPrefix(comment[0], "Windows") {
				return "Windows"
			} else if strings.HasPrefix(comment[0], "Symbian") {
				return "Symbian"
			} else if strings.HasPrefix(comment[0], "webOS") {
				return "webOS"
			} else if comment[0] == "BB10" {
				return "BlackBerry"
			}
			return comment[0]
		}
	}
	return ""
}

func (u *UserAgent) detectOS(s section) {
	if s.name == "Mozilla" {
		u.OS.Platform = getPlatform(s.comments)
		if u.OS.Platform == "Windows" && len(s.comments) > 0 {
			u.OS.Name = normalizeOS(s.comments[0])
		}

		switch u.Engine.Name {
		case "":
			u.undecided = true
		case "Gecko", "AppleWebKit", "Blink":
			webkitOrGecko(u, s.comments)
		case "Trident":
			trident(u, s.comments)
		}
	} else if s.name == "Opera" {
		if len(s.comments) > 0 {
			opera(u, s.comments)
		}
	} else if s.name == "Dalvik" {
		if len(s.comments) > 0 {
			dalvik(u, s.comments)
		}
	} else if s.name == "okhttp" {
		u.Device.Type = Mobile
		u.Browser.Name = "OkHttp"
		u.Browser.Version = s.version
	} else {
		u.undecided = true
	}

	// Special case for iPhone weirdness
	os := strings.Replace(u.OS.Name, "like Mac Os X", "", 1)
	os = strings.Replace(os, "Cpu", "", 1)
	os = strings.Trim(os, " ")

	osSplit := strings.Split(os, " ")

	if os == "Windows XP x64 Edition" {
		osSplit = osSplit[:len(osSplit)-2]
	}

	name, version := osName(osSplit)

	if strings.Contains(name, "/") {
		s := strings.Split(name, "/")
		name = s[0]
		version = s[1]
	}

	version = strings.Replace(version, "_", ".", -1)

	u.OS.Name = name
	u.OS.Version = version
}

func webkitOrGecko(p *UserAgent, comment []string) {
	if p.OS.Platform == "webOS" {
		p.Browser.Name = p.OS.Platform
		p.OS.Name = "Palm"
		p.Device.Type = Mobile
	} else if p.OS.Platform == "Symbian" {
		p.Device.Type = Mobile
		p.Browser.Name = p.OS.Platform
		p.OS.Name = comment[0]
	} else if p.OS.Platform == "Linux" {
		p.Device.Type = Mobile
		if p.Browser.Name == "Safari" {
			p.Browser.Name = "Android"
		}
		if len(comment) > 1 {
			if comment[1] == "U" || comment[1] == "arm_64" {
				if len(comment) > 2 {
					p.OS.Name = comment[2]
				} else {
					p.OS.Name = comment[0]
				}
			} else {
				p.OS.Name = comment[1]
			}
		}
		if len(comment) == 3 {
			_ = p.googleOrBingBot()
		}
	} else if len(comment) > 0 {
		if strings.HasPrefix(comment[0], "Windows NT") {
			p.OS.Name = normalizeOS(comment[0])
		} else if len(comment) < 2 {
			// Do nothing
		} else if len(comment) < 3 {
			if !p.googleOrBingBot() && !p.iMessagePreview() {
				p.OS.Name = normalizeOS(comment[1])
			}
		} else {
			p.OS.Name = normalizeOS(comment[2])
		}
		if p.OS.Platform == "BlackBerry" {
			p.Browser.Name = p.OS.Platform
			if p.OS.Name == "Touch" {
				p.OS.Name = p.OS.Platform
			}
		}
	}

	if p.OS.Platform == "Macintosh" && p.Engine.Name == "AppleWebKit" && p.Browser.Name == "Firefox" {
		p.OS.Platform = "iPad"
		p.Device.Type = Mobile
	}

	// Add support for HarmonyOS
	if strings.Contains(p.OS.Name, "HarmonyOS") {
		p.OS.Name = "HarmonyOS"
		p.Device.Type = Mobile
	}
}

//func (u *UserAgent) detectOS(s section) {
//	if s.name == "Mozilla" {
//		// Get the platform here. Be aware that IE11 provides a new format
//		// that is not backwards-compatible with previous versions of IE.
//		u.Os.Platform = getPlatform(s.comments)
//		if u.Os.Platform == "Windows" && len(s.comments) > 0 {
//			u.Os.Name = normalizeOS(s.comments[0])
//		}
//
//		// And finally get the Os depending on the engine.
//		switch u.Engine.Name {
//		case "":
//			u.undecided = true
//		case "Gecko":
//			gecko(u, s.comments)
//		case "AppleWebKit":
//			webkit(u, s.comments)
//		case "Blink":
//			webkit(u, s.comments)
//		case "Trident":
//			trident(u, s.comments)
//		}
//	} else if s.name == "Opera" {
//		if len(s.comments) > 0 {
//			opera(u, s.comments)
//		}
//	} else if s.name == "Dalvik" {
//		if len(s.comments) > 0 {
//			dalvik(u, s.comments)
//		}
//	} else if s.name == "okhttp" {
//		u.Device.Type = Mobile
//		u.Browser.Name = "OkHttp"
//		u.Browser.Version = s.version
//	} else {
//		// Check whether this is a bot or just a weird browser.
//		u.undecided = true
//	}
//
//	// Special case for iPhone weirdness
//	os := strings.Replace(u.Os.Name, "like Mac Os X", "", 1)
//	os = strings.Replace(os, "Cpu", "", 1)
//	os = strings.Trim(os, " ")
//
//	osSplit := strings.Split(os, " ")
//
//	// Special case for x64 edition of Windows
//	if os == "Windows XP x64 Edition" {
//		osSplit = osSplit[:len(osSplit)-2]
//	}
//
//	name, version := osName(osSplit)
//
//	// Special case for names that contain a forward slash version separator.
//	if strings.Contains(name, "/") {
//		s := strings.Split(name, "/")
//		name = s[0]
//		version = s[1]
//	}
//
//	// Special case for versions that use underscores
//	version = strings.Replace(version, "_", ".", -1)
//
//	u.Os.Name = name
//	u.Os.Version = version
//}

// Return Os name and version from a slice of strings created from the full name of the Os.
func osName(osSplit []string) (name, version string) {
	if len(osSplit) == 1 {
		name = osSplit[0]
		version = ""
	} else {
		// Assume version is stored in the last part of the array.
		nameSplit := osSplit[:len(osSplit)-1]
		version = osSplit[len(osSplit)-1]

		// Nicer looking Mac Os X
		if len(nameSplit) >= 2 && nameSplit[0] == "Intel" && nameSplit[1] == "Mac" {
			nameSplit = nameSplit[1:]
		}
		name = strings.Join(nameSplit, " ")

		if strings.Contains(version, "x86") || strings.Contains(version, "i686") {
			// x86_64 and i868 are not Linux versions but architectures
			version = ""
		} else if version == "X" && name == "Mac Os" {
			// X is not a version for Mac Os.
			name = name + " " + version
			version = ""
		}
	}
	return name, version
}
