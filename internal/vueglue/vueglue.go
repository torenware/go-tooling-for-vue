package vueglue

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"reflect"
	"strconv"
	"strings"
)

const (
	AssetsDir = "dist/assets"
	RootItem  = "src/main.ts" // set in vite.config.ts
)

type VueGlue struct {
	MainModule string
	Imports    []string
	CSSModule  []string
	DistFS     *embed.FS
}

type ManifestTarget struct {
	File    string   `json:"file"`
	Source  string   `json:"src"`
	IsEntry bool     `json:"isEntry"`
	Imports []string `json:"imports"`
	CSS     []string `json:"css"`
}

type JSONLookup map[string]reflect.StructField

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

// @see https://yourbasic.org/golang/json-example/
func (m ManifestTarget) parseWithoutReflection(jsonData []byte) {
	var v interface{}
	json.Unmarshal(jsonData, &v)
	m.siftCollections("", "", v)
}

func (m ManifestTarget) siftCollections(indent, key string, v interface{}) {
	data, ok := v.(map[string]interface{})
	if ok {
		for k, v := range data {
			m.processInterface(indent, k, v)
		}
	} else if arrayData, ok := v.([]interface{}); ok {
		for i, v := range arrayData {
			m.processInterface(indent, strconv.Itoa(i), v)
		}
	} else {
		m.processInterface(indent, key, v)
	}
}

// call this for recurisve structures.
func (m ManifestTarget) processInterface(indent, k string, v interface{}) {

	switch v := v.(type) {
	case string:
		fmt.Println(indent, k, "=", v, "(string)")
	case float64:
		fmt.Println(indent, k, "=", v, "(float64)")
	case bool:
		fmt.Println(indent, k, "=", v, "(bool)")
	case []interface{}:
		fmt.Println(indent, k, "=", "(array):")
		m.siftCollections(indent+"    ", k, v)
	case map[string]interface{}:
		fmt.Println(indent, k, "=", "(map):")
		m.siftCollections(indent+"    ", k, v)
	default:
		fmt.Printf("%s %s ?? %T (unknown)", indent, k, v)
	}
}

type ManifestMaps map[string]interface{}

func NewVueGlue(dist *embed.FS) (*VueGlue, error) {
	var glue VueGlue
	// prefix := "/dist"
	// removeChars := len(prefix) - 1

	if !fs.ValidPath(AssetsDir) {
		return nil, errors.New("vite dist directory not found")
	}
	glue.DistFS = dist

	dir, err := dist.ReadDir(".")
	if err != nil {
		return nil, errors.New("could not read dir")
	}
	for _, entry := range dir {
		log.Println(entry.Name())
	}

	// Get the manifest file
	contents, err := dist.ReadFile("dist/manifest.json")
	if err != nil {
		return nil, err
	}

	var testRslt ManifestTarget

	// temp: lookup for json fields.
	lookup := testRslt.LookupFields()
	fmt.Println(lookup)

	testRslt.parseWithoutReflection(contents)
	return &glue, nil

}
