package gogo

import (
	"os"
	"testing"
	"time"
)

func Test(t *testing.T) {
	Letsgogo()
	cc := false
	gsc := false
	RegisterCommand("cc", func(command string) {
		t.Log("cc")
		cc = true
	})
	RegisterCommand("gsc", func(command string) {
		t.Log("gsc")
		gsc = true
	})

	os.Stdin.Write([]byte("cc\n"))
	os.Stdin.Write([]byte("cc:1\n"))
	os.Stdin.Write([]byte("gsc\n"))
	os.Stdin.Write([]byte("gsc:1\n"))

	<-time.After(1000)
	if !cc {
		t.FailNow()
	}
	if !gsc {
		t.FailNow()
	}
}
