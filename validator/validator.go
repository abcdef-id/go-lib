package validator

import (
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

var lock = &sync.Mutex{}
var validate *validator.Validate

//GetValidator Initiatilize validator in singleton way
func GetValidator() *validator.Validate {
	lock.Lock()
	defer lock.Unlock()

	if validate == nil {
		validate = validator.New()
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}

	return validate
}
