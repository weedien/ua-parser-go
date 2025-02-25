package uaparser

import "strings"

func (u *UserAgent) detectDevice(s section) {
	if u.Device.Type != Mobile {
		u.Device.Type = Desktop
		return
	}
	if u.OS.Platform == "iPhone" || u.OS.Platform == "iPad" {
		u.Device.Model = u.OS.Platform
		return
	}
	// Android model
	if s.name == "Mozilla" && u.OS.Platform == "Linux" && len(s.comments) > 2 {
		mostAndroidModel := s.comments[2]
		if strings.Contains(mostAndroidModel, "Android") || strings.Contains(mostAndroidModel, "Linux") {
			mostAndroidModel = s.comments[len(s.comments)-1]
		}
		tmp := strings.Split(mostAndroidModel, "Build")
		if len(tmp) > 0 {
			u.Device.Model = strings.Trim(tmp[0], " ")
			return
		}
	}
	// traverse all item
	for _, v := range s.comments {
		if strings.Contains(v, "Build") {
			tmp := strings.Split(v, "Build")
			u.Device.Model = strings.Trim(tmp[0], " ")
		}
	}
}
