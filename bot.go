package uaparser

import (
	"regexp"
	"strings"
)

var botFromSiteRegexp = regexp.MustCompile(`https?://.+\.\w+`)

// Get the name of the bot from the website that may be in the given comments. If
// there is no website in the comments, then an empty string is returned.
func getFromSite(comment []string) string {
	if len(comment) == 0 {
		return ""
	}

	// Where we should check the website.
	idx := 2
	if len(comment) < 3 {
		idx = 0
	} else if len(comment) == 4 {
		idx = 3
	}

	// Pick the site.
	results := botFromSiteRegexp.FindStringSubmatch(comment[idx])
	if len(results) == 1 {
		// If it's a simple comments, just return the name of the site.
		if idx == 0 {
			return results[0]
		}

		// This is a large comments, usually the name will be in the previous
		// field of the comments.
		return strings.TrimSpace(comment[idx-1])
	}
	return ""
}

// Returns true if the info that we currently have corresponds to the Google
// or Bing mobile bot. This function also modifies some attributes in the receiver
// accordingly.
func (u *UserAgent) googleOrBingBot() bool {
	// This is a hackish way to detect
	// Google's mobile bot (Googlebot, AdsBot-Google-mobile, etc.)
	// (See https://support.google.com/webmasters/answer/1061943)
	// and Bing's mobile bot
	// (See https://www.bing.com/webmaster/help/which-crawlers-does-bing-use-8c184ec0)
	if strings.Contains(u.UA, "Google") || strings.Contains(u.UA, "bingbot") {
		//u.platform = ""
		u.undecided = true
	}
	return u.undecided
}

// Returns true if we think that it is iMessage-Preview. This function also
// modifies some attributes in the receiver accordingly.
func (u *UserAgent) iMessagePreview() bool {
	// iMessage-Preview doesn't advertise itself. We have a to rely on a hack
	// to detect it: it impersonates both facebook and twitter bots.
	// See https://medium.com/@siggi/apples-imessage-impersonates-twitter-facebook-bots-when-scraping-cef85b2cbb7d
	if !strings.Contains(u.UA, "facebookexternalhit") {
		return false
	}
	if !strings.Contains(u.UA, "Twitterbot") {
		return false
	}
	//u.bot = true
	u.Browser.Name = "iMessage-Preview"
	u.Engine.Name = ""
	u.Engine.Version = ""
	// We don't set the mobile flag because iMessage can be on iOS (mobile) or macOS (not mobile).
	return true
}

// Set the attributes of the receiver as given by the parameters. All the other
// parameters are set to empty.
func (u *UserAgent) setSimple(name, version string, bot bool) {
	if bot {
		u.Device.Type = Bot
	}
	u.Browser.Name = name
	u.Browser.Version = version
	u.Engine.Name = ""
	u.Engine.Version = ""
	//u.Os = ""
	//u.localization = ""
}

// Fix some values for some weird browsers.
func (u *UserAgent) fixOther(sections []section) {
	if len(sections) > 0 {
		u.Browser.Name = sections[0].name
		u.Browser.Version = sections[0].version
		//u.mozilla = ""
	}
}

var botRegex = regexp.MustCompile("(?i)(bot|crawler|sp([iy])der|search|worm|fetch|nutch)")

// Check if we're dealing with a bot or with some weird browser. If that is the
// case, the receiver will be modified accordingly.
func (u *UserAgent) checkBot(sections []section) {
	// If there's only one element, and it's doesn't have the Mozilla string,
	// check whether this is a bot or not.
	if len(sections) == 1 && sections[0].name != "Mozilla" {
		//u.mozilla = ""

		// Check whether the name has some suspicious "bot" or "crawler" in his name.
		if botRegex.Match([]byte(sections[0].name)) {
			u.setSimple(sections[0].name, "", true)
			return
		}

		// Tough luck, let's try to see if it has a website in his comments.
		if name := getFromSite(sections[0].comments); name != "" {
			// First of all, this is a bot. Moreover, since it doesn't have the
			// Mozilla string, we can assume that the name and the version are
			// the ones from the first section.
			u.setSimple(sections[0].name, sections[0].version, true)
			return
		}

		// At this point we are sure that this is not a bot, but some weirdo.
		u.setSimple(sections[0].name, sections[0].version, false)
	} else {
		// Let's iterate over the available comments and check for a website.
		for _, v := range sections {
			if name := getFromSite(v.comments); name != "" {
				// Ok, we've got a bot name.
				results := strings.SplitN(name, "/", 2)
				version := ""
				if len(results) == 2 {
					version = results[1]
				}
				u.setSimple(results[0], version, true)
				return
			}
		}

		// We will assume that this is some other weird browser.
		u.fixOther(sections)
	}
}
