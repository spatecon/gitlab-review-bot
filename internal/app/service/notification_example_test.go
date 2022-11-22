package service

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// This test can be used with testdata/user_template.go.tpl for generation json string to put into mongodb.
func TestExample_BuildUserNotification(t *testing.T) {
	f, err := os.Open("testdata/user_template.go.tpl")
	require.NoError(t, err)
	defer f.Close()

	str, err := io.ReadAll(f)
	require.NoError(t, err)

	val, err := json.Marshal(string(str))
	require.NoError(t, err)

	fmt.Println(string(val))
}

// This test can be used with testdata/channel_template.go.tpl for generation json string to put into mongodb.
func TestExample_BuildChannelNotification(t *testing.T) {
	f, err := os.Open("testdata/channel_template.go.tpl")
	require.NoError(t, err)
	defer f.Close()

	str, err := io.ReadAll(f)
	require.NoError(t, err)

	val, err := json.Marshal(string(str))
	require.NoError(t, err)

	fmt.Println(string(val))
}
