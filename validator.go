package clop

import (
	"reflect"
	"sort"
	"strings"
	"sync"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	en_translations "github.com/go-playground/validator/v10/translations/en"

	"github.com/go-playground/validator/v10"
)

var valid *defaultValidator = &defaultValidator{}

type defaultValidator struct {
	once     sync.Once
	validate *validator.Validate
	trans    ut.Translator
}

func (v *defaultValidator) ValidateStruct(obj interface{}) error {

	if kindOfData(obj) == reflect.Struct {

		v.lazyinit()

		if err := v.validate.Struct(obj); err != nil {
			return err
		}
	}

	return nil
}

func (v *defaultValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func showShortLongUsage(clopName string, tagName string) string {
	var usage []string

	opt := strings.Split(clopName, ";")

	for _, o := range opt {

		o = strings.TrimSpace(o)
		if len(o) == 0 {
			continue
		}

		switch o {
		case "short":
			usage = append(usage, "-"+strings.ToLower(string(tagName[0])))

			continue
		case "long":
			if len(tagName) > 1 {
				longName, _ := gnuOptionName(tagName)
				usage = append(usage, "--"+longName)
			}
			continue
		}

		if o[0] != '-' {
			continue
		}

		usage = append(usage, o)
	}

	sort.Slice(usage, func(i, j int) bool {
		return len(usage[i]) < len(usage[j])
	})

	return strings.Join(usage, ";")
}

func (v *defaultValidator) lazyinit() {
	v.once.Do(func() {
		en := en.New()
		uni := ut.New(en, en)

		v.validate = validator.New()
		v.validate.SetTagName("valid")
		v.trans, _ = uni.GetTranslator("en")
		en_translations.RegisterDefaultTranslations(v.validate, v.trans)

		v.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			return "error: " + showShortLongUsage(fld.Tag.Get("clop"), fld.Name)
		})

		v.validate.RegisterTranslation("required", v.trans, func(ut ut.Translator) error {
			return ut.Add("required", "{0} must have a value!", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("required", fe.Field())

			return t
		})

		// add any custom validations etc. here
	})
}

func kindOfData(data interface{}) reflect.Kind {

	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}
