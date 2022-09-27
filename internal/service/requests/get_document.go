package requests

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/spf13/cast"
)

type GetDocumentByIDRequest struct {
	ID int64
}

func NewGetDocumentID(r *http.Request) (GetDocumentByIDRequest, error) {

	request := GetDocumentByIDRequest{}

	id := r.URL.Query().Get("id")
	if _, err := strconv.Atoi(id); err != nil {
		return request, errors.New("id is not an integer")
	}

	request.ID = cast.ToInt64(id)
	return request, nil
}
