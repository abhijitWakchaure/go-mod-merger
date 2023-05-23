package semvar

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Compare will compare semvar of 2 modules and return the higher version
// The funcion supports following type of versions:
// v1.3.10
// v0.9.4-hf.1
// v0.9.1-0.20190603184501-d845e1d612f8
// v0.9.0-rc.1.0.20190509204259-4246269fb68e
func Compare(v1Str, v2Str string) (string, error) {
	if !strings.HasPrefix(v1Str, "v") || !strings.HasPrefix(v2Str, "v") {
		return "", fmt.Errorf("version does not follow semvar guidelines")
	}
	// fast check if versions are same
	if v1Str == v2Str {
		return v1Str, nil
	}
	v1Split := strings.SplitN(v1Str, ".", 3)
	v2Split := strings.SplitN(v2Str, ".", 3)
	// Compare major versions
	if v1Split[0] != v2Split[0] {
		return "", fmt.Errorf("major version mismatched")
	}
	// Compare minor versions
	if v1Split[1] != v2Split[1] {
		v1, _ := strconv.Atoi(v1Split[1])
		v2, _ := strconv.Atoi(v2Split[1])
		if v1 > v2 {
			return v1Str, nil
		}
		return v2Str, nil
	}
	// Minor versions are same, check patch versions
	// check if patch version is just a number
	p1 := getPatchDigit(v1Split[2])
	p2 := getPatchDigit(v2Split[2])
	if p1 == p2 {
		if len(v1Str) > len(v2Str) {
			return v1Str, nil
		}
		return v2Str, nil
	}
	if p1 > p2 {
		return v1Str, nil
	}
	return v2Str, nil
}

func getPatchDigit(pStr string) int {
	p, err := strconv.Atoi(pStr)
	if err == nil {
		return p
	}
	var index int
	for i, v := range pStr {
		if unicode.IsDigit(v) {
			continue
		}
		index = i
		break
	}
	p, _ = strconv.Atoi(pStr[:index])
	return p
}
