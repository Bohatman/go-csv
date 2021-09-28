package go_csv

import (
	"bufio"
	"errors"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
	"os"
	"strings"
)

type CsvInfo struct {
	nilCharacter        string
	separatorCharacter  string
	fileEncode          string
	valueLeftTrim       bool
	valueRightTrim      bool
	enableHeader        bool
	columnCaseSensitive bool
	columnSerialize     []string
	haltIfError         bool
	funcToCall          func(map[string]string)
	funcToCallWhenError func(string)
}

func Make(propertiesPath string, toCall func(map[string]string)) *CsvInfo {
	prop, err := ReadProperties(propertiesPath)
	if err != nil {
		panic("error while reading properties file")
	}
	info := CsvInfo{}
	info.nilCharacter = getConfigString(prop, NIL_CHARACTER_KEY, NIL_CHARACTER_DEFAULT)
	info.separatorCharacter = getConfigString(prop, SEPARATOR_CHARACTER_KEY, SEPARATOR_CHARACTER_DEFAULT)
	info.fileEncode = getConfigString(prop, FILE_ENCODE_KEY, FILE_ENCODE_DEFAULT)
	info.valueLeftTrim = getConfigBool(prop, VALUE_LEFT_TRIM_KEY, VALUE_LEFT_TRIM_DEFAULT)
	info.valueRightTrim = getConfigBool(prop, VALUE_RIGHT_TRIM_KEY, VALUE_RIGHT_TRIM_DEFAULT)
	info.enableHeader = getConfigBool(prop, ENABLE_HEADER_KEY, ENABLE_HEADER_DEFAULT)
	info.columnCaseSensitive = getConfigBool(prop, COLUMN_CASE_SENTITIVE_KEY, COLUMN_CASE_SENTITIVE_DEFAULT)
	info.columnSerialize = getConfigStringArray(prop, COLUMN_SERIALIZE_KEY, "")
	info.haltIfError = getConfigBool(prop, HALT_IF_ERROR_KEY, HALT_IF_ERROR_DEFAULT)
	info.funcToCall = toCall
	return &info
}
func (info *CsvInfo) ReadAndCall(fileToRead string) {
	file, err := os.Open(fileToRead)
	if err != nil {
		panic("can not open file " + fileToRead)
	}
	encoding, _ := charset.Lookup(info.fileEncode)
	if encoding == nil {
		panic("can not find encode " + info.fileEncode)
	}
	reader := transform.NewReader(file, encoding.NewDecoder())
	scanner := bufio.NewScanner(reader)
	if info.enableHeader {
		scanner.Scan()
		headers := strings.Split(scanner.Text(), info.separatorCharacter)
		info.columnSerialize = processDetectHeader(headers, info.columnCaseSensitive)
	}
	for scanner.Scan() {
		raw := scanner.Text()
		if len(info.columnSerialize) == 0 {
			text := strings.Split(raw, info.separatorCharacter)
			info.columnSerialize = processAutoHeader(len(text))
		}
		if info.haltIfError {
			info.processHalt(raw)
		} else {
			info.processDontCare(raw)
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
func processAutoHeader(size int) []string {
	toReturn := make([]string, size)
	for i := 0; i < size; i++ {
		toReturn[i] = "COL_" + string(i)
	}
	return toReturn
}
func processDetectHeader(headers []string, isCaseSensitive bool) []string {
	toReturn := make([]string, len(headers))
	for i, header := range headers {
		tmp := strings.Trim(header, " ")
		if isCaseSensitive {
			tmp = strings.ToUpper(tmp)
		}
		toReturn[i] = tmp
	}
	return toReturn
}
func (info *CsvInfo) processHalt(raw string) {
	text := strings.Split(raw, info.separatorCharacter)
	result, err := process(info.columnSerialize, text)
	if err != nil {
		info.funcToCallWhenError(raw)
		panic(err)
	}
	info.funcToCall(result)
}
func (info *CsvInfo) processDontCare(raw string) {
	text := strings.Split(raw, info.separatorCharacter)
	result, err := process(info.columnSerialize, text)
	if err != nil {
		info.funcToCallWhenError(raw)
	} else {
		info.funcToCall(result)
	}
}
func process(columns []string, values []string) (map[string]string, error) {
	if len(columns) != len(values) {
		return nil, errors.New("columns(" + string(len(columns)) + ") not match to values(" + string(len(values)) + ")")
	}
	toReturn := make(map[string]string)
	for i := 0; i < len(columns); i++ {
		toReturn[columns[i]] = values[i]
	}
	return toReturn, nil
}
