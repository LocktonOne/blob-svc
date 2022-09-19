package config

import (
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/copus"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/kit/pgdb"
	doormanConfig "gitlab.com/tokene/doorman/connector/config"
)

type Config interface {
	comfig.Logger
	pgdb.Databaser
	types.Copuser
	comfig.Listenerer
	AWSConfiger
	doormanConfig.DoormanConfiger
}

type config struct {
	comfig.Logger
	pgdb.Databaser
	types.Copuser
	comfig.Listenerer
	getter kv.Getter
	AWSConfiger
	doormanConfig.DoormanConfiger
}

func New(getter kv.Getter) Config {
	return &config{
		getter:          getter,
		Databaser:       pgdb.NewDatabaser(getter),
		Copuser:         copus.NewCopuser(getter),
		Listenerer:      comfig.NewListenerer(getter),
		AWSConfiger:     NewAWSConfiger(getter),
		Logger:          comfig.NewLogger(getter, comfig.LoggerOpts{}),
		DoormanConfiger: doormanConfig.NewDoormanConfiger(getter),
	}
}
