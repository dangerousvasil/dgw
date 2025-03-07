package main

import (
	"log"

	"github.com/BurntSushi/toml"
)

var defaultTypeMapCfg PgTypeMapConfig

func init() {
	if _, err := toml.Decode(typeMap, &defaultTypeMapCfg); err != nil {
		log.Fatal(err)
	}
}

const typeMap = `
[string]
db_types = ["character", "character varying", "text", "money"]
notnull_go_type = "string"
nullable_go_type = "sql.NullString"

[time]
db_types = [
    "time with time zone", "time without time zone",
    "timestamp without time zone", "timestamp with time zone", "date"
]
notnull_go_type = "time.Time"
nullable_go_type = "*time.Time"

[bool]
db_types = ["boolean"]
notnull_go_type = "bool"
nullable_go_type = "bool"

[smallint]
db_types = ["smallint"]
notnull_go_type = "int"
nullable_go_type = "sql.NullInt32"

[integer]
db_types = ["integer"]
notnull_go_type = "int64"
nullable_go_type = "sql.NullInt32"

[bigint]
db_types = ["bigint"]
notnull_go_type = "int64"
nullable_go_type = "sql.NullInt64"

[smallserial]
db_types = ["smallserial"]
notnull_go_type = "uint"
nullable_go_type = "sql.NullInt64"

[serial]
db_types = ["serial"]
notnull_go_type = "uint32"
nullable_go_type = "sql.NullInt64"

[real]
db_types = ["real"]
notnull_go_type = "float32"
nullable_go_type = "sql.NullFloat64"

[numeric]
db_types = ["numeric", "double precision"]
notnull_go_type = "float64"
nullable_go_type = "sql.NullFloat64"

[bytea]
db_types = ["bytea"]
notnull_go_type = "[]byte"
nullable_go_type = "[]byte"

[json]
db_types = ["json", "jsonb"]
notnull_go_type = "[]byte"
nullable_go_type = "[]byte"

[xml]
db_types = ["xml"]
notnull_go_type = "[]byte"
nullable_go_type = "[]byte"

[interval]
db_types = ["interval"]
notnull_go_type = "time.Duration"
nullable_go_type = "*time.Duration"

[bit]
db_types = ["bit"]
notnull_go_type = "string"
nullable_go_type = "sql.NullString"

[default]
db_types = ["*"]
notnull_go_type = "interface{}"
nullable_go_type = "interface{}"
`
