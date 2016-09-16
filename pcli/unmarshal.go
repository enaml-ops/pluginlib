package pcli

import (
	"bufio"
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"unicode"

	cli "gopkg.in/urfave/cli.v2"
)

const omgTagName = "omg"

// UnmarshalFlags populates obj from the specified CLI context.
// By default, it will populate a struct field named "FooBar"
// with the contents of the flag "foo-bar".
// You can override this behavior and specify the flag name
// with the `omg:"flag-name"` struct tag.
//
// Warning - UnmarshalFlags may zero-out any fields that were
// not specified in the context.  If you want UnmarshalFlags
// to leave a field alone under all circumstances, then annotate
// it with the `omg:"-"` tag.
//
// By default, UnmarshalFlags assumes that all fields in obj that
// are not explicitly skipped over with `omg:"-"` are required flags.
// It will return an error if a required flag is missing.
//
// If a flag is not required and you don't want an error if it's
// missing, the flag can be annotated with `omg:"flag-name,optional"`.
// Note that the flag name is required in addition to 'optional' in
// this case.
func UnmarshalFlags(obj interface{}, c *cli.Context) error {
	typ := reflect.TypeOf(obj)
	k := typ.Kind()
	if k != reflect.Ptr || (k == reflect.Ptr && typ.Elem().Kind() != reflect.Struct) {
		panic("unmarshal: obj must be a pointer to a struct type")
	}
	typ = typ.Elem()

	var missingFlags []string

	structVal := reflect.ValueOf(obj)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		flagName := field.Tag.Get(omgTagName)
		optional := false

		// skip over fields tagged with "-"
		if flagName == "-" {
			continue
		}

		if flagName == "" {
			flagName = deriveFlagName(field.Name)
		} else {
			flagName, optional = parseFlagName(flagName)
		}

		var value interface{}
		switch field.Type.Kind() {
		case reflect.Bool:
			value = c.Bool(flagName)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
			value = c.Int(flagName)
		case reflect.Int64:
			value = c.Int64(flagName)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
			value = c.Uint(flagName)
		case reflect.Uint64:
			value = c.Uint64(flagName)
		case reflect.Float64:
			value = c.Float64(flagName)
		case reflect.String:
			value = c.String(flagName)
		case reflect.Slice:
			switch field.Type.Elem().Kind() {
			case reflect.String:
				value = c.StringSlice(flagName)
			case reflect.Int:
				value = c.IntSlice(flagName)
			default:
				// unsupported slice type, skip over this field
				continue
			}
		default:
			// unsupported field type
			continue
		}

		desiredValue := reflect.ValueOf(value)
		structVal.Elem().Field(i).Set(desiredValue)

		// Check for missing [required] flags.
		// Note: if a flag has a default value and was not specified on the command line,
		// c.IsSet(flag) will return false. So in order to check if a flag is truly missing,
		// we check that its default value is non-zero.
		if !optional && !c.IsSet(flagName) {
			if reflect.Zero(reflect.TypeOf(value)).Interface() == desiredValue.Interface() {
				missingFlags = append(missingFlags, flagName)
			}
		}
	}

	if len(missingFlags) > 0 {
		return fmt.Errorf("unmarshal: missing flags %v", missingFlags)
	}
	return nil
}

func parseFlagName(tag string) (flag string, optional bool) {
	pieces := strings.Split(tag, ",")
	flag = pieces[0]
	if len(pieces) == 2 {
		optional = pieces[1] == "optional"
	}
	return
}

func deriveFlagName(fieldName string) string {
	buf := &bytes.Buffer{}
	w := bufio.NewWriter(buf)
	for index, runeValue := range fieldName {
		if unicode.IsUpper(runeValue) {
			if index != 0 {
				w.WriteRune('-')
			}
			w.WriteRune(unicode.ToLower(runeValue))
		} else {
			w.WriteRune(runeValue)
		}
	}
	w.Flush()
	return buf.String()
}
