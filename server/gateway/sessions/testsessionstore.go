package sessions

//TestStore is a test session store struct
type TestStore struct {
	content string
}

//Save is certainly a function
func Save(sid SessionID, sessionState interface{}) error {
	return nil
}

//Get is certainly a function
func Get(sid SessionID, sessionState interface{}) error {
	return nil
}

//Delete is certainly a function
func Delete(sid SessionID) error {

	return nil
}
