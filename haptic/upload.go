package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"

	nan "nanocloud.com/core/lib/libnan"
)

var completedFiles = make(chan string, 100)

type ByChunk []os.FileInfo

func (a ByChunk) Len() int      { return len(a) }
func (a ByChunk) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByChunk) Less(i, j int) bool {
	ai, _ := strconv.Atoi(a[i].Name())
	aj, _ := strconv.Atoi(a[j].Name())
	return ai < aj
}

func uploadHandler(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		return streamingReader(w, r)
	} else if r.Method == "GET" {
		return continueUpload(w, r)
	} else {
		return errors.New("Not found")
	}
}

type streamHandler func(http.ResponseWriter, *http.Request) error

func (fn streamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func continueUpload(w http.ResponseWriter, r *http.Request) error {
	chunkDirPath := nan.Config().CommonBaseDir + "/incomplete/" + r.FormValue("flowFilename") + "/" + r.FormValue("flowChunkNumber")
	if _, err := os.Stat(chunkDirPath); err != nil {
		w.WriteHeader(206)
		return err // streamingReader(w, r)
	}
	return nil
}

func streamingReader(w http.ResponseWriter, r *http.Request) error {
	buf := new(bytes.Buffer)
	reader, err := r.MultipartReader()
	// Part 1: Chunk Number
	// Part 4: Total Size (bytes)
	// Part 6: File Name
	// Part 8: Total Chunks
	// Part 9: Chunk Data
	if err != nil {
		return err
	}

	part, err := reader.NextPart() // 1
	if err != nil {
		return err
	}
	io.Copy(buf, part)
	chunkNo := buf.String()
	buf.Reset()

	for i := 0; i < 3; i++ { // 2 3 4
		// move through unused parts
		part, err = reader.NextPart()
		if err != nil {
			return err
		}
	}

	io.Copy(buf, part)
	flowTotalSize := buf.String()
	buf.Reset()

	for i := 0; i < 2; i++ { // 5 6
		// move through unused parts
		part, err = reader.NextPart()
		if err != nil {
			return err
		}
	}

	io.Copy(buf, part)
	fileName := buf.String()
	buf.Reset()

	for i := 0; i < 3; i++ { // 7 8 9
		// move through unused parts
		part, err = reader.NextPart()
		if err != nil {
			return err
		}
	}

	chunkDirPath := nan.Config().CommonBaseDir + "/incomplete/" + fileName
	err = os.MkdirAll(chunkDirPath, 02750)
	if err != nil {
		return err
	}

	dst, err := os.Create(chunkDirPath + "/" + chunkNo)
	if err != nil {
		return err
	}
	defer dst.Close()
	io.Copy(dst, part)

	fileInfos, err := ioutil.ReadDir(chunkDirPath)
	if err != nil {
		return err
	}

	if flowTotalSize == strconv.Itoa(int(totalSize(fileInfos))) {
		assemble(chunkDirPath, fileName)
		go SyncUploadedFile(nan.Config().CommonBaseDir + "uploads/" + fileName)
	}

	return nil
}

func totalSize(fileInfos []os.FileInfo) int64 {
	var sum int64
	for _, fi := range fileInfos {
		sum += fi.Size()
	}
	return sum
}

func chunkedReader(w http.ResponseWriter, r *http.Request) error {
	r.ParseMultipartForm(25)

	fileName := r.FormValue("flowFilename")
	chunkDirPath := nan.Config().CommonBaseDir + "/incomplete/" + fileName
	err := os.MkdirAll(chunkDirPath, 02750)
	if err != nil {
		return err
	}

	for _, fileHeader := range r.MultipartForm.File["file"] {
		src, err := fileHeader.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		dst, err := os.Create(chunkDirPath + "/" + r.FormValue("flowChunkNumber"))
		if err != nil {
			return err
		}
		defer dst.Close()
		io.Copy(dst, src)

		fileInfos, err := ioutil.ReadDir(chunkDirPath)
		if err != nil {
			return err
		}

		if r.FormValue("flowTotalSize") == strconv.Itoa(int(totalSize(fileInfos))) {
			assemble(chunkDirPath, fileName)
			go SyncUploadedFile(nan.Config().CommonBaseDir + "uploads/" + fileName)
		}
	}
	return nil
}

func assemble(path string, fileName string) {
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}

	// create final file to write to
	err = os.MkdirAll(nan.Config().CommonBaseDir+"uploads", 02750)
	if err != nil {
		return
	}

	dst, err := os.Create(nan.Config().CommonBaseDir + "uploads/" + fileName)
	if err != nil {
		return
	}
	defer dst.Close()

	sort.Sort(ByChunk(fileInfos))
	for _, fs := range fileInfos {
		src, err := os.Open(path + "/" + fs.Name())
		if err != nil {
			return
		}
		defer src.Close()
		io.Copy(dst, src)
	}
	os.RemoveAll(path)
}
