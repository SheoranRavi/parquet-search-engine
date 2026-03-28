package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/SheoranRavi/parquet-search-engine/internal/model"
	"github.com/parquet-go/parquet-go"
)

func main() {
	pathStr, err := os.Executable()
	if err != nil {
		panic(err)
	}
	path := filepath.Dir(pathStr)
	pFilesPath := filepath.Join(path, "..", "parquet_files")
	files, err := os.ReadDir(pFilesPath)
	if err != nil {
		panic(err)
	}
	var oneRow model.Message
	for _, f := range files {
		filePath := filepath.Join(pFilesPath, f.Name())
		rows, _ := parquet.ReadFile[model.Message](filePath)
		oneRow = rows[0]
		fmt.Printf("Number of records: %d, file: %s\n", len(rows), f.Name())
		fHandle, _ := os.Open(filePath)
		stat, _ := fHandle.Stat()
		pf, _ := parquet.OpenFile(fHandle, stat.Size())
		schema := pf.Schema()
		fmt.Println(schema.String())
	}
	fmt.Printf("Message: %+v\n", oneRow)
}

type RowType struct {
	Message    string
	MessageRaw string
}
