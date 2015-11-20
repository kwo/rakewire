package bolt

import (
	"bytes"
	"errors"
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
	timeFormat = "2006-01-02T15:04:05.000"
	empty      = ""
)

type metadata struct {
	name  string
	key   string
	index map[string][]string
	value reflect.Value
}

// TODO: expand marshal to include the following steps:
// 	- unmarshal (get previous),
//  - marshal
//  - clean old keys
//  - remove old indexes
//  - add new index entries
// TODO: store indexes in Index/entity-name/index-name
// TODO: create function to rebuild indexes

// TODO: need a control table to keep track of schema version,
// need functions to convert from one schema version to the next,
// and that is how to rename fields.

// marshal saves the object with the given ID to the specified bucket.
func marshal(object interface{}, tx *bolt.Tx) error {

	meta, err := getMetadata(object)
	if err != nil {
		return err
	}

	logger.Debugf("Bucket name: %s", meta.name)
	b := tx.Bucket([]byte(meta.name))

	// loop thru object fields
	for i := 0; i < meta.value.NumField(); i++ {

		field := meta.value.Field(i)
		fieldName := meta.value.Type().Field(i).Name
		key := []byte(fmt.Sprintf("%s%s%s", meta.key, chSep, fieldName))

		valueStr, err := getValue(field, fieldName)
		if err != nil {
			return err
		}

		value := []byte(valueStr)
		if len(value) > 0 {
			if err := b.Put(key, value); err != nil {
				return err
			}
		}

		// TODO: add each valid value-holding field to list

	} // for loop

	// TODO: loop with cursor thru record, remove database fields not in field list

	return nil

}

// unmarshal retrieves the object with the given ID from the specified bucket.
func unmarshal(object interface{}, tx *bolt.Tx) error {

	meta, err := getMetadata(object)
	if err != nil {
		return err
	}

	logger.Debugf("Bucket name: %s", meta.name)
	b := tx.Bucket([]byte(meta.name))

	// seek to the min key for the given ID
	// loop through cursor until ID changes
	c := b.Cursor()
	min := []byte(meta.key + chMin)
	max := []byte(meta.key + chMax)
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {

		// extract field name from key
		// key format: ID/fieldname
		key := string(k)
		keyParts := strings.SplitN(key, chSep, 2)
		if len(keyParts) != 2 {
			return errors.New("Malformatted key: " + key)
		}
		fieldName := keyParts[1]

		// lookup up field in value
		field := meta.value.FieldByName(fieldName)
		if !field.IsValid() {
			return errors.New("Invalid fieldname in database: " + fieldName)
		}

		if err = setValue(field, fieldName, string(v)); err != nil {
			return err
		}

	} // for loop

	return nil

}

func getMetadata(object interface{}) (*metadata, error) {

	result := &metadata{
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
		return nil, errors.New("Cannot get name of object")
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
					return nil, errors.New("Duplicate primary key defined for " + result.name)
				}

			} else if strings.HasPrefix(tagField, "index") {
				// populate indexes
				elements := strings.SplitN(tagField, ":", 2)
				if len(elements) != 2 {
					return nil, errors.New("Invalid index definition: " + tagField)
				}
				indexName := elements[0][5:] // remove prefix from 0
				indexPosition, err := strconv.Atoi(elements[1])
				if err != nil {
					return nil, errors.New("Index position is not an integer: " + tagField)
				} else if indexPosition < 1 {
					return nil, errors.New("Index positions are one-based: " + tagField)
				}
				indexElements := result.index[indexName]
				for len(indexElements) < indexPosition {
					indexElements = append(indexElements, "")
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
		return nil, errors.New("Empty primary key!")
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
			return empty, errors.New("Will not get value of struct for " + fieldName)
		}

	default:
		return empty, errors.New("Unknown field type when getting value for " + fieldName)

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
			return errors.New("Will not set value for struct: " + fieldName)
		}

	default:
		return errors.New("Unknown field type when setting value for " + fieldName)

	} // switch

	return nil

}
