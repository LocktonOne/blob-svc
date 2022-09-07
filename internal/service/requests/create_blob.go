package requests

import (
	"encoding/json"
	"net/http"

	"gitlab.com/tokene/blob-svc/internal/types"
	"gitlab.com/tokene/blob-svc/resources"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
)

func NewCreateBlobRequest(r *http.Request) (resources.Blob, error) {
	request := resources.BlobResponse{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {

		return request.Data, errors.Wrap(err, "failed to unmarshal")
	}

	return request.Data, validate(request.Data)
}

func validate(r resources.Blob) error {
	return validation.Errors{
		"/data/type":                        validation.Validate(&r.Type, validation.Required, validation.In(resources.BLOB)),
		"/data/attributes/blob":             validation.Validate(&r.Attributes.Blob, validation.Required),
		"/data/relationships/owner/data/id": validation.Validate(&r.Relationships.Owner.Data.ID, validation.Required, validation.Match(types.AddressRegexp)),
		"/data/attributes/purpose":          validation.Validate(&r.Attributes.Purpose, validation.Required, types.IsPurpose),
	}.Filter()
}
