package uaparser

import (
	"regexp"
	"strconv"
	"strings"
)

type UserAgent struct {
	UA        string   `json:"ua"`
	Browser   IBrowser `json:"browser"`
	Engine    IEngine  `json:"engine"`
	Device    IDevice  `json:"device"`
	OS        IOs      `json:"os"`
	CPU       ICpu     `json:"cpu"`
	undecided bool
}

func (u *UserAgent) WithClientHints(ch map[string]string) *UserAgent {
	uaCh := NewClientHints(ch)

	// Detect the browser and version
	var brands []IBrand
	if uaCh.fullVerList != nil {
		brands = uaCh.fullVerList
	} else {
		brands = uaCh.brands
	}
	var prevName string
	if brands != nil {
		for _, brand := range brands {
			brandName := strip("(Google|Microsoft) ", brand.Name)
			brandVersion := brand.Version
			if !regexp.MustCompile("not.a.brand").MatchString(strings.ToLower(brandName)) &&
				(prevName == "" ||
					(strings.Contains(strings.ToLower(prevName), "chrom") &&
						!strings.Contains(strings.ToLower(brandName), "chromi"))) {
				u.Browser.Name = brandName
				u.Browser.Version = brandVersion
				u.Browser.Major = majorize(brandVersion)
				prevName = brandName
			}
		}
	}

	// Detect the device
	if uaCh.mobile {
		u.Device.Type = Mobile
	}
	if uaCh.model != "" {
		u.Device.Model = uaCh.model
	}
	if uaCh.model == "Xbox" {
		u.Device.Type = Console
		u.Device.Vendor = Microsoft
	}

	// Detect the os
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
		u.OS.Name = osName
		u.OS.Version = osVersion
	}
	if u.OS.Name == Windows && uaCh.model == "Xbox" {
		u.OS.Name = "Xbox"
		u.OS.Version = ""
	}
	return u
}

// Parse the user agent string and print the sections.
func Parse(ua string) *UserAgent {
	u := &UserAgent{}

	var sections []section

	u.UA = ua
	for index, limit := 0, len(ua); index < limit; {
		s := parseSection(ua, &index)
		sections = append(sections, s)
	}

	if len(sections) > 0 {
		u.detectBrowser(sections)
		u.detectEngine()
		u.detectDevice(sections[0])
		u.detectOS(sections[0])
		u.detectCPU(sections[0])
	}

	return u
}

type section struct {
	name     string
	version  string
	comments []string
	comment  string
}

func (s section) String() string {
	return s.name + "/" + s.version + " (" + strings.Join(s.comments, "; ") + ")"
}

// Read from the given string until the given delimiter or the
// end of the string have been reached.
//
// The first argument is the user agent string being parsed. The second
// argument is a reference pointing to the current index of the user agent
// string. The delimiter argument specifies which character is the delimiter
// and the cat argument determines whether nested '(' should be ignored or not.
//
// Returns a string containing what has been read.
func readUntil(ua string, index *int, delimiter byte, cat bool) string {
	var builder strings.Builder
	catalan := 0

	for i := *index; i < len(ua); i++ {
		if ua[i] == delimiter {
			if catalan == 0 {
				*index = i + 1
				return builder.String()
			}
			catalan--
		} else if cat && ua[i] == '(' {
			catalan++
		}
		builder.WriteByte(ua[i])
	}
	*index = len(ua)
	return builder.String()
}

// Parse the given product, that is, just a name or a string
// formatted as name/version.
//
// It returns two strings. The first string is the name of the product and the
// second string contains the version of the product.
func parseProduct(product string) (string, string) {
	prod := strings.SplitN(product, "/", 2)
	if len(prod) == 2 {
		return prod[0], prod[1]
	}
	return product, ""
}

// Parse a section. A section is typically formatted as follows
// "name/version (comments)". Both, the comments and the version are optional.
//
// The first argument is the user agent string being parsed. The second
// argument is a reference pointing to the current index of the user agent
// string.
//
// Returns a section containing the information that we could extract
// from the last parsed section.
func parseSection(ua string, index *int) (s section) {
	// Check for empty products
	if *index < len(ua) && ua[*index] != '(' && ua[*index] != '[' {
		buffer := readUntil(ua, index, ' ', false)
		s.name, s.version = parseProduct(buffer)
	}

	if *index < len(ua) && ua[*index] == '(' {
		*index++
		buffer := readUntil(ua, index, ')', true)
		s.comment = buffer
		s.comments = strings.Split(buffer, "; ")
		*index++
	}

	// Discards any trailing data within square brackets
	if *index < len(ua) && ua[*index] == '[' {
		*index++
		_ = readUntil(ua, index, ']', true)
		*index++
	}
	return s
}
