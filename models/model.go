package model

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hamidOyeyiola/mail-service-api/interfaces"
)

//SQLQueryToInsert xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
type SQLQueryToInsert string

//SQLQueryToSelect xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
type SQLQueryToSelect string

//SQLQueryToDelete xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
type SQLQueryToDelete string

//SQLQueryToUpdate xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
type SQLQueryToUpdate struct {
	Stmts  []string
	Values []string
	ID     string
}

//JSONObject string xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
type JSONObject string

//Model string xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
type Model interface {
	Name() string       //Name() refers to the table in the database
	PrimaryKey() string //Primarykey() refers to the name of primaryKey field
	Validate(interfaces.Iterator, io.Reader) (string, bool, SQLQueryToDelete)
	FromQueryResult(interfaces.Iterator) (JSONObject, int)
	FromQueryResultArray([]interfaces.Iterator) (JSONObject, int)
	InsertInto(jsonObject io.Reader) ([]SQLQueryToInsert, []SQLQueryToSelect, bool)
	InsertIntoWhere(jsonObject io.Reader, value string) (SQLQueryToInsert, SQLQueryToSelect, bool)
	InsertIntoIf(jsonObject io.Reader, value string) (SQLQueryToInsert, SQLQueryToSelect, bool)
	UpdateWhere(jsonObject io.Reader, value string) (SQLQueryToUpdate, SQLQueryToSelect, bool)
	UpdateIf(jsonObject io.Reader, value string) (SQLQueryToUpdate, SQLQueryToSelect, bool)
	Update(jsonObject io.Reader) ([]SQLQueryToUpdate, []SQLQueryToSelect, bool)
	SelectFromWhere(value string) (SQLQueryToSelect, bool)
	DeleteFromWhere(value string) (SQLQueryToDelete, bool)
}

//GetQueryToSelect string xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
func GetQueryToSelect(m Model, id string) (q SQLQueryToSelect, ok bool) {
	q = SQLQueryToSelect(fmt.Sprintf("SELECT * FROM %s WHERE id IN (%s)", m.Name(), id))
	return q, true
}

//GetQueryToSelectAll string xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
func GetQueryToSelectAll(m Model) (q SQLQueryToSelect, ok bool) {
	q = SQLQueryToSelect(fmt.Sprintf("SELECT * FROM %s ", m.Name()))
	return q, true
}

//GetQueryToDelete string xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
func GetQueryToDelete(m Model, id string) (q SQLQueryToDelete, ok bool) {
	q = SQLQueryToDelete(fmt.Sprintf("DELETE FROM %s WHERE id IN (%s)", m.Name(), id))
	return q, ok
}

//GetParamFromRequest string xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
func GetParamFromRequest(req *http.Request, param string) (string, bool) {
	vars := mux.Vars(req)
	id, ok := vars[param]
	return id, ok
}
