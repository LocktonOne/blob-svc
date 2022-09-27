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

	documentID, err := requests.NewGetDocumentID(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("invalid request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	document, err := helpers.DocumentsQ(r).FilterByID(documentID).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get document from DB")
		ape.Render(w, problems.InternalError())
		return
	}

	err = helpers.Authorization(r, document.OwnerAddress)
	if err != nil {
		helpers.Log(r).WithError(err).Info("invalid auth token")
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	if document == nil {
		ape.Render(w, problems.NotFound())
		return
	}

	url, err := helpers.GetItemURL(helpers.NewAwsSession(r), document.Name, *helpers.AwsConfig(r))

	result := resources.DocumentResponse{
		Data: newDocumentModel(*document, url),
	}

	ape.Render(w, result)
}
