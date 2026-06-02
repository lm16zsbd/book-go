package binder

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

const tagName = "query"

func BindQuery(r *http.Request, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("binder: BindQuery requires a pointer to struct")
	}

	elem := rv.Elem()
	t := elem.Type()
	query := r.URL.Query()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := elem.Field(i)

		tag := field.Tag.Get(tagName)
		if tag == "" {
			continue
		}

		if !value.CanSet() {
			continue
		}

		parts := strings.Split(tag, ",")
		name := parts[0]

		if name == "" {
			key := toSnakeCase(field.Name)
			if len(parts) > 1 {
				key = parts[1]
			}
			if err := setField(value, query, key, field.Type); err != nil {
				return err
			}
			continue
		}

		defaultVal := ""
		if len(parts) > 1 {
			defaultVal = parts[1]
		}

		if err := setFieldWithDefault(value, query, name, field.Type, defaultVal); err != nil {
			return err
		}
	}

	return nil
}

func setFieldWithDefault(field reflect.Value, query map[string][]string, name string, typ reflect.Type, defaultVal string) error {
	vals, ok := query[name]
	val := ""
	if ok && len(vals) > 0 {
		val = vals[len(vals)-1]
	}
	if val == "" && defaultVal != "" {
		val = defaultVal
	}
	return assignValue(field, val, typ)
}

func setField(field reflect.Value, query map[string][]string, name string, typ reflect.Type) error {
	vals, ok := query[name]
	val := ""
	if ok && len(vals) > 0 {
		val = vals[len(vals)-1]
	}
	return assignValue(field, val, typ)
}

func assignValue(field reflect.Value, val string, typ reflect.Type) error {
	if val == "" {
		return nil
	}

	switch typ.Kind() {
	case reflect.String:
		field.SetString(val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return fmt.Errorf("binder: cannot parse %q as int for field %s", val, field.Type().Name())
		}
		field.SetInt(n)
	case reflect.Bool:
		b, err := strconv.ParseBool(val)
		if err != nil {
			return fmt.Errorf("binder: cannot parse %q as bool for field %s", val, field.Type().Name())
		}
		field.SetBool(b)
	case reflect.Slice:
		if typ.Elem().Kind() == reflect.String {
			parts := parseCommaSep(val)
			slice := reflect.MakeSlice(typ, len(parts), len(parts))
			for i, p := range parts {
				slice.Index(i).SetString(p)
			}
			field.Set(slice)
		}
	default:
		return fmt.Errorf("binder: unsupported field type %s for field %s", typ, field.Type().Name())
	}
	return nil
}

func parseCommaSep(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func toSnakeCase(name string) string {
	var b strings.Builder
	for i, r := range name {
		if i > 0 && r >= 'A' && r <= 'Z' {
			b.WriteByte('_')
		}
		b.WriteRune(r)
	}
	return strings.ToLower(b.String())
}
