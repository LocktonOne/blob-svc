package pg

import (
	"database/sql"
	"fmt"

	"gitlab.com/tokene/blob-svc/internal/data"

	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const blobsTableName = "blobs"

func NewBlobsQ(db *pgdb.DB) data.BlobsQ {
	return &BlobsQ{
		db:  db.Clone(),
		sql: sq.Select("b.*").From(fmt.Sprintf("%s as b", blobsTableName)),
	}
}

type BlobsQ struct {
	db  *pgdb.DB
	sql sq.SelectBuilder
}

func (q *BlobsQ) New() data.BlobsQ {
	return NewBlobsQ(q.db)
}

func (q *BlobsQ) Get() (*data.Blob, error) {
	var result data.Blob
	err := q.db.Get(&result, q.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (q *BlobsQ) Select() ([]data.Blob, error) {
	var result []data.Blob
	err := q.db.Select(&result, q.sql)
	return result, err
}

func (q *BlobsQ) Insert(value data.Blob) (data.Blob, error) {
	clauses := structs.Map(value)

	var result data.Blob
	stmt := sq.Insert(blobsTableName).SetMap(clauses).Suffix("returning *")
	err := q.db.Get(&result, stmt)

	return result, err
}
func (q *BlobsQ) Page(pageParams pgdb.OffsetPageParams) data.BlobsQ {
	q.sql = pageParams.ApplyTo(q.sql, "id")
	return q
}

func (q *BlobsQ) FilterByAddress(ids ...string) data.BlobsQ {
	q.sql = q.sql.Where(sq.Eq{"b.owner_address": ids})
	return q
}
func (q *BlobsQ) DelById(ids ...int64) error {
	s := sq.Delete(blobsTableName).Where(sq.Eq{"id": ids})
	err := q.db.Exec(s)
	return err
}
func (q *BlobsQ) FilterByID(ids ...int64) data.BlobsQ {
	q.sql = q.sql.Where(sq.Eq{"b.id": ids})
	return q
}
