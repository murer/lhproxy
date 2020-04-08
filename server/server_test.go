package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"os"
	"io/ioutil"

	"github.com/stretchr/testify/assert"

	"github.com/murer/lhproxy/util"
)

func TestVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(Handle))
	defer server.Close()
	t.Logf("URL: %s", server.URL)
	resp, err := http.Get(server.URL + "/version.txt")
	util.Check(err)
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, util.Version, util.ReadAllString(resp.Body))
}

func TestSelf(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(Handle))
	defer server.Close()
	t.Logf("URL: %s", server.URL)
	resp, err := http.Get(server.URL + "/self/lhproxy")
	util.Check(err)
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)
	path, err := os.Executable()
	util.Check(err)
	localData, err := ioutil.ReadFile(path)
	util.Check(err)
	assert.Equal(t, util.SHA256(localData), util.SHA256(util.ReadAll(resp.Body)))
}
