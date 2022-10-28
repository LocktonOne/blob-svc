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

func newBlobsList(blobs []data.Blob) []resources.Blob {
	result := make([]resources.Blob, len(blobs))
	for i, blob := range blobs {
		result[i] = newBlobModel(blob)
	}
	return result
}

func GetBlobsByOwnerAddress(w http.ResponseWriter, r *http.Request) {

	req, err := requests.NewGetBLobsListRequest(r)
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

	blobs, err := helpers.BlobsQ(r).FilterByAddress(req.Owner).Select()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get blobs from DB")
		ape.Render(w, problems.InternalError())
		return
	}

	result := resources.BlobListResponse{
		Data: newBlobsList(blobs),
	}

	ape.Render(w, result)
}
