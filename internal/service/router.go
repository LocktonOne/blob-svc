package service

import (
	"gitlab.com/tokene/blob-svc/internal/config"
	"gitlab.com/tokene/blob-svc/internal/data/pg"
	"gitlab.com/tokene/blob-svc/internal/service/handlers"
	"gitlab.com/tokene/blob-svc/internal/service/helpers"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
)

func (s *service) router(cfg config.Config) chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			helpers.CtxLog(s.log),
			helpers.CtxBlobsQ(pg.NewBlobsQ(cfg.DB())),
			helpers.CtxDocumentsQ(pg.NewDocumentsQ(cfg.DB())),
			helpers.CtxAwsConfig(cfg.AWSConfig()),
			helpers.CtxDoormanConnector(cfg.DoormanConnector()),
		),
	)
	r.Route("/integrations/storage", func(r chi.Router) {
		r.Route("/blobs", func(r chi.Router) {
			r.Post("/", handlers.CreateBlob)
			r.Get("/", handlers.GetBlobsByOwnerAddress)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", handlers.GetBlobByID)
				r.Delete("/", handlers.DeleteBlob)
			})
		})
		r.Route("/documents", func(r chi.Router) {
			r.Post("/", handlers.CreateDocument)
			r.Get("/", handlers.GetDocumentsByOwnerAddress)
			r.Route("/{id}", func(r chi.Router) {
				r.Delete("/", handlers.DeleteDocument)
				r.Get("/", handlers.GetDocumentByID)
			})
		})
	})
	return r
}
