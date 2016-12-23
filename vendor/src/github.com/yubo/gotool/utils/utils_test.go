package utils

import (
	"fmt"
	"os/user"
	"testing"
)

func TestGetUGroups(t *testing.T) {
	u, _ := user.Current()
	if gs, err := GetUGroups(u.Username); err != nil {
		t.Fatalf("GetUGroups: %v", err)
	} else {
		fmt.Printf("username:%s groups:%v\n", u.Username, gs)
	}
}
