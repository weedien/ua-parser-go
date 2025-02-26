package uaparser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewClientHints(t *testing.T) {
	headers := map[string]string{
		"Sec-Ch-Ua":                   "\"Microsoft Edge\";v=\"131\", \"Chromium\";v=\"131\", \"Not_A IBrand\";v=\"24\"",
		"Sec-Ch-Ua-Arch":              "\"x86\"",
		"Sec-Ch-Ua-bitness":           "\"64\"",
		"Sec-Ch-Ua-Full-Version-List": "\"Microsoft Edge\";v=\"131.0.2903.51\", \"Chromium\";v=\"131.0.6778.70\", \"Not_A IBrand\";v=\"24.0.0.0\"",
		"Sec-Ch-Ua-mobile":            "?0",
		"Sec-Ch-Ua-model":             "",
		"Sec-Ch-Ua-platform":          "\"Windows\"",
		"Sec-Ch-Ua-platform-Version":  "\"15.0.0\"",
	}
	ch := NewClientHints(headers)
	t.Logf("client hints: %+v", ch)
}

func BenchmarkNewClientHints(b *testing.B) {
	for range b.N {
		headers := map[string]string{
			"Sec-Ch-Ua":                   "\"Microsoft Edge\";v=\"131\", \"Chromium\";v=\"131\", \"Not_A IBrand\";v=\"24\"",
			"Sec-Ch-Ua-Arch":              "\"x86\"",
			"Sec-Ch-Ua-bitness":           "\"64\"",
			"Sec-Ch-Ua-Full-Version-List": "\"Microsoft Edge\";v=\"131.0.2903.51\", \"Chromium\";v=\"131.0.6778.70\", \"Not_A IBrand\";v=\"24.0.0.0\"",
			"Sec-Ch-Ua-mobile":            "?0",
			"Sec-Ch-Ua-model":             "",
			"Sec-Ch-Ua-platform":          "\"Windows\"",
			"Sec-Ch-Ua-platform-Version":  "\"15.0.0\"",
		}
		NewClientHints(headers)
	}
}

func TestMapUACHHeaders(t *testing.T) {
	headers := map[string]string{
		"Sec-Ch-Ua":                   "\"Chromium\";v=\"93\", \"Google Chrome\";v=\"93\", \" Not;A Brand\";v=\"99\"",
		"Sec-Ch-Ua-Arch":              "\"arm\"",
		"Sec-Ch-Ua-bitness":           "\"64\"",
		"Sec-Ch-Ua-Full-Version-List": "\"Chromium\";v=\"93.0.1.2\", \"Google Chrome\";v=\"93.0.1.2\", \" Not;A Brand\";v=\"99.0.1.2\"",
		"Sec-Ch-Ua-mobile":            "?1",
		"Sec-Ch-Ua-model":             "Pixel 99",
		"Sec-Ch-Ua-platform":          "\"Windows\"",
		"Sec-Ch-Ua-platform-Version":  "\"13\"",
		"User-Agent":                  "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36",
	}

	uap := NewUAParser(headers["User-Agent"], headers, nil).Result()
	browser := NewUAParser(headers["User-Agent"], headers, nil).Browser()
	cpu := NewUAParser(headers["User-Agent"], headers, nil).CPU()
	device := NewUAParser(headers["User-Agent"], headers, nil).Device()
	engine := NewUAParser(headers["User-Agent"], headers, nil).Engine()
	os := NewUAParser(headers["User-Agent"], headers, nil).Os()

	t.Run("Can read from client-hints headers using `withClientHints()`", func(t *testing.T) {
		assert.Equal(t, "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36", uap.UA)
		assert.Equal(t, "Chrome", uap.Browser.Name)
		assert.Equal(t, "93.0.1.2", uap.Browser.Version)
		assert.Equal(t, "93", uap.Browser.Major)
		assert.Equal(t, "Chrome", browser.Name)
		assert.Equal(t, "93.0.1.2", browser.Version)
		assert.Equal(t, "93", browser.Major)
		assert.Equal(t, "arm64", uap.Cpu.Architecture)
		assert.Equal(t, "arm64", cpu.Architecture)
		assert.Equal(t, "mobile", uap.Device.Type)
		assert.Equal(t, "Pixel 99", uap.Device.Model)
		assert.Equal(t, "Google", uap.Device.Vendor)
		assert.Equal(t, "mobile", device.Type)
		assert.Equal(t, "Pixel 99", device.Model)
		assert.Equal(t, "Google", device.Vendor)
		assert.Equal(t, "Blink", uap.Engine.Name)
		assert.Equal(t, "93.0.1.2", uap.Engine.Version)
		assert.Equal(t, "Blink", engine.Name)
		assert.Equal(t, "93.0.1.2", engine.Version)
		assert.Equal(t, "Windows", uap.Os.Name)
		assert.Equal(t, "11", uap.Os.Version)
		assert.Equal(t, "Windows", os.Name)
		assert.Equal(t, "11", os.Version)
	})

	t.Run("Only read from user-agent header when called without `withClientHints()`", func(t *testing.T) {
		uap = NewUAParser(headers["User-Agent"], nil, nil).Result()
		browser = NewUAParser(headers["User-Agent"], nil, nil).Browser()
		cpu = NewUAParser(headers["User-Agent"], nil, nil).CPU()
		device = NewUAParser(headers["User-Agent"], nil, nil).Device()
		engine = NewUAParser(headers["User-Agent"], nil, nil).Engine()
		os = NewUAParser(headers["User-Agent"], nil, nil).Os()

		assert.Equal(t, "Chrome", uap.Browser.Name)
		assert.Equal(t, "110.0.0.0", uap.Browser.Version)
		assert.Equal(t, "110", uap.Browser.Major)
		assert.Equal(t, "amd64", uap.Cpu.Architecture)
		assert.Equal(t, "", uap.Device.Type)
		assert.Equal(t, "", uap.Device.Model)
		assert.Equal(t, "", uap.Device.Vendor)
		assert.Equal(t, "Blink", uap.Engine.Name)
		assert.Equal(t, "110.0.0.0", uap.Engine.Version)
		assert.Equal(t, "Linux", uap.Os.Name)
		assert.Equal(t, "", uap.Os.Version)
	})

	t.Run("Fallback to user-agent header when using `withClientHints()` but found no client hints-related headers", func(t *testing.T) {
		headers2 := map[string]string{
			"sec-ch-ua-mobile": "?1",
			"user-agent":       "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36",
		}
		uap = NewUAParser(headers2["user-agent"], headers2, nil).Result()

		assert.Equal(t, "Chrome", uap.Browser.Name)
		assert.Equal(t, "110.0.0.0", uap.Browser.Version)
		assert.Equal(t, "110", uap.Browser.Major)
		assert.Equal(t, "amd64", uap.Cpu.Architecture)
		assert.Equal(t, "mobile", uap.Device.Type)
		assert.Equal(t, "", uap.Device.Model)
		assert.Equal(t, "", uap.Device.Vendor)
		assert.Equal(t, "Blink", uap.Engine.Name)
		assert.Equal(t, "110.0.0.0", uap.Engine.Version)
		assert.Equal(t, "Linux", uap.Os.Name)
		assert.Equal(t, "", uap.Os.Version)
	})

	t.Run("Can detect Apple silicon from client hints data", func(t *testing.T) {
		httpHeadersFromAppleSilicon := map[string]string{
			"sec-ch-ua-arch":     "arm",
			"sec-ch-ua-platform": "macOS",
			"sec-ch-ua-mobile":   "?0",
			"sec-ch-ua":          "\"Google Chrome\";v=\"111\", \"Not(A:Brand\";v=\"8\", \"Chromium\";v=\"111\"",
			"user-agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:97.0) Gecko/20100101 Firefox/97.0",
		}

		uap = NewUAParser(httpHeadersFromAppleSilicon["user-agent"], httpHeadersFromAppleSilicon, nil).Result()

		assert.Equal(t, true, uap.Os.Name == "macOS")
		assert.Equal(t, true, uap.Cpu.Architecture == "arm")
		assert.Equal(t, false, uap.Device.Type == "mobile")
		assert.Equal(t, false, uap.Device.Type == "tablet")
	})

	t.Run("Can detect form-factors from client-hints", func(t *testing.T) {
		headersVR := map[string]string{
			"sec-ch-ua-form-factors": "\"VR\"",
		}
		headersEInk := map[string]string{
			"sec-ch-ua-form-factors": "\"Tablet\", \"EInk\"",
		}
		headersUnknown := map[string]string{
			"sec-ch-ua-form-factors": "\"Unknown\"",
		}

		uapVR := NewUAParser("", headersVR, nil).Result()
		assert.Equal(t, "xr", uapVR.Device.Type)

		uapEInk := NewUAParser("", headersEInk, nil).Result()
		assert.Equal(t, "tablet", uapEInk.Device.Type)

		uapUnknown := NewUAParser("", headersUnknown, nil).Result()
		assert.Equal(t, "", uapUnknown.Device.Type)
	})

	t.Run("Avoid error on headers variation", func(t *testing.T) {
		headers = map[string]string{
			"sec-ch-ua":                   "\"Google Chrome\";v=\"119\", \"Chromium\";v=\"119\", \"Not?A_Brand\";v=\"24\"",
			"sec-ch-ua-full-version-list": "\"Google Chrome\", \"Chromium\", \"Not?A_Brand\";v=\"24.0.0.0\"",
			"sec-ch-ua-full-version":      "\"\"",
			"sec-ch-ua-mobile":            "?0",
			"sec-ch-ua-bitness":           "\"\"",
			"sec-ch-ua-model":             "\"\"",
			"sec-ch-ua-platform":          "\"Windows\"",
			"sec-ch-ua-platform-version":  "\"\"",
			"sec-ch-ua-wow64":             "?0",
		}

		uap = NewUAParser(headers["user-agent"], headers, nil).Result()

		assert.Equal(t, "Chrome", uap.Browser.Name)
		assert.Equal(t, "", uap.Browser.Version)
		assert.Equal(t, "", uap.Browser.Major)
	})

	t.Run("Prioritize more specific brand name regardless the order", func(t *testing.T) {
		headersList := []map[string]string{
			{"sec-ch-ua-full-version-list": "\"Not_A Brand;v=8, Chromium;v=120.0.6099.131, Google Chrome;v=120.0.6099.132\""},
			{"sec-ch-ua-full-version-list": "\"Chromium;v=120.0.6099.131, Not_A Brand;v=8, Google Chrome;v=120.0.6099.132\""},
			{"sec-ch-ua-full-version-list": "\"Google Chrome;v=120.0.6099.132, Chromium;v=120.0.6099.131, Not_A Brand;v=8\""},
			{"sec-ch-ua-full-version-list": "\"Microsoft Edge;v=120.0.6099.133, Google Chrome;v=120.0.6099.132, Chromium;v=120.0.6099.131, Not_A Brand;v=8\""},
			{"sec-ch-ua-full-version-list": "\"Chromium;v=120.0.6099.131, Google Chrome;v=120.0.6099.132, Microsoft Edge;v=120.0.6099.133, Not_A Brand;v=8\""},
			{"sec-ch-ua-full-version-list": "\"Not_A Brand;v=8, Microsoft Edge;v=120.0.6099.133, Google Chrome;v=120.0.6099.132, Chromium;v=120.0.6099.131\""},
		}

		expectedResults := []struct {
			name    string
			version string
		}{
			{"Chrome", "120.0.6099.132"},
			{"Chrome", "120.0.6099.132"},
			{"Chrome", "120.0.6099.132"},
			{"Edge", "120.0.6099.133"},
			{"Edge", "120.0.6099.133"},
			{"Edge", "120.0.6099.133"},
		}

		for i, headers := range headersList {
			uap := NewUAParser("", headers, nil).Result()
			assert.Equal(t, expectedResults[i].name, uap.Browser.Name)
			assert.Equal(t, expectedResults[i].version, uap.Browser.Version)
		}
	})
}

func TestUACHHeaders(t *testing.T) {
	type Headers map[string]string

	type Expect struct {
		Browser IBrowser `json:"browser"`
	}

	testCases := []struct {
		Headers Headers
		Expect  Expect
	}{
		{
			Headers: Headers{
				"sec-ch-ua": `"Avast Secure Browser";v="131", "Chromium";v="131", "Not_A Brand";v="24"`,
			},
			Expect: Expect{
				Browser: IBrowser{
					Name:    "Avast Secure Browser",
					Version: "131",
					Major:   "131",
				},
			},
		},
		{
			Headers: Headers{
				"sec-ch-ua": `"Not A(Brand";v="8", "Chromium";v="132", "Brave";v="132"`,
			},
			Expect: Expect{
				Browser: IBrowser{
					Name:    "Brave",
					Version: "132",
					Major:   "132",
				},
			},
		},
		{
			Headers: Headers{
				"sec-ch-ua": `"Google Chrome";v="111", "Not(A:Brand";v="8", "Chromium";v="111"`,
			},
			Expect: Expect{
				Browser: IBrowser{
					Name:    "Chrome",
					Version: "111",
					Major:   "111",
				},
			},
		},
		{
			Headers: Headers{
				"sec-ch-ua": `"Chromium";v="124", "HeadlessChrome";v="124", "Not-A.Brand";v="99"`,
			},
			Expect: Expect{
				Browser: IBrowser{
					Name:    "Chrome Headless",
					Version: "124",
					Major:   "124",
				},
			},
		},
		{
			Headers: Headers{
				"sec-ch-ua": `"Android WebView";v="123", "Not:A-Brand";v="8", "Chromium";v="123"`,
			},
			Expect: Expect{
				Browser: IBrowser{
					Name:    "Chrome WebView",
					Version: "123",
					Major:   "123",
				},
			},
		},
		{
			Headers: Headers{
				"sec-ch-ua": `"DuckDuckGo";v="131", "Chromium";v="131", "Not_A Brand";v="24"`,
			},
			Expect: Expect{
				Browser: IBrowser{
					Name:    "DuckDuckGo",
					Version: "131",
					Major:   "131",
				},
			},
		},
		{
			Headers: Headers{
				"sec-ch-ua": `"Not_A Brand";v="8", "Chromium";v="120", "Microsoft Edge";v="120"`,
			},
			Expect: Expect{
				Browser: IBrowser{
					Name:    "Edge",
					Version: "120",
					Major:   "120",
				},
			},
		},
		{
			Headers: Headers{
				"sec-ch-ua": `"Not.A/Brand";v="8", "Chromium";v="114", "HuaweiBrowser";v="114"`,
			},
			Expect: Expect{
				Browser: IBrowser{
					Name:    "Huawei Browser",
					Version: "114",
					Major:   "114",
				},
			},
		},
		{
			Headers: Headers{
				"sec-ch-ua": `"Miui Browser";v="123", "Not:A-Brand";v="8", "Chromium";v="123"`,
			},
			Expect: Expect{
				Browser: IBrowser{
					Name:    "MIUI Browser",
					Version: "123",
					Major:   "123",
				},
			},
		},
		{
			Headers: Headers{
				"sec-ch-ua": `"Chromium";v="130", "Oculus Browser";v="36", "Not?A_Brand";v="99"`,
			},
			Expect: Expect{
				Browser: IBrowser{
					Name:    "Oculus Browser",
					Version: "36",
					Major:   "36",
				},
			},
		},
		{
			Headers: Headers{
				"sec-ch-ua": `"Opera";v="116", "Chromium";v="131", "Not_A Brand";v="24"`,
			},
			Expect: Expect{
				Browser: IBrowser{
					Name:    "Opera",
					Version: "116",
					Major:   "116",
				},
			},
		},
		{
			Headers: Headers{
				"sec-ch-ua": `"Chromium";v="128", "Not;A=Brand";v="24", "Opera GX";v="114"`,
			},
			Expect: Expect{
				Browser: IBrowser{
					Name:    "Opera GX",
					Version: "114",
					Major:   "114",
				},
			},
		},
		{
			Headers: Headers{
				"sec-ch-ua": `"OperaMobile";v="86", ";Not A Brand";v="99", "Opera";v="115", "Chromium";v="130"`,
			},
			Expect: Expect{
				Browser: IBrowser{
					Name:    "Opera Mobi",
					Version: "86",
					Major:   "86",
				},
			},
		},
		{
			Headers: Headers{
				"sec-ch-ua": `"Chromium";v="132", "OperaMobile";v="87", "Opera";v="117", " Not A;Brand";v="99"`,
			},
			Expect: Expect{
				Browser: IBrowser{
					Name:    "Opera Mobi",
					Version: "87",
					Major:   "87",
				},
			},
		},
		{
			Headers: Headers{
				"sec-ch-ua": `"Chromium";v="125", "Not.A/Brand";v="24", "Samsung Internet";v="27.0"`,
			},
			Expect: Expect{
				Browser: IBrowser{
					Name:    "Samsung Internet",
					Version: "27.0",
					Major:   "27",
				},
			},
		},
		{
			Headers: Headers{
				"sec-ch-ua": `"Chromium";v="130", "YaBrowser";v="24.12", "Not?A_Brand";v="99", "Yowser";v="2.5"`,
			},
			Expect: Expect{
				Browser: IBrowser{
					Name:    "Yandex",
					Version: "24.12",
					Major:   "24",
				},
			},
		},
	}

	for _, test := range testCases {
		browser := NewUAParser("", test.Headers, nil).Browser()
		assert.Equal(t, test.Expect.Browser, browser)
	}
}

func TestIdentifyVendorAndTypeOfDeviceFromGivenModelName(t *testing.T) {
	type Expect struct {
		Vendor string `json:"vendor"`
		Type   string `json:"type"`
	}

	testCases := []struct {
		Model  string `json:"model"`
		Expect Expect `json:"expect"`
	}{
		{Model: "220733SG", Expect: Expect{Vendor: "Xiaomi", Type: "mobile"}},
		{Model: "5087Z", Expect: Expect{Vendor: "TCL", Type: "mobile"}},
		{Model: "9137W", Expect: Expect{Vendor: "TCL", Type: "tablet"}},
		{Model: "BE2015", Expect: Expect{Vendor: "OnePlus", Type: "mobile"}},
		{Model: "CPH2389", Expect: Expect{Vendor: "OPPO", Type: "mobile"}},
		{Model: "Infinix X669C", Expect: Expect{Vendor: "Infinix", Type: "mobile"}},
		{Model: "itel L6502", Expect: Expect{Vendor: "itel", Type: "mobile"}},
		{Model: "Lenovo TB-X606F", Expect: Expect{Vendor: "Lenovo", Type: "tablet"}},
		{Model: "LM-Q720", Expect: Expect{Vendor: "LG", Type: "mobile"}},
		{Model: "M2003J15SC", Expect: Expect{Vendor: "Xiaomi", Type: "mobile"}},
		{Model: "MAR-LX1A", Expect: Expect{Vendor: "Huawei", Type: "mobile"}},
		{Model: "moto g(20)", Expect: Expect{Vendor: "Motorola", Type: "mobile"}},
		{Model: "Nokia C210", Expect: Expect{Vendor: "Nokia", Type: "mobile"}},
		{Model: "Pixel 8", Expect: Expect{Vendor: "Google", Type: "mobile"}},
		{Model: "Redmi Note 9S", Expect: Expect{Vendor: "Xiaomi", Type: "mobile"}},
		{Model: "RMX3830", Expect: Expect{Vendor: "Realme", Type: "mobile"}},
		{Model: "SM-S536DL", Expect: Expect{Vendor: "Samsung", Type: "mobile"}},
		{Model: "SM-S546VL", Expect: Expect{Vendor: "Samsung", Type: "mobile"}},
		{Model: "SM-T875", Expect: Expect{Vendor: "Samsung", Type: "tablet"}},
		{Model: "STK-L21", Expect: Expect{Vendor: "Huawei", Type: "mobile"}},
		{Model: "T430W", Expect: Expect{Vendor: "TCL", Type: "mobile"}},
		{Model: "TECNO KI5k", Expect: Expect{Vendor: "TECNO", Type: "mobile"}},
		{Model: "vivo 1820", Expect: Expect{Vendor: "Vivo", Type: "mobile"}},
		{Model: "Xbox", Expect: Expect{Vendor: "Microsoft", Type: "console"}},
	}

	for _, test := range testCases {
		headers := map[string]string{
			"Sec-Ch-Ua-Model": test.Model,
		}
		device := NewUAParser("", headers, nil).Device()
		assert.Equal(t, test.Model, device.Model)
		assert.Equal(t, test.Expect.Vendor, device.Vendor)
		assert.Equal(t, test.Expect.Type, device.Type)
	}
}
