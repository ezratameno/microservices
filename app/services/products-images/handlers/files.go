package handlers

import (
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/ezratameno/microservices/app/services/products-images/files"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

// Files is a handler for reading and writing files
type Files struct {
	log   hclog.Logger
	store files.Storage
}

// NewFiles creates a new File handler
func NewFiles(s files.Storage, l hclog.Logger) *Files {
	return &Files{store: s, log: l}
}

// UploadRest implements the http.Handler interface
func (f *Files) UploadRest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	filename := vars["filename"]

	f.log.Info("Handle POST", "id", id, "filename", filename)

	// no need to check for invalid id or filename as the mux router will not send requests
	// here unless they have the correct parameters

	f.saveFile(id, filename, w, r.Body)
}

func (f *Files) UploadMultipart(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(128 * 1024)
	if err != nil {
		f.log.Error("Bad request", err)
		http.Error(w, "Expected multipart form data", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		f.log.Error("Bad request", err)
		http.Error(w, "Expected integer id", http.StatusBadRequest)
		return
	}

	f.log.Info("Process form for id", "id", id)

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		f.log.Error("Bad request", err)
		http.Error(w, "Expected file", http.StatusBadRequest)
		return
	}

	f.saveFile(r.FormValue("id"), fileHeader.Filename, w, file)
}

func (f *Files) invalidURI(uri string, w http.ResponseWriter) {
	f.log.Error("Invalid path", "path", uri)
	http.Error(w, "Invalid file path should be in the format: /[id]/[filepath]", http.StatusBadRequest)
}

// saveFile saves the contents of the request to a file
func (f *Files) saveFile(id, path string, w http.ResponseWriter, r io.ReadCloser) {
	f.log.Info("Save file for product", "id", id, "path", path)

	fp := filepath.Join(id, path)
	err := f.store.Save(fp, r)
	if err != nil {
		f.log.Error("Unable to save file", "error", err)
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
	}
}
