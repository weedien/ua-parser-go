package uaparser

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

type TestCase struct {
	Desc   string            `json:"desc"`
	Ua     string            `json:"ua"`
	Expect map[string]string `json:"expect"`
}

func loadJson(path string) []TestCase {
	var tests []TestCase
	if fi, err := os.Stat(path); err == nil && fi.IsDir() {
		files, _ := os.ReadDir(path)
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			tests = append(tests, loadJsonFile(path+"/"+file.Name())...)
		}
	} else {
		tests = loadJsonFile(path)
	}
	return tests
}

func loadJsonFile(path string) []TestCase {
	var tests []TestCase
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return tests
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	data, err := io.ReadAll(f)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return tests
	}

	if err := json.Unmarshal(data, &tests); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return tests
	}
	return tests
}

// Helper function to check for undefined values
func isUndefined(val string) bool {
	return val == "undefined"
}

// Helper function to handle reflection and comparison
func compareField(field reflect.Value, key, expectedValue string, label string, ua string) error {
	if field.Kind() == reflect.Struct {
		fieldName := strings.Title(key)
		actualVal := field.FieldByName(fieldName)
		if actualVal.IsValid() && actualVal.Kind() == reflect.String {
			actual := actualVal.String()
			if strings.ToLower(actual) != strings.ToLower(expectedValue) {
				return fmt.Errorf("\033[35m%s\u001B[0m, [%s] \033[34mkey: %s\033[0m, \033[32mexpect: %s\033[0m, \033[33mactual: %s\033[0m",
					label, ua, key, expectedValue, actual)
			}
		}
	}
	return nil
}

func TestUAParser_Result(t *testing.T) {
	alltestcases := []struct {
		label string
		list  []TestCase
	}{
		{
			label: "Browser",
			list:  loadJson("./data/ua/browser/browser-all.json"),
		},
		{
			label: "Cpu",
			list:  loadJson("./data/ua/cpu/cpu-all.json"),
		},
		{
			label: "Device",
			list:  loadJson("./data/ua/device/"),
		},
		{
			label: "Engine",
			list:  loadJson("./data/ua/engine/engine-all.json"),
		},
		{
			label: "Os",
			list:  loadJson("./data/ua/os/"),
		},
	}

	// Loop through each test case group (Browser, Cpu, etc.)
	for _, singleCase := range alltestcases {
		for _, tc := range singleCase.list {
			// Get the result from UAParser
			r := NewUAParser(tc.Ua).Result()

			// Reflect to find the correct field by label (Browser, Cpu, etc.)
			field := reflect.ValueOf(r).FieldByName(singleCase.label)

			// Iterate over all expected values in the TestCase
			for key, val := range tc.Expect {
				if isUndefined(val) {
					continue
				}

				// Compare the expected value with the actual parsed value
				if err := compareField(field, key, val, singleCase.label, tc.Ua); err != nil {
					t.Error(err.Error())
				}
			}
		}
	}
}

func TestUAParser_Concurrent(t *testing.T) {
	tc := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36 Edg/132.0.0.0"
	parser := NewUAParser(tc)

	var wg sync.WaitGroup
	numGoroutines := 100

	// Function to test a single method concurrently
	testMethod := func(method func() interface{}) {
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_ = method()
			}()
		}
	}

	testMethod(func() interface{} { return parser.Result() })

	wg.Wait()
}

func BenchmarkUAParser_Concurrent(b *testing.B) {
	tc := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36 Edg/132.0.0.0"
	numGoroutines := 100

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		for j := 0; j < numGoroutines; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_ = NewUAParser(tc).Result()
			}()
		}
		wg.Wait()
	}
}

// < 250ms/op
func BenchmarkUAParser_Result(b *testing.B) {
	// Load test cases
	tests := loadJson("./data/ua/browser/browser-all.json")

	// Benchmark the Result() method
	for i := 0; i < b.N; i++ {
		for _, tc := range tests {
			_ = NewUAParser(tc.Ua).Result()
		}
	}
}

// ~> 900μs/op
func BenchmarkUAParser_Single(b *testing.B) {
	tc := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36 Edg/132.0.0.0"

	for i := 0; i < b.N; i++ {
		_ = NewUAParser(tc).Result()
	}
}

func BenchmarkUAParser_Methods(b *testing.B) {
	tc := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36 Edg/132.0.0.0"

	b.Run("Result", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = NewUAParser(tc).Result()
		}
	})

	b.Run("Browser", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = NewUAParser(tc).Browser()
		}
	})

	b.Run("Os", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = NewUAParser(tc).Os()
		}
	})

	b.Run("CPU", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = NewUAParser(tc).CPU()
		}
	})

	// 耗时最长，超出其他方法的总和一个数量级
	b.Run("Device", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = NewUAParser(tc).Device()
		}
	})

	b.Run("Engine", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = NewUAParser(tc).Engine()
		}
	})
}

func FuzzUAParser_Result(f *testing.F) {
	f.Add("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36 Edg/132.0.0.0")
	f.Add("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.6261.119 Safari/537.36 XiaoMi/MiuiBrowser/19.0.120123")
	f.Add("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36 NetType/WIFI MicroMessenger/7.0.20.1781(0x6700143B) WindowsWechat(0x63090c25) XWEB/11581 Flue")

	f.Fuzz(func(t *testing.T, ua string) {
		start := time.Now()
		_ = NewUAParser(ua).Result()
		elapsed := time.Since(start)
		if elapsed > 100*time.Millisecond {
			t.Fatalf("[Potential ReDos] Takes %v. User-Agent: %s", elapsed, ua)
		}
	})
}

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

	uap := NewUAParser(headers["User-Agent"]).WithHeaders(headers).Result()
	browser := NewUAParser(headers["User-Agent"]).WithHeaders(headers).Browser()
	cpu := NewUAParser(headers["User-Agent"]).WithHeaders(headers).CPU()
	device := NewUAParser(headers["User-Agent"]).WithHeaders(headers).Device()
	engine := NewUAParser(headers["User-Agent"]).WithHeaders(headers).Engine()
	_os := NewUAParser(headers["User-Agent"]).WithHeaders(headers).Os()

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
		assert.Equal(t, "Windows", _os.Name)
		assert.Equal(t, "11", _os.Version)
	})

	t.Run("Only read from user-agent header when called without `withClientHints()`", func(t *testing.T) {
		uap = NewUAParser(headers["User-Agent"]).Result()
		browser = NewUAParser(headers["User-Agent"]).Browser()
		cpu = NewUAParser(headers["User-Agent"]).CPU()
		device = NewUAParser(headers["User-Agent"]).Device()
		engine = NewUAParser(headers["User-Agent"]).Engine()
		_os = NewUAParser(headers["User-Agent"]).Os()

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
		uap = NewUAParser(headers2["user-agent"]).WithHeaders(headers2).Result()

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

		uap = NewUAParser(httpHeadersFromAppleSilicon["user-agent"]).WithHeaders(httpHeadersFromAppleSilicon).Result()

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

		uapVR := NewUAParser("").WithHeaders(headersVR).Result()
		assert.Equal(t, "xr", uapVR.Device.Type)

		uapEInk := NewUAParser("").WithHeaders(headersEInk).Result()
		assert.Equal(t, "tablet", uapEInk.Device.Type)

		uapUnknown := NewUAParser("").WithHeaders(headersUnknown).Result()
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

		uap = NewUAParser(headers["user-agent"]).WithHeaders(headers).Result()

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
			uap := NewUAParser("").WithHeaders(headers).Result()
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
		browser := NewUAParser("").WithHeaders(test.Headers).Browser()
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
		device := NewUAParser("").WithHeaders(headers).Device()
		assert.Equal(t, test.Model, device.Model)
		assert.Equal(t, test.Expect.Vendor, device.Vendor)
		assert.Equal(t, test.Expect.Type, device.Type)
	}
}
