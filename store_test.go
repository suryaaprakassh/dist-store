package main

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathTransform(t *testing.T) {
	path := CASPathTransformFunc("test")
	//TODO: fix hardcoded value in test
	expectedPath := "a94a8/fe5cc/b19ba/61c4c/0873d/391e9/87982/fbbd3"
	assert.NotNil(t, path)
	assert.Equal(t, expectedPath, path.Path)
	path = CASPathTransformFunc("testies")
	assert.NotEqual(t, expectedPath, path)
}

func TestStore(t *testing.T) {
	opts := StorageOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	store := NewStore(opts)
	key := "testKey"
	data := []byte("this is a test")
	dataReader := bytes.NewReader(data)

	//Create
	assert.NoError(t, store.Write(key, dataReader))

	//Read
	file, err := store.Read(key)
	assert.Nilf(t, err, "File with key %v should exist", key)

	fileData ,err :=io.ReadAll(file)
	assert.Nilf(t, err, "File with key %v should be read without errors", key)

	assert.Equalf(t, data,fileData , "Contents of the file with key %v should be equal", key)

	//Delete
	assert.Truef(t, store.Has(key), "File with key %v should exist", key)

	assert.NoErrorf(t, store.Delete(key), "Should delete file with key %v", key)

	assert.Falsef(t, store.Has(key), "File with key %v should not exist", key)
}
