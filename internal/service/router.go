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
		),
	)
	r.Route("/integrations/blobs", func(r chi.Router) {
		r.Get("/", handlers.GetBlobs)
		r.Post("/", handlers.CreateBlob)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handlers.GetBlobByID)
			r.Delete("/", handlers.DeleteBlob)
		})
	})

	return r
}
