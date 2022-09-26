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

func NewGetDocumentID(r *http.Request) (int64, error) {

	id := r.URL.Query().Get("id")
	if _, err := strconv.Atoi(id); err != nil {
		return 0, errors.New("id is not an integer")
	}

	return cast.ToInt64(id), nil
}
