package tg_template

import (
	"fmt"
	"regexp"
)

func ExtractReText(input string, reg string) string {
	re := regexp.MustCompile(reg)
	matches := re.FindStringSubmatch(input)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func ExtractReTextArr(input string, reg string) []string {
	re := regexp.MustCompile(reg)
	matches := re.FindAllString(input, -1)
	return matches
}

func DelReText(input string, reg string) string {
	re := regexp.MustCompile(reg)
	return re.ReplaceAllString(input, "")
}

func ReplaceReText(input string, reg string, replace any) string {
	re := regexp.MustCompile(reg)
	return re.ReplaceAllString(input, fmt.Sprintf("%v", replace))
}
