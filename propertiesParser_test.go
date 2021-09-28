package go_csv

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func TestReadFileProperties(t *testing.T) {
	path := filepath.Join("example/config", "config.properties")
	result, err := ReadProperties(path)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "false", result[ENABLE_HEADER_KEY])
	assert.Equal(t, "", result[NIL_CHARACTER_KEY])
	assert.Equal(t, ",", result[SEPARATOR_CHARACTER_KEY])
	assert.Equal(t, "UTF-8", result[FILE_ENCODE_KEY])
	assert.Equal(t, "true", result[VALUE_LEFT_TRIM_KEY])
	assert.Equal(t, "true", result[VALUE_RIGHT_TRIM_KEY])
	assert.Equal(t, "false", result[COLUMN_CASE_SENTITIVE_KEY])
	assert.Equal(t, "COL_A,COL_B,COL_C", result[COLUMN_SERIALIZE_KEY])

}
func TestReadUnResolveFileProperties(t *testing.T) {
	path := filepath.Join("example/config", "404.properties")
	_, err := ReadProperties(path)
	if err != nil {
		assert.Error(t, err)
	}
}
func TestReadFilePropertiesWithMalformedFormat(t *testing.T) {
	path := filepath.Join("example/config", "error.properties")
	_, err := ReadProperties(path)
	if err != nil {
		t.Log(err)
		assert.Error(t, err)
	}
}
func TestGetConfigStringToReturnDefaultValue(t *testing.T) {
	empty := make(map[string]string)
	defaultValue := "Hi I'm Mark"
	assert.Equal(t, defaultValue, getConfigString(empty, "404", defaultValue))
}
func TestGetConfigStringToReturnMapValue(t *testing.T) {
	mapProp := make(map[string]string)
	value := "Oh Page not found maybe rerun again"
	mapProp["404"] = value
	assert.Equal(t, value, getConfigString(mapProp, "404", ""))
}
func TestGetConfigStringArrayToReturnDefaultValue(t *testing.T) {
	empty := make(map[string]string)
	defaultValue := []string{"I'm", "hungry", "I write", "this", "code", "at", "00:00"}
	arrayToCheck := getConfigStringArray(empty, "404", defaultValue)
	assert.Equal(t, defaultValue, arrayToCheck)
}
func TestGetConfigStringArrayToReturnMapValue(t *testing.T) {
	mapProp := make(map[string]string)
	value := "COFFEE,PAY,ME"
	mapProp["404"] = value
	arrayToCheck := getConfigStringArray(mapProp, "404", make([]string, 0))
	assert.Equal(t, 3, len(arrayToCheck))
	assert.Equal(t, "COFFEE", arrayToCheck[0])
	assert.Equal(t, "PAY", arrayToCheck[1])
	assert.Equal(t, "ME", arrayToCheck[2])
}
func TestGetConfigBoolToReturnDefaultValue(t *testing.T) {
	empty := make(map[string]string)
	assert.Equal(t, true, getConfigBool(empty, "404", true))
}
func TestGetConfigBoolToReturnMapValue(t *testing.T) {
	mapProp := make(map[string]string)
	mapProp["404"] = "false"
	assert.Equal(t, false, getConfigBool(mapProp, "404", true))
}
func TestGetConfigBoolThatInvalidFormat(t *testing.T) {
	mapProp := make(map[string]string)
	mapProp["404"] = "Hello"
	defaultValue := true
	panicF := func() {
		getConfigBool(mapProp, "404", defaultValue)
	}
	require.Panics(t, panicF)
}
