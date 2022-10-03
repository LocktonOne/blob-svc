package requests

import (
	"encoding/json"
	"net/http"
	"strings"

	"gitlab.com/tokene/blob-svc/internal/types"
	"gitlab.com/tokene/blob-svc/resources"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
)

type CreateBlobRequest struct {
	Data resources.Blob `json:"data"`
}

func NewCreateBlobRequest(r *http.Request) (resources.Blob, error) {
	var request CreateBlobRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {

		return request.Data, errors.Wrap(err, "failed to unmarshal")
	}
	request.Data.Relationships.Owner.Data.ID = strings.ToLower(request.Data.Relationships.Owner.Data.ID)

	return request.Data, request.validate()
}

func (r CreateBlobRequest) validate() error {
	return validation.Errors{
		"/data/type":                        validation.Validate(&r.Data.Type, validation.Required, validation.In(resources.BLOB)),
		"/data/attributes/blob":             validation.Validate(&r.Data.Attributes.Blob, validation.Required),
		"/data/relationships/owner/data/id": validation.Validate(&r.Data.Relationships.Owner.Data.ID, validation.Required, validation.Match(types.AddressRegexp)),
		"/data/attributes/purpose":          validation.Validate(&r.Data.Attributes.Purpose, validation.Required, types.IsPurpose),
	}.Filter()
}
