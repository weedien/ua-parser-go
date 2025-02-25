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

var testcases2 = []struct {
	title    string
	ua       string
	expected string
}{

	// 1. 主流的桌面端浏览器
	// 2. 小众一些的主要用于移动端的浏览器 夸克UCQQ搜狗百度 （可能也有桌面端但用的少）
	// 3. 应用内置浏览器 微信/QQ
	// 平台主要考虑 Windows Android 然后是 Mac Linux iOS HarmonyOS

	// Microsoft Edge
	{
		title:    "EdgeDesktop",
		ua:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36 Edg/132.0.0.0",
		expected: "Browser:Edge-132.0.0.0 Engine:Blink-132.0.0.0 Device:desktop Os:Windows-10.0 Cpu:amd64",
	},
	{
		title:    "EdgeMobile",
		ua:       "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Mobile Safari/537.36 EdgA/132.0.0.0",
		expected: "Browser:Edge-132.0.0.0 Engine:Blink-132.0.0.0 Device:mobile Os:Android-10",
	},

	// Quark
	{
		title:    "QuarkDesktop",
		ua:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 QuarkPC/2.0.5.220",
		expected: "Browser:Quark-2.0.5.220 Engine:Blink-112.0.0.0 Device:desktop Os:Windows-10 Cpu:amd64",
	},
	{
		title:    "QuarkMobile",
		ua:       "Mozilla/5.0 (Linux; U; Android 13; zh-CN; M2102J2SC Build/TKQ1.221114.001) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/123.0.6312.80 Quark/7.7.5.740 Mobile Safari/537.36",
		expected: "Browser:Quark-7.7.5.740 Engine:Blink-123.0.6312.80 Device:mobile-M2102J2SC Os:Android-13",
	},

	// MiuiBrowser
	{
		title:    "MiuiBrowserMobile",
		ua:       "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.6261.119 Safari/537.36 XiaoMi/MiuiBrowser/19.0.120123",
		expected: "Browser:MiuiBrowser-19.0.120123 Engine:Blink-122.0.6261.119 Device:mobile-XiaoMi Os:Android Cpu:amd64",
	},

	// WechatBrowser
	{
		title:    "WechatDesktop",
		ua:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36 NetType/WIFI MicroMessenger/7.0.20.1781(0x6700143B) WindowsWechat(0x63090c25) XWEB/11581 Flue",
		expected: "Browser:WeChat-7.0.20.1781 Engine:Blink-122.0.0.0 Device:desktop Os:Windows-10.0 Cpu:amd64",
	},
	{
		title:    "WechatMobile",
		ua:       "Mozilla/5.0 (Linux; Android 13; M2102J2SC Build/TKQ1.221114.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/130.0.6723.103 Mobile Safari/537.36 XWEB/1300333 MMWEBSDK/20241202 MMWEBID/9403 MicroMessenger/8.0.56.2800(0x28003855) WeChat/arm64 Weixin NetType/4G Language/zh_CN ABI/arm64",
		expected: "Browser:WeChat-8.0.56.2800 Engine:Blink-130.0.6723.103 Device:mobile-M2102J2SC Os:Android-13",
	},

	// HuaweiBrowser
	{
		title:    "HuaweiBrowserMobile",
		ua:       "Mozilla/5.0 (Linux; Android 10; HarmonyOS; VOG-AL00; HMSCore 6.15.0.302; GMSCore 20.15.16) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.5735.196 HuaweiBrowser/15.0.10.302 Mobile Safari/537.36",
		expected: "Browser:HuaweiBrowser-15.0.10.302 Engine:Blink-114.0.5735.196 Device:mobile-VOG-AL00 Os:HarmonyOS-10",
	},
}

var testcases_aigen = []struct {
	title    string
	ua       string
	expected string
}{
	{
		title:    "EdgeDesktop",
		ua:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36 Edg/132.0.0.0",
		expected: "Browser:Edge-132.0.0.0 Engine:Blink-132.0.0.0 Device:desktop Os:Windows-10.0 Cpu:amd64",
	},
	{
		title:    "EdgeMac",
		ua:       "Mozilla/5.0 (Macintosh; Intel Mac Os X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36 Edg/132.0.0.0",
		expected: "Browser:Edge-132.0.0.0 Engine:Blink-132.0.0.0 Device:desktop Os:Mac-10.15.7 Cpu:x64",
	},
	{
		title:    "EdgeLinux",
		ua:       "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Edge/132.0.0.0 Safari/537.36",
		expected: "Browser:Edge-132.0.0.0 Engine:Blink-132.0.0.0 Device:desktop Os:Linux Cpu:x64",
	},
	{
		title:    "EdgeMobile",
		ua:       "Mozilla/5.0 (Linux; Android 12; Pixel 6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Mobile Safari/537.36 Edg/132.0.0.0",
		expected: "Browser:Edge-132.0.0.0 Engine:Blink-132.0.0.0 Device:mobile Os:Android-12 Cpu:arm64",
	},
	{
		title:    "EdgeiOS",
		ua:       "Mozilla/5.0 (iPhone; Cpu iPhone Os 18_3 like Mac Os X) AppleWebKit/605.1.15 (KHTML, like Gecko) EdgiOS/132.0.0.0 Version/18.0 Mobile/15E148 Safari/604.1",
		expected: "Browser:Edge-132.0.0.0 Engine:WebKit-132.0.0.0 Device:mobile Os:iOS-18.3 Cpu:arm64",
	},
	{
		title:    "ChromeDesktop",
		ua:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36",
		expected: "Browser:Chrome-132.0.0.0 Engine:Blink-132.0.0.0 Device:desktop Os:Windows-10.0 Cpu:amd64",
	},
	{
		title:    "ChromeMac",
		ua:       "Mozilla/5.0 (Macintosh; Intel Mac Os X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36",
		expected: "Browser:Chrome-132.0.0.0 Engine:Blink-132.0.0.0 Device:desktop Os:Mac-10.15.7 Cpu:x64",
	},
	{
		title:    "ChromeLinux",
		ua:       "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36",
		expected: "Browser:Chrome-132.0.0.0 Engine:Blink-132.0.0.0 Device:desktop Os:Linux Cpu:x64",
	},
	{
		title:    "ChromeMobile",
		ua:       "Mozilla/5.0 (Linux; Android 12; Pixel 6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Mobile Safari/537.36",
		expected: "Browser:Chrome-132.0.0.0 Engine:Blink-132.0.0.0 Device:mobile Os:Android-12 Cpu:arm64",
	},
	{
		title:    "FirefoxDesktop",
		ua:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:132.0) Gecko/20100101 Firefox/132.0",
		expected: "Browser:Firefox-132.0 Engine:Gecko-132.0 Device:desktop Os:Windows-10.0 Cpu:amd64",
	},
	{
		title:    "FirefoxMac",
		ua:       "Mozilla/5.0 (Macintosh; Intel Mac Os X 10_15_7; rv:132.0) Gecko/20100101 Firefox/132.0",
		expected: "Browser:Firefox-132.0 Engine:Gecko-132.0 Device:desktop Os:Mac-10.15.7 Cpu:x64",
	},
	{
		title:    "FirefoxLinux",
		ua:       "Mozilla/5.0 (X11; Linux x86_64; rv:132.0) Gecko/20100101 Firefox/132.0",
		expected: "Browser:Firefox-132.0 Engine:Gecko-132.0 Device:desktop Os:Linux Cpu:x64",
	},
	{
		title:    "FirefoxMobile",
		ua:       "Mozilla/5.0 (Android 12; Mobile; rv:132.0) Gecko/132.0 Firefox/132.0",
		expected: "Browser:Firefox-132.0 Engine:Gecko-132.0 Device:mobile Os:Android-12 Cpu:arm64",
	},
	{
		title:    "SafariDesktop",
		ua:       "Mozilla/5.0 (Macintosh; Intel Mac Os X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/132.0.0.0 Safari/605.1.15",
		expected: "Browser:Safari-132.0.0.0 Engine:WebKit-132.0.0.0 Device:desktop Os:Mac-10.15.7 Cpu:x64",
	},
	{
		title:    "SafariMobile",
		ua:       "Mozilla/5.0 (iPhone; Cpu iPhone Os 18_3 like Mac Os X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/132.0.0.0 Mobile/15E148 Safari/604.1",
		expected: "Browser:Safari-132.0.0.0 Engine:WebKit-132.0.0.0 Device:mobile Os:iOS-18.3 Cpu:arm64",
	},
	{
		title:    "SafariiPad",
		ua:       "Mozilla/5.0 (iPad; Cpu Os 18_3 like Mac Os X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/132.0.0.0 Mobile/15E148 Safari/604.1",
		expected: "Browser:Safari-132.0.0.0 Engine:WebKit-132.0.0.0 Device:tablet Os:iOS-18.3 Cpu:arm64",
	},

	{
		title:    "360安全浏览器_Windows",
		ua:       "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36 QIHU 360SE",
		expected: "Browser:360安全浏览器 Engine:Blink Device:desktop Os:Windows-10.0 Cpu:x86",
	},
	{
		title:    "360极速浏览器_Windows",
		ua:       "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36 QIHU 360EE",
		expected: "Browser:360极速浏览器 Engine:Blink Device:desktop Os:Windows-10.0 Cpu:x86",
	},
	{
		title:    "QQ浏览器_Windows",
		ua:       "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.25 Safari/537.36 Core/1.70.3722.400 QQBrowser/10.5.3739.400",
		expected: "Browser:QQ浏览器 Engine:Blink Device:desktop Os:Windows-10.0 Cpu:x86",
	},
	{
		title:    "UC浏览器_Windows",
		ua:       "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 UBrowser/6.2.3964.2 Safari/537.36",
		expected: "Browser:UC浏览器 Engine:Blink Device:desktop Os:Windows-10.0 Cpu:x86",
	},
	{
		title:    "搜狗浏览器_Windows",
		ua:       "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36 SE 2.X MetaSr 1.0",
		expected: "Browser:搜狗浏览器 Engine:Blink Device:desktop Os:Windows-10.0 Cpu:x86",
	},
	{
		title:    "百度浏览器_Windows",
		ua:       "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.25 Safari/537.36 Core/1.70.3722.400 baidubrowser/10.5.3739.400",
		expected: "Browser:百度浏览器 Engine:Blink Device:desktop Os:Windows-10.0 Cpu:x86",
	},
	{
		title:    "夸克浏览器_Windows",
		ua:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36 Quark/1.0.0",
		expected: "Browser:夸克浏览器 Engine:Blink Device:desktop Os:Windows-10.0 Cpu:x64",
	},

	{
		title:    "360浏览器_Android",
		ua:       "Mozilla/5.0 (Linux; Android 7.1.1; OPPO R9sk Build/NMF26F; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/62.0.3202.97 Mobile Safari/537.36",
		expected: "Browser:360浏览器 Engine:Blink Device:mobile Os:Android-7.1.1 Cpu:arm",
	},
	{
		title:    "QQ浏览器_Android",
		ua:       "Mozilla/5.0 (Linux; U; Android 9; zh-cn; V1816A Build/PKQ1.180819.001) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/66.0.3359.126 MQQBrowser/10.1 Mobile Safari/537.36",
		expected: "Browser:QQ浏览器 Engine:Blink Device:mobile Os:Android-9 Cpu:arm",
	},
	{
		title:    "UC浏览器_Android",
		ua:       "Mozilla/5.0 (Linux; U; Android 8.1.0; zh-CN; MI 5X Build/OPM1.171019.019) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/57.0.2987.108 UCBrowser/12.1.2.992 Mobile Safari/537.36",
		expected: "Browser:UC浏览器 Engine:Blink Device:mobile Os:Android-8.1.0 Cpu:arm",
	},
	{
		title:    "搜狗浏览器_Android",
		ua:       "Mozilla/5.0 (Linux; Android 10; HMA-AL00 Build/HUAWEIHMA-AL00; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/76.0.3809.89 Mobile Safari/537.36 SE 2.X MetaSr 1.0",
		expected: "Browser:搜狗浏览器 Engine:Blink Device:mobile Os:Android-10 Cpu:arm",
	},
	{
		title:    "百度浏览器_Android",
		ua:       "Mozilla/5.0 (Linux; Android 10; HMA-AL00 Build/HUAWEIHMA-AL00; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/76.0.3809.89 Mobile Safari/537.36 T7/11.20 baidubrowser/11.20.2.2",
		expected: "Browser:百度浏览器 Engine:Blink Device:mobile Os:Android-10 Cpu:arm",
	},
	{
		title:    "夸克浏览器_Android",
		ua:       "Mozilla/5.0 (Linux; Android 10; HMA-AL00 Build/HUAWEIHMA-AL00) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Mobile Safari/537.36 Quark/1.0.0",
		expected: "Browser:夸克浏览器 Engine:Blink Device:mobile Os:Android-10 Cpu:arm",
	},
}

var testcases = []string{
	// Opera
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36 OPR/26.0.1656.60",
	"Opera/8.0 (Windows NT 5.1; U; en)",
	"Mozilla/5.0 (Windows NT 5.1; U; en; rv:1.8.1) Gecko/20061208 Firefox/2.0.0 Opera 9.50",
	"Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; en) Opera 9.50",
	// Firefox
	"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:34.0) Gecko/20100101 Firefox/34.0",
	"Mozilla/5.0 (X11; U; Linux x86_64; zh-CN; rv:1.9.2.10) Gecko/20100922 Ubuntu/10.10 (maverick) Firefox/3.6.10",
	// Safari
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/534.57.2 (KHTML, like Gecko) Version/5.1.7 Safari/534.57.2",
	// Chrome
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11",
	"Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/534.16 (KHTML, like Gecko) Chrome/10.0.648.133 Safari/534.16",
	// Edge
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36 Edg/132.0.0.0",
	// 360
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30.0.1599.101 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko",
	// 淘宝浏览器
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/536.11 (KHTML, like Gecko) Chrome/20.0.1132.11 TaoBrowser/2.0 Safari/536.11",
	// QQ浏览器
	"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; .NET4.0C; .NET4.0E; QQBrowser/7.0.3698.400)",
	"Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1; QQDownload 732; .NET4.0C; .NET4.0E)",
	// UC浏览器
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/38.0.2125.122 UBrowser/4.0.3214.0 Safari/537.36",
	// 百度浏览器
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.57 BIDUBrowser/2.x Safari/537.36",
	// 猎豹浏览器
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.57 Safari/537.36 LBBROWSER",
	"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; .NET4.0C; .NET4.0E; LBBROWSER)",
	"Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1; QQDownload 732; .NET4.0C; .NET4.0E; LBBROWSER)",
	// 搜狗浏览器
	"Mozilla/5.0 (Windows NT 5.1) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.84 Safari/535.11 SE 2.X MetaSr 1.0",
	"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Trident/4.0; SV1; QQDownload 732; .NET4.0C; .NET4.0E; SE 2.X MetaSr 1.0)",
	// 夸克浏览器
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 QuarkPC/2.0.5.220",
	// 鸿蒙设备
	"Mozilla/5.0 (Phone; OpenHarmony 5.0; HarmonyOS 5.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 ArkWeb/4.1.6.1 Mobile",
}

func beautify(u *UserAgent) (s string) {
	if u.Browser.Name != "" {
		s += "Browser:" + u.Browser.Name
		if u.Browser.Version != "" {
			s += "-" + u.Browser.Version
		}
	}
	if u.Engine.Name != "" {
		s += " Engine:" + u.Engine.Name
		if u.Engine.Version != "" {
			s += "-" + u.Engine.Version
		}
	}
	if u.Device.Type != "" {
		s += " Device:" + u.Device.Type
		if u.Device.Model != "" {
			s += "-" + u.Device.Model
		}
		if u.Device.Vendor != "" {
			s += "-" + u.Device.Vendor
		}
	}
	if u.OS.Name != "" {
		s += " Os:" + u.OS.Name
		if u.OS.Version != "" {
			s += "-" + u.OS.Version
		}
	}
	if u.CPU.Architecture != "" {
		s += " Cpu:" + u.CPU.Architecture
	}
	return
}

func Test(t *testing.T) {
	alltestcases := []struct {
		title string
		label string
		list  []TestCase
	}{
		{
			title: "getBrowser()",
			label: "Browser",
			list:  loadJson("./data/ua/browser/browser-all.json"),
		},
		{
			title: "getCPU()",
			label: "Cpu",
			list:  loadJson("./data/ua/cpu/cpu-all.json"),
		},
		{
			title: "getDevice()",
			label: "Device",
			list:  loadJson("./data/ua/device/"),
		},
		{
			title: "getEngine()",
			label: "Engine",
			list:  loadJson("./data/ua/engine/engine-all.json"),
		},
		{
			title: "getOS()",
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
								singleCase.title, tc.Ua, key, val, actual)
						}
					}
				}
			}
		}
	}
}

func Test2(t *testing.T) {
	ua := "Mozilla/5.0 (Linux; Android 11; itel P651L Build/RP1A.201005.001) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.5672.76 Mobile Safari/537.36"
	r := NewUAParser(ua, nil, nil).Result()
	t.Logf("UA: %s\n%s\n", ua, r)
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
	defer f.Close()

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
