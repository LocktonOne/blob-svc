package requests

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/spf13/cast"
)

type DeleteBlobRequest struct {
	BlobID int64 `url:"-"`
}

func NewDeleteBlobRequest(r *http.Request) (DeleteBlobRequest, error) {
	request := DeleteBlobRequest{}

	request.BlobID = cast.ToInt64(chi.URLParam(r, "id"))

	return request, nil
}
