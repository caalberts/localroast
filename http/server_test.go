package http

import (
	"testing"

	"github.com/caalberts/localghost"

	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	port := "8888"
	server := NewServer(port, localghost.Schema{Path: "/"})
	assert.Equal(t, ":8888", server.Addr)
}
