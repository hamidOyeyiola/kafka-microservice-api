package model

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/hamidOyeyiola/mail-service-api/interfaces"
	"github.com/hamidOyeyiola/mail-service-api/utils"
)

//User xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
type User struct {
	FirstName string             `json:"firstname"`
	LastName  string             `json:"lastname"`
	Email     utils.EmailAddress `json:"email"`
	PhoneNo   string             `json:"phoneno"`
	Password  string             `json:"password"`
	MessageID string             `json:"-"`
	CreatedOn string             `json:"-"`
	UpdatedOn string             `json:"-"`
	ID        int                `json:"-"`
}

//utils.NewDate()
func (u User) String() string {
	//"INSERT INTO cities(name, population) VALUES ('Moscow', 12506000)"
	return fmt.Sprintf("User(firstname, lastname, email, phoneno, password, messageID, createdOn, updatedOn) VALUES ('%s','%s','%s','%s','%s','%s','%s','%s')",
		u.FirstName, u.LastName, u.Email, u.PhoneNo, u.Password, u.MessageID, u.CreatedOn, u.UpdatedOn)
}

//Name xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
func (u User) Name() string {
	return "users"
}

//PrimaryKey   xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
func (u User) PrimaryKey() string {
	return "email"
}

//Validate   xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
func (u User) Validate(s interfaces.Iterator, r io.Reader) (string, bool, SQLQueryToDelete) {
	v := struct {
		Email    utils.EmailAddress
		Password string
	}{}
	err := json.NewDecoder(r).Decode(&v)
	//fmt.Println(v)
	if s.Next() {
		fmt.Println("ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ")
		err = s.Scan(&u.FirstName, &u.LastName, &u.Email, &u.PhoneNo, &u.Password, &u.MessageID, &u.CreatedOn, &u.UpdatedOn, &u.ID)
		if err != nil {
			return "", false, ""
		}
	}
	return string(v.Email), (u.Password == utils.EncryptPassword(v.Password) && v.Email == u.Email), ""
}

//InsertIntoIf xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
func (u User) InsertIntoIf(r io.Reader, value string) (ins SQLQueryToInsert, sel SQLQueryToSelect, ok bool) {
	err := json.NewDecoder(r).Decode(&u)
	if u.Email.IsValid() && string(u.Email) == value {
		ins = SQLQueryToInsert(fmt.Sprintf("INSERT INTO users(firstname, lastname, email, phoneno, password, messageID, createdOn,updatedOn,id) VALUES ('%s','%s','%s','%s','%s','%s','%s','%s',%d)",
			u.FirstName, u.LastName, u.Email, u.PhoneNo, utils.EncryptPassword(u.Password), u.MessageID, utils.NewDate(), "", 0))
		//sel = SQLQueryToSelect(fmt.Sprintf("SELECT * FROM %s WHERE %s='%s'", u.Name(), u.PrimaryKey(), value))
	}
	return ins, sel, err == nil
}

//InsertInto xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
func (u User) InsertInto(r io.Reader) (ins []SQLQueryToInsert, sel []SQLQueryToSelect, ok bool) {
	v := []User{}

	err := json.NewDecoder(r).Decode(&v)

	for _, u = range v {
		//fmt.Println(u.Password)
		if u.Email.IsValid() {
			ins = append(ins, SQLQueryToInsert(fmt.Sprintf("INSERT INTO users(firstname, lastname, email, phoneno, password, messageID, createdOn,updatedOn,id) VALUES ('%s','%s','%s','%s','%s','%s','%s','%s',%d)",
				u.FirstName, u.LastName, u.Email, u.PhoneNo, utils.EncryptPassword(u.Password), u.MessageID, utils.NewDate(), "", 0)))
			//sel = append(sel, SQLQueryToSelect(fmt.Sprintf("SELECT * FROM %s WHERE %s='%s'", u.Name(), u.PrimaryKey(), u.Email)))
		}
	}
	sel = append(sel, "")
	return ins, sel, err == nil
}

//InsertIntoWhere xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
func (u User) InsertIntoWhere(r io.Reader, value string) (ins SQLQueryToInsert, sel SQLQueryToSelect, ok bool) {
	err := json.NewDecoder(r).Decode(&u)
	if utils.EmailAddress(value).IsValid() {
		ins = SQLQueryToInsert(fmt.Sprintf("INSERT INTO users(firstname, lastname, email, phoneno, password, messageID, createdOn,updatedOn,id) VALUES ('%s','%s','%s','%s','%s','%s','%s','%s',%d)",
			u.FirstName, u.LastName, value, u.PhoneNo, utils.EncryptPassword(u.Password), u.MessageID, utils.NewDate(), "", 0))
		//sel = append(sel, SQLQueryToSelect(fmt.Sprintf("SELECT * FROM %s WHERE %s='%s'", u.Name(), u.PrimaryKey(), u.Email)))
	}
	return ins, sel, err == nil
}

//Update xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
func (u User) Update(r io.Reader) (upt []SQLQueryToUpdate, sel []SQLQueryToSelect, ok bool) {
	v := []User{}
	err := json.NewDecoder(r).Decode(&v)
	upt = make([]SQLQueryToUpdate, len(v))
	sel = make([]SQLQueryToSelect, len(v))
	//fmt.Println(v)
	i := 0
	for _, u = range v {
		upt[i].ID = string(u.Email)
		if u.FirstName != "" {
			upt[i].Stmts = append(upt[i].Stmts, string("UPDATE users SET firstname=? WHERE email=?"))
			upt[i].Values = append(upt[i].Values, u.FirstName)
		}
		if u.LastName != "" {
			upt[i].Stmts = append(upt[i].Stmts, string("UPDATE users SET lastname=? WHERE email=?"))
			upt[i].Values = append(upt[i].Values, u.LastName)
		}
		if u.PhoneNo != "" {
			upt[i].Stmts = append(upt[i].Stmts, string("UPDATE users SET phoneno=? WHERE email=?"))
			upt[i].Values = append(upt[i].Values, u.PhoneNo)
		}
		if u.Password != "" {
			upt[i].Stmts = append(upt[i].Stmts, string("UPDATE users SET password=? WHERE email=?"))
			upt[i].Values = append(upt[i].Values, utils.EncryptPassword(u.Password))
		}
		if u.MessageID != "" {
			upt[i].Stmts = append(upt[i].Stmts, string("UPDATE users SET messageID=? WHERE email=?"))
			upt[i].Values = append(upt[i].Values, string(u.MessageID))
		}
		upt[i].Stmts = append(upt[i].Stmts, string("UPDATE users SET updatedOn=? WHERE email=?"))
		date := fmt.Sprintf("%s", utils.NewDate())
		upt[i].Values = append(upt[i].Values, date)
		sel[i] = SQLQueryToSelect(fmt.Sprintf("SELECT * FROM %s WHERE %s='%s'", u.Name(), u.PrimaryKey(), u.Email))
		i++
	}
	return upt, sel, err == nil
}

//UpdateWhere xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
func (u User) UpdateWhere(r io.Reader, value string) (upt SQLQueryToUpdate, sel SQLQueryToSelect, ok bool) {
	err := json.NewDecoder(r).Decode(&u)
	upt.ID = value
	if u.FirstName != "" {
		upt.Stmts = append(upt.Stmts, string("UPDATE users SET firstname=? WHERE email=?"))
		upt.Values = append(upt.Values, u.FirstName)
	}
	if u.LastName != "" {
		upt.Stmts = append(upt.Stmts, string("UPDATE users SET lastname=? WHERE email=?"))
		upt.Values = append(upt.Values, u.LastName)
	}
	if u.PhoneNo != "" {
		upt.Stmts = append(upt.Stmts, string("UPDATE users SET phoneno=? WHERE email=?"))
		upt.Values = append(upt.Values, u.PhoneNo)
	}
	if u.Password != "" {
		upt.Stmts = append(upt.Stmts, string("UPDATE users SET password=? WHERE email=?"))
		upt.Values = append(upt.Values, utils.EncryptPassword(u.Password))
	}
	if u.MessageID != "" {
		upt.Stmts = append(upt.Stmts, string("UPDATE users SET messageID=? WHERE email=?"))
		upt.Values = append(upt.Values, string(u.MessageID))
	}
	upt.Stmts = append(upt.Stmts, string("UPDATE users SET updatedOn=? WHERE email=?"))
	date := fmt.Sprintf("%s", utils.NewDate())
	upt.Values = append(upt.Values, date)
	sel = SQLQueryToSelect(fmt.Sprintf("SELECT * FROM %s WHERE %s='%s'", u.Name(), u.PrimaryKey(), u.Email))

	return upt, sel, err == nil
}

//UpdateIf xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
func (u User) UpdateIf(r io.Reader, value string) (upt SQLQueryToUpdate, sel SQLQueryToSelect, ok bool) {
	err := json.NewDecoder(r).Decode(&u)
	if string(u.Email) == value {
		upt.ID = string(u.Email)
		if u.FirstName != "" {
			upt.Stmts = append(upt.Stmts, string("UPDATE users SET firstname=? WHERE email=?"))
			upt.Values = append(upt.Values, u.FirstName)
		}
		if u.LastName != "" {
			upt.Stmts = append(upt.Stmts, string("UPDATE users SET lastname=? WHERE email=?"))
			upt.Values = append(upt.Values, u.LastName)
		}
		if u.PhoneNo != "" {
			upt.Stmts = append(upt.Stmts, string("UPDATE users SET phoneno=? WHERE email=?"))
			upt.Values = append(upt.Values, u.PhoneNo)
		}
		if u.Password != "" {
			upt.Stmts = append(upt.Stmts, string("UPDATE users SET password=? WHERE email=?"))
			upt.Values = append(upt.Values, utils.EncryptPassword(u.Password))
		}
		if u.MessageID != "" {
			upt.Stmts = append(upt.Stmts, string("UPDATE users SET messageID=? WHERE email=?"))
			upt.Values = append(upt.Values, string(u.MessageID))
		}
		upt.Stmts = append(upt.Stmts, string("UPDATE users SET updatedOn=? WHERE email=?"))
		date := fmt.Sprintf("%s", utils.NewDate())
		upt.Values = append(upt.Values, date)
		sel = SQLQueryToSelect(fmt.Sprintf("SELECT * FROM %s WHERE %s='%s'", u.Name(), u.PrimaryKey(), u.Email))
	}
	return upt, sel, err == nil
}

//FromQueryResult xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
func (u User) FromQueryResult(s interfaces.Iterator) (JSONObject, int) {
	//fmt.Printf("%v", u)

	v := []User{}
	var n int
	var id int
	for n = 0; s.Next(); n++ {
		//fmt.Println("yyyyyyyyyyyyyyyyyyyyyyyyyyyyyy")
		err := s.Scan(&u.FirstName, &u.LastName, &u.Email, &u.PhoneNo, &u.Password, &u.MessageID, &u.CreatedOn, &u.UpdatedOn, &id)
		//fmt.Println(u.FirstName)
		if err != nil {
			break
		}
		v = append(v, u)
	}
	result := struct {
		N    int    `json:"rows"`
		Data []User `json:"data"`
	}{N: n, Data: v}
	o, _ := json.MarshalIndent(result, "", "  ")
	return JSONObject(o), n
}

//FromQueryResultArray xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
func (u User) FromQueryResultArray(ss []interfaces.Iterator) (JSONObject, int) {
	//fmt.Printf("%v", u)

	v := []User{}
	var id int
	n := 0
	for _, s := range ss {
		for ; s.Next(); n++ {
			//fmt.Println("yyyyyyyyyyyyyyyyyyyyyyyyyyyyyy")
			err := s.Scan(&u.FirstName, &u.LastName, &u.Email, &u.PhoneNo, &u.Password, &u.MessageID, &u.CreatedOn, &u.UpdatedOn, &id)
			//fmt.Println(u.FirstName)
			if err != nil {
				break
			}
			v = append(v, u)
		}
	}
	result := struct {
		N    int    `json:"rows"`
		Data []User `json:"data"`
	}{N: n, Data: v}
	o, _ := json.MarshalIndent(result, "", "  ")
	return JSONObject(o), n
}

//SelectFromWhere xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
func (u User) SelectFromWhere(value string) (q SQLQueryToSelect, ok bool) {
	q = SQLQueryToSelect("SELECT * FROM " + u.Name())
	q = q + SQLQueryToSelect(fmt.Sprintf(" WHERE %s = '%s'", u.PrimaryKey(), value))
	return q, true
}

//SelectFromRequest xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
/*func (u User) SelectFromRequest(req *http.Request, param string) (q SQLQueryToSelect, ok bool) {
	q = SQLQueryToSelect("SELECT * FROM " + u.Name())
	id, ok := GetParamFromRequest(req, param)
	if ok {
		q = q + SQLQueryToSelect(fmt.Sprintf(" WHERE %s = '%s'", u.PrimaryKey(), id))
	}
	return q, ok
}*/

//DeleteFromWhere xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
func (u User) DeleteFromWhere(key string) (q SQLQueryToDelete, ok bool) {
	q = SQLQueryToDelete(fmt.Sprintf("DELETE FROM %s WHERE %s = '%s'", u.Name(), u.PrimaryKey(), key))
	return q, ok
}
