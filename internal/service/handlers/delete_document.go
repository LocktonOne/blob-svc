package handlers

import (
	"net/http"

	"gitlab.com/tokene/blob-svc/internal/service/helpers"
	"gitlab.com/tokene/blob-svc/internal/service/requests"
	"gitlab.com/tokene/doorman/connector"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func DeleteDocument(w http.ResponseWriter, r *http.Request) {

	delReq, err := requests.NewGetDocumentID(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("invalid request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	doorman := connector.NewConnectorMockKyc(helpers.DoormanConfig(r).ServiceUrl) //remove mockremove

	token, err := doorman.GetAuthToken(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("invalid auth token")
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	document, err := helpers.DocumentsQ(r).FilterByID(delReq.DocumentID).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get document from DB")
		ape.Render(w, problems.InternalError())
		return
	}
	if document == nil {
		ape.Render(w, problems.NotFound())
		return
	}
	//TODO add check permission

	validation, err := doorman.ValidateJwt(token, document.OwnerAddress)
	if err != nil || !validation {
		helpers.Log(r).WithError(err).Info("invalid auth token")
		ape.RenderErr(w, problems.Unauthorized())
		return
	}
	helpers.DeleteItem(helpers.NewAwsSession(r), &helpers.AwsConfig(r).Bucket, &document.Name)
	if err != nil {
		helpers.Log(r).WithError(err).Info("failed to delete document from s3 bucket")
		ape.Render(w, problems.InternalError())
		return
	}
	err = helpers.DocumentsQ(r).DelById(delReq.DocumentID)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to delete document from DB")
		ape.Render(w, problems.InternalError())
		return
	}
	ape.Render(w, http.StatusOK)
}
