package handlers

import (
	"net/http"

	"gitlab.com/tokene/blob-svc/internal/service/helpers"
	"gitlab.com/tokene/blob-svc/internal/service/requests"
	"gitlab.com/tokene/blob-svc/resources"
	"gitlab.com/tokene/doorman/connector"

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
	doorman := connector.NewConnectorMockKyc(helpers.DoormanConfig(r).ServiceUrl)

	token, err := doorman.GetAuthToken(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("invalid auth token")
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	blob, err := helpers.BlobsQ(r).FilterByID(req.BlobID).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get blob from DB")
		ape.Render(w, problems.InternalError())
		return
	}
	if blob == nil {
		ape.Render(w, problems.NotFound())
		return
	}
	//TODO add check permission

	validation, err := doorman.ValidateJwt(token, blob.OwnerAddress)
	if err != nil || !validation {
		helpers.Log(r).WithError(err).Info("invalid auth token")
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	result := resources.BlobResponse{
		Data: newBlobModel(*blob),
	}
	ape.Render(w, result)
}
