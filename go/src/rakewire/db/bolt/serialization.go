package bolt

import (
	"bytes"
	"errors"
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

// Marshal saves the object with the given ID to the specified bucket.
// TODO: use cursor instead of bucket
func (z *Database) Marshal(v interface{}, ID string, b *bolt.Bucket) error {

	// if err = b.Put([]byte(formatFeedLogKey(id, &entry.StartTime)), data); err != nil {
	// 	return err
	// }

	return nil
}

// Unmarshal retrieves the object with the given ID from the specified bucket.
func (z *Database) Unmarshal(object interface{}, ID string, c *bolt.Cursor) error {

	var obj reflect.Value
	if reflect.TypeOf(object).Name() != "" {
		obj = reflect.ValueOf(object)
	} else {
		obj = reflect.ValueOf(object).Elem()
	}

	// seek to the min key for the given ID
	// loop through cursor until ID changes
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
				return errors.New("Will not unmarshall struct for " + fieldName)
			}
			break

		default:
			return errors.New("Unknown field type when unmarshalling value for " + fieldName)

		} // switch

	} // for loop

	return nil

}
