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
	{{.Name}}Entity    = "{{.Name}}"
	{{range $name, $fields := .Indexes}}{{$structure.Name}}Index{{$name}} = "{{$name}}"
	{{end}}
)

const (
{{range $index, $field := .Fields}}{{$structure.NameLower}}{{.Name}} = "{{.Name}}"
{{end}}
)

// GetName return the name of the entity.
func (z *{{.Name}}) GetName() string {
	return {{.Name}}Entity
}

// GetID return the primary key of the object.
func (z *{{.Name}}) GetID() uint64 {
	return z.ID
}

// SetID sets the primary key of the object.
func (z *{{.Name}}) SetID(id uint64) {
	z.ID = id
}

// Clear reset all fields to zero/empty
func (z *{{.Name}}) Clear() {
{{range $index, $field := .Fields}}z.{{.Name}} = {{.EmptyValue}}
{{end}}
}

// Serialize serializes an object to a list of key-values.
func (z *{{.Name}}) Serialize() map[string]string {
	result := make(map[string]string)
	{{range $index, $field := .Fields}}
		if {{.ZeroTest}} {
			result[{{$structure.NameLower}}{{.Name}}] = {{.SerializeCommand}}
		}
	{{end}}
	return result
}

// Deserialize serializes an object to a list of key-values.
func (z *{{.Name}}) Deserialize(values map[string]string) error {
	var errors []error
	{{range $index, $field := .Fields}}
		z.{{.Name}} = {{.DeserializeCommand}}
	{{end}}
	if len(errors) > 0 {
		return errors[0]
	}
	return nil
}

// IndexKeys returns the keys of all indexes for this object.
func (z *{{.Name}}) IndexKeys() map[string][]string {
	{{$struct := .}}
	result := make(map[string][]string)
	{{if ne (len .Indexes) 0}}
	data := z.Serialize()
	{{end}}
	{{range $name, $fields := .Indexes}}
	result[{{$structure.Name}}Index{{$name}}] = []string{
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
	result, err := strconv.ParseBool(values[fieldName])
	if err != nil {
		errors = append(errors, err)
		return false
	}
	return result
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
		t, err := time.Parse(time.RFC3339Nano, value)
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

var tplSerializeDefault = `z.{{.Name}}`
var tplSerializeBool = `fmt.Sprintf("%t", z.{{.Name}})`
var tplSerializeFloat = `fmt.Sprintf("%f", z.{{.Name}})`
var tplSerializeInt = `fmt.Sprintf("%d", z.{{.Name}})`
var tplSerializeIntKey = `fmt.Sprintf("%05d", z.{{.Name}})`
var tplSerializeTime = `z.{{.Name}}.Format(time.RFC3339Nano)`
var tplSerializeDuration = `z.{{.Name}}.String()`

var tplZeroTestDefault = `z.{{.Name}} != {{.EmptyValue}}`
var tplZeroTestBool = `z.{{.Name}}`
var tplZeroTestTime = `!z.{{.Name}}.IsZero()`
