package generator

// KVTemplate is the top-level template
type KVTemplate struct {
	PackageName string
	Imports     map[string]bool
	Structures  []*StructInfo
}

var kvTemplateText = `
package {{.PackageName}}

import (
	{{range $k := .Imports}}"{{$k}}"
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
{{range $index, $field := .Fields}}{{$structure.Name}}{{.Name}} = "{{.Name}}"
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
		{{if eq .Type "bool"}}
			if z.{{.Name}} {
				result[{{$structure.Name}}{{.Name}}] = z.{{.Name}}
			}
		{{else if eq .Type "uint64"}}
			if z.{{.Name}} != {{.EmptyValue}} {
				result[{{$structure.Name}}{{.Name}}] = strconv.FormatUint(uint64(z.{{.Name}}), 10)
			}
		{{else}}
			if z.{{.Name}} != {{.EmptyValue}} {
				result[{{$structure.Name}}{{.Name}}] = z.{{.Name}}
			}
		{{end}}
	{{end}}
	return result
}

// Deserialize serializes an object to a list of key-values.
func (z *{{.Name}}) Deserialize(values map[string]string) error {
	var errors []error
	z.ID = getUint(uID, values, errors)
	z.Username = getString(uUsername, values, errors)
	z.PasswordHash = getString(uPasswordHash, values, errors)
	z.FeverHash = getString(uFeverHash, values, errors)
	if len(errors) > 0 {
		return errors[0]
	}
	return nil
}

// IndexKeys returns the keys of all indexes for this object.
func (z *{{.Name}}) IndexKeys() map[string][]string {
	result := make(map[string][]string)
	{{range $name, $fields := .Indexes}}
	result[{{$structure.Name}}Index{{$name}}] = []string{
		{{range $index, $f := $fields}} z.{{$f}}, {{end}}
	}
	{{end}}
	return result
}

{{end}}

`
