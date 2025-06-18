package rest

import (
	"CookFinder.Backend/pkg/puberr"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io"

	"github.com/go-playground/mold/v4/modifiers"
)

var (
	modifier = modifiers.New()
	validate = validator.New()
)

// MapJSON кроме маршалинга строки в объкт, также проводит валидацию, модификацию и обработку частых ошибок
func MapJSON(r io.Reader, target any) error {
	err := json.NewDecoder(r).Decode(target)
	if err != nil {
		var errSyntax *json.SyntaxError
		var errUnmarshal *json.UnmarshalTypeError

		switch {
		case errors.Is(err, io.EOF), errors.Is(err, io.ErrUnexpectedEOF):
			return puberr.NewPubErr("unexpected EOF").SetCause(err)
		case errors.As(err, &errSyntax):
			return puberr.NewPubErr(err.Error()).SetCause(err)
		case errors.As(err, &errUnmarshal):
			msg := fmt.Sprintf("incorrect JSON type for field %q at %d", errUnmarshal.Field, errUnmarshal.Offset)
			return puberr.NewPubErr(msg).SetCause(err)
		default:
			return puberr.ErrGeneric
		}
	}

	if err = modifier.Struct(context.Background(), target); err != nil {
		return puberr.NewPubErr(err.Error())
	}

	if err = validate.Struct(target); err != nil {
		return puberr.NewPubErr(err.Error())
	}

	return nil
}
