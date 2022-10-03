package handlers

import (
	"net/http"

	"gitlab.com/tokene/blob-svc/internal/service/helpers"
	"gitlab.com/tokene/blob-svc/internal/service/requests"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func DeleteBlob(w http.ResponseWriter, r *http.Request) {

	req, err := requests.NewGetBlobIDRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("invalid request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	blob, err := helpers.BlobsQ(r).FilterByID(req.ID).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get blob from DB")
		ape.Render(w, problems.InternalError())
		return
	}
	if blob == nil {
		ape.Render(w, problems.NotFound())
		return
	}

	err = helpers.Authorization(r, blob.OwnerAddress)
	if err != nil {
		helpers.Log(r).WithError(err).Info("user does not have permission")
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	err = helpers.BlobsQ(r).DelById(req.ID)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to delete blob from DB")
		ape.Render(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
