package model

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/hamidOyeyiola/mail-service-api/interfaces"
	"github.com/hamidOyeyiola/mail-service-api/utils"
)

//Session xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
type Session struct {
	Sessionkey string `json:"sessionkey"`
	UserKey    string `json:"userkey"`
	Deadline   string `json:"-"`
}

//utils.NewDate()
func (se Session) String() string {
	//"INSERT INTO cities(name, population) VALUES ('Moscow', 12506000)"
	return fmt.Sprintf("User(sessionkey, userkey, deadline) VALUES ('%s','%s','%s')",
		se.Sessionkey, se.UserKey, se.Deadline)
}

//Name xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
func (se Session) Name() string {
	return "session"
}

//PrimaryKey   xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
func (se Session) PrimaryKey() string {
	return "sessionkey"
}

//Validate   xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
func (se Session) Validate(s interfaces.Iterator, r io.Reader) (string, bool, SQLQueryToDelete) {
	//v := struct{ Password string }{}
	//err := json.NewDecoder(r).Decode(&v)
	var id int
	if s.Next() {
		err := s.Scan(&se.Sessionkey, &se.UserKey, &se.Deadline, &id)
		if err != nil {
			return "", false, ""
		}
	}
	//fmt.Printf("%s", u)
	d, _ := se.DeleteFromWhere(se.Sessionkey)
	return se.UserKey, utils.IsNotDeadline(se.Deadline), d
}

//InsertInto xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
func (se Session) InsertInto(r io.Reader) (ins []SQLQueryToInsert, sel []SQLQueryToSelect, ok bool) {
	v := []struct {
		Email string
	}{}
	err := json.NewDecoder(r).Decode(&v)
	for _, s := range v {
		ins = append(ins, SQLQueryToInsert(fmt.Sprintf("INSERT INTO session(sessionkey, userkey, deadline) VALUES ('%s','%s','%s')",
			utils.GetSessionToken(), s.Email, utils.MakeDeadline(time.Minute*5))))
	}
	sel = append(sel, "")
	return ins, sel, err == nil
}

//InsertIntoWhere xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
func (se Session) InsertIntoWhere(r io.Reader, value string) (ins SQLQueryToInsert, sel SQLQueryToSelect, ok bool) {

	ins = SQLQueryToInsert(fmt.Sprintf("INSERT INTO session(sessionkey, userkey, deadline) VALUES ('%s','%s','%s')",
		utils.GetSessionToken(), value, utils.MakeDeadline(time.Minute*5)))
	return ins, sel, true
}

//InsertIntoIf xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
func (se Session) InsertIntoIf(r io.Reader, value string) (ins SQLQueryToInsert, sel SQLQueryToSelect, ok bool) {
	ins = SQLQueryToInsert(fmt.Sprintf("INSERT INTO session(sessionkey, userkey, deadline) VALUES ('%s','%s','%s')",
		utils.GetSessionToken(), value, utils.MakeDeadline(time.Minute*5)))
	return ins, sel, true
}

//Update xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
func (se Session) Update(r io.Reader) (upt []SQLQueryToUpdate, sel []SQLQueryToSelect, ok bool) {
	return upt, sel, true
}

//UpdateWhere xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
func (se Session) UpdateWhere(r io.Reader, value string) (upt SQLQueryToUpdate, sel SQLQueryToSelect, ok bool) {
	return upt, sel, true
}

//UpdateIf xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
func (se Session) UpdateIf(r io.Reader, value string) (upt SQLQueryToUpdate, sel SQLQueryToSelect, ok bool) {

	return upt, sel, true
}

//FromQueryResult xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
func (se Session) FromQueryResult(s interfaces.Iterator) (JSONObject, int) {
	//fmt.Printf("%v", u)

	v := []Session{}
	var n int
	var id int
	for n = 0; s.Next(); n++ {
		err := s.Scan(&se.Sessionkey, &se.UserKey, &se.Deadline, &id)
		if err != nil {
			break
		}
		v = append(v, se)
	}
	result := struct {
		N    int       `json:"rows"`
		Data []Session `json:"data"`
	}{N: n, Data: v}
	o, _ := json.MarshalIndent(result, "", "  ")
	return JSONObject(o), n
}

//FromQueryResultArray xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
func (se Session) FromQueryResultArray(ss []interfaces.Iterator) (JSONObject, int) {
	//fmt.Printf("%v", u)

	v := []Session{}
	var n int
	var id int
	for _, s := range ss {
		for n = 0; s.Next(); n++ {
			err := s.Scan(&se.Sessionkey, &se.Deadline, &id)
			if err != nil {
				break
			}
			v = append(v, se)
		}
	}
	result := struct {
		N    int       `json:"rows"`
		Data []Session `json:"data"`
	}{N: n, Data: v}
	o, _ := json.MarshalIndent(result, "", "  ")
	return JSONObject(o), n
}

//SelectFromWhere xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
func (se Session) SelectFromWhere(value string) (q SQLQueryToSelect, ok bool) {
	q = SQLQueryToSelect("SELECT * FROM " + se.Name())
	q = q + SQLQueryToSelect(fmt.Sprintf(" WHERE %s = '%s'", se.PrimaryKey(), value))
	return q, true
}

//DeleteFromWhere xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
func (se Session) DeleteFromWhere(value string) (q SQLQueryToDelete, ok bool) {
	q = SQLQueryToDelete(fmt.Sprintf("DELETE FROM %s WHERE %s = '%s'", se.Name(), se.PrimaryKey(), value))
	return q, true
}
