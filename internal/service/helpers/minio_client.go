package helpers

import (
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinioClient() (*minio.Client, error) { //add config
	//config
	endpoint := "s3.amazonaws.com"
	accessKeyID := "AKIA3M562ZQ54BSD3V74"
	secretAccessKey := "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"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""), //config creds
		Secure: useSSL,                                                    //config
	})
	if err != nil {
		log.Fatalln(err)
		return minioClient, err
	}

	return minioClient, nil
}
