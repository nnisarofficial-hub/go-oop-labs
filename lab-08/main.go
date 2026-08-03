package main

import (
	"errors"
	"fmt"
)

type Reader interface {
	Read(id string) (string, error)
}

type Writer interface {
	Write(id, content string) error
}

type Deleter interface {
	Delete(id string) error
}

type ReadWriter interface {
	Reader
	Writer
}

type ReadWriteDeleter interface {
	Reader
	Writer
	Deleter
}

type DocumentStore struct {
	docs map[string]string
}

func NewDocumentStore() *DocumentStore {
	return &DocumentStore{
		docs: make(map[string]string),
	}
}

func (d *DocumentStore) Read(id string) (string, error) {
	doc, ok := d.docs[id]
	if !ok {
		return "", errors.New("document not found")
	}
	return doc, nil
}

func (d *DocumentStore) Write(id, content string) error {
	d.docs[id] = content
	return nil
}

func (d *DocumentStore) Delete(id string) error {
	_, ok := d.docs[id]
	if !ok {
		return fmt.Errorf("document with id %s not found", id)
	}
	delete(d.docs, id)
	return nil
}

func DisplayDocument(r Reader, id string) {
	content, err := r.Read(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Document '%s': %s\n", id, content)
}

func BackupDocument(rw ReadWriter, fromID, toID string) {
	content, err := rw.Read(fromID)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = rw.Write(toID, content)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Backed up '%s' → '%s'\n", fromID, toID)
	fmt.Printf("%s: %s\n", toID, content)
}

func ArchiveDocument(rwd ReadWriteDeleter, id, archiveID string) {
	content, err := rwd.Read(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = rwd.Write(archiveID, content)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = rwd.Delete(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Archived '%s': moved to '%s', original deleted\n", id, archiveID)
}

func main() {
	store := NewDocumentStore()
	store.Write("doc-1", "Hello, this is document one.")
	store.Write("doc-2", "Hello, this is document two.")
	DisplayDocument(store, "doc-1")
	BackupDocument(store, "doc-1", "doc-1-backup")
	ArchiveDocument(store, "doc-2", "archive-doc-2")
	_, err := store.Read("doc-2")
	if err != nil {
		fmt.Printf("Read deleted doc: %s\n", err)
	}
}
