package handlers

import (
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"

	"gitlab.com/tokene/blob-svc/internal/data"
	"gitlab.com/tokene/blob-svc/internal/service/helpers"
	"gitlab.com/tokene/blob-svc/internal/service/requests"
	"gitlab.com/tokene/blob-svc/resources"
)

func newDocumentModel(document data.Document) resources.Document {
	result := resources.Document{
		Key:           resources.NewKeyInt64(document.ID, resources.ResourceType(document.Type)),
		Attributes:    resources.DocumentAttributes{Purpose: document.Purpose, Url: document.ImageUrl},
		Relationships: resources.DocumentRelationships{Owner: resources.Relation{Data: &resources.Key{ID: document.OwnerAddress}}},
	}
	return result
}

func CreatDocument(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewСreateDocumentRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	file, h, err := r.FormFile("Image")
	awsCfg := helpers.AwsConfig(r)

	//Create session
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("eu-central-1"),
		Credentials: credentials.NewStaticCredentials(awsCfg.AccessKeyID, awsCfg.SecretKeyID, ""),
		DisableSSL:  &awsCfg.SslDisable,
	}))

	uploader := s3manager.NewUploader(sess)

	//File data
	contentType := h.Header.Get("Content-Type")
	objectName := uuid.New().String() + "." + strings.Split(contentType, "/")[1]

	// Upload the file to S3.
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(awsCfg.Bucket),
		Key:    aws.String(objectName),
		Body:   file,
	})
	if err != nil {
		helpers.Log(r).WithError(err).Info("cannot upload file")
		ape.Render(w, problems.InternalError())
		return
	}

	//Get url
	svc := s3.New(sess)
	getObjectReq, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(awsCfg.Bucket),
		Key:    aws.String(objectName),
	})
	url, err := getObjectReq.Presign(awsCfg.Expiration)
	if err != nil {
		helpers.Log(r).WithError(err).Info("cannot get object's url")
		ape.Render(w, problems.InternalError())
		return
	}

	//Insert to db
	docImage := data.Document{
		Type:         string(req.Type),
		ImageUrl:     url,
		OwnerAddress: req.Relationships.Owner.Data.ID,
		Purpose:      req.Attributes.Purpose,
		Name:         objectName,
	}

	doc, err := helpers.DocumentsQ(r).Insert(docImage)
	if err != nil {
		helpers.Log(r).WithError(err).Info("cannot insert document to db")
		ape.Render(w, problems.InternalError())
		return
	}

	resp := resources.DocumentResponse{
		Data: newDocumentModel(doc),
	}
	ape.Render(w, resp)
}
