package validators

import (
	"encoding/json"
	"reflect"

	"github.com/go-playground/validator/v10"
)

func JSONValidator(fl validator.FieldLevel) bool {
	v := fl.Field()

	// Если nil, или zero, возвращаем false
	if !v.IsValid() {
		return false
	}

	// Пример: если это интерфейс / указатель, разворачиваем
	kind := v.Kind()
	// Если интерфейс или указатель — получаем подлежащие значения
	if kind == reflect.Interface || kind == reflect.Ptr {
		v = v.Elem()
		kind = v.Kind()
	}

	var data []byte

	switch kind {
	case reflect.String:
		// если строка — используем String()
		data = []byte(v.String())
	case reflect.Slice:
		// если слайс байтов — допустим, RawMessage
		// Нужно проверить, что это []uint8
		if v.Type().Elem().Kind() == reflect.Uint8 {
			data = v.Bytes()
		} else {
			// другой слайс — не то, не валидируем как JSON
			return false
		}
	default:
		// Тип не поддерживается
		return false
	}

	// Проверяем, валидный ли JSON
	var js json.RawMessage

	return json.Unmarshal(data, &js) == nil
}
