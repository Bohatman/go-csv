package go_csv

const NIL_CHARACTER_KEY = "parser.nil.character"
const NIL_CHARACTER_DEFAULT = ""
const NIL_CHARACTER_DOC = "If match this character value will turn to nil"

const SEPARATOR_CHARACTER_KEY = "parser.separator.character"
const SEPARATOR_CHARACTER_DEFAULT = ","
const SEPARATOR_CHARACTER_DOC = "Character for separate column"

//const SKIP_LINE_KEY = "parser.skip.line"
//const SKIP_LINE_DEFAULT = 0
//const SKIP_LINE_DOC = "Skip line that will not process"

const FILE_ENCODE_KEY = "file.encode"
const FILE_ENCODE_DEFAULT = "UTF-8"
const FILE_ENCODE_DOC = "File encode for reader"

const VALUE_LEFT_TRIM_KEY = "value.left.trim"
const VALUE_LEFT_TRIM_DEFAULT = false
const VALUE_LEFT_TRIM_DOC = "left trim value"

const VALUE_RIGHT_TRIM_KEY = "value.right.trim"
const VALUE_RIGHT_TRIM_DEFAULT = false
const VALUE_RIGHT_TRIM_DOC = "right trim value"

const ENABLE_HEADER_KEY = "parser.enable.header"
const ENABLE_HEADER_DEFAULT = false
const ENABLE_HEADER_DOC = "Enable header for column name"

const COLUMN_CASE_SENTITIVE_KEY = "column.case.sensitive"
const COLUMN_CASE_SENTITIVE_DEFAULT = false
const COLUMN_CASE_SENTITIVE_DOC = "If set to false all column will change to uppercase"

const COLUMN_SERIALIZE_KEY = "column.serialize"
const COLUMN_SERIALIZE_DOC = "If parser.enable.header set to false will use this to define column name"

const HALT_IF_ERROR_KEY = "parser.halt.if.error"
const HALT_IF_ERROR_DEFAULT = true
const HALT_IF_ERROR_DOC = "If set to true process will return error and stop process if set to false it will ignore and skip line"
