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
	structVal := reflect.ValueOf(obj)

	var missingFlags []string
	missingFlags = append(missingFlags, unmarshal(structVal, typ, c)...)

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

func unmarshal(structVal reflect.Value, typ reflect.Type, c *cli.Context) (missingFlags []string) {
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		if field.Anonymous {
			// recurse into embedded fields
			var embeddedField reflect.Value
			if structVal.Kind() == reflect.Ptr {
				embeddedField = structVal.Elem().Field(i)
			} else {
				embeddedField = structVal.Field(i)
			}
			missingFlags = append(missingFlags, unmarshal(embeddedField, embeddedField.Type(), c)...)
			continue
		}

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
			optional = true // bool flags aren't required to be present on the command line
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
		if structVal.Kind() == reflect.Ptr {
			structVal.Elem().Field(i).Set(desiredValue)
		} else {
			structVal.Field(i).Set(desiredValue)
		}

		// Check for missing [required] flags.
		if !optional && !c.IsSet(flagName) {
			if flagIsMissing(value, desiredValue) { //reflect.Zero(reflect.TypeOf(value)).Interface() == desiredValue.Interface() {
				missingFlags = append(missingFlags, flagName)
			}
		}
	}
	return
}

// flagIsMissing checks whether the flag is missing.
//
// Note: if a flag has a default value and was not specified on the command line,
// c.IsSet(flag) will return false. So in order to check if a flag is truly missing,
// we check that its default value is non-zero.
//
// For non-slice types, we do a simple comparison.
// However, equality is not defined for slice types, so for slices
// we just check if they are empty.
func flagIsMissing(contextValue interface{}, desiredValue reflect.Value) bool {
	if desiredValue.Kind() == reflect.Slice {
		return sliceIsEmpty(desiredValue)
	} else {
		return reflectValueIsZero(contextValue, desiredValue)
	}
}

func reflectValueIsZero(contextValue interface{}, desiredValue reflect.Value) bool {
	return reflect.Zero(reflect.TypeOf(contextValue)).Interface() == desiredValue.Interface()
}

func sliceIsEmpty(v reflect.Value) bool {
	return v.Len() == 0
}
