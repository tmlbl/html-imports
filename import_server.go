package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const importStr = "<!-- import"

type ImportServer struct {
	root       string
	fileServer http.Handler
}

func NewImportServer(root string) *ImportServer {
	return &ImportServer{
		root:       root,
		fileServer: http.FileServer(http.Dir(root)),
	}
}

func (s *ImportServer) isIndexPath(path string) bool {
	p := filepath.Join(s.root, path, "index.html")
	_, err := os.Stat(p)
	return err == nil
}

func (s *ImportServer) executeTemplate(path, indexPath string) ([]byte, error) {
	p := filepath.Join(s.root, path, "index.html")
	data, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}
	outLines := []string{}
	lines := strings.Split(string(data), "\n")
	for _, ln := range lines {
		importIx := strings.Index(ln, importStr)
		if importIx != -1 {
			endIx := strings.Index(ln, "-->")
			importPath := filepath.Join(indexPath, strings.Trim(ln[importIx+len(importStr):endIx], " "))
			data, err := ioutil.ReadFile(importPath)
			if err != nil {
				return nil, err
			}
			outLines = append(outLines, string(data))
		} else {
			outLines = append(outLines, ln)
		}
	}
	return []byte(strings.Join(outLines, "\n")), nil
}

func (s *ImportServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.isIndexPath(r.URL.Path) {
		data, err := s.executeTemplate(r.URL.Path, path.Dir(r.URL.Path)[1:])
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(`
			<h1 style="color:red">An error occurred</h1>
			<h3>` + err.Error() + "</h3>"))
			return
		}
		w.WriteHeader(200)
		w.Write(data)
		w.Header().Add("Content-Type", "text/html")
		return
	}
	s.fileServer.ServeHTTP(w, r)
}

func main() {
	http.ListenAndServe(":8000", NewImportServer("."))
}
