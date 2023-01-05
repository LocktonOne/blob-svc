package handlers

import (
	"net/http"

	"gitlab.com/tokene/blob-svc/internal/service/helpers"
	"gitlab.com/tokene/blob-svc/internal/service/requests"
	"gitlab.com/tokene/blob-svc/resources"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetDocumentByID(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewGetDocumentID(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("invalid request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	document, err := helpers.DocumentsQ(r).FilterByID(req.ID).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get document from DB")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if document == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	err = helpers.Authorization(r, document.OwnerAddress)
	if err != nil {
		helpers.Log(r).WithError(err).Info("invalid auth token")
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	session := helpers.NewAwsSession(r)

	url, err := helpers.GetItemURL(session, document.Name, *helpers.AwsConfig(r))
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get url for document")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	result := resources.DocumentResponse{
		Data: newDocumentModel(*document, url),
	}

	ape.Render(w, result)
}
