package requests

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/spf13/cast"
)

type GetBlobIDRequest struct {
	BlobID int64 `url:"-"`
}

func NewGetBlobIDRequest(r *http.Request) (GetBlobIDRequest, error) {
	request := GetBlobIDRequest{}

	request.BlobID = cast.ToInt64(chi.URLParam(r, "id"))
	return request, nil
}
