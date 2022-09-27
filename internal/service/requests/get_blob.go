package requests

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/spf13/cast"
)

type GetBlobByIDRequest struct {
	ID int64
}

func NewGetBlobIDRequest(r *http.Request) (GetBlobByIDRequest, error) {
	request := GetBlobByIDRequest{}

	id := chi.URLParam(r, "id")
	if _, err := strconv.Atoi(id); err != nil {
		return request, errors.New("id is not an integer")
	}

	request.ID = cast.ToInt64(id)
	return request, nil
}
