package ngxnet

import (
	"strings"
)

func SplitStr(s string, sep string) []string {
	return strings.Split(s, sep)
}

func StrSplit(s string, sep string) []string {
	return strings.Split(s, sep)
}

func SplitStrN(s string, sep string, n int) []string {
	return strings.SplitN(s, sep, n)
}

func StrSplitN(s string, sep string, n int) []string {
	return strings.SplitN(s, sep, n)
}

func StrFind(s string, f string) int {
	return strings.Index(s, f)
}

func FindStr(s string, f string) int {
	return strings.Index(s, f)
}

func ReplaceStr(s, old, new string) string {
	return strings.Replace(s, old, new, -1)
}

func StrReplace(s, old, new string) string {
	return strings.Replace(s, old, new, -1)
}

func TrimStr(s string) string {
	return strings.TrimSpace(s)
}

func StrTrim(s string) string {
	return strings.TrimSpace(s)
}

func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func JoinStr(a []string, sep string) string {
	return strings.Join(a, sep)
}

func StrJoin(a []string, sep string) string {
	return strings.Join(a, sep)
}

func StrToLower(s string) string {
	return strings.ToLower(s)
}

func ToLowerStr(s string) string {
	return strings.ToLower(s)
}

func StrToUpper(s string) string {
	return strings.ToUpper(s)
}

func ToUpperStr(s string) string {
	return strings.ToUpper(s)
}
