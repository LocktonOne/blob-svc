package requests

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/spf13/cast"
)

func NewGetBlobIDRequest(r *http.Request) (int64, error) {

	id := chi.URLParam(r, "id")
	if _, err := strconv.Atoi(id); err != nil {
		return 0, errors.New("id is not an integer")
	}

	return cast.ToInt64(id), nil
}
