package requests

import (
	"encoding/json"
	"net/http"

	"gitlab.com/tokene/blob-svc/internal/types"

	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/tokene/blob-svc/resources"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

func NewСreateDocumentRequest(r *http.Request) (resources.Document, error) {
	var request resources.DocumentResponse
	err := r.ParseMultipartForm(1 << 32)
	if err != nil {
		return request.Data, errors.Wrap(err, "multipart limit to unmarshal")
	}
	if err := json.Unmarshal([]byte(r.FormValue("Document")), &request); err != nil {
		return request.Data, errors.Wrap(err, "multipart limit to unmarshal")
	}
	return request.Data, ValidatePutDocumentRequest(request.Data)
}

func ValidatePutDocumentRequest(r resources.Document) error {
	return validation.Errors{
		"/data/type":                        validation.Validate(&r.Type, validation.Required),
		"/data/relationships/owner/data/id": validation.Validate(&r.Relationships.Owner.Data.ID, validation.Required, validation.Match(types.AddressRegexp)),
		"/data/attributes/purpose":          validation.Validate(&r.Attributes.Purpose, validation.Required, types.IsPurpose),
	}.Filter()
}
