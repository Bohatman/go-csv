## Simple Golang CSV(comma-separated values) reader

### Supported Go version

Go-csv is available as a Go module.For 1.14+ 

### Feature Overview

1. Config and run for read CSV file.
2. Simple code
3. Error handle as your choice

### Guide

#### Installation

> go get github.com/Bohatman/go-cvs

#### Example

```go
package main

import (
	"fmt"
	"github.com/Bohatman/go-csv"
)
func main() {
	prop,_ := go_csv.ReadProperties("example/prop.properties")
	prop[go_csv.ENABLE_HEADER_KEY] = "true"
	info:= go_csv.Make(prop,process,processBadMessage)
	info.ReadAndCallByFile("example/data")
}
func process(row map[string]interface{}){
	fmt.Println(row)
}
func processBadMessage(row string){
	fmt.Println( "Oops error row detect: " + row)
}
```

### Properties

1. parser.enable.header
    - ***Possible value:*** true,false
    - ***Description:*** If set to "true" the first row will be determined as column name else will auto generated column name as prefix "COL_"
    - ***Default value:*** false
    
2. column.case.sensitive
    - ***Possible value:*** true,false
    - ***Description:*** If set to false all column will change to uppercase
    - ***Default value:*** false    

3. column.serialize
    - ***Possible value:*** Any string
    - ***Example:*** field_1,filed_2,field_3
    - ***Description:*** If parser.enable.header set to false will use this to define column name , if set this column.case.sensitive flag will be ignored.
        

4. parser.halt.if.error
    - ***Possible value:*** true,false
    - ***Description:*** If set to true process will return error and stop process if set to false it will ignore and skip line
    - ***Default value:*** true
    
5. value.left.trim
   - ***Possible value:*** true,false
   - ***Description:*** left trim value
   - ***Default value:*** false    

6. value.right.trim
    - ***Possible value:*** true,false
    - ***Description:*** right trim value
    - ***Default value:*** false
    
7. file.encode
    - ***Description:*** File encode for reader
    - ***Default value:*** UTF-8
    
8. parser.nil.character
    - ***Possible value:*** Any string
    - ***Description:*** If match this character value will turn to nil
    - ***Default value:*** (EMPTY STRING)

9. parser.separator.character
    - ***Possible value:*** Any string
    - ***Description:*** Character for separate column
    - ***Default value:*** ,
    
    