package requests

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/spf13/cast"
)

type DeleteBlobRequest struct {
	BlobID int64 `url:"-"`
}

func NewDeleteBlobRequest(r *http.Request) (DeleteBlobRequest, error) {
	request := DeleteBlobRequest{}
	id := chi.URLParam(r, "id")
	if _, err := strconv.Atoi(id); err != nil {
		return request, errors.New("id is not an integer")
	}

	request.BlobID = cast.ToInt64(id)

	return request, nil
}
