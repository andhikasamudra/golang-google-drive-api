package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"google.golang.org/api/drive/v3"
	"log"
	"os"
)

func main() {
	gDrive, err := NewGDriveService(GoogleDriveScopeDrive, GoogleDriveScopeDriveFile)
	if err != nil {
		fmt.Println(err)
	}

	folderId := os.Getenv("GoogleDriveFolderID")

	fileList, err := gDrive.GetFileList()
	if err != nil {
		fmt.Println(err)
	}

	for _, item := range fileList.Files {
		if item.Name == "test.csv" {
			err := gDrive.DeleteFile(item.Id)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	testCSV := [][]string{
		{"test1", "ganti update"},
		{"test2", "ganti baris 2"},
	}

	b := new(bytes.Buffer)
	w := csv.NewWriter(b)

	w.WriteAll(testCSV)

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}

	driveMetaData := drive.File{
		Name:     "test.csv",
		MimeType: MimeTypeCSV,
		Parents:  []string{folderId},
	}

	_, err = gDrive.CreateFile(&driveMetaData, b.Bytes())
	if err != nil {
		fmt.Println(err)
	}
}
