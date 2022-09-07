package handlers

import (
	"bytes"
	"io"
	"net/http"
	"os"

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
		Attributes:    resources.DocumentAttributes{Purpose: document.Purpose},
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
	docsDB := helpers.IamgesQ(r)

	file, h, err := r.FormFile("Image")
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	//Convert file to bytes
	buffer := bytes.NewBuffer(nil)
	_, err = io.Copy(buffer, file)
	if err != nil {
		helpers.Log(r).WithError(err).Info("cannot parse file")
		ape.Render(w, problems.InternalError())
		return
	}
	//Create image
	docImage := data.Document{
		Type:         string(req.Type),
		Image:        buffer.Bytes(),
		OwnerAddress: req.Relationships.Owner.Data.ID,
		Format:       h.Header.Get("Content-Type"),
		DocumentName: h.Filename,
		Purpose:      req.Attributes.Purpose,
	}

	doc, err := docsDB.Insert(docImage)
	if err != nil {
		helpers.Log(r).WithError(err).Info("cannot insert document to db")
		ape.Render(w, problems.InternalError())
		return
	}
	tmpfile, err := os.Create("./" + h.Filename)
	defer tmpfile.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tmpfile.Write(doc.Image)

	result := resources.DocumentResponse{
		Data: newDocumentModel(doc),
	}
	ape.Render(w, result)
}
