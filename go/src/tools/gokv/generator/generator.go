package generator

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"strconv"
	"strings"
	"text/template"
)

// StructInfo describes a struct
type StructInfo struct {
	Name    string
	Imports map[string]bool
	Fields  []*FieldInfo
	Indexes map[string][]string
}

// FieldInfo describes a field
type FieldInfo struct {
	Name       string
	Type       string
	Tags       map[string]string
	Comment    string
	EmptyValue string
}

// Finalize populates index and key information from tags.
func (z *StructInfo) Finalize() error {

	hasTime := false

	for _, field := range z.Fields {

		if field.Type == "time.Time" {
			hasTime = true
		}

		tagFields := strings.Split(field.Tags["kv"], ",")
		for _, tagField := range tagFields {
			if tagField == "" {
				continue
			}

			elements := strings.Split(tagField, ":")
			if len(elements) != 2 {
				return fmt.Errorf("Invalid index definition: %s", tagField)
			}
			indexName := elements[0]
			indexPosition, err := strconv.Atoi(elements[1])
			if err != nil {
				return fmt.Errorf("Index position is not an integer: %s", tagField)
			} else if indexPosition < 1 {
				return fmt.Errorf("Index positions are one-based: %s", tagField)
			}

			index := z.Indexes[indexName]
			for len(index) < indexPosition {
				index = append(index, "")
			}
			index[indexPosition-1] = field.Name
			z.Indexes[indexName] = index

		} // tagFields

	} // fields

	if hasTime {
		z.Imports["time"] = true
	}

	return nil

}

// Generate ${filename}_kv.go for the given file
func Generate(filename, kvFilename string) error {

	packageName, structInfos, err := ExtractStructs(filename)
	if err != nil {
		return err
	}

	imports := make(map[string]bool)
	for _, s := range structInfos {
		for k := range s.Imports {
			imports[k] = true
		}
	}

	kvTemplate := &KVTemplate{
		PackageName: packageName,
		Imports:     imports,
		Structures:  structInfos,
	}

	t := template.Must(template.New("kvTemplate").Parse(kvTemplateText))

	buf := new(bytes.Buffer)
	err = t.Execute(buf, kvTemplate)
	if err != nil {
		return err
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	f, err := os.OpenFile(kvFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(formatted)
	return err

}

// ExtractStructs retrieves the structs from the given go file
func ExtractStructs(filename string) (string, []*StructInfo, error) {

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return "", nil, err
	}

	packageName := f.Name.String()

	structs := []*StructInfo{}

	for k, d := range f.Scope.Objects {
		if d.Kind == ast.Typ {
			structInfo := &StructInfo{
				Name:    k,
				Imports: make(map[string]bool),
				Indexes: make(map[string][]string),
			}

			ast.Inspect(d.Decl.(ast.Node), func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.StructType:
					structs = append(structs, structInfo)
					for _, field := range x.Fields.List {
						f := &FieldInfo{
							Name:    field.Names[0].String(),
							Type:    types.ExprString(field.Type),
							Tags:    extractTags(field.Tag),
							Comment: extractCommentTexts(field.Comment),
						}
						f.EmptyValue = getEmptyValue(f.Type)
						if f.Tags["kv"] != "-" {
							structInfo.Fields = append(structInfo.Fields, f)
						}
					}
				}
				return true
			}) // ast.Inspect

			if err := structInfo.Finalize(); err != nil {
				return "", nil, err
			}

		} // ast.Typ

	} // scope.Objects

	return packageName, structs, nil

}

func extractCommentTexts(comments *ast.CommentGroup) string {

	list := []string{}

	if comments != nil && comments.List != nil {
		for _, comment := range comments.List {
			list = append(list, comment.Text)
		}
	}

	return strings.Join(list, ",")

}

func extractTags(literal *ast.BasicLit) map[string]string {
	result := make(map[string]string)
	if literal != nil {
		if tagstr, err := strconv.Unquote(literal.Value); err == nil {
			tagfields := strings.Fields(tagstr)
			for _, tagfield := range tagfields {
				tag := strings.SplitN(tagfield, ":", 2)
				if len(tag) == 2 {
					if tagvalue, err := strconv.Unquote(tag[1]); err == nil {
						result[tag[0]] = tagvalue
					}
				}
			}
		}
	}
	return result
}

func getEmptyValue(typ string) string {

	switch typ {
	case "string":
		return "\"\""

	case "bool":
		return "false"

	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "float32", "float64":
		return "0"

	case "time.Time":
		return "time.Time{}"

	default:
		fmt.Printf("Unknown type: %s\n", typ)
		return "0"

	}

}
