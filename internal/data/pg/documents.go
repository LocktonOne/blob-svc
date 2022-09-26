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

func NewDocumentsQ(db *pgdb.DB) data.DocumentsQ {
	return &DocumentsQ{
		db:  db.Clone(),
		sql: sq.Select("b.*").From(fmt.Sprintf("%s as b", documentsTableName)),
	}
}

type DocumentsQ struct {
	db  *pgdb.DB
	sql sq.SelectBuilder
}

func (q *DocumentsQ) New() data.DocumentsQ {
	return NewDocumentsQ(q.db)
}

func (q *DocumentsQ) Get() (*data.Document, error) {
	var result data.Document
	err := q.db.Get(&result, q.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (q *DocumentsQ) Select() ([]data.Document, error) {
	var result []data.Document
	err := q.db.Select(&result, q.sql)
	return result, err
}

func (q *DocumentsQ) Insert(value data.Document) (int64, error) {
	clauses := structs.Map(value)

	var result data.Document
	stmt := sq.Insert(documentsTableName).SetMap(clauses).Suffix("returning *")
	err := q.db.Get(&result, stmt)

	return result.ID, err
}
func (q *DocumentsQ) Page(pageParams pgdb.OffsetPageParams) data.DocumentsQ {
	q.sql = pageParams.ApplyTo(q.sql, "id")
	return q
}

func (q *DocumentsQ) FilterByAddress(ids ...string) data.DocumentsQ {
	q.sql = q.sql.Where(sq.Eq{"b.owner_address": ids})
	return q
}
func (q *DocumentsQ) DelById(ids ...int64) error {
	s := sq.Delete(documentsTableName).Where(sq.Eq{"id": ids})
	err := q.db.Exec(s)
	return err
}
func (q *DocumentsQ) FilterByID(ids ...int64) data.DocumentsQ {
	q.sql = q.sql.Where(sq.Eq{"b.id": ids})
	return q
}
