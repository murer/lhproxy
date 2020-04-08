package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"os"
	"io/ioutil"
	"encoding/base64"

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
	path, err := os.Executable()
	util.Check(err)
	localData, err := ioutil.ReadFile(path)
	util.Check(err)

	resp, err := http.Get(server.URL + "/self/lhproxy.txt")
	util.Check(err)
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)
	body := util.ReadAll(resp.Body)
	content, err := base64.StdEncoding.DecodeString(string(body))
	util.Check(err)
	assert.Equal(t, util.SHA256(localData), util.SHA256(content))
	assert.Equal(t, resp.Header["X-Sec"][0], util.SHA256Hex(append(body, util.Secret()...)))
}
