package users

import (
	"errors"
	"net/http"
	"time"
)

//FakeUserStore exists for testing auth
type FakeUserStore struct {
}

//NewFakeUserStore function for testing auth
func NewFakeUserStore() *FakeUserStore {
	return &FakeUserStore{}
}

//GetByID function for testing auth
func (f *FakeUserStore) GetByID(id int64) (*User, error) {
	if id == 1 {
		return &User{
			ID: 1,
		}, nil
	}

	return nil, errors.New("only id 1 can be found")
}

//GetByEmail function for testing auth
func (f *FakeUserStore) GetByEmail(email string) (*User, error) {
	if email == "valid" {
		u := &User{
			ID:    1,
			Email: "valid",
		}
		u.SetPassword("1")
		return u, nil
	}

	return nil, errors.New("email != valid")
}

//GetByUserName function for testing auth
func (f *FakeUserStore) GetByUserName(username string) (*User, error) {
	if username == "valid" {
		return &User{
			ID:       1,
			UserName: "valid",
		}, nil
	}

	return nil, errors.New("userName != valid")
}

//Insert function for testing auth
func (f *FakeUserStore) Insert(user *User) (*User, error) {
	return user, nil
}

//Update function for testing auth
func (f *FakeUserStore) Update(id int64, updates *Updates) (*User, error) {
	return &User{ID: 1}, nil
}

//Delete function for testing auth
func (f *FakeUserStore) Delete(id int64) error {
	return nil
}

//Track function for testing auth
func (f *FakeUserStore) Track(r *http.Request, id int64, now time.Time) error {
	return nil
}
