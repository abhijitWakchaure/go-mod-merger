package semvar

import (
	"fmt"
	"strconv"
	"strings"
	"time"
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
	if p1 > p2 {
		return v1Str, nil
	} else if p2 < p1 {
		return v2Str, nil
	}
	// Patch digit is same, need to compare pseudo-version numbers
	// Format: baseVersionPrefix-timestamp-revisionIdentifier
	// v0.0.0-20190607162005-6e2aefe19808
	// v0.0.0-20190319203104-747b5cae041f
	// v0.9.0-rc.1.0.20190509204259-4246269fb68e
	// v0.9.0-beta.2

	// Check if patch version is NOT a standard pseudo-version numeber using it's length i.e. 28
	if len(v1Split[2]) != 28+len(fmt.Sprintf("%d", p1)) {
		return "", fmt.Errorf("version [%s] is not standard pseudo-version number. Human intervention is required", v1Str)
	}

	if len(v2Split[2]) != 28+len(fmt.Sprintf("%d", p2)) {
		return "", fmt.Errorf("version [%s] is not standard pseudo-version number. Human intervention is required", v2Str)
	}

	p1Timestamp := strings.Split(v1Split[2], "-")[1]
	p2Timestamp := strings.Split(v2Split[2], "-")[1]

	var p1Time, p2Time time.Time
	var err error
	p1Time, err = convertTimestamp(p1Timestamp)
	if err != nil {
		return "", fmt.Errorf("failed to parse timestamp [%s] for version %s due to %s", p1Timestamp, v1Str, err.Error())
	}
	p2Time, err = convertTimestamp(p2Timestamp)
	if err != nil {
		return "", fmt.Errorf("failed to parse timestamp [%s] for version %s due to %s", p2Timestamp, v2Str, err.Error())
	}
	if p1Time.Equal(p2Time) {
		return "", fmt.Errorf("cannot choose between versions [%s] and [%s] as both have same timestamp in pseudo-version number. Human intervention is required", v1Str, v2Str)
	}
	if p1Time.Before(p2Time) {
		return v2Str, nil
	}
	return v1Str, nil
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

func convertTimestamp(timestamp string) (time.Time, error) {
	// Format: yyyymmddhhmmss
	layout := "20060102150405"
	return time.Parse(layout, timestamp)
}
