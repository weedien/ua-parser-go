package uaparser

import (
	"encoding/json"
	"fmt"
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

func TestUAParser_Result(t *testing.T) {
	// Helper function to check for undefined values
	isUndefined := func(val string) bool {
		return val == "undefined"
	}

	// Helper function to handle reflection and comparison
	compareField := func(field reflect.Value, key, expectedValue string, label string, ua string) bool {
		if field.Kind() == reflect.Struct {
			fieldName := strings.Title(key)
			actualVal := field.FieldByName(fieldName)
			if actualVal.IsValid() && actualVal.Kind() == reflect.String {
				actual := actualVal.String()
				if strings.ToLower(actual) != strings.ToLower(expectedValue) {
					t.Errorf("\033[35m%s\u001B[0m, [%s] \033[34mkey: %s\033[0m, \033[32mexpect: %s\033[0m, \033[33mactual: %s\033[0m",
						label, ua, key, expectedValue, actual)
					return false
				}
			}
		}
		return true
	}

	// Test cases with their respective JSON file paths
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
			r := NewUAParser(tc.Ua, nil, nil).Result()

			// Reflect to find the correct field by label (Browser, Cpu, etc.)
			field := reflect.ValueOf(r).FieldByName(singleCase.label)

			// Iterate over all expected values in the TestCase
			for key, val := range tc.Expect {
				if isUndefined(val) {
					continue
				}

				// Compare the expected value with the actual parsed value
				if !compareField(field, key, val, singleCase.label, tc.Ua) {
					break
				}
			}
		}
	}
}

func TestUAParser_Concurrent(t *testing.T) {
	tc := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36 Edg/132.0.0.0"
	parser := NewUAParser(tc, nil, nil)

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
				_ = NewUAParser(tc, nil, nil).Result()
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
			_ = NewUAParser(tc.Ua, nil, nil).Result()
		}
	}
}

// ~> 900μs/op
func BenchmarkUAParser_Single(b *testing.B) {
	tc := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36 Edg/132.0.0.0"

	for i := 0; i < b.N; i++ {
		_ = NewUAParser(tc, nil, nil).Result()
	}
}

func BenchmarkUAParser_Methods(b *testing.B) {
	tc := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36 Edg/132.0.0.0"

	b.Run("Result", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = NewUAParser(tc, nil, nil).Result()
		}
	})

	b.Run("Browser", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = NewUAParser(tc, nil, nil).Browser()
		}
	})

	b.Run("Os", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = NewUAParser(tc, nil, nil).Os()
		}
	})

	b.Run("CPU", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = NewUAParser(tc, nil, nil).CPU()
		}
	})

	// 耗时最长，超出其他方法的总和一个数量级
	b.Run("Device", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = NewUAParser(tc, nil, nil).Device()
		}
	})

	b.Run("Engine", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = NewUAParser(tc, nil, nil).Engine()
		}
	})
}

func FuzzUAParser_Result(f *testing.F) {
	f.Add("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36 Edg/132.0.0.0")
	f.Add("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.6261.119 Safari/537.36 XiaoMi/MiuiBrowser/19.0.120123")
	f.Add("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36 NetType/WIFI MicroMessenger/7.0.20.1781(0x6700143B) WindowsWechat(0x63090c25) XWEB/11581 Flue")

	f.Fuzz(func(t *testing.T, ua string) {
		start := time.Now()
		_ = NewUAParser(ua, nil, nil).Result()
		elapsed := time.Since(start)
		if elapsed > 100*time.Millisecond {
			t.Fatalf("[Potential ReDos] Takes %v. User-Agent: %s", elapsed, ua)
		}
	})
}
