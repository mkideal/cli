package cli

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

const DefaultEditor = "vim"

// GetEditor sets callback to get editor program
var GetEditor func() (string, error)

func getEditor() (string, error) {
	if GetEditor != nil {
		return GetEditor()
	}
	return exec.LookPath(DefaultEditor)
}

func launchEditor(editor string) (content []byte, err error) {
	buf := make([]byte, 16)
	_, err = rand.Read(buf)
	if err != nil {
		return
	}
	filename := fmt.Sprintf(".%x", buf)

	return launchEditorWithFilename(editor, filename)
}

func launchEditorWithFilename(editor, filename string) (content []byte, err error) {
	cmd := exec.Command(editor, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return
	}
	defer os.Remove(filename)
	content, err = ioutil.ReadFile(filename)
	if err != nil {
		return []byte{}, nil
	}
	return
}
