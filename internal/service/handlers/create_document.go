package handlers

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
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

	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	//Create minio client
	minioClient, err := helpers.NewMinioClient()
	if err != nil {
		helpers.Log(r).WithError(err).Info("cannot create client")
		ape.RenderErr(w, problems.InternalError())
	}

	bucketName := "documents-kyc0123" // TODO cfg name
	location := "eu-central-1"        // TODO cfg location

	//Check of bucket
	err = minioClient.MakeBucket(r.Context(), bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket
		exists, errBucketExists := minioClient.BucketExists(r.Context(), bucketName)
		if errBucketExists != nil || !exists {
			helpers.Log(r).WithError(err).Info("cannot connect to the bucket")
			ape.Render(w, problems.InternalError())
			return
		}
	}
	helpers.Log(r).Debug("connected to the bucket")

	//document data
	fileName := uuid.New().String()
	contentType := h.Header.Get("Content-Type")
	objectName := fileName + "." + strings.Split(contentType, "/")[1]

	// Upload the document file with PutObject
	_, err = minioClient.PutObject(r.Context(), bucketName, objectName, file, h.Size, minio.PutObjectOptions{ContentType: contentType})

	if err != nil {
		helpers.Log(r).WithError(err).Info("cannot connect to the bucket")
		ape.RenderErr(w, problems.InternalError())
	}

	//Get url
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename=\""+objectName+"\"")

	// Generates a presigned url which expires in a day.
	presignedURL, err := minioClient.PresignedGetObject(r.Context(), bucketName, objectName, time.Second*24*60*60, reqParams)

	if err != nil {
		helpers.Log(r).WithError(err).Info("cannot connect to the bucket")
		ape.RenderErr(w, problems.InternalError())
	}
	docImage := data.Document{
		Type:         string(req.Type),
		ImageUrl:     presignedURL.String(),
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
