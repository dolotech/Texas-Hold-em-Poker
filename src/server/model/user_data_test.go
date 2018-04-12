package model

import "testing"

func TestUserData_Login(t *testing.T) {
	u:= &UserData{AccountID:"123"}

	t.Log(u.Insert())
	u= &UserData{AccountID:"123"}
	t.Log(u.ExistByAccountID())
}


func TestSeq(t *testing.T) {
	userID, err := mongoDBNextSeq("users")
	if err != nil {
		t.Logf("get next users id error: %v", err)
	}

	t.Log(userID)
}
