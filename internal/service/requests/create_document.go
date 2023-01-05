package requests

import (
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
)

type CreateDocumentRequest struct {
	MimeType string `json:"mime_type"`
	Name     string `json:"name"`
}

func NewCreateDocumentRequest(r *http.Request) (CreateDocumentRequest, error) {
	request := struct {
		Data CreateDocumentRequest `json:"data"`
	}{}

	err := r.ParseMultipartForm(1 << 32)
	if err != nil {
		return CreateDocumentRequest{}, errors.Wrap(err, "multipart limit to unmarshal")
	}

	if err := json.Unmarshal([]byte(r.FormValue("Metadata")), &request); err != nil {
		return CreateDocumentRequest{}, errors.Wrap(err, "multipart limit to unmarshal")
	}

	return request.Data, nil
}
