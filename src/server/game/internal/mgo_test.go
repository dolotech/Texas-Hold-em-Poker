package internal

import (
	"testing"
)

func TestSeq(t *testing.T) {
	userID, err := mongoDBNextSeq("users")
	if err != nil {
		t.Logf("get next users id error: %v", err)
	}

	t.Log(userID)
}
