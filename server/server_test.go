package server

import (
	"testing"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"

	"github.com/murer/lhproxy/util"
)

func TestVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handle))
	defer server.Close()
	t.Logf("URL: %s", server.URL)
	resp, err := http.Get(server.URL + "/version.txt")
	util.Check(err)
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, util.Version, util.ReadAllString(resp.Body))
}
