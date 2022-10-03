package requests

import (
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/tokene/blob-svc/internal/types"
)

type GetBLobsListRequest struct {
	Owner string
}

func NewGetBLobsListRequest(r *http.Request) (GetBLobsListRequest, error) {
	request := GetBLobsListRequest{}

	request.Owner = strings.ToLower(r.URL.Query().Get("owner"))
	if err := validation.Validate(request.Owner, validation.Required, validation.Match(types.AddressRegexp)); err != nil {
		return request, err
	}

	return request, nil
}

func (r GetBLobsListRequest) validate() error {
	return validation.Validate(r.Owner, validation.Required, validation.Match(types.AddressRegexp))
}
