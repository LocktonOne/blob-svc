package config

import (
	"time"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type AWSConfiger interface {
	AWSConfig() *AWSConfig
}

type AWSConfig struct {
	AccessKeyID    string        `fig:"access_key,required"`
	SecretKeyID    string        `fig:"secret_key,required"`
	Bucket         string        `fig:"bucket,required"`
	Expiration     time.Duration `fig:"expiration,required"`
	SslDisable     bool          `fig:"ssldisable,required"`
	Region         string        `fig:"region,required"`
	Endpoint       string        `fig:"endpoint,required"`
	ForcePathStyle bool          `fig:"force_path_style,required"`
}

func NewAWSConfiger(getter kv.Getter) AWSConfiger {
	return &awsConfig{
		getter: getter,
	}
}

type awsConfig struct {
	getter kv.Getter
	once   comfig.Once
}

func (c *awsConfig) AWSConfig() *AWSConfig {
	return c.once.Do(func() interface{} {
		raw := kv.MustGetStringMap(c.getter, "aws")
		config := AWSConfig{}
		err := figure.Out(&config).From(raw).Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out"))
		}

		return &config
	}).(*AWSConfig)
}
