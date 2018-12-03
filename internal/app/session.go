package app

import (
	"encoding/gob"

	"github.com/gorilla/sessions"
)

var (
	Store *sessions.FilesystemStore
)

// Init -
func Init() error {
	Store = sessions.NewFilesystemStore("", []byte("gophers-secret"))
	gob.Register(map[string]interface{}{})
	return nil
}
