package handlers

import (
	"fmt"
	"net/http"

	"gitlab.com/tokene/blob-svc/internal/data"
	"gitlab.com/tokene/blob-svc/internal/service/requests"
	"gitlab.com/tokene/blob-svc/resources"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func newBlobModel(blob data.Blob) resources.Blob {
	owner := resources.Relation{
		Data: &resources.Key{
			ID:   blob.OwnerID,
			Type: "user",
		},
	}
	result := resources.Blob{
		Key: resources.NewKeyInt64(blob.ID, resources.BLOB),
		Attributes: resources.BlobAttributes{
			Blob:    []byte(blob.BlobContent),
			Purpose: blob.Purpose,
		},
		Relationships: resources.BlobRelationships{
			Owner: owner,
		},
	}

	return result
}

func CreateBlob(w http.ResponseWriter, r *http.Request) {

	req, err := requests.NewCreateBlobRequest(r)
	if err != nil {
		Log(r).WithError(err).Info("invalid request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	onwer_id := fmt.Sprintf("%v", req.Data.Relationships.Owner.Data.ID)

	c_blob := data.Blob{
		OwnerID:     onwer_id,
		BlobContent: string([]byte(req.Data.Attributes.Blob)),
		Purpose:     req.Data.Attributes.Purpose,
	}
	blob, err := BlobsQ(r).Insert(c_blob)
	if err != nil {
		Log(r).WithError(err).Error("failed to create blob in DB")
		ape.Render(w, problems.InternalError())
		return
	}
	result := resources.BlobResponse{
		Data: newBlobModel(blob),
	}
	ape.Render(w, result)
}
