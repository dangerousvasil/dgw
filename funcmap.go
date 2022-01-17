package main

import (
	"fmt"
	"strconv"
	"text/template"
)

var tmplFuncMap = template.FuncMap{
	"createMapFields":            createMapFields,
	"createMapIndexFields":       createMapIndexFields,
	"createMapDataFields":        createMapDataFields,
	"createUpdateSQL":            createUpdateSQL,
	"createDeleteSQL":            createDeleteSQL,
	"createDeletePkSQL":          createDeletePkSQL,
	"createInsertSQL":            createInsertSQL,
	"createInsertParams":         createInsertParams,
	"createInsertScan":           createInsertScan,
	"createCount":                createCount,
	"createSelect":               createSelect,
	"createFrom":                 createFrom,
	"createSelectByPkSQL":        createSelectByPkSQL,
	"createSelectByPkFuncParams": createSelectByPkFuncParams,
	"criteriaAddPkCriterion":     criteriaAddPkCriterion,
	"createSelectByPkSQLParams":  createSelectByPkSQLParams,
	"createSelectByPkScan":       createSelectByPkScan,
	"createNamedParams":          createNamedParams,
	"createNamedPkParams":        createNamedPkParams,
	"createSelectByPkSQLWHERE":   createSelectByPkSQLWHERE,

	"snippetCheckOrder":       snippetCheckOrder,
	"snippetBuildOrder":       snippetBuildOrder,
	"snippetCheckCriteria":    snippetCheckCriteria,
	"snippetBuildCriteria":    snippetBuildCriteria,
	"snippetBuildLimitOffset": snippetBuildLimitOffset,
	"snippetBuildSqlSelect":   snippetBuildSqlSelect,
}

func snippetBuildLimitOffset() string {
	return `if t.Limit.GetLimit() > 0 {
       limitOffset += " LIMIT " + strconv.FormatInt(t.Limit.GetLimit(), 10)
    }
    if t.Limit.GetOffset() > 0 {
       limitOffset += " OFFSET " + strconv.FormatInt(t.Limit.GetOffset(), 10)
    }
`
}
func snippetBuildOrder() string {
	return `if t.Order.GetName() != "" {
        orderBy = " ORDER BY " + t.Order.GetName()
        if t.Order.IsDesc() {
            orderBy += " DESC "
        } else {
            orderBy += " ASC "
        }
    }
`
}
func snippetCheckOrder() string {
	return `if  _, exist := t.columns[t.Order.GetName()]; t.Order.GetName() != "" && !exist {
  		err = errors.New("order column not exist: " + t.Order.GetName())
  	}`
}
func snippetCheckCriteria() string {
	return `for _, crit := range t.Criteria.GetCriteria() {
        if _, exist := t.columns[crit.GetField()]; !exist {
            err = errors.New("Criteria not exist: " + crit.GetName()+":"+crit.GetField())
            break
        }
    }`
}

func snippetBuildCriteria() string {
	return `for _, crit := range t.Criteria.GetCriteria() {
        if crit.GetRaw() != "" {
			sqlWhere = append(sqlWhere, crit.GetField() + crit.GetOperand() + crit.GetRaw())
		} else {
			sqlWhere = append(sqlWhere, crit.GetField() + crit.GetOperand() + ":" + crit.GetName())
        }
		filterWhere[crit.GetName()] = crit.GetValue()
    }`
}

func snippetBuildSqlSelect(st *Struct) string {
	return `count := len(t.iColumn) - 1
        for i := 0; i <= count; i++ {
                name := t.iColumn[i]
                if val, ok := t.Select[name]; ok && val != "" {
                        sqlSelect += val
                } else {
                        sqlSelect += name
                }
                if i < count {
                        sqlSelect += ", "
                }
        }`
}

func createMapFields(st *Struct) string {
	var colMap []string
	for _, c := range st.Table.Columns {

		colMap = append(colMap, `"`+c.Name+`": true`)
	}
	return `map[string]bool{` + flatten(colMap, ", ") + `}`
}

func createMapIndexFields(st *Struct) string {
	var colMap []string
	for i, c := range st.Table.Columns {

		colMap = append(colMap, strconv.Itoa(i)+`:"`+c.Name+`"`)
	}
	return `map[int]string{` + flatten(colMap, ", ") + `}`
}

func createMapDataFields(st *Struct) string {
	var colMap []string
	var key int
	for _, c := range st.Table.Columns {
		if !c.IsPrimaryKey {
			colMap = append(colMap, strconv.Itoa(key)+`:"`+c.Name+`"`)
			key++
		}
	}
	return `map[int]string{` + flatten(colMap, ", ") + `}`
}

func createCount(st *Struct) string {
	var sql string
	sql = "SELECT COUNT(*) FROM " + st.Table.Schema + "." + st.Table.Name
	return sql
}
func createFrom(st *Struct) string {
	var sql string
	sql = " FROM " + st.Table.Schema + "." + st.Table.Name
	return sql
}

func createSelect(st *Struct) string {
	var sql string
	var colNames []string
	var pkNames []string
	for _, c := range st.Table.Columns {
		if c.IsPrimaryKey {
			pkNames = append(pkNames, c.Name)
		}
		colNames = append(colNames, c.Name)
	}
	sql = "SELECT " + flatten(colNames, ", ") + " FROM " + st.Table.Schema + "." + st.Table.Name
	return sql
}

func createSelectByPkSQL(st *Struct) string {
	var sql string
	var colNames []string
	var pkNames []string
	for _, c := range st.Table.Columns {
		if c.IsPrimaryKey {
			pkNames = append(pkNames, c.Name)
		}
		colNames = append(colNames, c.Name)
	}
	sql = "SELECT " + flatten(colNames, ", ") + " FROM " + st.Table.Schema + "." + st.Table.Name + " WHERE "
	for i, c := range pkNames {
		placeHolder := i + 1
		if i == 0 {
			sql = sql + c + fmt.Sprintf(" = $%d", placeHolder)
		} else {
			sql = sql + " AND " + c + fmt.Sprintf(" = $%d", placeHolder)
		}
	}
	return sql
}

func createSelectByPkSQLWHERE(st *Struct) string {
	var sql string
	var colNames []string
	var pkNames []string
	for _, c := range st.Table.Columns {
		if c.IsPrimaryKey {
			pkNames = append(pkNames, c.Name)
		}
		colNames = append(colNames, c.Name)
	}
	sql = " WHERE "
	for i, c := range pkNames {
		placeHolder := i + 1
		if i == 0 {
			sql = sql + c + fmt.Sprintf(" = $%d", placeHolder)
		} else {
			sql = sql + " AND " + c + fmt.Sprintf(" = $%d", placeHolder)
		}
	}
	return sql
}

func createSelectByPkScan(st *Struct) string {
	var s []string
	for _, f := range st.Fields {
		s = append(s, fmt.Sprintf("&r.%s", f.Name))
	}
	return flatten(s, ", ")
}

func createSelectByPkSQLParams(st *Struct) string {
	var fs []string
	for i, f := range st.Fields {
		if f.Column.IsPrimaryKey {
			fs = append(fs, fmt.Sprintf("pk%d", i))
		}
	}
	return flatten(fs, ", ")
}

func createSelectByPkFuncParams(st *Struct) string {
	var fs []string
	for i, f := range st.Fields {
		if f.Column.IsPrimaryKey {
			fs = append(fs, fmt.Sprintf("pk%d ", i)+f.Type)
		}
	}
	return flatten(fs, ", ")
}
func criteriaAddPkCriterion(st *Struct) string {
	var fs []string
	for i, f := range st.Fields {
		if f.Column.IsPrimaryKey {
			fs = append(fs, fmt.Sprintf(`criteria.AddCriterion(NewCriterion("%s","%s","=",pk%d))`, f.Column.Name, f.Column.Name, i))
		}
	}
	return flatten(fs, "\r\n")
}

func createNamedParams(st *Struct) string {
	var fs []string
	for _, f := range st.Fields {
		fs = append(fs, `"`+f.Column.Name+`": r.`+f.Name)
	}
	return `map[string]interface{}{` + flatten(fs, ", ") + `}`
}

func createNamedPkParams(st *Struct) string {
	var fs []string
	for _, f := range st.Fields {
		if f.Column.IsPrimaryKey {
			fs = append(fs, `"`+f.Column.Name+`": r.`+f.Name)
		}
	}
	return `map[string]interface{}{` + flatten(fs, ", ") + `}`
}

func createInsertScan(st *Struct) string {
	var fs []string
	for _, f := range st.Fields {
		if f.Column.IsPrimaryKey && st.Table.AutoGenPk {
			fs = append(fs, "&r."+f.Name)
		}
	}
	return flatten(fs, ", ")
}

func createInsertParams(st *Struct) string {
	var fs []string
	for _, f := range st.Fields {
		if f.Column.IsPrimaryKey && st.Table.AutoGenPk {
			continue
		} else {
			fs = append(fs, "&r."+f.Name)
		}
	}
	return flatten(fs, ", ")
}

func flatten(elems []string, sep string) string {
	var str string
	for i, e := range elems {
		if i == 0 {
			str = str + e
		} else {
			str = str + sep + e
		}
	}
	return str
}

func placeholders(l []string) string {
	var ph string
	var j int
	for i := range l {
		j = i + 1
		if i == 0 {
			ph = ph + fmt.Sprintf("$%d", j)
		} else {
			ph = ph + fmt.Sprintf(", $%d", j)
		}
	}
	return ph
}

func createInsertSQL(st *Struct) string {
	var sql string
	sql = "INSERT INTO " + st.Table.Schema + "." + st.Table.Name + " ("

	if len(st.Table.Columns) == 1 && st.Table.Columns[0].IsPrimaryKey && st.Table.AutoGenPk {
		sql = sql + st.Table.Columns[0].Name + ") VALUES (DEFAULT)"
	} else {
		var colNames []string
		var fieldNames []string
		for _, c := range st.Table.Columns {
			if c.IsPrimaryKey && st.Table.AutoGenPk {
				continue
			} else {
				colNames = append(colNames, c.Name)
				fieldNames = append(fieldNames, `t.BuildFieldValue("`+c.Name+`")`)
			}
		}
		sql = sql + flatten(colNames, ", ") + ") VALUES (`+" + flatten(fieldNames, "+`, `+") + "+`)"
	}
	sql = sql + " RETURNING *"
	return sql
}

func createUpdateSQL(st *Struct) string {
	var sql string
	sql = "UPDATE " + st.Table.Schema + "." + st.Table.Name + " SET "

	if len(st.Table.PrimaryKeys) < 1 {
		return ""
	} else {
		var colNames []string
		var fieldNames []string
		for _, c := range st.Table.Columns {
			if c.IsPrimaryKey {
				fieldNames = append(fieldNames, c.Name+"=:"+c.Name)
			} else {
				colNames = append(colNames, c.Name+"=`+t.BuildFieldValue(`"+c.Name+"`)+`")
			}
		}
		sql = sql + flatten(colNames, ", ") + " WHERE " + flatten(fieldNames, "AND ")
	}

	return sql
}

func createDeletePkSQL(st *Struct) string {
	var sql string
	sql = "DELETE FROM " + st.Table.Schema + "." + st.Table.Name

	if len(st.Table.PrimaryKeys) < 1 {
		return ""
	} else {
		var fieldNames []string
		for _, c := range st.Table.Columns {
			if c.IsPrimaryKey {
				fieldNames = append(fieldNames, c.Name+"=:"+c.Name)
			}
		}
		sql = sql + " WHERE " + flatten(fieldNames, "AND ")
	}

	return sql
}

func createDeleteSQL(st *Struct) string {
	var sql string
	sql = "DELETE FROM " + st.Table.Schema + "." + st.Table.Name
	return sql
}
