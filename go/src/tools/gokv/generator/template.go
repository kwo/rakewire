package generator

// KVTemplate is the top-level template
type KVTemplate struct {
	PackageName string
	Imports     map[string]bool
	Structures  []*StructInfo
}

var kvTemplateText = `
package {{.PackageName}}

/*
 *  CODE GENERATED AUTOMATICALLY WITH gokv.
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 */

import (
	{{range $k, $v := .Imports}}"{{$k}}"
	{{end}}
)

{{range $index, $structure := .Structures}}

// index names
const (
	{{.NameLower}}Entity    = "{{.Name}}"
	{{range $name, $fields := .Indexes}}{{$structure.NameLower}}Index{{$name}} = "{{$name}}"
	{{end}}
)

const (
{{range $index, $field := .Fields}}{{$structure.NameLower}}{{.Name}} = "{{.Name}}"
{{end}}
)

var (
	{{$structure.NameLower}}AllFields = []string{
	  {{range $index, $field := .Fields}}{{$structure.NameLower}}{{.Name}},{{end}}
	}
)

// GetID return the primary key of the object.
func (z *{{.Name}}) getID() uint64 {
	return z.ID
}

// SetID sets the primary key of the object.
func (z *{{.Name}}) setID(id uint64) {
	z.ID = id
}

// Clear reset all fields to zero/empty
func (z *{{.Name}}) clear() {
{{range $index, $field := .Fields}}z.{{.Name}} = {{.EmptyValue}}
{{end}}
}

// Serialize serializes an object to a list of key-values.
// An optional flag, when set, will serialize all fields to the resulting map, not just the non-zero values.
func (z *{{.Name}}) serialize(flags ...bool) map[string]string {
	flagNoZeroCheck := len(flags) > 0 && flags[0]
	result := make(map[string]string)
	{{range $index, $field := .Fields}}
		if flagNoZeroCheck || {{.ZeroTest}} {
			result[{{$structure.NameLower}}{{.Name}}] = {{.SerializeCommand}}
		}
	{{end}}
	return result
}

// Deserialize serializes an object to a list of key-values.
// An optional flag, when set, will return an error if unknown keys are contained in the values.
func (z *{{.Name}}) deserialize(values map[string]string, flags ...bool) error {
	flagUnknownCheck := len(flags) > 0 && flags[0]
	{{$struct := .}}
	var errors []error
	var missing []string
	var unknown []string
	{{range $index, $field := .Fields}}
		z.{{.Name}} = {{.DeserializeCommand}}
		{{if .Required}}
		if !({{.ZeroTest}}) {
			missing = append(missing, {{$struct.NameLower}}{{.Name}})
		}
		{{end}}
	{{end}}
	if flagUnknownCheck {
		for fieldname := range values {
			if !isStringInArray(fieldname, {{$struct.NameLower}}AllFields) {
				unknown = append(unknown, fieldname)
			}
		}
	}
	return newDeserializationError({{$struct.NameLower}}Entity, errors, missing, unknown)
}

// IndexKeys returns the keys of all indexes for this object.
func (z *{{.Name}}) indexKeys() map[string][]string {
	{{$struct := .}}
	result := make(map[string][]string)
	{{if ne (len .Indexes) 0}}
	data := z.serialize(true)
	{{end}}
	{{range $name, $fields := .Indexes}}
	result[{{$structure.NameLower}}Index{{$name}}] = []string{
		{{range $index, $f := $fields}}
			{{if eq $f.Filter "lower"}}
				strings.ToLower(data[{{$struct.NameLower}}{{$f.Field}}]),
			{{else}}
				data[{{$struct.NameLower}}{{$f.Field}}],
			{{end}}
	  {{end}}
	}
	{{end}}
	return result
}

{{end}}

`

var tplDeserializeDefault = `values[{{.StructNameLower}}{{.Name}}]`
var tplDeserializeBool = `func (fieldName string, values map[string]string, errors []error) bool {
	if value, ok := values[fieldName]; ok {
		return value == "1"
	}
	return false
}({{.StructNameLower}}{{.Name}}, values, errors)
`
var tplDeserializeFloat = `func (fieldName string, values map[string]string, errors []error) {{.Type}} {
	result, err := strconv.ParseFloat(values[fieldName], 64)
	if err != nil {
		errors = append(errors, err)
		return 0
	}
	return {{.Type}}(result)
}({{.StructNameLower}}{{.Name}}, values, errors)
`
var tplDeserializeInt = `func (fieldName string, values map[string]string, errors []error) {{.Type}} {
	result, err := strconv.ParseInt(values[fieldName], 10, 64)
	if err != nil {
		errors = append(errors, err)
		return 0
	}
	return {{.Type}}(result)
}({{.StructNameLower}}{{.Name}}, values, errors)
`
var tplDeserializeUint = `func (fieldName string, values map[string]string, errors []error) {{.Type}} {
	result, err := strconv.ParseUint(values[fieldName], 10, 64)
	if err != nil {
		errors = append(errors, err)
		return 0
	}
	return {{.Type}}(result)
}({{.StructNameLower}}{{.Name}}, values, errors)
`
var tplDeserializeTime = `func (fieldName string, values map[string]string, errors []error) time.Time {
	result := time.Time{}
	if value, ok := values[fieldName]; ok {
		t, err := time.Parse(time.RFC3339, value)
		if err != nil {
			errors = append(errors, err)
		} else {
			result = t
		}
	}
	return result
}({{.StructNameLower}}{{.Name}}, values, errors)
`

var tplDeserializeDuration = `func (fieldName string, values map[string]string, errors []error) time.Duration {
	var result time.Duration
	if value, ok := values[fieldName]; ok {
		t, err := time.ParseDuration(value)
		if err != nil {
			errors = append(errors, err)
		} else {
			result = t
		}
	}
	return result
}({{.StructNameLower}}{{.Name}}, values, errors)
`

var tplDeserializeUintArray = `func (fieldName string, values map[string]string, errors []error) {{.Type}} {
	var result {{.Type}}
	if value, ok := values[fieldName]; ok {
		elements := strings.Fields(value)
		for _, element := range elements {
			value, err := strconv.ParseUint(element, 10, 64)
			if err != nil {
				errors = append(errors, err)
				break
			}
			result = append(result, value)
		}
	}
	return result
}({{.StructNameLower}}{{.Name}}, values, errors)
`

var tplSerializeDefault = `z.{{.Name}}`
var tplSerializeBool = `func(value {{.Type}}) string {
	if value {
		return "1"
	}
	return "0"
}(z.{{.Name}})`
var tplSerializeFloat = `fmt.Sprintf("%f", z.{{.Name}})`
var tplSerializeInt = `fmt.Sprintf("%d", z.{{.Name}})`
var tplSerializeIntKey = `fmt.Sprintf("%05d", z.{{.Name}})`
var tplSerializeTime = `z.{{.Name}}.UTC().Format(time.RFC3339)`
var tplSerializeDuration = `z.{{.Name}}.String()`
var tplSerializeIntArray = `func(values {{.Type}}) string {
	var buffer bytes.Buffer
  for i, value := range values {
		if i > 0 {
			buffer.WriteString(" ")
		}
		buffer.WriteString(fmt.Sprintf("%d", value))
	}
	return buffer.String()
}(z.{{.Name}})`

var tplZeroTestDefault = `z.{{.Name}} != {{.EmptyValue}}`
var tplZeroTestBool = `z.{{.Name}}`
var tplZeroTestTime = `!z.{{.Name}}.IsZero()`
var tplZeroTestArray = `len(z.{{.Name}}) > 0`
