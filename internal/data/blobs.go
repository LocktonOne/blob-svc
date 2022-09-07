package data

import "gitlab.com/distributed_lab/kit/pgdb"

type BlobsQ interface {
	New() BlobsQ

	Get() (*Blob, error)
	Select() ([]Blob, error)

	Insert(data Blob) (Blob, error)
	Page(pageParams pgdb.OffsetPageParams) BlobsQ
	DelById(id ...int64) error
	FilterByAddress(id ...string) BlobsQ
	FilterByID(id ...int64) BlobsQ
}

type Blob struct {
	ID           int64  `db:"id" structs:"-"`
	OwnerAddress string `db:"owner_address" structs:"owner_address"`
	BlobContent  string `db:"blob_content" structs:"blob_content"`
	Purpose      string `db:"purpose" structs:"purpose"`
}
