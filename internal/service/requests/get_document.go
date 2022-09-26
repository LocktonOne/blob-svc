package requests

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/spf13/cast"
)

type GetDocumentRequest struct {
	DocumentID int64 `url:"-"`
}

func NewGetDocumentID(r *http.Request) (GetDocumentRequest, error) {
	request := GetDocumentRequest{}

	id := r.URL.Query().Get("id")
	if _, err := strconv.Atoi(id); err != nil {
		return request, errors.New("id is not an integer")
	}

	request.DocumentID = cast.ToInt64(id)

	return request, nil
}
