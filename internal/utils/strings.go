package utils

import "unsafe"

// BytesToString converts byte slice to string via Go 1.20+ optimal unsafeness.
func BytesToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return unsafe.String(&b[0], len(b))
}

// StringToBytes converts string to byte slice via Go 1.20+ optimal unsafeness.
func StringToBytes(s string) []byte {
	if len(s) == 0 {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

type StringSort []string

func (s StringSort) Len() int {
	return len(s)
}

func (s StringSort) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s StringSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
