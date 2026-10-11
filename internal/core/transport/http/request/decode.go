package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	core_errors "gitverse.ru/IvanVerse/go-todoApp/internal/core/errors"
)

var requestValidator = validator.New()

type validatable interface {
	Validate() error
}

func DecodeAndValidateRequest(r *http.Request, dest any) error {

	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode json: %v:%w", err, core_errors.ErrInvalidArgument)
	}

	v, ok := dest.(validatable)
	if ok {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("request validation:%v:%w", err, core_errors.ErrInvalidArgument)
		}
	} else {
		if err := requestValidator.Struct(dest); err != nil {
			return fmt.Errorf("request validation:%v:%w", err, core_errors.ErrInvalidArgument)
		}
	}

	return nil
}
