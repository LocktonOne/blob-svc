package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
)

type GetBLobsListRequest struct {
	pgdb.OffsetPageParams
	FilterAddress []string `filter:"address"`
}

func NewGetBLobsListRequest(r *http.Request) (GetBLobsListRequest, error) {
	request := GetBLobsListRequest{}

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	return request, nil
}
