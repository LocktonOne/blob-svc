package handlers

import (
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"

	"gitlab.com/tokene/blob-svc/internal/data"
	"gitlab.com/tokene/blob-svc/internal/service/helpers"
	"gitlab.com/tokene/blob-svc/internal/service/requests"
	"gitlab.com/tokene/blob-svc/resources"
)

func newDocumentModel(document data.Document, url string) resources.Document {
	result := resources.Document{
		Key:           resources.NewKeyInt64(document.ID, resources.ResourceType(document.Type)),
		Attributes:    resources.DocumentAttributes{Purpose: document.Purpose, Url: url},
		Relationships: resources.DocumentRelationships{Owner: resources.Relation{Data: &resources.Key{ID: document.OwnerAddress}}},
	}
	return result
}

func CreateDocument(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewСreateDocumentRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	err = helpers.Authorization(r, req.Relationships.Owner.Data.ID)
	if err != nil {
		helpers.Log(r).WithError(err).Debug("user does not have permission")
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	file, h, err := r.FormFile("Image")
	awsCfg := helpers.AwsConfig(r)

	//Create session
	sess := helpers.NewAwsSession(r)

	//Create updloader
	uploader := s3manager.NewUploader(sess)

	//Check document extenstion
	contentType := h.Header.Get("Content-Type")
	fileExtension := strings.Split(contentType, "/")[1]
	if err := helpers.CheckFileExtension(fileExtension); err != nil {
		helpers.Log(r).WithError(err).Debug("invalid file type")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	//Generate key
	fileName := uuid.New().String() + "." + fileExtension

	// Upload the file to S3.
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(awsCfg.Bucket),
		Key:    aws.String(fileName),
		Body:   file,
	})
	if err != nil {
		helpers.Log(r).WithError(err).Info("cannot upload file")
		ape.Render(w, problems.InternalError())
		return
	}

	//Get url
	url, err := helpers.GetItemURL(sess, fileName, *awsCfg)
	if err != nil {
		helpers.Log(r).WithError(err).Info("cannot get object's url")
		ape.Render(w, problems.InternalError())
		return
	}

	//Insert into db
	document := data.Document{
		Type:         string(req.Type),
		OwnerAddress: req.Relationships.Owner.Data.ID,
		Purpose:      req.Attributes.Purpose,
		Name:         fileName,
	}

	document.ID, err = helpers.DocumentsQ(r).Insert(document)
	if err != nil {
		helpers.Log(r).WithError(err).Info("cannot insert document to db")
		ape.Render(w, problems.InternalError())
		return
	}

	resp := resources.DocumentResponse{
		Data: newDocumentModel(document, url),
	}
	ape.Render(w, resp)
}
