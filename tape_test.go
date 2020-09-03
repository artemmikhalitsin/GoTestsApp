package poker_test

import (
	"io/ioutil"
	"testing"

	poker "github.com/artemmikhalitsin/GoTestsApp"
)

func TestTape_Write(t *testing.T) {
	file, clean := createTempFile(t, "12345")
	defer clean()

	tape := &poker.Tape{file}
	tape.Write([]byte("abc"))

	file.Seek(0, 0)
	newFileContents, _ := ioutil.ReadAll(file)

	got := string(newFileContents)
	want := "abc"

	if got != want {
		t.Errorf("Expected file to contain %q, but got %q", want, got)
	}
}
