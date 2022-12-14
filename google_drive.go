package main

import (
	"bytes"
	"context"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"os"
)

const (
	GoogleDriveScopeDrive     = "https://www.googleapis.com/auth/drive"
	GoogleDriveScopeDriveFile = "https://www.googleapis.com/auth/drive.file"

	MimeTypePNG     = "image/png"
	MimeTypeJPG     = "image/jpg"
	MimeTypeJPEG    = "image/jpeg"
	MimeTypeCSV     = "text/csv"
	MimeOctetStream = "application/octet-stream"
)

type GDrive struct {
	service *drive.Service
}

func NewGDriveService(scope ...string) (*GDrive, error) {
	b, err := os.ReadFile(os.Getenv("GoogleDriveServicePath"))
	if err != nil {
		return nil, fmt.Errorf("unable to read google drive client secret file : %s", err)
	}

	credentials, err := google.CredentialsFromJSON(context.Background(), b, scope...)
	if err != nil {
		return nil, fmt.Errorf("unable to acquire google drive credentials : %s", err)
	}

	service, err := drive.NewService(context.Background(), option.WithCredentials(credentials))
	if err != nil {
		return nil, err
	}

	return &GDrive{
		service: service,
	}, nil
}

func (d *GDrive) GetFileList() (*drive.FileList, error) {
	fileList, err := d.service.Files.
		List().
		PageSize(10).
		SupportsAllDrives(true).
		IncludeItemsFromAllDrives(true).
		Fields("files(id, name)").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to get file list : %s", err)
	}

	return fileList, nil
}

func (d *GDrive) DeleteFile(fileID string) error {
	err := d.service.Files.
		Delete(fileID).
		Do()
	if err != nil {
		return err
	}
	return nil
}

func (d *GDrive) CreateFile(metaData *drive.File, fileBuff []byte) (*drive.File, error) {
	file, err := d.service.Files.
		Create(metaData).
		Media(bytes.NewReader(fileBuff)).
		Do()
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (d *GDrive) UpdateFile(existingFile drive.File, fileBuff []byte) (*drive.File, error) {
	file, err := d.service.Files.
		Update(existingFile.Id, &existingFile).
		Media(bytes.NewReader(fileBuff)).
		Do()
	if err != nil {
		return nil, err
	}

	return file, nil
}
