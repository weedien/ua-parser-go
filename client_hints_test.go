package uaparser

import "testing"

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
