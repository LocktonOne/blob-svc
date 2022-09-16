package handlers

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"

	"gitlab.com/tokene/blob-svc/internal/service/helpers"
	"gitlab.com/tokene/blob-svc/internal/service/requests"
)

func CreatDocument(w http.ResponseWriter, r *http.Request) {
	getReq, err := requests.NewGetDocumentID(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("invalid request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	docDb := helpers.DocumentsQ(r)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to delete blob from DB")
		ape.Render(w, problems.InternalError())
		return
	}
	document, err := docDb.FilterByID(getReq.DocumentID).Get()
	awsCfg := helpers.AwsConfig(r)

	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to delete blob from DB")
		ape.Render(w, problems.InternalError())
		return
	}

	//Create session
	sess := helpers.NewAwsSession(r)

	ape.Render(w, http.StatusOK)
	svc := s3.New(sess)
	getObjectReq, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(awsCfg.Bucket),
		Key:    aws.String(document.Name),
	})
	document.ImageUrl, err = getObjectReq.Presign(awsCfg.Expiration)
	if err != nil {
		helpers.Log(r).WithError(err).Info("cannot get object's url")
		ape.Render(w, problems.InternalError())
		return
	}
}
