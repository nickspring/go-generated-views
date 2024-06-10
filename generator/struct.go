package generator

import (
	"fmt"
	"go/ast"
	"go/types"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type (

	// StructTag describes tag value for each view, default view has * key
	StructTag map[string]string

	// StructField describes one field on struct
	StructField struct {
		Name string
		Type string
		Tag  *StructTag
	}

	// Struct describes one struct
	Struct struct {
		Name   string
		Views  []string
		Fields []StructField
	}
)

// NewStructTagFromString creates StructTag from source struct tag
func NewStructTagFromString(sourceTag string) *StructTag {
	result := StructTag{}
	temp := map[string]map[string]string{commonView: {}}
	tags := parseTag(sourceTag)
	for _, tagSlice := range tags {
		tagName, tagValue := strings.TrimSpace(tagSlice[0]), strings.TrimSpace(tagSlice[1])
		matches := specialTagRegexp.FindStringSubmatch(tagName)
		if matches == nil {
			temp[commonView][tagName] = tagValue
			continue
		}
		realTagName := matches[1]
		views := strings.Split(matches[2], ",")
		for _, view := range views {
			view = cases.Title(language.English, cases.NoLower).String(strings.TrimSpace(view))
			if temp[view] == nil {
				temp[view] = make(map[string]string)
			}
			temp[view][realTagName] = tagValue
		}
	}
	for viewName, viewTags := range temp {
		newMap := make(map[string]string, len(temp[commonView]))
		for key, value := range temp[commonView] {
			newMap[key] = value
		}
		if viewName != commonView {
			for key, value := range viewTags {
				newMap[key] = value
			}
		}
		if len(newMap) == 0 {
			continue
		}

		keys := make([]string, len(newMap))
		i := 0
		for k := range newMap {
			keys[i] = k
			i++
		}
		sort.Strings(keys)

		result[viewName] = ""
		for _, key := range keys {
			result[viewName] += fmt.Sprintf(`%s:"%s" `, key, strings.ReplaceAll(newMap[key], `"`, `\"`))
		}
		result[viewName] = strings.TrimSpace(result[viewName])
	}
	return &result
}

// NewStructList creates a list of Struct from the provided ast.File node.
func NewStructList(node *ast.File) []Struct {
	var structs []Struct
	for _, decl := range node.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if structType, ok := typeSpec.Type.(*ast.StructType); ok {
						var structInfo Struct
						viewKeys := map[string]struct{}{}
						structInfo.Name = typeSpec.Name.Name
						for _, field := range structType.Fields.List {
							tag := ""
							if field.Tag != nil {
								tag = strings.Trim(field.Tag.Value, "`")
							}
							structTag := NewStructTagFromString(tag)
							for viewName := range *structTag {
								if viewName == commonView {
									continue
								}
								if _, exists := viewKeys[viewName]; exists {
									continue
								}
								viewKeys[viewName] = struct{}{}
								structInfo.Views = append(structInfo.Views, viewName)
							}
							sf := StructField{
								Name: field.Names[0].Name,
								Type: types.ExprString(field.Type),
								Tag:  structTag,
							}
							structInfo.Fields = append(structInfo.Fields, sf)
						}
						sort.Strings(structInfo.Views)
						structs = append(structs, structInfo)
					}
				}
			}
		}
	}
	return structs
}

// parseTag returns slice of tags pairs (name, value) from tags string
// original algorithm from reflect package is used
func parseTag(tag string) [][]string {
	var result [][]string
	for tag != "" {
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// Scan to colon. A space, a quote or a control character is a syntax error.
		// Strictly speaking, control chars include the range [0x7f, 0x9f], not just
		// [0x00, 0x1f], but in practice, we ignore the multi-byte control characters
		// as it is simpler to inspect the tag's bytes than the tag's runes.
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := tag[:i]
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := tag[:i+1]
		tag = tag[i+1:]

		value, err := strconv.Unquote(qvalue)
		if err == nil {
			result = append(result, []string{name, value})
		}
	}
	return result
}
