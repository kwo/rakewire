package serial

import (
	"fmt"
	"rakewire/logging"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	// TagName ist the name of the struct tag used by this package
	TagName = "db"
)

const (
	chMin      = " "
	chMax      = "~"
	chSep      = "/"
	empty      = ""
	timeFormat = time.RFC3339
)

var (
	logger = logging.New("serial")
)

// Metadata contains field names.
type Metadata struct {
	// Name of the struct.
	Name string
	// Key holds the field used as the primary key.
	Key string
	// Indexes holds the fields used to construct the index keys, mapped by index name.
	Indexes map[string][]string
}

// Data contains the actual data.
type Data struct {
	// Name of the struct.
	Name string
	// Key holds the value of the primary key.
	Key string
	// Indexes holds the composite index key values, mapped by index name.
	Indexes map[string][]string
	// Data contains the key/values.
	Values map[string]string
}

// Encode a pointer to a struct into key/value pairs.
func Encode(object interface{}) (*Metadata, *Data, error) {

	wrapper := reflect.ValueOf(object)
	if wrapper.Kind() != reflect.Ptr {
		return nil, nil, fmt.Errorf("Cannot decode non-pointer object")
	}
	wrapper = wrapper.Elem()

	meta := &Metadata{
		Name:    wrapper.Type().Name(),
		Indexes: make(map[string][]string),
	}
	data := &Data{
		Name:    meta.Name,
		Indexes: make(map[string][]string),
		Values:  make(map[string]string),
	}

	pkeyFound := false

	// loop thru wrapper fields
	for i := 0; i < wrapper.NumField(); i++ {

		field := wrapper.Field(i)
		typeField := wrapper.Type().Field(i)
		fieldName := typeField.Name

		fieldValue, err := getValue(field, fieldName)
		if err != nil {
			return nil, nil, err
		}
		data.Values[fieldName] = fieldValue

		tag := typeField.Tag
		tagFields := strings.Split(tag.Get(TagName), ",")

		for _, tagFieldX := range tagFields {
			tagField := strings.TrimSpace(tagFieldX)

			if tagField == "primary-key" {
				if meta.Key == empty {
					value, err := getValue(field, fieldName)
					if err != nil {
						return nil, nil, err
					}
					meta.Key = fieldName
					data.Key = value
					pkeyFound = true
				} else {
					return nil, nil, fmt.Errorf("Duplicate primary key defined for %s.", meta.Name)
				}

			} else if strings.HasPrefix(tagField, "index") {
				// populate indexes
				elements := strings.SplitN(tagField, ":", 2)
				if len(elements) != 2 {
					return nil, nil, fmt.Errorf("Invalid index definition: %s.", tagField)
				}
				indexName := elements[0][5:] // remove prefix from 0
				indexPosition, err := strconv.Atoi(elements[1])
				if err != nil {
					return nil, nil, fmt.Errorf("Index position is not an integer: %s.", tagField)
				} else if indexPosition < 1 {
					return nil, nil, fmt.Errorf("Index positions are one-based: %s.", tagField)
				}

				metaIndex := meta.Indexes[indexName]
				for len(metaIndex) < indexPosition {
					metaIndex = append(metaIndex, empty)
				}
				metaIndex[indexPosition-1] = fieldName
				meta.Indexes[indexName] = metaIndex

				dataIndex := data.Indexes[indexName]
				for len(dataIndex) < indexPosition {
					dataIndex = append(dataIndex, empty)
				}
				positionValue, err := getValue(field, fieldName)
				if err != nil {
					return nil, nil, err
				}
				dataIndex[indexPosition-1] = positionValue
				data.Indexes[indexName] = dataIndex

			}

		} // loop tag fields

	} // loop fields

	// if primary key is not found, use the field named ID
	if !pkeyFound {
		fieldName := "ID"
		idValue := wrapper.FieldByName(fieldName)
		if idValue.IsValid() {
			value, err := getValue(idValue, fieldName)
			if err != nil {
				return nil, nil, err
			}
			meta.Key = fieldName
			data.Key = value
		}
	}

	// validate that primary key is not an empty string
	if meta.Key == empty {
		return nil, nil, fmt.Errorf("Empty primary key for %s.", meta.Name)
	}

	// validate contiguous indexes
	for name := range meta.Indexes {
		index := meta.Indexes[name]
		for _, field := range index {
			if field == empty {
				return nil, nil, fmt.Errorf("Non-contiguous index names for entity %s, index %s.", meta.Name, name)
			}
		}
	}

	return meta, data, nil

}

// Decode a pointer to a struct from key/value pairs.
func Decode(object interface{}, values map[string]string) error {

	wrapper := reflect.ValueOf(object)
	if wrapper.Kind() != reflect.Ptr {
		return fmt.Errorf("Cannot decode non-pointer object")
	}
	wrapper = wrapper.Elem()

	// loop thru object fields
	for fieldName, v := range values {

		// lookup up field in value
		field := wrapper.FieldByName(fieldName)
		if !field.IsValid() {
			return fmt.Errorf("Invalid fieldname: %s.", fieldName)
		}

		if err := setValue(field, fieldName, v); err != nil {
			return err
		}

	}

	return nil

}

// DataFrom constructs a Data object from Metadata and values (key/value pairs).
func DataFrom(metadata *Metadata, values map[string]string) *Data {

	data := &Data{
		Name:    metadata.Name,
		Indexes: make(map[string][]string),
		Values:  values,
	}

	// set primary key
	data.Key = values[metadata.Key]

	// set indexes
	for index := range metadata.Indexes {
		var elements []string
		fieldNames := metadata.Indexes[index]
		for _, fieldName := range fieldNames {
			elements = append(elements, values[fieldName])
		}
		data.Indexes[index] = elements
	}

	return data

}

func getValue(field reflect.Value, fieldName string) (string, error) {

	switch field.Kind() {

	case reflect.String:
		return field.String(), nil

	case reflect.Bool:
		v := field.Bool()
		if v {
			return strconv.FormatBool(v), nil
		}

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		v := field.Int()
		if v != 0 {
			return strconv.FormatInt(v, 10), nil
		}

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		v := field.Uint()
		if v != 0 {
			return strconv.FormatUint(v, 10), nil
		}

	case reflect.Float32:
		v := field.Float()
		if v != 0 {
			return strconv.FormatFloat(v, 'f', -1, 32), nil
		}

	case reflect.Float64:
		v := field.Float()
		if v != 0 {
			return strconv.FormatFloat(v, 'f', -1, 64), nil
		}

	case reflect.Struct:
		if field.Type() == reflect.TypeOf(time.Time{}) {
			v := field.Interface().(time.Time)
			if !v.IsZero() {
				// dates must be in UTC so that indexes sort properly
				return v.UTC().Format(timeFormat), nil
			}
		} else {
			return empty, fmt.Errorf("Will not get value of struct for %s.", fieldName)
		}

	default:
		return empty, fmt.Errorf("Unknown field type when getting value for %s.", fieldName)

	} // switch

	return empty, nil

}

func setValue(field reflect.Value, fieldName string, val string) error {

	switch field.Kind() {

	case reflect.String:
		field.SetString(val)

	case reflect.Bool:
		value, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		field.SetBool(value)

	case reflect.Int8:
		value, err := strconv.ParseInt(val, 10, 8)
		if err != nil {
			return err
		}
		field.SetInt(value)

	case reflect.Int16:
		value, err := strconv.ParseInt(val, 10, 16)
		if err != nil {
			return err
		}
		field.SetInt(value)

	case reflect.Int32:
		value, err := strconv.ParseInt(val, 10, 32)
		if err != nil {
			return err
		}
		field.SetInt(value)

	case reflect.Int, reflect.Int64:
		value, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(value)

	case reflect.Uint8:
		value, err := strconv.ParseUint(val, 10, 8)
		if err != nil {
			return err
		}
		field.SetUint(value)

	case reflect.Uint16:
		value, err := strconv.ParseUint(val, 10, 16)
		if err != nil {
			return err
		}
		field.SetUint(value)

	case reflect.Uint32:
		value, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			return err
		}
		field.SetUint(value)

	case reflect.Uint, reflect.Uint64:
		value, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(value)

	case reflect.Float32:
		value, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return err
		}
		field.SetFloat(value)

	case reflect.Float64:
		value, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		field.SetFloat(value)

	case reflect.Ptr, reflect.Struct:
		if field.Type() == reflect.TypeOf(time.Time{}) {
			value, err := time.Parse(timeFormat, val)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(value.Truncate(time.Second)))
		} else {
			return fmt.Errorf("Will not set value for struct: %s.", fieldName)
		}

	default:
		return fmt.Errorf("Unknown field type when setting value for %s.", fieldName)

	} // switch

	return nil

}
