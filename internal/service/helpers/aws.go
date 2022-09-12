package helpers

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
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
