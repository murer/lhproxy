package server

import (
	"net/http"
	"os"
	"path/filepath"
	"io/ioutil"
	"fmt"
	"strconv"
	"github.com/murer/lhproxy/util"
)

func handleSelf(w http.ResponseWriter, r *http.Request) {
	path, err := os.Executable()
	util.Check(err)
	localData, err := ioutil.ReadFile(path)
	util.Check(err)
	filename := filepath.Base(r.URL.Path)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Content-Length", strconv.Itoa(len(localData)))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(localData)
}
