package go_csv

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
	"io"
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
	funcToCall          func(map[string]interface{})
	funcToCallWhenError func(string)
}

func Make(prop map[string]string, toCall func(map[string]interface{}), toCallWhenError func(string)) *CsvInfo {

	info := CsvInfo{}
	info.nilCharacter = getConfigString(prop, NIL_CHARACTER_KEY, NIL_CHARACTER_DEFAULT)
	info.separatorCharacter = getConfigString(prop, SEPARATOR_CHARACTER_KEY, SEPARATOR_CHARACTER_DEFAULT)
	info.fileEncode = getConfigString(prop, FILE_ENCODE_KEY, FILE_ENCODE_DEFAULT)
	info.valueLeftTrim = getConfigBool(prop, VALUE_LEFT_TRIM_KEY, VALUE_LEFT_TRIM_DEFAULT)
	info.valueRightTrim = getConfigBool(prop, VALUE_RIGHT_TRIM_KEY, VALUE_RIGHT_TRIM_DEFAULT)
	info.enableHeader = getConfigBool(prop, ENABLE_HEADER_KEY, ENABLE_HEADER_DEFAULT)
	info.columnCaseSensitive = getConfigBool(prop, COLUMN_CASE_SENTITIVE_KEY, COLUMN_CASE_SENTITIVE_DEFAULT)
	info.columnSerialize = getConfigStringArray(prop, COLUMN_SERIALIZE_KEY, make([]string, 0))
	info.haltIfError = getConfigBool(prop, HALT_IF_ERROR_KEY, HALT_IF_ERROR_DEFAULT)
	info.funcToCallWhenError = toCallWhenError
	info.funcToCall = toCall
	return &info
}
func (info *CsvInfo) ReadAndCallByFile(fileToRead string) {
	handle, err := os.Open(fileToRead)
	if err != nil {
		panic("can not open file " + fileToRead)
	}
	defer handle.Close()
	info.ReadAndCall(handle)
}
func (info *CsvInfo) ReadAndCall(handle io.Reader) {
	encoding, _ := charset.Lookup(info.fileEncode)
	if encoding == nil {
		panic("can not find encode " + info.fileEncode)
	}
	reader := transform.NewReader(handle, encoding.NewDecoder())
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
			info.columnSerialize = generateColumnsName(len(text))
		}
		info.process(raw)
	}
}
func generateColumnsName(size int) []string {
	toReturn := make([]string, size)
	for i := 0; i < size; i++ {
		toReturn[i] = "COL_" + fmt.Sprint(i)
	}
	return toReturn
}
func processDetectHeader(headers []string, isCaseSensitive bool) []string {
	toReturn := make([]string, len(headers))
	for i, header := range headers {
		tmp := strings.Trim(header, " ")
		if !isCaseSensitive {
			tmp = strings.ToUpper(tmp)
		}
		toReturn[i] = tmp
	}
	return toReturn
}
func (info *CsvInfo) process(raw string) {
	text := strings.Split(raw, info.separatorCharacter)
	result, err := process(info.columnSerialize, text, info.nilCharacter, info.valueLeftTrim, info.valueRightTrim)
	if err != nil {
		info.funcToCallWhenError(raw)
		if info.haltIfError {
			panic(err)
		}
	} else {
		info.funcToCall(result)
	}
}
func process(columns []string, values []string, nilChar string, isLtrim bool, isRtrim bool) (map[string]interface{}, error) {
	if len(columns) != len(values) {
		return nil, errors.New("columns(" + fmt.Sprint(len(columns)) + ") not match to values(" + fmt.Sprint(len(values)) + ")")
	}
	toReturn := make(map[string]interface{})
	for i := 0; i < len(columns); i++ {
		toSet := trimValue(values[i], isLtrim, isRtrim)
		toReturn[columns[i]] = checkValueNull(toSet, nilChar)
	}
	return toReturn, nil
}
func checkValueNull(value string, nilChar string) interface{} {
	if value == nilChar {
		return nil
	} else {
		return value
	}
}
func trimValue(data string, isLtrim bool, isRtrim bool) string {
	if isLtrim {
		data = strings.TrimLeft(data, " ")
	}
	if isRtrim {
		data = strings.TrimRight(data, " ")
	}
	return data
}
