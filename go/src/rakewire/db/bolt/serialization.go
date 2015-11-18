package bolt

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	m "rakewire/model"
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
func marshal(object m.Identifiable, tx *bolt.Tx) error {

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
		key := []byte(fmt.Sprintf("%s%s%s", object.GetID(), chSep, fieldType.Name))

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
func unmarshal(object m.Identifiable, tx *bolt.Tx) error {

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
	min := []byte(object.GetID() + chMin)
	max := []byte(object.GetID() + chMax)
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
