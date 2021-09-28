package go_csv

import (
	"testing"
)

func TestReadProperties(t *testing.T) {
	exampleConfig := "D:\\dev\\go-csv\\example\\config\\config.properties"
	result, err := ReadProperties(exampleConfig)
	if err != nil {
		t.Log(err)
	}
	t.Log(result)
}
