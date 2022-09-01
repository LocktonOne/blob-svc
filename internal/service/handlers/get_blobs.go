package handlers

import (
	"net/http"

	"gitlab.com/tokene/blob-svc/internal/data"

	"gitlab.com/tokene/blob-svc/internal/service/requests"
	"gitlab.com/tokene/blob-svc/resources"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func newBlobsList(blobs []data.Blob) []resources.Blob {
	result := make([]resources.Blob, len(blobs))
	for i, blob := range blobs {
		result[i] = newBlobModel(blob)
	}
	return result
}
func GetBlobs(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewGetBLobsListRequest(r)
	if err != nil {
		Log(r).WithError(err).Info("invalid request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	q := BlobsQ(r)
	ApplyFilter(q, req)
	blobs, err := q.Select()
	if err != nil {
		Log(r).WithError(err).Error("failed to get blobs")
		ape.Render(w, problems.InternalError())
		return
	}
	result := resources.BlobListResponse{
		Data: newBlobsList(blobs),
	}
	ape.Render(w, result)
}
func ApplyFilter(q data.BlobsQ, request requests.GetBLobsListRequest) {
	q.Page(request.OffsetPageParams)

	if len(request.FilterOwnerID) > 0 {
		q.FilterByOwnerID(request.FilterOwnerID...)
	}

}
