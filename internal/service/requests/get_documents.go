package requests

import (
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/tokene/blob-svc/internal/types"
)

type GetDocumentsListRequest struct {
	Owner string
}

func NewGetDocumentsListRequest(r *http.Request) (GetDocumentsListRequest, error) {
	request := GetDocumentsListRequest{}

	request.Owner = strings.ToLower(r.URL.Query().Get("owner"))

	return request, request.validate()
}

func (r GetDocumentsListRequest) validate() error {
	return validation.Validate(r.Owner, validation.Required, validation.Match(types.AddressRegexp))
}
