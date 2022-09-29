package requests

import (
	"encoding/json"
	"net/http"
	"strings"

	"gitlab.com/tokene/blob-svc/internal/types"

	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/tokene/blob-svc/resources"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

type CreateDocumentRequest struct {
	Data resources.Document `json:"data"`
}

func NewСreateDocumentRequest(r *http.Request) (resources.Document, error) {
	var request CreateDocumentRequest
	err := r.ParseMultipartForm(1 << 32)
	if err != nil {
		return request.Data, errors.Wrap(err, "multipart limit to unmarshal")
	}

	if err := json.Unmarshal([]byte(r.FormValue("Document")), &request); err != nil {
		return request.Data, errors.Wrap(err, "multipart limit to unmarshal")
	}

	request.Data.Relationships.Owner.Data.ID = strings.ToLower(request.Data.Relationships.Owner.Data.ID)

	return request.Data, request.validate()
}

func (r CreateDocumentRequest) validate() error {
	return validation.Errors{
		"/data/type":                        validation.Validate(&r.Data.Type, validation.Required),
		"/data/relationships/owner/data/id": validation.Validate(&r.Data.Relationships.Owner.Data.ID, validation.Required, validation.Match(types.AddressRegexp)),
		"/data/attributes/purpose":          validation.Validate(&r.Data.Attributes.Purpose, validation.Required, types.IsPurpose),
	}.Filter()
}
