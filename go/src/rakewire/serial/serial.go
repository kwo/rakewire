package serial

import (
	"fmt"
	"rakewire/logging"
	"reflect"
	"strconv"
	"time"
)

const (
	chMin      = " "
	chMax      = "~"
	chSep      = "/"
	empty      = ""
	timeFormat = "2006-01-02T15:04:05.000"
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
	// Index holds the fields used to construct the index key, mapped by index name.
	Index map[string][]string
}

// Summary contains the actual data.
type Summary struct {
	// Name of the struct.
	Name string
	// Key holds the value of the primary key.
	Key string
	// Index holds the composite index key values, mapped by index name.
	Index map[string][]string
	// Data contains the key/values.
	Data map[string]string
}

func mar(o interface{}) {

	// encode (new) object
	// encode(o) -> (name string, pkey string, data map[string]string, indexes map[string][]string, error)

	// get existing data
	// load(name, pkey) -> data map[string]string

	// get metadata from object - output from encode?
	// construct summary from metadata and data for old

	// now I have enough info to delete old indexes
	// save
	// index

}

/*

	- Encode(object interface{}) (metadata, summary, error)
	- Decode(object interface{}, data map[string]string) error
	- Summarize(metadata, data) (summary, error)
	- getValue(field reflect.Value, fieldName string) (string, error)
	- setValue(field reflect.Value, fieldName string, val string) error

*/

// Encode a struct into key/value pairs.
func Encode(object interface{}) (*Metadata, *Summary, error) {
	return nil, nil, nil
}

// Decode a struct from key/value pairs.
func Decode(object interface{}, data map[string]string) error {

	value := reflect.ValueOf(object).Elem()

	// loop thru object fields
	for fieldName, v := range data {

		// lookup up field in value
		field := value.FieldByName(fieldName)
		if !field.IsValid() {
			return fmt.Errorf("Invalid fieldname: %s.", fieldName)
		}

		if err := setValue(field, fieldName, v); err != nil {
			return err
		}

	}

	return nil

}

// Summarize constructs a Summary object from Metadata and data (key/value pairs).
func Summarize(metadata *Metadata, data map[string]string) *Summary {

	summary := &Summary{
		Name:  metadata.Name,
		Index: make(map[string][]string),
		Data:  data,
	}

	// set primary key
	summary.Key = data[metadata.Key]

	// set indexes
	for indexName := range metadata.Index {
		var elements []string
		fieldNames := metadata.Index[indexName]
		for _, fieldName := range fieldNames {
			elements = append(elements, data[fieldName])
		}
		summary.Index[indexName] = elements
	}

	return summary

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
			value, err := time.ParseInLocation(timeFormat, val, time.UTC)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(value.Truncate(time.Millisecond)))
		} else {
			return fmt.Errorf("Will not set value for struct: %s.", fieldName)
		}

	default:
		return fmt.Errorf("Unknown field type when setting value for %s.", fieldName)

	} // switch

	return nil

}
