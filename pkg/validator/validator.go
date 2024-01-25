package validator

import (
	errors "coffee-shop/pkg/custom_errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type validatorStruct struct {
	v *validator.Validate
}

type ValidatorInterface interface {
	Struct(interface{}) error
	Map(map[string]interface{}, map[string]map[string]string) error
}

func NewValidator() ValidatorInterface {
	return &validatorStruct{
		v: validator.New(),
	}
}

func (vs *validatorStruct) Struct(value interface{}) error {
	err := vs.v.Struct(value)
	if err != nil {
		return err
	}
	return nil
}

func (vs *validatorStruct) Map(inputMap map[string]interface{}, validationMap map[string]map[string]string) error {
	for key, value := range inputMap {
		validations, ok := validationMap[key]
		if !ok {
			return errors.BadRequest.Newf(nil, false, "Invalid key :%s", key)
		}
		dataType := fmt.Sprintf("%T", value)
		if dataType == "[]interface {}" {
			temp := value.([]interface{})
			if len(temp) == 0 {
				return errors.BadRequest.Newf(nil, false, "Invalid value for key %s", key)
			}
			dataType = fmt.Sprintf("[]%T", temp[0])
		}

		if !strings.Contains(validations["data_type"], dataType) {
			return errors.BadRequest.Newf(nil, false, "Invalid type of value for key %s", key)
		}

		if dataType == "bool" {
			_, ok := value.(bool)
			if !ok {
				return errors.BadRequest.Newf(nil, false, "Invalid value for key %s", key)
			}
		} else if dataType == "<nil>" {
			if value != nil {
				return errors.BadRequest.Newf(nil, false, "Invalid value for key %s", key)
			}
		} else {
			err := vs.v.Var(value, validations["rules"])
			if err != nil {
				return errors.BadRequest.Wrapf(nil, false, err, "Invalid value for key %s", key)
			}
		}
	}

	return nil
}
