package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestLoadKazaamTransformWithMissingFile(t *testing.T) {
	_, err := loadKazaamTransform("")

	if err == nil {
		t.Error("Should have errored for missing file")
	}
}

func TestLoadKazaamTransformWithNoFile(t *testing.T) {
	_, err := loadKazaamTransform("doesnt-exist")

	if err == nil {
		t.Error("Should have errored for non-existent file")
	}
}

func TestLoadKazaamTransformWithInvalidFile(t *testing.T) {
	fd, err := ioutil.TempFile("", "kz-main-test-")
	if err != nil {
		t.Error("Unable to create tmpfile for test", err)
	}
	_, _ = fd.WriteString(`{"invalid json"}`)
	defer os.Remove(fd.Name())
	defer fd.Close()

	_, err = loadKazaamTransform(fd.Name())
	if err == nil {
		t.Error("Should have errored for empty file")
	}
}

func TestLoadKazaamTransform(t *testing.T) {
	fd, err := ioutil.TempFile("", "kz-main-test-")
	if err != nil {
		t.Error("Unable to create tmpfile for test", err)
	}
	_, _ = fd.WriteString(`[{"operation": "pass"}]`)
	fd.Close()
	defer os.Remove(fd.Name())

	_, err = loadKazaamTransform(fd.Name())
	if err != nil {
		t.Error("Shouldn't have errored with valid transform", err)
	}
}

func TestGetInputByFilename(t *testing.T) {
	filenameTestData := `testByFilename`
	fd, err := ioutil.TempFile("", "kz-main-test-")
	if err != nil {
		t.Error("Unable to create tmpfile for test", err)
	}
	_, _ = fd.WriteString(filenameTestData)
	_ = fd.Close()
	defer os.Remove(fd.Name())

	data, err := getInput(fd.Name(), nil)
	if err != nil {
		t.Error("Unexpected error reading file", err)
	}
	if data != filenameTestData {
		t.Error("Unexpected file contents")
	}
}

func TestGetInputByFileHandle(t *testing.T) {
	fileHandleTestData := `testByFileHandle`
	fd, err := ioutil.TempFile("", "kz-main-test-")
	if err != nil {
		t.Error("Unable to create tmpfile for test", err)
	}
	_, _ = fd.WriteString(fileHandleTestData)
	_, _ = fd.Seek(0, 0)
	defer os.Remove(fd.Name())

	data, err := getInput("", fd)
	if err != nil {
		t.Error("Unexpected error reading file", err)
	}
	if data != fileHandleTestData {
		t.Error("Unexpected file contents")
	}
}

func TestGetInputByClosedFileHandle(t *testing.T) {
	fd, err := ioutil.TempFile("", "kz-main-test-")
	if err != nil {
		t.Error("Unable to create tmpfile for test", err)
	}
	_ = fd.Close()
	defer os.Remove(fd.Name())

	_, err = getInput("", fd)
	if err == nil {
		t.Error("Should have thrown error for unreadable file", err)
	}
}

func TestGetInputPriority(t *testing.T) {
	filenameTestData := `testByFilename`
	fileHandleTestData := `testByFileHandle`
	fdA, _ := ioutil.TempFile("", "kz-main-test-")
	fdB, _ := ioutil.TempFile("", "kz-main-test-")
	_, _ = fdA.WriteString(filenameTestData)
	defer os.Remove(fdA.Name())
	_, _ = fdB.WriteString(fileHandleTestData)
	_, _ = fdB.Seek(0, 0)
	defer os.Remove(fdA.Name())

	data, err := getInput(fdA.Name(), fdB)
	if err != nil {
		t.Error("Unexpected error reading file", err)
	}
	if data != filenameTestData {
		t.Error("Unexpected file contents")
	}
}
