package generator

import (
	"bytes"
	"flag"
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

var (
	templates       = make(map[string]*template.Template)
	flagUnformatted = flag.Bool("u", false, "Leave generated source code unformatted.")
)

// StructInfo describes a struct
type StructInfo struct {
	Name      string
	NameLower string
	Imports   map[string]bool
	Fields    []*FieldInfo
	Indexes   map[string][]string
}

// FieldInfo describes a field
type FieldInfo struct {
	Name          string
	NameLower     string
	Type          string
	Tags          map[string]string
	Comment       string
	EmptyValue    string
	ZeroTest      string
	FormatCommand string
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

// Finalize populates index and key information from tags.
func (z *FieldInfo) Finalize() error {

	switch z.Type {
	case "string":
		z.EmptyValue = "\"\""
		z.ZeroTest = executeTemplate("ZeroTestDefault", z)
		z.FormatCommand = executeTemplate("FormatDefault", z)

	case "bool":
		z.EmptyValue = "false"
		z.ZeroTest = executeTemplate("ZeroTestBool", z)
		z.FormatCommand = executeTemplate("FormatBool", z)

	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "float32", "float64":
		z.EmptyValue = "0"
		z.ZeroTest = executeTemplate("ZeroTestDefault", z)
		z.FormatCommand = executeTemplate("FormatNumeric", z)

	case "uint64":
		z.EmptyValue = "0"
		z.ZeroTest = executeTemplate("ZeroTestDefault", z)
		z.FormatCommand = executeTemplate("FormatNumericKey", z)

	case "time.Time":
		z.EmptyValue = "time.Time{}"
		z.ZeroTest = executeTemplate("ZeroTestTime", z)
		z.FormatCommand = executeTemplate("FormatTime", z)

	default:
		fmt.Printf("Unknown type for fmt verb: %s\n", z.Type)
		z.EmptyValue = "0"
		z.ZeroTest = executeTemplate("ZeroTestDefault", z)
		z.FormatCommand = executeTemplate("FormatDefault", z)

	}

	return nil

}

// Generate ${filename}_kv.go for the given file
func Generate(filename, kvFilename string) error {

	templates["FormatDefault"] = template.Must(template.New("FormatDefault").Parse(tplFormatDefault))
	templates["FormatBool"] = template.Must(template.New("FormatBool").Parse(tplFormatBool))
	templates["FormatNumeric"] = template.Must(template.New("FormatNumeric").Parse(tplFormatNumeric))
	templates["FormatNumericKey"] = template.Must(template.New("FormatNumericKey").Parse(tplFormatNumericKey))
	templates["FormatTime"] = template.Must(template.New("FormatTime").Parse(tplFormatTime))
	templates["ZeroTestDefault"] = template.Must(template.New("ZeroTestDefault").Parse(tplZeroTestDefault))
	templates["ZeroTestBool"] = template.Must(template.New("ZeroTestBool").Parse(tplZeroTestBool))
	templates["ZeroTestTime"] = template.Must(template.New("ZeroTestTime").Parse(tplZeroTestTime))

	packageName, structInfos, err := ExtractStructs(filename)
	if err != nil {
		return err
	}

	imports := make(map[string]bool)
	imports["fmt"] = true
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

	resultBytes := buf.Bytes()
	if !*flagUnformatted {
		formatted, err := format.Source(buf.Bytes())
		if err != nil {
			return err
		}
		resultBytes = formatted
	}

	f, err := os.OpenFile(kvFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(resultBytes)
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
				Name:      k,
				NameLower: strings.ToLower(k),
				Imports:   make(map[string]bool),
				Indexes:   make(map[string][]string),
			}

			ast.Inspect(d.Decl.(ast.Node), func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.StructType:
					structs = append(structs, structInfo)
					for _, field := range x.Fields.List {
						f := &FieldInfo{
							Name:      field.Names[0].String(),
							NameLower: strings.ToLower(field.Names[0].String()),
							Type:      types.ExprString(field.Type),
							Tags:      extractTags(field.Tag),
							Comment:   extractCommentTexts(field.Comment),
						}
						f.Finalize()
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

func executeTemplate(tplName string, data interface{}) string {
	t := templates[tplName]
	if t == nil {
		fmt.Printf("Template does not exist: %s\n", tplName)
		return ""
	}
	buf := new(bytes.Buffer)
	err := t.Execute(buf, data)
	if err != nil {
		fmt.Printf("Error executing template: %s\n", err.Error())
	}
	return string(buf.Bytes())
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
