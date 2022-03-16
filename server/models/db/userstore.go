package db

import (
	"chat/server/models/users"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type UserStore struct {
	conn *Connection
}

func (c *Connection) InitUserStore() *UserStore {
	return &UserStore{conn: c}
}

//GetByID returns the User with the given ID
func (store *UserStore) GetByID(id int64) (*users.User, error) {
	db := store.conn.db

	q := "select userid, email, firstname, lastname, username, photourl, passhash from users where id = ?"
	row := db.QueryRow(q, id)

	u := &users.User{}
	err := row.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.UserName, &u.PhotoURL, &u.PassHash)

	if err != nil {
		return nil, err
	}

	return u, nil
}

//GetByEmail returns the User with the given email
func (store *UserStore) GetByEmail(email string) (*users.User, error) {
	db := store.conn.db

	q := "select userid, email, firstname, lastname, username, photourl, passhash from users where email = ?"
	row := db.QueryRow(q, email)

	u := &users.User{}
	err := row.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.UserName, &u.PhotoURL, &u.PassHash)

	if err != nil {
		return nil, err
	}

	return u, nil
}

//GetByUserName returns the User with the given Username
func (store *UserStore) GetByUserName(username string) (*users.User, error) {

	db := store.conn.db
	q := "select userid, email, firstname, lastname, username, photourl, passhash from users where username = ?"
	row := db.QueryRow(q, username)

	u := &users.User{}
	err := row.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.UserName, &u.PhotoURL, &u.PassHash)

	if err != nil {
		return nil, err
	}
	return u, nil
}

//Insert inserts the user into the database, and returns
//the newly-inserted User, complete with the DBMS-assigned ID
func (store *UserStore) Insert(user *users.User) (*users.User, error) {
	db := store.conn.db
	insq := "insert into users(email, firstname, lastname, username, photourl, passhash) values(?,?,?,?,?,?)"

	res, err := db.Exec(insq, user.Email, user.FirstName, user.LastName, user.UserName, user.PhotoURL, user.PassHash)
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
func (store *UserStore) Update(id int64, updates *users.Updates) error {

	db := store.conn.db
	queryF := "update users set firstname = ? , lastname = ? where userid = ?"

	_, err := db.Exec(queryF, updates.FirstName, updates.LastName, id)

	return err

}

//Delete deletes the user with the given ID
func (store *UserStore) Delete(id int64) error {
	db := store.conn.db
	ex := "delete from users where userid = ?"

	_, err := db.Exec(ex, id)

	if err != nil {
		return err
	}

	return nil
}

//Track tracks sessions
func (store *UserStore) Track(r *http.Request, id int64, now time.Time) error {
	db := store.conn.db
	insq := "insert into userlog(userid, timeinitiated, ip) values(?,?,?)"
	//insq = insert query

	ip := r.RemoteAddr

	if ip2 := r.Header.Get("X-FORWARDED-FOR"); len(ip2) != 0 {
		ips := strings.Split(ip2, ", ")
		ip = strings.TrimSpace(ips[0])
	}

	res, err := db.Exec(insq, id, now, ip)

	if err != nil {
		fmt.Printf("error inserting new row: %v\n", err)
		return err
	}

	_, err = res.LastInsertId()

	if err != nil {
		return err
	}

	return nil
}
