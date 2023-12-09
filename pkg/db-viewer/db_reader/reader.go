package db_reader

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"reflect"
	"strings"
	"time"
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
	http.HandleFunc("/table", reader.displayTable)
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

	stmt := `SELECT table_name FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE='BASE TABLE' AND table_schema=?;`

	if rows, err = rd.db.Query(stmt, rd.dbName); err != nil {
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
		stmt = fmt.Sprintf("describe %s;", table)
		if rows, err = rd.db.Query(stmt); err != nil {
			rd.log.Error(err.Error())
			return nil, err
		}

		for rows.Next() {
			var record, trash []uint8
			if err = rows.Scan(&record, &trash, &trash, &trash, &trash, &trash); err != nil {
				return nil, err
			}
			result[table] = append(result[table], string(record))
		}
	}

	return result, nil
}

func (rd *Reader) selectAll(tableName string) [][]string {
	var s string
	var err error
	var rows *sql.Rows
	var result [][]string
	caption := rd.tables[tableName]
	numCols := len(caption)
	log := rd.log.With("op", "db_reader.selectAll")

	// при построении запроса неизвестно количество запрашиваемых столбцов, поэтому запрос строится небезопасным методом
	stmt := "SELECT " + strings.Join(caption, ", ") + " FROM " + tableName + ";"
	if rows, err = rd.db.Query(stmt); err != nil {
		log.Error(err.Error())
		return [][]string{}
	}

	for rows.Next() {
		t := make([]interface{}, numCols)
		for i := range t {
			t[i] = &t[i]
		}

		err = rows.Scan(t...)
		if err != nil {
			log.Error(err.Error())
			return nil
		}
		result = append(result, make([]string, numCols))
		for i := range t {
			if reflect.ValueOf(t[i]).Type().String() == "time.Time" {
				timeVal := t[i].(time.Time)
				s = timeVal.String()
			} else {
				anotherVal := t[i].([]uint8)
				s = string(anotherVal)
			}
			result[len(result)-1][i] = s
		}
	}

	return result
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
		rows = append(rows, []string{makeLink(name, fmt.Sprintf("table?name=%s", name), linkStyle)})
	}

	if _, err := io.WriteString(w, makeTable([]string{"Таблицы"}, rows)); err != nil {
		return
	}
}

// displayTable хендлер для отображения таблицы. В GET параметре name передается название таблицы, которую необходимо
// отобразить
func (rd *Reader) displayTable(w http.ResponseWriter, r *http.Request) {
	tableName := r.URL.Query().Get("name")
	caption := rd.tables[tableName]

	if len(tableName) == 0 || len(caption) == 0 {
		_, _ = io.WriteString(w, "Нет таблицы для отображения")
		return
	}
	data := rd.selectAll(tableName)
	table := makeTable(caption, data)
	if _, err := io.WriteString(w, table); err != nil {
		return
	}
}
