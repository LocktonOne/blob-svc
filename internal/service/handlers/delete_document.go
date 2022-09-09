package handlers

import (
	"net/http"

	"gitlab.com/tokene/blob-svc/internal/service/helpers"
	"gitlab.com/tokene/blob-svc/internal/service/requests"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func DeleteDocument(w http.ResponseWriter, r *http.Request) {

	del_req, err := requests.NewGetDocumentID(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("invalid request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	err = helpers.BlobsQ(r).DelById(del_req.DocumentID)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to delete blob from DB")
		ape.Render(w, problems.InternalError())
		return
	}
	ape.Render(w, http.StatusOK)
}
