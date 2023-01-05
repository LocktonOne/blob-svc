package data

import (
	"gitlab.com/distributed_lab/kit/pgdb"
)

type DocumentsQ interface {
	New() DocumentsQ

	Get() (*Document, error)
	Select() ([]Document, error)

	Insert(data Document) (int64, error)
	Page(pageParams pgdb.OffsetPageParams) DocumentsQ
	DelById(id ...int64) error
	FilterByAddress(id ...string) DocumentsQ
	FilterByID(id ...int64) DocumentsQ
}

type Document struct {
	ID           int64  `db:"id" structs:"-"`
	Name         string `db:"name" structs:"name"`
	OwnerAddress string `db:"owner_address" structs:"owner_address"`
	FileKey      string `db:"file_key" structs:"file_key"`
	MimeType     string `db:"mime_type" structs:"mime_type"`
}
