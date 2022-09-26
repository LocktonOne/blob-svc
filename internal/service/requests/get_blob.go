package requests

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/spf13/cast"
)

type GetBlobIDRequest struct {
	BlobID int64 `url:"-"`
}

func NewGetBlobIDRequest(r *http.Request) (GetBlobIDRequest, error) {
	request := GetBlobIDRequest{}

	id := r.URL.Query().Get("id")
	if _, err := strconv.Atoi(id); err != nil {
		return request, errors.New("id is not an integer")
	}

	request.BlobID = cast.ToInt64(id)

	return request, nil
}
