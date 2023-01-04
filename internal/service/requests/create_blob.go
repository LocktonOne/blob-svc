package requests

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"gitlab.com/tokene/blob-svc/internal/types"
	"net/http"
)

type CreateBlobRequest struct {
	Blob    json.RawMessage `json:"blob"`
	Purpose string          `json:"purpose"`
	Owner   string          `json:"owner"`
}

func NewCreateBlobRequest(r *http.Request) (CreateBlobRequest, error) {
	request := struct {
		Data CreateBlobRequest `json:"data"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {

		return request.Data, errors.Wrap(err, "failed to unmarshal")
	}

	return request.Data, request.Data.validate()
}

func (r CreateBlobRequest) validate() error {
	return validation.Errors{
		"/data/attributes/blob":             validation.Validate(&r.Blob, validation.Required),
		"/data/relationships/owner/data/id": validation.Validate(&r.Owner, validation.Required, validation.Match(types.AddressRegexp)),
		"/data/attributes/purpose":          validation.Validate(&r.Purpose, validation.Required),
	}.Filter()
}
