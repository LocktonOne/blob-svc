package requests

import (
	"encoding/json"
	"gitlab.com/tokene/blob-svc/resources"
	"net/http"
	"strings"

	"gitlab.com/tokene/blob-svc/internal/types"

	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type CreateDocumentRequest struct {
	ContentType string                 `json:"content_type"`
	Owner       string                 `json:"owner"`
	Type        resources.ResourceType `json:"type"`
}

func NewCreateDocumentRequest(r *http.Request) (CreateDocumentRequest, error) {
	request := struct {
		Data CreateDocumentRequest `json:"data"`
	}{}

	err := r.ParseMultipartForm(1 << 32)
	if err != nil {
		return CreateDocumentRequest{}, errors.Wrap(err, "multipart limit to unmarshal")
	}

	if err := json.Unmarshal([]byte(r.FormValue("Document")), &request); err != nil {
		return CreateDocumentRequest{}, errors.Wrap(err, "multipart limit to unmarshal")
	}

	request.Data.Owner = strings.ToLower(request.Data.Owner)

	return request.Data, request.Data.validate()
}

func (r CreateDocumentRequest) validate() error {
	return validation.Errors{
		"/data/type":                        validation.Validate(&r.Type, validation.Required),
		"/data/relationships/owner/data/id": validation.Validate(&r.Owner, validation.Required, validation.Match(types.AddressRegexp)),
	}.Filter()
}
