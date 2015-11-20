package bolt

import (
	"bytes"
	"fmt"
	"github.com/boltdb/bolt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	chMin      = " "
	chMax      = "~"
	chSep      = "/"
	empty      = ""
	timeFormat = "2006-01-02T15:04:05.000"
)

// metadata contains field names
type metadata struct {
	// name of the struct.
	name string
	// key holds the field used as the primary key.
	key string
	// index holds the fields used to construct the index key, mapped by index name.
	index map[string][]string
}

// summary contains the actual data
type summary struct {
	// name of the struct.
	name string
	// key holds the value of the primary key.
	key string
	// index holds the composite index key values, mapped by index name.
	index map[string][]string
	// deprecated
	value reflect.Value
	// data contains the key/values.
	data map[string]string
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

	- save(name string, pkey string, data map[string]string, tx) error
	- load(name string, pkey string, tx) (map[string]string, error)
	- index(name string, indexesNew map[string][]string, indexesOld map[string][]string, tx) error

	- encode(object interface{}) (metadata, summary, error)
	- decode(object interface{}, data map[string]string) error
	- summarize(metadata, data) (summary, error)
	- getValue(field reflect.Value, fieldName string) (string, error)
	- setValue(field reflect.Value, fieldName string, val string) error

*/

// TODO: extract serialization from database with map[string]string as interface
// add summary struct to describe struct by naming field not holding data (just key and index)
// rename summary to extract

// TODO: expand marshal to include the following steps:
// 	- unmarshal (get previous),
//  - marshal
//  - remove old indexes
//  - add new index entries
// TODO: create function to rebuild indexes

// TODO: need a control table to keep track of schema version,
// need functions to convert from one schema version to the next,
// and that is how to rename fields.

// marshal saves the object with the given ID to the specified bucket.
func marshal(object interface{}, tx *bolt.Tx) error {

	summary, err := getSummary(object)
	if err != nil {
		return err
	}

	keyList := make(map[string]bool)

	logger.Debugf("Bucket name: %s", summary.name)
	bucketData := tx.Bucket([]byte(bucketData))
	b := bucketData.Bucket([]byte(summary.name))

	// loop thru object fields
	for i := 0; i < summary.value.NumField(); i++ {

		field := summary.value.Field(i)
		fieldName := summary.value.Type().Field(i).Name
		keyStr := fmt.Sprintf("%s%s%s", summary.key, chSep, fieldName)

		valueStr, err := getValue(field, fieldName)
		if err != nil {
			return err
		}

		key := []byte(keyStr)
		value := []byte(valueStr)
		if len(value) > 0 {
			keyList[keyStr] = true
			if err := b.Put(key, value); err != nil {
				return err
			}
		}

	} // for loop object fields

	// loop with cursor thru record, remove database fields not in field list
	c := b.Cursor()
	min := []byte(summary.key + chMin)
	max := []byte(summary.key + chMax)
	// logger.Debugf("marshall cursor min/max: %s / %s", min, max)
	for k, _ := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, _ = c.Next() {
		if !keyList[string(k)] {
			// logger.Debugf("Removing key: %s", k)
			c.Delete()
		}
	}

	return nil

}

// unmarshal retrieves the object with the given ID from the specified bucket.
func unmarshal(object interface{}, tx *bolt.Tx) error {

	summary, err := getSummary(object)
	if err != nil {
		return err
	}

	logger.Debugf("Bucket name: %s", summary.name)
	bucketData := tx.Bucket([]byte(bucketData))
	b := bucketData.Bucket([]byte(summary.name))

	// seek to the min key for the given ID
	// loop through cursor until ID changes
	c := b.Cursor()
	min := []byte(summary.key + chMin)
	max := []byte(summary.key + chMax)
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {

		// extract field name from key
		// key format: ID/fieldname
		key := string(k)
		keyParts := strings.SplitN(key, chSep, 2)
		if len(keyParts) != 2 {
			return fmt.Errorf("Malformatted key: %s.", key)
		}
		fieldName := keyParts[1]

		// lookup up field in value
		field := summary.value.FieldByName(fieldName)
		if !field.IsValid() {
			return fmt.Errorf("Invalid fieldname in database: %s.", fieldName)
		}

		if err = setValue(field, fieldName, string(v)); err != nil {
			return err
		}

	} // for loop

	return nil

}

func index(summary *summary, oldsummary *summary, tx *bolt.Tx) error {

	indexBucket := tx.Bucket([]byte(bucketIndex)).Bucket([]byte(summary.name))

	// delete old indexes
	for indexName := range oldsummary.index {
		b := indexBucket.Bucket([]byte(indexName))
		indexElements := oldsummary.index[indexName]
		keyStr := strings.Join(indexElements, chSep)
		if err := b.Delete([]byte(keyStr)); err != nil {
			return err
		}
	}

	// add new indexes
	for indexName := range summary.index {
		b := indexBucket.Bucket([]byte(indexName))
		indexElements := summary.index[indexName]
		keyStr := strings.Join(indexElements, chSep)
		if err := b.Put([]byte(keyStr), []byte(summary.key)); err != nil {
			return err
		}
	}

	return nil

}

func getSummary(object interface{}) (*summary, error) {

	result := &summary{
		index: make(map[string][]string),
	}

	result.name = reflect.TypeOf(object).Name()
	if result.name != "" {
		result.value = reflect.ValueOf(object)
	} else {
		result.value = reflect.ValueOf(object).Elem()
		result.name = result.value.Type().Name()
	}

	if result.name == "" {
		return nil, fmt.Errorf("Cannot get name of object: %v.", object)
	}

	pkeyFound := false

	// loop thru object fields
	for i := 0; i < result.value.NumField(); i++ {

		field := result.value.Field(i)
		typeField := result.value.Type().Field(i)
		fieldName := typeField.Name
		tag := typeField.Tag
		tagFields := strings.Split(tag.Get("db"), ",")

		for _, tagFieldX := range tagFields {
			tagField := strings.TrimSpace(tagFieldX)

			if tagField == "primary-key" {
				if result.key == empty {
					value, err := getValue(field, fieldName)
					if err != nil {
						return nil, err
					}
					result.key = value
					pkeyFound = true
				} else {
					return nil, fmt.Errorf("Duplicate primary key defined for %s.", result.name)
				}

			} else if strings.HasPrefix(tagField, "index") {
				// populate indexes
				elements := strings.SplitN(tagField, ":", 2)
				if len(elements) != 2 {
					return nil, fmt.Errorf("Invalid index definition: %s.", tagField)
				}
				indexName := elements[0][5:] // remove prefix from 0
				indexPosition, err := strconv.Atoi(elements[1])
				if err != nil {
					return nil, fmt.Errorf("Index position is not an integer: %s.", tagField)
				} else if indexPosition < 1 {
					return nil, fmt.Errorf("Index positions are one-based: %s.", tagField)
				}
				indexElements := result.index[indexName]
				for len(indexElements) < indexPosition {
					indexElements = append(indexElements, empty)
				}
				positionValue, err := getValue(field, fieldName)
				if err != nil {
					return nil, err
				}
				indexElements[indexPosition-1] = positionValue
				result.index[indexName] = indexElements
			}

		} // loop tag fields

	} // loop fields

	// if primary key is not found, use the field named ID
	if !pkeyFound {
		fieldName := "ID"
		idValue := result.value.FieldByName(fieldName)
		if idValue.IsValid() {
			value, err := getValue(idValue, fieldName)
			if err != nil {
				return nil, err
			}
			result.key = value
		}
	}

	// validate that primary key is not an empty string
	if result.key == empty {
		return nil, fmt.Errorf("Empty primary key for %s.", result.name)
	}

	return result, nil

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
