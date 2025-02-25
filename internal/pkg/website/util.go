package website

import (
	"fmt"
	"strconv"
	"strings"
)

func getYear(href string) (string, error) {
	parts := strings.Split(href, "/")
	if len(parts) != 4 {
		return "", fmt.Errorf("Invalid href year %s", href)
	}

	return parts[2], nil
}

func incrementYear(year string) (string, error) {
	parts := strings.Split(year, "-")
	if len(parts) != 2 {
		return "", fmt.Errorf("Invalid year %s", year)
	}

	start, err1 := strconv.Atoi(parts[0])
	end, err2 := strconv.Atoi(parts[1])

	if err1 != nil || err2 != nil {
		return "", fmt.Errorf("Invalid year %s", year)
	}

	return fmt.Sprintf("%d-%d", start+1, end+1), nil
}
