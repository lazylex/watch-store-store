package db_reader

import (
	"context"
	"database/sql"
)

type Reader struct {
	db *sql.DB
}

func New(db *sql.DB) *Reader {
	return &Reader{db: db}
}

func (r *Reader) ReadTables(databaseName string) (map[string][]string, error) {
	var rows *sql.Rows
	var err error
	result := make(map[string][]string)
	ctx := context.Background()

	stmt := `SELECT table_name FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE='BASE TABLE' AND table_schema=?;`

	if rows, err = r.db.QueryContext(ctx, stmt, databaseName); err != nil {
		return nil, err
	}
	for rows.Next() {
		var record []uint8
		if err = rows.Scan(&record); err != nil {
			return nil, err
		}
		result[string(record)] = make([]string, 0)
	}

	for table, _ := range result {
		stmt = `SELECT column_name FROM information_schema.columns WHERE table_schema=? AND table_name=?;`

		if rows, err = r.db.QueryContext(ctx, stmt, databaseName, table); err != nil {
			return nil, err
		}
		for rows.Next() {
			var record []uint8
			if err = rows.Scan(&record); err != nil {
				return nil, err
			}
			result[table] = append(result[table], string(record))
		}
	}

	return result, nil
}
