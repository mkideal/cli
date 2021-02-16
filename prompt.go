package cli

import (
	"errors"
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

var (
	errRequiredMissing = errors.New("required missing")
	errInvalidBoolean  = errors.New("invalid boolean")
)

type readerWriter struct {
	io.Reader
	io.Writer
}

func doPrompt(text string, password bool) (string, error) {
	stdin := os.Stdin
	stdout := os.Stdout
	term := terminal.NewTerminal(readerWriter{stdin, stdout}, text)
	stdinFD := int(stdin.Fd())
	stdinState, err := terminal.MakeRaw(stdinFD)
	if err != nil {
		return "", err
	}
	defer terminal.Restore(stdinFD, stdinState)

	var line string
	if password {
		line, err = term.ReadPassword(text)
	} else {
		line, err = term.ReadLine()
	}
	if err != nil {
		return "", err
	}
	return line, nil
}

func prompt(text string, required bool) (string, error) {
	line, err := doPrompt(text, false)
	if err != nil {
		return line, err
	}
	if required && line == "" {
		return line, errRequiredMissing
	}
	return line, err
}

func promptDefault(text string, dft string) (string, error) {
	line, err := doPrompt(text, false)
	if err != nil {
		return line, err
	}
	if line == "" {
		return dft, nil
	}
	return line, err
}

func password(text string) (string, error) {
	return doPrompt(text, true)
}

func ask(question string, dft bool) (bool, error) {
	line, err := doPrompt(question, false)
	if err != nil {
		return false, err
	}
	if line == "" {
		return dft, nil
	}
	return line == "y" || line == "Y" || line == "yes" || line == "Yes" || line == "YES", nil
}
