package go_csv

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

func ReadProperties(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	toReturn := make(map[string]string)
	count := 0
	for scanner.Scan() {
		count += 1
		line := strings.TrimLeft(scanner.Text(), " ")
		// If length of line is zero or start with # will skip
		if len(line) == 0 || string(line[0]) == "#" {
			continue
		}
		separatorIndex := strings.Index(line, "=")
		// If can not find "=" or len invalid will return error
		if separatorIndex <= 0 || len(line) < separatorIndex+1 {
			return nil, errors.New("line is invalid format" + " (line " + strconv.Itoa(count) + ")" + ": " + line)
		}
		key := strings.Trim(line[0:separatorIndex], " ")
		value := line[separatorIndex+1:]
		toReturn[key] = value
	}
	return toReturn, nil
}

func getConfigString(prop map[string]string, key string, defaultValue string) string {
	var result string
	if value, ok := prop[key]; ok {
		result = value
	} else {
		result = defaultValue
	}
	return result
}
func getConfigStringArray(prop map[string]string, key string, defaultValue []string) []string {
	var result []string
	if value, ok := prop[key]; ok {
		result = strings.Split(value, ",")
	} else {
		result = defaultValue
	}
	return result
}
func getConfigBool(prop map[string]string, key string, defaultValue bool) bool {
	var result bool
	if value, ok := prop[key]; ok {
		tmp, err := strconv.ParseBool(value)
		if err != nil {
			panic(key + " can not parse to boolean(" + value + ")")
		}
		result = tmp
	} else {
		result = defaultValue
	}
	return result
}
