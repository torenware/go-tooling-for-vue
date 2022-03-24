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

type ManifestNode struct {
	Key      string
	Type     reflect.Kind
	Value    reflect.Value
	Children []*ManifestNode
}

type ManifestTarget struct {
	File    string   `json:"file"`
	Source  string   `json:"src"`
	IsEntry bool     `json:"isEntry"`
	Imports []string `json:"imports"`
	CSS     []string `json:"css"`
	Nodes   []*ManifestNode
}

func (n *ManifestNode) subKey(key string) *ManifestNode {
	if len(n.Children) == 0 {
		return nil
	}
	for _, leaf := range n.Children {
		if leaf.Key == key {
			return leaf
		}
	}
	return nil
}

// @see https://yourbasic.org/golang/json-example/
func (m *ManifestTarget) parseWithoutReflection(jsonData []byte) (*VueGlue, error) {
	var v interface{}
	json.Unmarshal(jsonData, &v)
	topNode := ManifestNode{
		Key: "top",
	}
	m.Nodes = append(m.Nodes, &topNode)
	m.siftCollections(&topNode, "", "", v)

	// Get entry point
	entry := (*ManifestNode)(nil)
	glue := &VueGlue{}

	for _, leaf := range topNode.Children {
		if leaf.subKey("isEntry") != nil {
			entry = leaf
			glue.MainModule = leaf.subKey("file").Value.String()
			break
		}
	}
	if entry == nil {
		return nil, errors.New("manifest lacked entry point")
	}

	imports := entry.subKey("imports")
	if imports == nil || len(imports.Children) == 0 {
		return nil, errors.New("expected code to have js dependencies")
	}

	for _, child := range imports.Children {
		// these have a level of indirection for some reason
		deref := topNode.subKey(child.Value.String())
		if deref == nil {
			return nil, errors.New("expected details for import")
		}
		item := deref.subKey("file")
		if item == nil {
			return nil, errors.New("expected path for import")
		}
		glue.Imports = append(glue.Imports, item.Value.String())
	}

	css := entry.subKey("css")
	if css == nil || len(css.Children) == 0 {
		// not an error, since CSS is optional
		return glue, nil
	}

	for _, child := range css.Children {
		glue.CSSModule = append(glue.CSSModule, child.Value.String())
	}

	return glue, nil
}

func (m *ManifestTarget) siftCollections(leaf *ManifestNode, indent, key string, v interface{}) {
	data, ok := v.(map[string]interface{})
	if ok {
		leaf.Type = reflect.Map
		for k, v := range data {
			child := &ManifestNode{
				Key: k,
			}
			leaf.Children = append(leaf.Children, child)
			m.processInterface(child, indent, k, v)
		}
	} else if arrayData, ok := v.([]interface{}); ok {
		leaf.Type = reflect.Slice
		for i, v := range arrayData {
			child := &ManifestNode{}
			leaf.Children = append(leaf.Children, child)
			m.processInterface(child, indent, strconv.Itoa(i), v)
		}
	} else {
		m.processInterface(leaf, indent, key, v)
	}
}

// call this for recurisve structures.
func (m *ManifestTarget) processInterface(leaf *ManifestNode, indent, k string, v interface{}) {

	switch v := v.(type) {
	case string:
		leaf.Type = reflect.String
		leaf.Value = reflect.ValueOf(v)
		fmt.Println(indent, k, "=", v, "(string)")
	case float64:
		leaf.Type = reflect.Float64
		leaf.Value = reflect.ValueOf(v)
		fmt.Println(indent, k, "=", v, "(float64)")
	case bool:
		leaf.Type = reflect.Bool
		leaf.Value = reflect.ValueOf(v)
		fmt.Println(indent, k, "=", v, "(bool)")
	case []interface{}:
		fmt.Println(indent, k, "=", "(array):")
		m.siftCollections(leaf, indent+"    ", k, v)
	case map[string]interface{}:
		fmt.Println(indent, k, "=", "(map):")
		m.siftCollections(leaf, indent+"    ", k, v)
	default:
		fmt.Printf("%s %s ?? %T (unknown)", indent, k, v)
	}
}

type ManifestMaps map[string]interface{}

func NewVueGlue(dist *embed.FS) (*VueGlue, error) {

	if !fs.ValidPath(AssetsDir) {
		return nil, errors.New("vite dist directory not found")
	}

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
	glue, err := testRslt.parseWithoutReflection(contents)
	if err != nil {
		return nil, err
	}
	glue.DistFS = dist

	output, _ := json.MarshalIndent(glue, "", "  ")
	fmt.Println(string(output))

	return glue, nil

}
