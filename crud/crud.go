package crud

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type Crud struct {
	db     *sql.DB
	schema string
}

func NewCrud(db *sql.DB, schema string) *Crud {
	return &Crud{db: db, schema: schema}
}

func StructToMap(input any) map[string]any {
	result := make(map[string]any)
	v := reflect.ValueOf(input)
	t := reflect.TypeOf(input)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		tag := field.Tag.Get("db")
		if tag == "" {
			continue
		}

		parts := strings.Split(tag, ",")
		column := parts[0]

		for _, p := range parts[1:] {
			if p == "pk" {
				goto skip
			}
		}

		result[column] = value.Interface()
	skip:
	}

	return result
}

func (c *Crud) CreateStruct(table string, model any) error {
	v := reflect.ValueOf(model)
	t := reflect.TypeOf(model)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	var columns []string
	var values []string
	var args []any

	argIndex := 1

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		tag := field.Tag.Get("db")
		if tag == "" {
			continue
		}

		parts := strings.Split(tag, ",")
		column := parts[0]

		var isPK bool
		var seq string

		for _, p := range parts[1:] {
			if p == "pk" {
				isPK = true
			}
			if strings.HasPrefix(p, "seq=") {
				seq = strings.TrimPrefix(p, "seq=")
			}
		}

		// 🔑 PK com sequence
		if isPK && seq != "" {
			columns = append(columns, column)
			values = append(values, seq+".NEXTVAL")
			continue
		}

		// ignora PK sem sequence
		if isPK {
			continue
		}

		columns = append(columns, column)
		values = append(values, fmt.Sprintf(":%d", argIndex))
		args = append(args, value.Interface())
		argIndex++
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		table,
		strings.Join(columns, ", "),
		strings.Join(values, ", "),
	)

	_, err := c.db.Exec(query, args...)
	return err
}

func (c *Crud) CreateStructReturningID(table string, model any, idDest *int64) error {
	v := reflect.ValueOf(model)
	t := reflect.TypeOf(model)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	var columns []string
	var values []string
	var args []any

	argIndex := 1
	var hasReturning bool

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		tag := field.Tag.Get("db")
		if tag == "" {
			continue
		}

		parts := strings.Split(tag, ",")
		column := parts[0]

		var isPK bool
		var seq string

		for _, p := range parts[1:] {
			if p == "pk" {
				isPK = true
			}
			if strings.HasPrefix(p, "seq=") {
				seq = strings.TrimPrefix(p, "seq=")
			}
		}

		// PK com sequence
		if isPK && seq != "" {
			columns = append(columns, column)
			values = append(values, seq+".NEXTVAL")
			hasReturning = true
			continue
		}

		if isPK {
			continue
		}

		columns = append(columns, column)
		values = append(values, fmt.Sprintf(":%d", argIndex))
		args = append(args, value.Interface())
		argIndex++
	}

	// OUT parameter
	if hasReturning {
		args = append(args, sql.Out{Dest: idDest})
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING ID INTO :%d",
		table,
		strings.Join(columns, ", "),
		strings.Join(values, ", "),
		argIndex,
	)

	_, err := c.db.Exec(query, args...)
	return err
}

func (c *Crud) UpdateStruct(table string, model any) error {
	v := reflect.ValueOf(model)
	t := reflect.TypeOf(model)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	var sets []string
	var args []any

	var pkColumn string
	var pkValue any
	argIndex := 1

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		tag := field.Tag.Get("db")
		if tag == "" {
			continue
		}

		parts := strings.Split(tag, ",")
		column := parts[0]

		var isPK bool
		for _, p := range parts[1:] {
			if p == "pk" {
				isPK = true
			}
		}

		if isPK {
			pkColumn = column
			pkValue = value.Interface()
			continue
		}

		sets = append(sets, fmt.Sprintf("%s = :%d", column, argIndex))
		args = append(args, value.Interface())
		argIndex++
	}

	if pkColumn == "" {
		return fmt.Errorf("pk não encontrada no model")
	}

	args = append(args, pkValue)

	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s = :%d",
		table,
		strings.Join(sets, ", "),
		pkColumn,
		argIndex,
	)

	_, err := c.db.Exec(query, args...)
	return err
}

func (c *Crud) DeleteByID(table string, pkColumn string, id any) error {
	query := fmt.Sprintf(
		"DELETE FROM %s WHERE %s = :1",
		table,
		pkColumn,
	)
	_, err := c.db.Exec(query, id)
	return err
}

func (c *Crud) DeleteByPK(table string, model any) error {
	v := reflect.ValueOf(model)
	t := reflect.TypeOf(model)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	var pkColumn string
	var pkValue any

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		tag := field.Tag.Get("db")
		if tag == "" {
			continue
		}

		parts := strings.Split(tag, ",")
		column := parts[0]

		for _, p := range parts[1:] {
			if p == "pk" {
				pkColumn = column
				pkValue = value.Interface()
				break
			}
		}
	}

	if pkColumn == "" {
		return fmt.Errorf("pk não encontrada no model")
	}

	query := fmt.Sprintf(
		"DELETE FROM %s WHERE %s = :1",
		table,
		pkColumn,
	)

	_, err := c.db.Exec(query, pkValue)
	return err
}

func (c *Crud) ListStruct(table string, dest any) error {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("dest deve ser ponteiro para slice")
	}

	sliceValue := v.Elem()
	elemType := sliceValue.Type().Elem()

	var columns []string
	var scanTargets []any

	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		tag := field.Tag.Get("db")
		if tag == "" {
			continue
		}

		column := strings.Split(tag, ",")[0]
		columns = append(columns, column)
	}

	query := fmt.Sprintf(
		"SELECT %s FROM %s ORDER BY %s",
		strings.Join(columns, ", "),
		table,
		columns[0],
	)

	rows, err := c.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		elemPtr := reflect.New(elemType)
		elem := elemPtr.Elem()

		scanTargets = scanTargets[:0]

		for i := 0; i < elemType.NumField(); i++ {
			field := elemType.Field(i)
			tag := field.Tag.Get("db")
			if tag == "" {
				continue
			}
			scanTargets = append(scanTargets, elem.Field(i).Addr().Interface())
		}

		if err := rows.Scan(scanTargets...); err != nil {
			return err
		}

		sliceValue.Set(reflect.Append(sliceValue, elem))
	}

	return nil
}

func (c *Crud) FindByID(table string, dest any) error {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("dest deve ser ponteiro para struct")
	}

	elem := v.Elem()
	elemType := elem.Type()

	var columns []string
	var scanTargets []any

	var pkColumn string
	var pkValue any

	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		tag := field.Tag.Get("db")
		if tag == "" {
			continue
		}

		parts := strings.Split(tag, ",")
		column := parts[0]

		var isPK bool
		for _, p := range parts[1:] {
			if p == "pk" {
				isPK = true
			}
		}

		if isPK {
			pkColumn = column
			pkValue = elem.Field(i).Interface()
		}

		columns = append(columns, column)
		scanTargets = append(scanTargets, elem.Field(i).Addr().Interface())
	}

	if pkColumn == "" {
		return fmt.Errorf("pk não encontrada no model")
	}

	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s = :1",
		strings.Join(columns, ", "),
		table,
		pkColumn,
	)

	row := c.db.QueryRow(query, pkValue)
	return row.Scan(scanTargets...)
}
