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
)

// TODO: use tags to indicate which fields are indexed (how would it work?)

// marshal saves the object with the given ID to the specified bucket.
func marshal(object interface{}, ID string, tx *bolt.Tx) error {

	// get reflection info from object
	var obj reflect.Value
	objectName := reflect.TypeOf(object).Name()
	if objectName != "" {
		obj = reflect.ValueOf(object)
	} else {
		obj = reflect.ValueOf(object).Elem()
		objectName = obj.Type().Name()
	}

	if objectName == "" {
		return errors.New("Cannot get name of object")
	}

	logger.Debugf("Bucket name: %s", objectName)
	b := tx.Bucket([]byte(objectName))

	// loop thru object fields
	for i := 0; i < obj.NumField(); i++ {

		field := obj.Field(i)
		fieldType := obj.Type().Field(i)
		key := []byte(fmt.Sprintf("%s%s%s", ID, chSep, fieldType.Name))

		switch fieldType.Type.Kind() {

		case reflect.String:
			value := []byte(field.String())
			if err := b.Put(key, value); err != nil {
				return err
			}
			break

		case reflect.Bool:
			value := []byte(strconv.FormatBool(field.Bool()))
			if err := b.Put(key, value); err != nil {
				return err
			}
			break

		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			value := []byte(strconv.FormatInt(field.Int(), 10))
			if err := b.Put(key, value); err != nil {
				return err
			}
			break

		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
			value := []byte(strconv.FormatUint(field.Uint(), 10))
			if err := b.Put(key, value); err != nil {
				return err
			}
			break

		case reflect.Float32:
			value := []byte(strconv.FormatFloat(field.Float(), 'f', -1, 32))
			if err := b.Put(key, value); err != nil {
				return err
			}
			break

		case reflect.Float64:
			value := []byte(strconv.FormatFloat(field.Float(), 'f', -1, 64))
			if err := b.Put(key, value); err != nil {
				return err
			}
			break

		case reflect.Struct:
			if fieldType.Type == reflect.TypeOf(time.Time{}) {
				value := []byte(field.Interface().(time.Time).UTC().Format(timeFormat))
				if err := b.Put(key, value); err != nil {
					return err
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
func unmarshal(object interface{}, ID string, tx *bolt.Tx) error {

	// get reflection info from object
	var obj reflect.Value
	objectName := reflect.TypeOf(object).Name()
	if objectName != "" {
		obj = reflect.ValueOf(object)
	} else {
		obj = reflect.ValueOf(object).Elem()
		objectName = obj.Type().Name()
	}

	if objectName == "" {
		return errors.New("Cannot get name of object")
	}

	logger.Debugf("Bucket name: %s", objectName)
	b := tx.Bucket([]byte(objectName))

	// seek to the min key for the given ID
	// loop through cursor until ID changes
	c := b.Cursor()
	min := []byte(ID + chMin)
	max := []byte(ID + chMax)
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
		field := obj.FieldByName(fieldName)
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
