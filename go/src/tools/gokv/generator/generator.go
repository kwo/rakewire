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
	Indexes   map[string][]*IndexField
}

// FieldInfo describes a field
type FieldInfo struct {
	Name               string
	NameLower          string
	StructName         string
	StructNameLower    string
	Type               string
	Tags               map[string]string
	Required           bool
	Comment            string
	EmptyValue         string
	ZeroTest           string
	SerializeCommand   string
	DeserializeCommand string
}

// IndexField represents a single value in an index key
type IndexField struct {
	Field  string
	Filter string
}

// Finalize populates index and key information from tags.
func (z *StructInfo) Finalize() error {

	for _, field := range z.Fields {

		//fmt.Printf("%s/%s: %s\n", z.Name, field.Name, field.Type)

		if field.Name == "ID" {
			field.Required = true
		}

		if field.Type == "time.Time" {
			z.Imports["time"] = true
		}

		tagFields := strings.Split(field.Tags["kv"], ",")
		for _, tagField := range tagFields {

			if tagField == "" {
				continue
			} else if tagField == "+required" {
				// required field
				field.Required = true
				if field.Type == "bool" {
					return fmt.Errorf("%s.%s, Type: bool - cannot be required", z.Name, field.Name)
				}
			} else {

				// index definition

				elements := strings.Split(tagField, ":")
				var filter string
				switch len(elements) {
				case 2:
				case 3:
					filter = elements[2]
					z.Imports["strings"] = true
				default:
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
					index = append(index, nil)
				}
				index[indexPosition-1] = &IndexField{Field: field.Name, Filter: filter}
				z.Indexes[indexName] = index

			}

		} // tagFields

	} // fields

	return nil

}

// Finalize populates index and key information from tags.
func (z *FieldInfo) Finalize(s *StructInfo) error {

	z.StructName = s.Name
	z.StructNameLower = s.NameLower

	switch z.Type {
	case "string":
		z.EmptyValue = "\"\""
		z.ZeroTest = executeTemplate("ZeroTestDefault", z)
		z.SerializeCommand = executeTemplate("SerializeDefault", z)
		z.DeserializeCommand = executeTemplate("DeserializeDefault", z)

	case "bool":
		z.EmptyValue = "false"
		z.ZeroTest = executeTemplate("ZeroTestBool", z)
		z.SerializeCommand = executeTemplate("SerializeBool", z)
		z.DeserializeCommand = executeTemplate("DeserializeBool", z)

	case "int", "int8", "int16", "int32", "int64":
		z.EmptyValue = "0"
		z.ZeroTest = executeTemplate("ZeroTestDefault", z)
		z.SerializeCommand = executeTemplate("SerializeInt", z)
		z.DeserializeCommand = executeTemplate("DeserializeInt", z)

	case "uint", "uint8", "uint16", "uint32":
		z.EmptyValue = "0"
		z.ZeroTest = executeTemplate("ZeroTestDefault", z)
		z.SerializeCommand = executeTemplate("SerializeInt", z)
		z.DeserializeCommand = executeTemplate("DeserializeUint", z)

	case "uint64":
		z.EmptyValue = "0"
		z.ZeroTest = executeTemplate("ZeroTestDefault", z)
		z.SerializeCommand = executeTemplate("SerializeIntKey", z)
		z.DeserializeCommand = executeTemplate("DeserializeUint", z)

	case "float32", "float64":
		z.EmptyValue = "0"
		z.ZeroTest = executeTemplate("ZeroTestDefault", z)
		z.SerializeCommand = executeTemplate("SerializeFloat", z)
		z.DeserializeCommand = executeTemplate("DeserializeFloat", z)

	case "time.Time":
		z.EmptyValue = "time.Time{}"
		z.ZeroTest = executeTemplate("ZeroTestTime", z)
		z.SerializeCommand = executeTemplate("SerializeTime", z)
		z.DeserializeCommand = executeTemplate("DeserializeTime", z)

	case "time.Duration":
		z.EmptyValue = "0"
		z.ZeroTest = executeTemplate("ZeroTestDefault", z)
		z.SerializeCommand = executeTemplate("SerializeDuration", z)
		z.DeserializeCommand = executeTemplate("DeserializeDuration", z)

	case "[]uint64":
		z.EmptyValue = "[]uint64{}"
		z.ZeroTest = executeTemplate("ZeroTestArray", z)
		z.SerializeCommand = executeTemplate("SerializeIntArray", z)
		z.DeserializeCommand = executeTemplate("DeserializeUintArray", z)
		s.Imports["bytes"] = true
		s.Imports["strings"] = true

	default:
		fmt.Printf("Unknown type %s, field: %s, struct: %s\n", z.Type, z.Name, z.StructName)
		z.EmptyValue = "0"
		z.ZeroTest = executeTemplate("ZeroTestDefault", z)
		z.SerializeCommand = executeTemplate("SerializeDefault", z)
		z.DeserializeCommand = executeTemplate("DeserializeDefault", z)

	}

	return nil

}

// Generate ${filename}_kv.go for the given file
func Generate(filename, kvFilename string) error {

	templates["DeserializeDefault"] = template.Must(template.New("DeserializeDefault").Parse(tplDeserializeDefault))
	templates["DeserializeBool"] = template.Must(template.New("DeserializeBool").Parse(tplDeserializeBool))
	templates["DeserializeFloat"] = template.Must(template.New("DeserializeFloat").Parse(tplDeserializeFloat))
	templates["DeserializeInt"] = template.Must(template.New("DeserializeInt").Parse(tplDeserializeInt))
	templates["DeserializeUint"] = template.Must(template.New("DeserializeUint").Parse(tplDeserializeUint))
	templates["DeserializeTime"] = template.Must(template.New("DeserializeTime").Parse(tplDeserializeTime))
	templates["DeserializeDuration"] = template.Must(template.New("DeserializeDuration").Parse(tplDeserializeDuration))
	templates["DeserializeUintArray"] = template.Must(template.New("DeserializeUintArray").Parse(tplDeserializeUintArray))

	templates["SerializeDefault"] = template.Must(template.New("SerializeDefault").Parse(tplSerializeDefault))
	templates["SerializeBool"] = template.Must(template.New("SerializeBool").Parse(tplSerializeBool))
	templates["SerializeFloat"] = template.Must(template.New("SerializeFloat").Parse(tplSerializeFloat))
	templates["SerializeInt"] = template.Must(template.New("SerializeInt").Parse(tplSerializeInt))
	templates["SerializeIntKey"] = template.Must(template.New("SerializeIntKey").Parse(tplSerializeIntKey))
	templates["SerializeTime"] = template.Must(template.New("SerializeTime").Parse(tplSerializeTime))
	templates["SerializeDuration"] = template.Must(template.New("SerializeDuration").Parse(tplSerializeDuration))
	templates["SerializeIntArray"] = template.Must(template.New("SerializeIntArray").Parse(tplSerializeIntArray))

	templates["ZeroTestDefault"] = template.Must(template.New("ZeroTestDefault").Parse(tplZeroTestDefault))
	templates["ZeroTestBool"] = template.Must(template.New("ZeroTestBool").Parse(tplZeroTestBool))
	templates["ZeroTestTime"] = template.Must(template.New("ZeroTestTime").Parse(tplZeroTestTime))
	templates["ZeroTestArray"] = template.Must(template.New("ZeroTestArray").Parse(tplZeroTestArray))

	packageName, structInfos, err := ExtractStructs(filename)
	if err != nil {
		return err
	}

	imports := make(map[string]bool)
	imports["fmt"] = true
	imports["strconv"] = true
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
				Indexes:   make(map[string][]*IndexField),
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
						if f.Tags["kv"] != "-" {
							f.Finalize(structInfo)
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
