package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestSignature(t *testing.T) {
	newFile, error := os.Create("new.file")
	defer newFile.Close()
	defer os.Remove("new.file")
	defer os.Remove("new.sigfile")

	if error != nil {
		log.Println("TestSignature: Error on create file", error)
		t.Error("Errror:")
	}
	newFile.Write([]byte("This is a Fresh file"))

	Signature("new.file", "new.sigfile")
	if _, error := os.Open("new.sigfile"); error != nil {
		log.Println("TestSignature > Error on open sigfile >", error)
		t.Error("should create a sigFile")
	}
}

func TestDelta(t *testing.T) {
	newFile, error := os.Create("new.file")
	defer newFile.Close()
	defer os.Remove("new.file")
	defer os.Remove("new.sigfile")
	defer os.Remove("new.deltafile")

	if error != nil {
		log.Println("TestSignature: Error on create file", error)
		t.Error("Errror:")
	}

	newFile.Write([]byte("This is a Fresh file"))
	Signature("new.file", "new.sigfile")

	Delta("new.file", "new.sigfile", "new.deltafile")
	if _, error := os.Open("new.deltafile"); error != nil {
		log.Println("TestDelta > Error on open deltafile >", error)
		t.Error("should create a deltaFile")
	}
}

func TestPatch(t *testing.T) {
	oldFile, error := os.Create("old.file")
	defer oldFile.Close()
	defer os.Remove("old.file")
	defer os.Remove("old.sigfile")
	defer os.Remove("old.deltafile")

	if error != nil {
		log.Println("TestSignature: Error on create file", error)
		t.Error("Errror:")
	}

	oldFile.Write([]byte("This is the old file"))
	Signature("old.file", "old.sigfile")

	oldFile.Write([]byte(" with new content"))
	Delta("old.file", "old.sigfile", "old.deltafile")

	Patch("old.file", "old.deltafile", "new.patchedfile")
	patchedFile, error := ioutil.ReadFile("new.patchedfile")
	if content := string(patchedFile); content != "This is the old file with new content" {
		t.Error("should apply the patch to the old file", content)
	}
}
