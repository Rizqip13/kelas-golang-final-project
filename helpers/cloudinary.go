package helpers

import (
	"bytes"
	"context"
	"final-project/configs"
	"io"
	"mime/multipart"
	"path"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadFile(fileHeader *multipart.FileHeader) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cld, err := cloudinary.NewFromURL(configs.EnvCloudURL())
	if err != nil {
		return "", err
	}

	// convert file
	fileReader, err := convertFile(fileHeader)
	if err != nil {
		return "", err
	}

	fileName := RemoveExtension(fileHeader.Filename)

	// upload file
	uploadParam, err := cld.Upload.Upload(ctx, fileReader, uploader.UploadParams{
		PublicID: fileName,
		Folder:   configs.EnvCloudUploadFolder(),
	})
	if err != nil {
		return "", err
	}

	return uploadParam.SecureURL, nil
}

func convertFile(fileHeader *multipart.FileHeader) (*bytes.Reader, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the file content into an in-memory buffer
	buffer := new(bytes.Buffer)
	if _, err := io.Copy(buffer, file); err != nil {
		return nil, err
	}

	// Create a bytes.Reader from the buffer
	fileReader := bytes.NewReader(buffer.Bytes())
	return fileReader, nil
}

func RemoveExtension(filename string) string {
	return path.Base(filename[:len(filename)-len(path.Ext(filename))])
}
