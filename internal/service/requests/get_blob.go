package requests

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/spf13/cast"
)

type GetBlobIDRequest struct {
	BlobID int64 `url:"-"`
}

func NewGetBlobIDRequest(r *http.Request) (GetBlobIDRequest, error) {
	request := GetBlobIDRequest{}

	id := chi.URLParam(r, "id")
	if _, err := strconv.Atoi(id); err == nil {
		return request, errors.New("id is not a integer")
	}

	request.BlobID = cast.ToInt64(id)

	return request, nil
}
