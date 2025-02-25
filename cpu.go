package uaparser

import (
	"strings"
)

func containsAny(s string, substrs []string) bool {
	for _, substr := range substrs {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

func (u *UserAgent) detectCPU(s section) {
	if u.CPU.Architecture != "" {
		return
	}

	comment := s.comment
	if containsAny(comment, []string{"amd64", "x64", "x86-64", "x86_64", "wow64", "win64"}) {
		u.CPU.Architecture = "amd64"
		return
	}

	if containsAny(comment, []string{"ia32"}) || containsAny(comment, []string{"i386", "i486", "i586", "i686", "x86;"}) {
		u.CPU.Architecture = "ia32"
		return
	}

	if containsAny(comment, []string{"aarch64", "arm64", "armv8", "armv8l", "arm_64"}) {
		u.CPU.Architecture = "arm64"
		return
	}

	if containsAny(comment, []string{"armv6", "armv7", "armhf", "armht", "armn", "armfl", "armfp"}) {
		u.CPU.Architecture = "armhf"
		return
	}

	if strings.Contains(comment, "windows ce") || strings.Contains(comment, "windows mobile") || strings.Contains(comment, "ppc") {
		u.CPU.Architecture = "arm"
		return
	}

	if containsAny(comment, []string{"ppc", "powerpc", "ppc64", "powerpc64"}) {
		u.CPU.Architecture = "ppc"
		return
	}

	if strings.Contains(comment, "sun4") {
		u.CPU.Architecture = "sparc"
		return
	}

	if containsAny(comment, []string{"avr32", "ia64", "68k", "armv", "atmel", "irix", "mips", "sparc", "pa-risc"}) {
		u.CPU.Architecture = strings.ToLower(comment)
		return
	}
}

//import (
//	"regexp"
//	"strings"
//)
//
//var (
//	amd64Regex = regexp.MustCompile(`(?i)\b(amd|x|x86[-_]?|wow|win)64\b`)
//	ia32Regex  = []*regexp.Regexp{
//		regexp.MustCompile(`(?i)(ia32(?=;))`),
//		regexp.MustCompile(`(?i)((?:i[346]|x)86)[;)]`),
//	}
//	arm64Regex = regexp.MustCompile(`(?i)\b(aarch64|arm(v?8e?l?|_?64))\b`)
//	armhfRegex = regexp.MustCompile(`(?i)\b(arm(?:v[67])?ht?n?[fl]p?)\b`)
//	armRegex   = regexp.MustCompile(`(?i)windows (ce|mobile); ppc;`)
//	ppcRegex   = regexp.MustCompile(`(?i)((?:ppc|powerpc)(?:64)?)(?: mac|;|\))`)
//	sparcRegex = regexp.MustCompile(`(?i)(sun4\w)[;)]`)
//	otherRegex = regexp.MustCompile(`(?i)((?:avr32|ia64(?=;))|68k(?=\))|\barm(?=v(?:[1-7]|[5-7]1)l?|;|eabi)|(?=atmel )avr|(?:irix|mips|sparc)(?:64)?\b|pa-risc)`)
//)
//
//func (u *UserAgent) detectCPU(s section) {
//	if u.Cpu.Architecture != "" {
//		return
//	}
//
//	part := s.String()
//	if amd64Regex.MatchString(part) {
//		u.Cpu.Architecture = "amd64"
//		return
//	}
//
//	for _, r := range ia32Regex {
//		if r.MatchString(part) {
//			u.Cpu.Architecture = "ia32"
//			return
//		}
//	}
//
//	if arm64Regex.MatchString(part) {
//		u.Cpu.Architecture = "arm64"
//		return
//	}
//
//	if armhfRegex.MatchString(part) {
//		u.Cpu.Architecture = "armhf"
//		return
//	}
//
//	if armRegex.MatchString(part) {
//		u.Cpu.Architecture = "arm"
//		return
//	}
//
//	if ppcRegex.MatchString(part) {
//		u.Cpu.Architecture = "ppc"
//		return
//	}
//
//	if sparcRegex.MatchString(part) {
//		u.Cpu.Architecture = "sparc"
//		return
//	}
//
//	if otherRegex.MatchString(part) {
//		m := otherRegex.FindStringSubmatch(part)
//		u.Cpu.Architecture = strings.ToLower(m[1])
//		return
//	}
//}
