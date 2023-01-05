package helpers

import (
	"errors"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"gitlab.com/tokene/blob-svc/internal/config"
)

var AllowedFileExtensions = []string{"png", "jpg", "jpeg", "bmp"}

func NewAwsSession(r *http.Request) *session.Session {
	awsCfg := AwsConfig(r)
	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint:         &awsCfg.Endpoint,
		S3ForcePathStyle: &awsCfg.ForcePathStyle,
		Region:           aws.String(awsCfg.Region),
		Credentials:      credentials.NewStaticCredentials(awsCfg.AccessKeyID, awsCfg.SecretKeyID, ""),
		DisableSSL:       &awsCfg.SslDisable,
	}))
	return sess
}

func GetItemURL(sess *session.Session, itemName string, cfg config.AWSConfig) (string, error) {
	svc := s3.New(sess)
	getObjectReq, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(cfg.Bucket),
		Key:    aws.String(itemName),
	})
	return getObjectReq.Presign(cfg.Expiration)
}

func DeleteItem(sess *session.Session, item *string, cfg config.AWSConfig) error {
	svc := s3.New(sess)

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: &cfg.Bucket,
		Key:    item,
	})
	if err != nil {
		return err
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: &cfg.Bucket,
		Key:    item,
	})
	if err != nil {
		return err
	}

	return nil
}
func CheckFileExtension(ext string) error {
	for _, el := range AllowedFileExtensions {
		if el == ext {
			return nil
		}
	}
	return errors.New("invalid file extension")
}
