package api

import (
	// "github.com/go-playground/locales/currency"
	//TODO:: will use the go-currency package later
	"github.com/MacbotX/simplebank_v1/util"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool{
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		//  check currency is supported 
		return util.IsSupportedCurrency(currency)
	}
	return false
}