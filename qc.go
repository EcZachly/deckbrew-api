package main

import (
        "strings"
        "fmt"
)

// users := Select("*").From("users")
// active := users.Where(Eq{"deleted_at": nil})


const (
        PARAM = "$%d"
)

type Condition interface {
    ToSql() (string, []interface{}, error)
}

type Expression struct {
        Fields []string
        Tables []string
        Action string
        Condition Condition
}

type comparison struct {
        column string
        operator string
        val interface{}
}

func (c comparison) ToSql() (string, []interface{}, error) {
        sql := fmt.Sprintf("%s %s %s", c.column, c.operator, PARAM)
        return sql, []interface{}{c.val}, nil
}

func And(conds ...Condition) Condition {
        return comparison{}
}

func Or(conds ...Condition) Condition {
        return comparison{}
}

func Eq(column string, val interface{}) Condition {
        return comparison{column: column, operator: "=", val: val}
}

func Gt(column string, val interface{}) Condition {
        return comparison{}
}

func Lt(column string, val interface{}) Condition {
        return comparison{}
}

func Gte(column string, val interface{}) Condition {
        return comparison{}
}

func Lte(column string, val interface{}) Condition {
        return comparison{}
}

func (e Expression) From(table ...string) Expression {
        e.Tables = table
        return e
}

func (e Expression) Where(cond Condition) Expression {
        e.Condition = cond
        return e
}

func (e Expression) ToSql() (string, []interface{}, error) {
        params := []interface{}{}
        if len(e.Fields) == 0 {
                return "", params, fmt.Errorf("No fields in SQL expression")
        }

        if len(e.Tables) == 0 {
                return "", params, fmt.Errorf("No tables to query in SQL expression")
        }

        sql := e.Action + " " + strings.Join(e.Fields, ", ")
        sql = sql + " FROM " + strings.Join(e.Tables, ", ")

        if e.Condition != nil {
                csql, args, err := e.Condition.ToSql()

                if err != nil {
                    return sql, params, nil
                }

                sql = sql + " WHERE " + csql

                for _, arg := range args {
                    params = append(params, arg)
                }
        }

        counts := []interface{}{}

        for i, _ := range params {
                counts = append(counts, i + 1)
        }

        sql = fmt.Sprintf(sql, counts...)

        return sql, params, nil
}

func Select(fields ...string) Expression {
        return Expression{Action: "SELECT", Fields: fields}
}
