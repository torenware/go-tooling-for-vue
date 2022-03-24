package vueglue

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Not currently using this file, since I found a simpler way to parse
// arbitrary JSON that won't work with json.Unmarshal. But if I
// need to use golang reflection again, at least I have this to
// look back on.

// From gopl book
func Display(name string, x interface{}) {
	fmt.Printf("Display %s (%T):\n", name, x)
	display(name, reflect.ValueOf(x))
}

//!-Display

// formatAtom formats a value without inspecting its internal structure.
// It is a copy of the the function in gopl.io/ch11/format.
func formatAtom(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	// ...floating-point and complex cases omitted for brevity...
	case reflect.Bool:
		if v.Bool() {
			return "true"
		}
		return "false"
	case reflect.String:
		return strconv.Quote(v.String())
	case reflect.Chan, reflect.Func, reflect.Ptr,
		reflect.Slice, reflect.Map:
		return v.Type().String() + " 0x" +
			strconv.FormatUint(uint64(v.Pointer()), 16)
	default: // reflect.Array, reflect.Struct, reflect.Interface
		return v.Type().String() + " value"
	}
}

//!+display
func display(path string, v reflect.Value) {
	switch v.Kind() {
	case reflect.Invalid:
		fmt.Printf("%s = invalid\n", path)
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			display(fmt.Sprintf("%s[%d]", path, i), v.Index(i))
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fieldPath := fmt.Sprintf("%s.%s", path, v.Type().Field(i).Name)
			display(fieldPath, v.Field(i))
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			display(fmt.Sprintf("%s[%s]", path,
				formatAtom(key)), v.MapIndex(key))
		}
	case reflect.Ptr:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			display(fmt.Sprintf("(*%s)", path), v.Elem())
		}
	case reflect.Interface:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			fmt.Printf("%s.type = %s\n", path, v.Elem().Type())
			display(path+".value", v.Elem())
		}
	default: // basic types, channels, funcs
		fmt.Printf("%s = %s\n", path, formatAtom(v))
	}
}

// end from gopl

type JSONLookup map[string]reflect.StructField

// from
func (b ManifestTarget) LookupFields() JSONLookup {
	lookup := make(map[string]reflect.StructField)
	jsonTable := make(JSONLookup)

	// Get a lookup of the fields of the struct
	// and their corresponding reflect values.
	// @see https://stackoverflow.com/a/24337241/8600734
	structVal := reflect.ValueOf(b)
	for i := 0; i < structVal.Type().NumField(); i++ {
		field := structVal.Type().Field(i)
		lookup[field.Name] = field
		fmt.Println("field", field.Name)
	}

	// Find the json tags for the fields.
	for i := 0; i < structVal.Type().NumField(); i++ {
		t := structVal.Type().Field(i)
		// fieldName := t.Name
		jsonTag := t.Tag.Get("json")

		switch jsonTag {
		case "-", "":
			// skip; no usable tag.
			continue
		default:
			parts := strings.Split(jsonTag, ",")
			fieldName := t.Name
			name := parts[0]
			if name == "" {
				continue
			}
			field, ok := lookup[fieldName]
			if !ok {
				continue
			}
			jsonTable[name] = field
		}
	}

	return jsonTable

}

// end
