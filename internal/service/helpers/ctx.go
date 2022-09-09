package helpers

import (
	"context"
	"net/http"

	"gitlab.com/tokene/blob-svc/internal/config"
	"gitlab.com/tokene/blob-svc/internal/data"

	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	blobsQCtxKey
	imagesQCtxKey
	awsCfgKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}
func CtxBlobsQ(entry data.BlobsQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, blobsQCtxKey, entry)
	}
}
func CtxImagesQ(entry data.ImagesQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, imagesQCtxKey, entry)
	}
}
func CtxAwsConfig(entry *config.AWSConfig) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, awsCfgKey, entry)
	}
}
func AwsConfig(r *http.Request) *config.AWSConfig {
	return r.Context().Value(awsCfgKey).(*config.AWSConfig)
}
func DocumentsQ(r *http.Request) data.ImagesQ {
	return r.Context().Value(imagesQCtxKey).(data.ImagesQ).New()
}
func BlobsQ(r *http.Request) data.BlobsQ {
	return r.Context().Value(blobsQCtxKey).(data.BlobsQ).New()
}
func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}
