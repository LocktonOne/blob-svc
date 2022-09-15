package helpers

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func NewAwsSession(r *http.Request) *session.Session {
	awsCfg := AwsConfig(r)
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(awsCfg.Region),
		Credentials: credentials.NewStaticCredentials(awsCfg.AccessKeyID, awsCfg.SecretKeyID, ""),
		DisableSSL:  &awsCfg.SslDisable,
	}))
	return sess
}
func DeleteItem(sess *session.Session, bucket *string, item *string) error {
	svc := s3.New(sess)

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: bucket,
		Key:    item,
	})
	if err != nil {
		return err
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: bucket,
		Key:    item,
	})
	if err != nil {
		return err
	}

	return nil
}
