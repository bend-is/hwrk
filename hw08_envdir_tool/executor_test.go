package main

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		output, err := captureStdout(func() {
			cmd := []string{"/bin/bash", "testdata/echo.sh", "arg1"}
			env := Environment{"FOO": EnvValue{Value: "TEST"}}

			code := RunCmd(cmd, env)

			require.Equal(t, 0, code)
		})

		require.NoError(t, err)

		want := `HELLO is ()
BAR is ()
FOO is (TEST)
UNSET is ()
ADDED is ()
EMPTY is ()
arguments are arg1
`

		require.Equal(t, want, output)
	})

	t.Run("error from cmd", func(t *testing.T) {
		scriptName := "./script.sh"
		script := []byte(
			`#!/usr/bin/env bash
			>&2 echo "test error" && exit 5
		`)

		if err := os.WriteFile(scriptName, script, os.ModePerm); err != nil {
			t.Fatal(err)
		}
		defer os.Remove(scriptName)

		output, err := captureStderr(func() {
			code := RunCmd([]string{scriptName}, nil)

			require.Equal(t, 5, code)
		})

		require.NoError(t, err)
		require.Equal(t, "test error\n", output)
	})

	t.Run("empty command list", func(t *testing.T) {
		output, err := captureStderr(func() {
			code := RunCmd(nil, nil)

			require.Equal(t, 1, code)
		})

		require.NoError(t, err)
		require.Equal(t, "Empty command list\n", output)
	})
}

func captureStdout(f func()) (string, error) {
	return captureStdoutOrStderr(f, true)
}

func captureStderr(f func()) (string, error) {
	return captureStdoutOrStderr(f, false)
}

func captureStdoutOrStderr(f func(), needStdout bool) (string, error) {
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}
	defer r.Close()

	if needStdout {
		backup := os.Stdout
		os.Stdout = w

		defer func() {
			os.Stdout = backup
		}()
	} else {
		backup := os.Stderr
		os.Stderr = w

		defer func() {
			os.Stderr = backup
		}()
	}

	f()

	if err := w.Close(); err != nil {
		return "", err
	}

	b, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
