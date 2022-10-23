package main

import (
    "bytes"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "testing"
)

func Test_Help(t *testing.T) {
    t.Run("main help", func(t *testing.T) {
        c := createCLI([]string{"--help"})

        buf := bytes.Buffer{}
        c.ErrorWriter = &buf
        c.HelpWriter = &buf
        _, err := c.Run()
        require.NoError(t, err)

        assert.Contains(t, buf.String(), "equipment")
        assert.Contains(t, buf.String(), "login")
    })

    t.Run("equipment help", func(t *testing.T) {
        c := createCLI([]string{"equipment", "--help"})

        buf := bytes.Buffer{}
        c.ErrorWriter = &buf
        c.HelpWriter = &buf
        _, err := c.Run()
        require.NoError(t, err)

        assert.Contains(t, buf.String(), "show")
        assert.Contains(t, buf.String(), "list")
    })
}
