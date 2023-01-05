package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokene/blob-svc/internal/data"
	"gitlab.com/tokene/blob-svc/internal/service/helpers"
	"gitlab.com/tokene/blob-svc/internal/service/requests"
	"gitlab.com/tokene/blob-svc/resources"
)

func GetDocumentsByOwnerAddress(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewGetDocumentsListRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to parse request")
		ape.Render(w, problems.InternalError())
		return
	}

	err = helpers.Authorization(r, req.Owner)
	if err != nil {
		helpers.Log(r).WithError(err).Info("invalid auth token")
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	documents, err := helpers.DocumentsQ(r).FilterByAddress(req.Owner).Select()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get blobs from DB")
		ape.Render(w, problems.InternalError())
		return
	}

	session := helpers.NewAwsSession(r)
	urls := make([]string, len(documents))

	for i, document := range documents {
		urls[i], err = helpers.GetItemURL(session, document.Name, *helpers.AwsConfig(r))
		if err != nil {
			helpers.Log(r).WithError(err).Error("failed to get url for document")
			ape.Render(w, problems.InternalError())
			return
		}
	}

	result := resources.DocumentListResponse{
		Data: newDocumentsList(documents, urls),
	}

	ape.Render(w, result)
}

func newDocumentsList(documents []data.Document, urls []string) []resources.Document {
	result := make([]resources.Document, len(documents))
	for i, document := range documents {
		result[i] = newDocumentModel(document, urls[i])
	}
	return result
}
