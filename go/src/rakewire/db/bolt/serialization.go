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
	primaryKey = "primary-key"
)

type metadata struct {
	name    string
	key     string
	indexes map[string][]string
	value   reflect.Value
}

// TODO: also use tags to define indexes (`db:"#indexFeedByURL:1"`) for field 1 of named index
// TODO: write functions to deliver the pk and indexes array
// TODO: remember - the value of an index is always the primary-key
// TODO: expand marshal to include the following steps:
// 	- unmarshal (get previous),
//  - marshal
//  - clean old keys
//  - remove old indexes
//  - add new index entries
// TODO: methods to retrieve values by index must still be custom made
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
		fieldType := meta.value.Type().Field(i)
		key := []byte(fmt.Sprintf("%s%s%s", meta.key, chSep, fieldType.Name))

		switch fieldType.Type.Kind() {

		case reflect.String:
			value := []byte(field.String())
			if len(value) > 0 {
				if err := b.Put(key, value); err != nil {
					return err
				}
			}
			break

		case reflect.Bool:
			v := field.Bool()
			if v {
				value := []byte(strconv.FormatBool(v))
				if err := b.Put(key, value); err != nil {
					return err
				}
			}
			break

		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			v := field.Int()
			if v != 0 {
				value := []byte(strconv.FormatInt(v, 10))
				if err := b.Put(key, value); err != nil {
					return err
				}
			}
			break

		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
			v := field.Uint()
			if v != 0 {
				value := []byte(strconv.FormatUint(v, 10))
				if err := b.Put(key, value); err != nil {
					return err
				}
			}
			break

		case reflect.Float32:
			v := field.Float()
			if v != 0 {
				value := []byte(strconv.FormatFloat(v, 'f', -1, 32))
				if err := b.Put(key, value); err != nil {
					return err
				}
			}
			break

		case reflect.Float64:
			v := field.Float()
			if v != 0 {
				value := []byte(strconv.FormatFloat(v, 'f', -1, 64))
				if err := b.Put(key, value); err != nil {
					return err
				}
			}
			break

		case reflect.Struct:
			if fieldType.Type == reflect.TypeOf(time.Time{}) {
				v := field.Interface().(time.Time)
				if !v.IsZero() {
					value := []byte(v.UTC().Format(timeFormat))
					if err := b.Put(key, value); err != nil {
						return err
					}
				}
			} else {
				return errors.New("Will not marshal struct for " + fieldType.Name)
			}
			break

		default:
			return errors.New("Unknown field type when marshaling value for " + fieldType.Name)

		} // switch

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

		// get type of field
		// parse byte value and apply to object
		val := string(v)
		switch field.Kind() {

		case reflect.String:
			field.SetString(val)
			break

		case reflect.Bool:
			value, err := strconv.ParseBool(val)
			if err != nil {
				return err
			}
			field.SetBool(value)
			break

		case reflect.Int8:
			value, err := strconv.ParseInt(val, 10, 8)
			if err != nil {
				return err
			}
			field.SetInt(value)
			break

		case reflect.Int16:
			value, err := strconv.ParseInt(val, 10, 16)
			if err != nil {
				return err
			}
			field.SetInt(value)
			break

		case reflect.Int32:
			value, err := strconv.ParseInt(val, 10, 32)
			if err != nil {
				return err
			}
			field.SetInt(value)
			break

		case reflect.Int, reflect.Int64:
			value, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return err
			}
			field.SetInt(value)
			break

		case reflect.Uint8:
			value, err := strconv.ParseUint(val, 10, 8)
			if err != nil {
				return err
			}
			field.SetUint(value)
			break

		case reflect.Uint16:
			value, err := strconv.ParseUint(val, 10, 16)
			if err != nil {
				return err
			}
			field.SetUint(value)
			break

		case reflect.Uint32:
			value, err := strconv.ParseUint(val, 10, 32)
			if err != nil {
				return err
			}
			field.SetUint(value)
			break

		case reflect.Uint, reflect.Uint64:
			value, err := strconv.ParseUint(val, 10, 64)
			if err != nil {
				return err
			}
			field.SetUint(value)
			break

		case reflect.Float32:
			value, err := strconv.ParseFloat(val, 32)
			if err != nil {
				return err
			}
			field.SetFloat(value)
			break

		case reflect.Float64:
			value, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return err
			}
			field.SetFloat(value)
			break

		case reflect.Ptr, reflect.Struct:

			if field.Type() == reflect.TypeOf(time.Time{}) {
				value, err := time.ParseInLocation(timeFormat, val, time.UTC)
				if err != nil {
					return err
				}
				field.Set(reflect.ValueOf(value.Truncate(time.Millisecond)))
			} else {
				return errors.New("Will not unmarshal struct for " + fieldName)
			}
			break

		default:
			return errors.New("Unknown field type when unmarshaling value for " + fieldName)

		} // switch

	} // for loop

	return nil

}

func getMetadata(object interface{}) (metadata, error) {

	result := metadata{
		indexes: make(map[string][]string),
	}

	result.name = reflect.TypeOf(object).Name()
	if result.name != "" {
		result.value = reflect.ValueOf(object)
	} else {
		result.value = reflect.ValueOf(object).Elem()
		result.name = result.value.Type().Name()
	}

	if result.name == "" {
		return result, errors.New("Cannot get name of object")
	}

	// loop thru object fields
	for i := 0; i < result.value.NumField(); i++ {

		field := result.value.Field(i)
		fieldType := result.value.Type().Field(i)

		tag := fieldType.Tag
		dbFields := strings.Split(tag.Get("db"), ",")
		for _, dbField := range dbFields {
			if strings.TrimSpace(dbField) == primaryKey {
				// TODO: allow others types
				result.key = field.String()
				break
			}
		}

	}

	return result, nil

}
