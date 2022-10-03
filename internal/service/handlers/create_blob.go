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

func newBlobModel(blob data.Blob) resources.Blob {
	result := resources.Blob{
		Key: resources.NewKeyInt64(blob.ID, resources.BLOB),
		Attributes: resources.BlobAttributes{
			Blob:    []byte(blob.BlobContent),
			Purpose: blob.Purpose,
		},
		Relationships: resources.BlobRelationships{Owner: resources.Relation{Data: &resources.Key{ID: blob.OwnerAddress}}},
	}

	return result
}

func CreateBlob(w http.ResponseWriter, r *http.Request) {

	req, err := requests.NewCreateBlobRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("invalid request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	ownerAddress := req.Relationships.Owner.Data.ID

	err = helpers.Authorization(r, ownerAddress)
	if err != nil {
		helpers.Log(r).WithError(err).Info("user does not have permission")
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	blob := data.Blob{
		OwnerAddress: ownerAddress,
		BlobContent:  string([]byte(req.Attributes.Blob)),
		Purpose:      req.Attributes.Purpose,
	}

	blob.ID, err = helpers.BlobsQ(r).Insert(blob)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to create blob in DB")
		ape.Render(w, problems.InternalError())
		return
	}

	result := resources.BlobResponse{
		Data: newBlobModel(blob),
	}
	ape.Render(w, result)
}
