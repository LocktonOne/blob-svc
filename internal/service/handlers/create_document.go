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

func CreateDocument(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewCreateDocumentRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to parse request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	owner, err := helpers.ValidateJwt(r)
	if err != nil {
		helpers.Log(r).WithError(err).Debug("user does not have permission")
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	file, _, err := r.FormFile("Document")
	awsCfg := helpers.AwsConfig(r)

	//Create session
	sess := helpers.NewAwsSession(r)

	//Create uploader
	uploader := s3manager.NewUploader(sess)

	//Check document extension

	fileExtension := strings.Split(req.MimeType, "/")[1]
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
		ape.RenderErr(w, problems.InternalError())
		return
	}

	//Get url
	url, err := helpers.GetItemURL(sess, fileName, *helpers.AwsConfig(r))
	if err != nil {
		helpers.Log(r).WithError(err).Info("cannot get object's url")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	//Insert into db
	document := data.Document{
		Name:         req.Name,
		OwnerAddress: owner,
		MimeType:     req.MimeType,
		FileKey:      fileName,
	}

	document.ID, err = helpers.DocumentsQ(r).Insert(document)
	if err != nil {
		helpers.Log(r).WithError(err).Info("cannot insert document to db")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	resp := resources.DocumentResponse{
		Data: newDocumentModel(document, url),
	}
	ape.Render(w, resp)
}

func newDocumentModel(document data.Document, url string) resources.Document {
	result := resources.Document{
		Key: resources.NewKeyInt64(document.ID, resources.DOCUMENT),
		Attributes: resources.DocumentAttributes{
			MimeType: document.MimeType,
			Name:     document.Name,
			Link:     url,
			Owner:    document.OwnerAddress,
		},
	}
	return result
}
