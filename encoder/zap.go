package encoder

import (
	"reflect"
	"strings"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type CustomEncoder struct {
	zapcore.Encoder
}

func NewCustomEncoder(config zapcore.EncoderConfig) (zapcore.Encoder, error) {
	encoder := zapcore.NewJSONEncoder(config)
	return &CustomEncoder{encoder}, nil
}

// EncodeEntry
// https://stackoverflow.com/questions/73469128/hide-sensitive-fields-in-uber-zap-go
func (e *CustomEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	filtered := make([]zapcore.Field, 0, len(fields))
	for _, field := range fields {
		if strings.ToLower(field.Key) == "password" {
			filtered = append(filtered, zapcore.Field{
				Key:    field.Key,
				Type:   zapcore.StringType,
				String: "[MASKED]",
			})
		} else if field.Interface != nil {
			masked := e.MaskFields(field.Interface)
			field.Interface = masked
			filtered = append(filtered, field)
		} else {
			filtered = append(filtered, field)
		}
	}
	return e.Encoder.EncodeEntry(entry, filtered)
}

// MaskFields
// https://github.com/jinzhu/copier/tree/master
func (e *CustomEncoder) MaskFields(data any) any {
	value := reflect.ValueOf(data)

	switch value.Kind() {
	case reflect.Ptr:
		if value.IsNil() {
			return data
		}
		elem := value.Elem()
		newElem := reflect.New(elem.Type()).Elem()
		newElem.Set(reflect.ValueOf(e.MaskFields(elem.Interface())))
		return newElem.Addr().Interface()
	case reflect.Slice, reflect.Array:
		newData := reflect.MakeSlice(value.Type(), value.Len(), value.Len())
		for i := 0; i < value.Len(); i++ {
			elem := value.Index(i)
			if elem.Kind() == reflect.Ptr && !elem.IsNil() {
				elem = elem.Elem()
				newElem := reflect.New(elem.Type())
				newElem.Elem().Set(reflect.ValueOf(e.MaskFields(elem.Interface())))
				newData.Index(i).Set(newElem)
			} else {
				newData.Index(i).Set(reflect.ValueOf(e.MaskFields(elem.Interface())))
			}
		}
		return newData.Interface()
	case reflect.Map:
		newData := reflect.MakeMapWithSize(value.Type(), value.Len())
		for _, key := range value.MapKeys() {
			if key.Kind() == reflect.String && e.IsSensitive(key.String()) {
				newData.SetMapIndex(key, reflect.ValueOf(e.MaskSensitiveField(key.String())))
				continue
			}
			val := value.MapIndex(key)
			if val.Kind() == reflect.Ptr && !val.IsNil() {
				val = val.Elem()
			}
			maskedVal := e.MaskFields(val.Interface())

			elemType := value.Type().Elem()
			if elemType.Kind() == reflect.Ptr && elemType.Elem().Kind() == reflect.Struct {
				newElem := reflect.New(elemType.Elem())
				newElem.Elem().Set(reflect.ValueOf(maskedVal))
				newData.SetMapIndex(key, newElem)
			} else if elemType.Kind() == reflect.Struct {
				newElem := reflect.New(elemType).Elem()
				newElem.Set(reflect.ValueOf(maskedVal))
				newData.SetMapIndex(key, newElem)
			} else {
				newData.SetMapIndex(key, reflect.ValueOf(maskedVal))
			}
		}
		return newData.Interface()
	case reflect.Struct:
		newData := reflect.New(value.Type()).Elem()
		for i := 0; i < value.NumField(); i++ {
			field := value.Field(i)
			if field.Kind() == reflect.Ptr && !field.IsNil() {
				field = field.Elem()
			}
			newData.Field(i).Set(reflect.ValueOf(e.MaskFields(field.Interface())))
			fieldName := value.Type().Field(i).Name
			if field.Kind() == reflect.String && strings.Contains(strings.ToLower(fieldName), "password") {
				newData.Field(i).SetString("[MASKED]")
			}
		}
		return newData.Interface()
	}

	return data
}
