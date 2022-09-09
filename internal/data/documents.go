package data

import (
	"gitlab.com/distributed_lab/kit/pgdb"
)

type ImagesQ interface {
	New() ImagesQ

	Get() (*Document, error)
	Select() ([]Document, error)

	Insert(data Document) (Document, error)
	Page(pageParams pgdb.OffsetPageParams) ImagesQ
	DelById(id ...int64) error
	FilterByAddress(id ...string) ImagesQ
	FilterByID(id ...int64) ImagesQ
}

type Document struct {
	ID           int64  `db:"id" structs:"-"`
	Type         string `db:"type" structs:"type"`
	OwnerAddress string `db:"owner_address" structs:"owner_address"`
	Name         string `db:"name" structs:"name"`
	ImageUrl     string `db:"image_url" structs:"image_url"`
	Purpose      string `db:"purpose" structs:"purpose"`
}
