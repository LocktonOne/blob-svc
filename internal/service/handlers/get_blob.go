package handlers

import (
	"net/http"

	"gitlab.com/tokene/blob-svc/internal/service/helpers"
	"gitlab.com/tokene/blob-svc/internal/service/requests"
	"gitlab.com/tokene/blob-svc/resources"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetBlobByID(w http.ResponseWriter, r *http.Request) {

	req, err := requests.NewGetBlobIDRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("invalid request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	blob, err := helpers.BlobsQ(r).FilterByID(req.ID).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get blob from DB")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if blob == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	err = helpers.Authorization(r, blob.OwnerAddress)
	if err != nil {
		helpers.Log(r).WithError(err).Info("invalid auth token")
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	result := resources.BlobResponse{
		Data: newBlobModel(*blob),
	}
	ape.Render(w, result)
}
