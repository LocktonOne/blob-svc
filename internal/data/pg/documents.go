package pg

import (
	"database/sql"
	"fmt"

	"gitlab.com/tokene/blob-svc/internal/data"

	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const documentsTableName = "documents"

func NewImagesQ(db *pgdb.DB) data.ImagesQ {
	return &ImagesQ{
		db:  db.Clone(),
		sql: sq.Select("b.*").From(fmt.Sprintf("%s as b", documentsTableName)),
	}
}

type ImagesQ struct {
	db  *pgdb.DB
	sql sq.SelectBuilder
}

func (q *ImagesQ) New() data.ImagesQ {
	return NewImagesQ(q.db)
}

func (q *ImagesQ) Get() (*data.Document, error) {
	var result data.Document
	err := q.db.Get(&result, q.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (q *ImagesQ) Select() ([]data.Document, error) {
	var result []data.Document
	err := q.db.Select(&result, q.sql)
	return result, err
}

func (q *ImagesQ) Insert(value data.Document) (data.Document, error) {
	clauses := structs.Map(value)

	var result data.Document
	stmt := sq.Insert(documentsTableName).SetMap(clauses).Suffix("returning *")
	err := q.db.Get(&result, stmt)

	return result, err
}
func (q *ImagesQ) Page(pageParams pgdb.OffsetPageParams) data.ImagesQ {
	q.sql = pageParams.ApplyTo(q.sql, "id")
	return q
}

func (q *ImagesQ) FilterByAddress(ids ...string) data.ImagesQ {
	q.sql = q.sql.Where(sq.Eq{"b.owner_address": ids})
	return q
}
func (q *ImagesQ) DelById(ids ...int64) error {
	s := sq.Delete(documentsTableName).Where(sq.Eq{"id": ids})
	err := q.db.Exec(s)
	return err
}
func (q *ImagesQ) FilterByID(ids ...int64) data.ImagesQ {
	q.sql = q.sql.Where(sq.Eq{"b.id": ids})
	return q
}
