package main

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	rootDir := "testdata/env"

	badFileName := path.Join(rootDir, "BAD=name")
	if err := os.WriteFile(badFileName, []byte("NO_CHANCES"), os.ModePerm); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(badFileName)

	childDir := path.Join(rootDir, "TEST")
	if err := os.Mkdir(childDir, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(childDir)

	env, err := ReadDir(rootDir)

	require.NoError(t, err)

	cases := []struct {
		giveKey          string
		wantValue        string
		wantMeedToRemove bool
	}{
		{giveKey: "UNSET", wantValue: "", wantMeedToRemove: true},
		{giveKey: "EMPTY", wantValue: "", wantMeedToRemove: false},
		{giveKey: "BAR", wantValue: "bar", wantMeedToRemove: false},
		{giveKey: "HELLO", wantValue: `"hello"`, wantMeedToRemove: false},
		{giveKey: "FOO", wantValue: "   foo\nwith new line", wantMeedToRemove: false},
	}

	require.Len(t, env, len(cases)) // without file with bad name and directory.

	for _, tc := range cases {
		v, exist := env[tc.giveKey]

		require.True(t, exist)
		require.Equal(t, tc.wantValue, v.Value)
		require.Equal(t, tc.wantMeedToRemove, v.NeedRemove)
	}
}

func TestEnvironment_Apply(t *testing.T) {
	testEnvs := map[string]string{
		"TEST1": "VALUE1",
		"TEST2": "VALUE2",
		"TEST3": "VALUE3",
	}

	for k, v := range testEnvs {
		require.NoError(t, os.Setenv(k, v))
	}

	env := Environment{
		"TEST1": EnvValue{Value: "", NeedRemove: true},
		"TEST2": EnvValue{Value: "", NeedRemove: false},
		"TEST3": EnvValue{Value: "NEW_VALUE", NeedRemove: false},
	}

	require.NoError(t, env.Apply())

	_, exist := os.LookupEnv("TEST1")
	require.False(t, exist)

	v, exist := os.LookupEnv("TEST2")
	require.True(t, exist)
	require.Equal(t, "", v)

	v, exist = os.LookupEnv("TEST3")
	require.True(t, exist)
	require.Equal(t, "NEW_VALUE", v)
}
