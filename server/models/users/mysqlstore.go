package users

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	//comment to satisfy lintr
	_ "github.com/go-sql-driver/mysql"
)

//MySQLStore stores a pointer to a database
type MySQLStore struct {
	DB *sql.DB
}

//GetByID returns the User with the given ID
func (store *MySQLStore) GetByID(id int64) (*User, error) {
	q := "select userid, email, firstname, lastname, username, photourl, passhash from users where id = ?"
	row := store.DB.QueryRow(q, id)

	u := &User{}
	err := row.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.UserName, &u.PhotoURL, &u.PassHash)

	if err != nil {
		return nil, err
	}

	return u, nil
}

//GetByEmail returns the User with the given email
func (store *MySQLStore) GetByEmail(email string) (*User, error) {
	q := "select userid, email, firstname, lastname, username, photourl, passhash from users where email = ?"
	row := store.DB.QueryRow(q, email)

	u := &User{}
	err := row.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.UserName, &u.PhotoURL, &u.PassHash)

	if err != nil {
		return nil, err
	}

	return u, nil
}

//GetByUserName returns the User with the given Username
func (store *MySQLStore) GetByUserName(username string) (*User, error) {
	q := "select userid, email, firstname, lastname, username, photourl, passhash from users where username = ?"
	row := store.DB.QueryRow(q, username)

	u := &User{}
	err := row.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.UserName, &u.PhotoURL, &u.PassHash)

	if err != nil {
		return nil, err
	}
	return u, nil
}

//Insert inserts the user into the database, and returns
//the newly-inserted User, complete with the DBMS-assigned ID
func (store *MySQLStore) Insert(user *User) (*User, error) {
	insq := "insert into users(email, firstname, lastname, username, photourl, passhash) values(?,?,?,?,?,?)"

	res, err := store.DB.Exec(insq, user.Email, user.FirstName, user.LastName, user.UserName, user.PhotoURL, user.PassHash)
	if err != nil {
		fmt.Printf("error inserting new row: %v\n", err)
		return nil, err
	}

	//get the auto-assigned ID for the new row
	id, err := res.LastInsertId()
	if err != nil {
		fmt.Printf("error getting new ID: %v\n", id)
		return nil, err
	}

	return user, nil
}

//Update applies UserUpdates to the given user ID
//and returns the newly-updated user
func (store *MySQLStore) Update(id int64, updates *Updates) (*User, error) {
	queryF := "update users set firstname = ? , lastname = ? where userid = ?"

	if _, err := store.DB.Exec(queryF, updates.FirstName, updates.LastName, id); err != nil {
		return nil, err
	}

	return store.GetByID(id)
}

//Delete deletes the user with the given ID
func (store *MySQLStore) Delete(id int64) error {
	ex := "delete from users where userid = ?"

	_, err := store.DB.Exec(ex, id)

	if err != nil {
		return err
	}

	return nil
}

//Track tracks sessions
func (store *MySQLStore) Track(r *http.Request, id int64, now time.Time) error {
	insq := "insert into userlog(userid, timeinitiated, ip) values(?,?,?)"
	//insq = insert query

	ip := r.RemoteAddr

	if ip2 := r.Header.Get("X-FORWARDED-FOR"); len(ip2) != 0 {
		ips := strings.Split(ip2, ", ")
		ip = strings.TrimSpace(ips[0])
	}

	res, err := store.DB.Exec(insq, id, now, ip)

	if err != nil {
		fmt.Printf("error inserting new row: %v\n", err)
		return err
	}

	_, err2 := res.LastInsertId()

	if err2 != nil {
		return err
	}

	return nil
}
