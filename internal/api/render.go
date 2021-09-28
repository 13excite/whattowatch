package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/xid"
	"net/http"
)

type ErrResponse struct {
	Status  string `json:"status,omitempty"`
	Error   string `json:"error,omitempty"`
	ErrorID string `json:"error_id,omitempty"`
}

func RenderJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(v); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(b, `{"render_error":"%s"}`, errString(err))
	} else {
		w.WriteHeader(code)
	}
	_, _ = w.Write(b.Bytes())
}

func RenderErrResourceNotFound(w http.ResponseWriter, resource string) {
	RenderJSON(w, http.StatusNotFound, ErrResponse{Status: resource + " not found", Error: resource + " not found"})
}

func RenderErrUnauthorized(w http.ResponseWriter) {
	RenderJSON(w, http.StatusUnauthorized, ErrResponse{Status: "not authorized", Error: "not authorized"})
}

func RenderErrInvalidRequest(w http.ResponseWriter, err error) {
	RenderJSON(w, http.StatusBadRequest, ErrResponse{Status: "invalid request", Error: errString(err)})
}

func RenderErrInternal(w http.ResponseWriter, err error) {
	RenderJSON(w, http.StatusInternalServerError, ErrResponse{Status: "internal error", Error: errString(err)})
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func RenderErrInternalWithID(w http.ResponseWriter, err error) string {
	errID := xid.New().String()
	RenderJSON(w, http.StatusInternalServerError, ErrResponse{Status: "internal error", Error: errString(err), ErrorID: errID})
	return errID
}
