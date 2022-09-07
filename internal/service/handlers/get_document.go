package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokene/blob-svc/internal/service/helpers"
	"gitlab.com/tokene/blob-svc/internal/service/requests"
	"gitlab.com/tokene/blob-svc/resources"
)

func GetDocument(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewGetDocumentID(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("invalid request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	doc, err := helpers.IamgesQ(r).FilterByID(req.DocumentID).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get blob from DB")
		ape.Render(w, problems.InternalError())
		return
	}
	if doc == nil {
		ape.Render(w, problems.NotFound())
		return
	}
	result := resources.DocumentResponse{
		Data: newDocumentModel(*doc),
	}

	writer := multipart.NewWriter(w)

	// Write document info
	metadata, _ := json.Marshal(result)
	metadataHeader := textproto.MIMEHeader{}
	metadataHeader.Set("Content-Type", "application/json; charset=UTF-8")
	part, _ := writer.CreatePart(metadataHeader)
	part.Write(metadata)

	//Write image

	mediaHeader := textproto.MIMEHeader{}
	mediaHeader.Set("Content-Type", doc.Format)
	mediaHeader.Add("Content-Length", fmt.Sprint(len(doc.Image)))
	mediaHeader.Add("Content-Disposition", "form-data; name=\"Image\"; filename=\""+doc.DocumentName+"\"")
	mediaPart, err := writer.CreatePart(mediaHeader)
	io.Copy(mediaPart, bytes.NewReader(doc.Image))
	w.Header().Set("Content-Type", writer.FormDataContentType())

	writer.Close()
}
