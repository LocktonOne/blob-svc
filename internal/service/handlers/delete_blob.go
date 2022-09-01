package handlers

import (
	"net/http"

	"gitlab.com/tokene/blob-svc/internal/service/requests"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func DeleteBlob(w http.ResponseWriter, r *http.Request) {

	del_req, err := requests.NewDeleteBlobRequest(r)
	if err != nil {
		Log(r).WithError(err).Info("invalid request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	err = BlobsQ(r).DelById(del_req.BlobID)
	if err != nil {
		Log(r).WithError(err).Error("failed to delete blob from DB")
		ape.Render(w, problems.InternalError())
		return
	}
	ape.Render(w, 200)
}
