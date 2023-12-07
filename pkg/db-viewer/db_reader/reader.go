package db_reader

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"log/slog"
	"net/http"
)

type Reader struct {
	db     *sql.DB
	log    *slog.Logger
	dbName string
	tables map[string][]string
}

// Start выводит список таблиц БД по адресу localhost<:port>/tables с возможностью вывода их содержимого
func Start(db *sql.DB, log *slog.Logger, port string) {
	var err error
	reader := &Reader{db: db, log: log}
	reader.dbName = reader.readDBName()
	if reader.dbName == "" {
		log.Error("empty db name, so no tables will viewed")
		return
	}
	if reader.tables, err = reader.readTables(); err != nil {
		log.Error("can't read tables info, so no tables will viewed")
		return
	}

	http.HandleFunc("/tables", reader.displayTablesList)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			log.Error("server closed")
		} else if err != nil {
			log.Error("error starting server: %s", err)
		}
	}

}

// readDBName возвращает название БД в СУБД
func (rd *Reader) readDBName() string {
	var name string
	stmt := `SELECT DATABASE();`

	err := rd.db.QueryRow(stmt).Scan(&name)
	if err != nil {
		rd.log.Error("can't read db name")
		return ""
	}

	return name
}

// readTables возвращает названия таблиц и их столбцов в СУБД
func (rd *Reader) readTables() (map[string][]string, error) {
	var rows *sql.Rows
	var err error
	result := make(map[string][]string)
	ctx := context.Background()

	stmt := `SELECT table_name FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE='BASE TABLE' AND table_schema=?;`

	if rows, err = rd.db.QueryContext(ctx, stmt, rd.dbName); err != nil {
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

		if rows, err = rd.db.QueryContext(ctx, stmt, rd.dbName, table); err != nil {
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

// tableNames возвращает список с названиями таблиц
func (rd *Reader) tableNames() []string {
	var result []string
	for name, _ := range rd.tables {
		result = append(result, name)
	}
	return result
}

// displayTablesList отображает страницу со списком таблиц
func (rd *Reader) displayTablesList(w http.ResponseWriter, r *http.Request) {
	linkStyle := makeStyles(map[string]string{"color": "black", "text-decoration": "none"})
	var rows [][]string

	for _, name := range rd.tableNames() {
		rows = append(rows, []string{makeLink(name, makePath([]string{"tables", name}), linkStyle)})
	}

	_, err := io.WriteString(w, makeTable([]string{"Таблицы"}, rows))
	if err != nil {
		return
	}

}
