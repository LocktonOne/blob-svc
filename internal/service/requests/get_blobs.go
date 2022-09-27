package requests

import (
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/tokene/blob-svc/internal/types"
)

type GetBLobsListRequest struct {
	OwnerFilter string `filter:"owner"`
}

func NewGetBLobsListRequest(r *http.Request) (GetBLobsListRequest, error) {
	request := GetBLobsListRequest{}

	request.OwnerFilter = strings.ToLower(r.URL.Query().Get("owner"))
	if err := validation.Validate(request.OwnerFilter, validation.Required, validation.Match(types.AddressRegexp)); err != nil {
		return request, err
	}

	return request, nil
}
