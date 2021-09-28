package go_csv

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateColumnsNameWithSize(t *testing.T) {
	cols := generateColumnsName(6)
	assert.Equal(t, 6, len(cols))
	for i, col := range cols {
		assert.Equal(t, "COL_"+fmt.Sprint(i), col)
	}
}
func TestProcessEmptyAll(t *testing.T) {
	var cols []string
	var values []string
	result, err := process(cols, values, "", false, false)
	assert.Nil(t, err, "error must be nil")
	assert.Equal(t, 0, len(result), "It must be zero because no column and no value")
}
func TestProcessWithColsAndValuesSameSize(t *testing.T) {
	cols := []string{"name", "lastname", "age"}
	values := []string{"Puttipong", "Suwannapornkul", "25"}
	result, err := process(cols, values, "", false, false)
	assert.Nil(t, err, "error must be nil")
	assert.Equal(t, "Puttipong", result["name"])
	assert.Equal(t, "Suwannapornkul", result["lastname"])
	assert.Equal(t, "25", result["age"])
}
func TestProcessWithColsAndValuesDiffSize(t *testing.T) {
	cols := []string{"name", "status", "age"}
	values := []string{"Tidarut", "Be loved"}
	_, err := process(cols, values, "", false, false)
	assert.NotNil(t, err, "Columns and values not match should return error")
}
func TestReadCsvFileWithHeader(t *testing.T) {
	reader := strings.NewReader("id,firstname,lastname,email,email2,profession\n0,Genovera,Wildermuth,Genovera.Wildermuth@yopmail.com,Genovera.Wildermuth@gmail.com,worker")
	config := make(map[string]string)
	config[ENABLE_HEADER_KEY] = "true"
	csv := Make(config, func(m map[string]interface{}) {
		assert.Equal(t, "0", m["ID"])
		assert.Equal(t, "Genovera", m["FIRSTNAME"])
		assert.Equal(t, "Wildermuth", m["LASTNAME"])
		assert.Equal(t, "Genovera.Wildermuth@yopmail.com", m["EMAIL"])
		assert.Equal(t, "Genovera.Wildermuth@gmail.com", m["EMAIL2"])
		assert.Equal(t, "worker", m["PROFESSION"])
	}, func(s string) {
		t.Fatal("error function should not be called")
	})
	csv.ReadAndCall(reader)
}
func TestReadCsvFileWithHeaderAndCaseSensitive(t *testing.T) {
	reader := strings.NewReader("id,firstname,lastname,email,email2,profession\n0,Genovera,Wildermuth,Genovera.Wildermuth@yopmail.com,Genovera.Wildermuth@gmail.com,worker")
	config := make(map[string]string)
	config[ENABLE_HEADER_KEY] = "true"
	config[COLUMN_CASE_SENTITIVE_KEY] = "true"
	csv := Make(config, func(m map[string]interface{}) {
		assert.Equal(t, "0", m["id"])
		assert.Equal(t, "Genovera", m["firstname"])
		assert.Equal(t, "Wildermuth", m["lastname"])
		assert.Equal(t, "Genovera.Wildermuth@yopmail.com", m["email"])
		assert.Equal(t, "Genovera.Wildermuth@gmail.com", m["email2"])
		assert.Equal(t, "worker", m["profession"])
	}, func(s string) {
		t.Fatal("error function should not be called")
	})
	csv.ReadAndCall(reader)
}
func TestReadCsvFileWithAutoGenerateColumnName(t *testing.T) {
	reader := strings.NewReader("0,Genovera,Wildermuth,Genovera.Wildermuth@yopmail.com,Genovera.Wildermuth@gmail.com,worker")
	config := make(map[string]string)
	config[ENABLE_HEADER_KEY] = "false"
	csv := Make(config, func(m map[string]interface{}) {
		assert.Equal(t, "0", m["COL_0"])
		assert.Equal(t, "Genovera", m["COL_1"])
		assert.Equal(t, "Wildermuth", m["COL_2"])
		assert.Equal(t, "Genovera.Wildermuth@yopmail.com", m["COL_3"])
		assert.Equal(t, "Genovera.Wildermuth@gmail.com", m["COL_4"])
		assert.Equal(t, "worker", m["COL_5"])
	}, func(s string) {
		t.Fatal("error function should not be called")
	})
	csv.ReadAndCall(reader)
}
func TestReadCsvFileWithErrorMessageOnHaltIfError(t *testing.T) {
	// Error on line 3
	reader := strings.NewReader("id,firstname,lastname,email,email2,profession\n" +
		"0,Genovera,Wildermuth,Genovera.Wildermuth@yopmail.com,Genovera.Wildermuth@gmail.com,worker\n" +
		"1,Raina,Gibbeon,Raina.Gibbeon@yopmail.com,Raina.Gibbeon@gmail.com\n" +
		"2,Steffane,Codding,Steffane.Codding@yopmail.com,Steffane.Codding@gmail.com,firefighter\n")
	config := make(map[string]string)
	config[ENABLE_HEADER_KEY] = "true"
	config[COLUMN_CASE_SENTITIVE_KEY] = "true"
	config[HALT_IF_ERROR_KEY] = "true"
	count := 0
	csv := Make(config, func(m map[string]interface{}) {
		count += 1
	}, func(s string) {
		assert.Equal(t, 1, count)
	})
	panicF := func() {
		csv.ReadAndCall(reader)
	}
	require.Panics(t, panicF)
}
func TestReadCsvFileWithErrorMessageOnNotHaltIfError(t *testing.T) {
	// Error on line 3
	reader := strings.NewReader("id,firstname,lastname,email,email2,profession\n" +
		"0,Genovera,Wildermuth,Genovera.Wildermuth@yopmail.com,Genovera.Wildermuth@gmail.com,worker\n" +
		"1,Raina,Gibbeon,Raina.Gibbeon@yopmail.com,Raina.Gibbeon@gmail.com\n" +
		"2,Steffane,Codding,Steffane.Codding@yopmail.com,Steffane.Codding@gmail.com,firefighter\n")
	config := make(map[string]string)
	config[ENABLE_HEADER_KEY] = "true"
	config[COLUMN_CASE_SENTITIVE_KEY] = "true"
	config[HALT_IF_ERROR_KEY] = "false"
	successCount := 0
	failedCount := 0
	csv := Make(config, func(m map[string]interface{}) {
		successCount += 1
	}, func(s string) {
		failedCount += 1
		assert.Equal(t, "1,Raina,Gibbeon,Raina.Gibbeon@yopmail.com,Raina.Gibbeon@gmail.com", s)
	})
	csv.ReadAndCall(reader)
	assert.Equal(t, 2, successCount)
	assert.Equal(t, 1, failedCount)
}
func TestReadFileThenProcessWithoutError(t *testing.T) {
	path := filepath.Join("example/source", "data")
	config := make(map[string]string)
	config[ENABLE_HEADER_KEY] = "true"
	config[COLUMN_CASE_SENTITIVE_KEY] = "true"
	successCount := 0
	failedCount := 0
	csv := Make(config, func(m map[string]interface{}) {
		successCount += 1
	}, func(s string) {
		failedCount += 1
	})
	csv.ReadAndCallByFile(path)
	assert.Equal(t, 1000, successCount)
	assert.Equal(t, 0, failedCount)

}
func TestReadFileFailed(t *testing.T) {
	path := filepath.Join("example/source", "404")
	config := make(map[string]string)
	config[ENABLE_HEADER_KEY] = "true"
	config[COLUMN_CASE_SENTITIVE_KEY] = "true"
	successCount := 0
	failedCount := 0
	csv := Make(config, func(m map[string]interface{}) {
		successCount += 1
	}, func(s string) {
		failedCount += 1
	})
	panicF := func() {
		csv.ReadAndCallByFile(path)
	}
	require.Panics(t, panicF)
}
func TestReadFileWithEncodeNotFound(t *testing.T) {
	path := filepath.Join("example/source", "data")
	config := make(map[string]string)
	config[ENABLE_HEADER_KEY] = "true"
	config[COLUMN_CASE_SENTITIVE_KEY] = "true"
	config[FILE_ENCODE_KEY] = "UTF-99"

	successCount := 0
	failedCount := 0
	csv := Make(config, func(m map[string]interface{}) {
		successCount += 1
	}, func(s string) {
		failedCount += 1
	})
	panicF := func() {
		csv.ReadAndCallByFile(path)
	}
	require.Panics(t, panicF)
}
func TestProcessWithColsAndValuesAndRtrim(t *testing.T) {
	cols := []string{"name", "lastname", "age"}
	values := []string{" Puttipong   ", "  Suwannapornkul   ", "   25   "}
	result, err := process(cols, values, "", false, true)
	assert.Nil(t, err, "error must be nil")
	assert.Equal(t, " Puttipong", result["name"])
	assert.Equal(t, "  Suwannapornkul", result["lastname"])
	assert.Equal(t, "   25", result["age"])
}
func TestProcessWithColsAndValuesAndLtrim(t *testing.T) {
	cols := []string{"name", "lastname", "age"}
	values := []string{" Puttipong ", "  Suwannapornkul  ", "   25   "}
	result, err := process(cols, values, "", true, false)
	assert.Nil(t, err, "error must be nil")
	assert.Equal(t, "Puttipong ", result["name"])
	assert.Equal(t, "Suwannapornkul  ", result["lastname"])
	assert.Equal(t, "25   ", result["age"])
}
func TestProcessWithColsAndValuesAndBothRtrimLtrim(t *testing.T) {
	cols := []string{"name", "lastname", "age"}
	values := []string{" Puttipong ", "  Suwannapornkul  ", "   25   "}
	result, err := process(cols, values, "", true, true)
	assert.Nil(t, err, "error must be nil")
	assert.Equal(t, "Puttipong", result["name"])
	assert.Equal(t, "Suwannapornkul", result["lastname"])
	assert.Equal(t, "25", result["age"])
}
func TestProcessWithColsAndValuesCheckNull(t *testing.T) {
	cols := []string{"name", "lastname", "age"}
	values := []string{"Puttipong", "", "25"}
	result, err := process(cols, values, "", false, false)
	assert.Nil(t, err, "error must be nil")
	assert.Equal(t, "Puttipong", result["name"])
	assert.Equal(t, nil, result["lastname"])
	assert.Equal(t, "25", result["age"])
}
func TestProcessWithColsAndValuesTrimCheckNull(t *testing.T) {
	cols := []string{"name", "lastname", "age"}
	values := []string{" Puttipong ", "  ", " 25 "}
	result, err := process(cols, values, "", true, true)
	assert.Nil(t, err, "error must be nil")
	assert.Equal(t, "Puttipong", result["name"])
	assert.Equal(t, nil, result["lastname"])
	assert.Equal(t, "25", result["age"])
}
func TestProcessWithColsAndValuesTrimCheckNullWithChar(t *testing.T) {
	cols := []string{"name", "lastname", "age"}
	values := []string{" Puttipong ", " null ", " 25 "}
	result, err := process(cols, values, "null", true, true)
	assert.Nil(t, err, "error must be nil")
	assert.Equal(t, "Puttipong", result["name"])
	assert.Equal(t, nil, result["lastname"])
	assert.Equal(t, "25", result["age"])
}
