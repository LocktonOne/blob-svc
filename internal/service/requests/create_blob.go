package requests

import (
	"encoding/json"
	"net/http"

	"gitlab.com/tokene/blob-svc/internal/service/helpers"
	"gitlab.com/tokene/blob-svc/resources"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pkg/errors"
)

type CreateBlobRequest struct {
	Data resources.Blob `url:"-"`
}

func NewCreateBlobRequest(r *http.Request) (CreateBlobRequest, error) {
	request := CreateBlobRequest{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {

		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, request.validate()
}

func (r *CreateBlobRequest) validate() error {
	return mergeErrors(validation.Errors{
		"/data/attributes/blob":               validation.Validate(string([]byte(r.Data.Attributes.Blob)), validation.Required, is.JSON),
		"/data/attributes/user_address":       validation.Validate(string([]byte(r.Data.Attributes.UserAddress)), validation.Required, validation.Match(helpers.AddressRegexp)),
		"/data/type":                          validation.Validate(string([]byte(r.Data.Key.Type)), validation.In("blob")),
		"/data/relationships/owner/data/id":   validation.Validate(string([]byte(r.Data.Relationships.Owner.Data.ID)), validation.Required, validation.Length(1, 255)),
		"/data/relationships/owner/data/type": validation.Validate(string([]byte(r.Data.Relationships.Owner.Data.Type)), validation.In("user")),
	},
	).Filter()
}
func mergeErrors(validationErrors ...validation.Errors) validation.Errors {
	result := make(validation.Errors)
	for _, errs := range validationErrors {
		for key, err := range errs {
			result[key] = err
		}
	}
	return result
}
