package handlers

import (
	"net/http"

	"gitlab.com/tokene/blob-svc/internal/service/helpers"
	"gitlab.com/tokene/blob-svc/internal/service/requests"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func DeleteDocument(w http.ResponseWriter, r *http.Request) {

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
		helpers.Log(r).WithError(err).Info("user does not have permission")
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	err = helpers.DeleteItem(helpers.NewAwsSession(r), &document.Name, *helpers.AwsConfig(r))
	if err != nil {
		helpers.Log(r).WithError(err).Info("failed to delete document from s3 bucket")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	err = helpers.DocumentsQ(r).DelById(req.ID)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to delete document from DB")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
