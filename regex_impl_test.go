package uaparser

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

func Test(t *testing.T) {
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

	for _, singleCase := range alltestcases {
		for _, tc := range singleCase.list {
			r := NewUAParser(tc.Ua, nil, nil).Result()
			field := reflect.ValueOf(r).FieldByName(singleCase.label)
			for key, val := range tc.Expect {
				if val == "undefined" {
					continue
				}
				if field.Kind() == reflect.Struct {
					fieldName := strings.Title(key)
					actualVal := field.FieldByName(fieldName)
					if actualVal.IsValid() && actualVal.Kind() == reflect.String {
						if actual := actualVal.String(); strings.ToLower(actual) != strings.ToLower(val) {
							t.Errorf("%s, [%s] \033[34mkey: %s\033[0m, \033[32mexpect: %s\033[0m, \033[33mactual: %s\033[0m",
								singleCase.label, tc.Ua, key, val, actual)
						}
					}
				}
			}
		}
	}
}

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
