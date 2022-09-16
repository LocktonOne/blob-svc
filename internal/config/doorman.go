package config

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokene/doorman/connector"
)

type DoormanConfiger interface {
	DoormanConfig() *DoormanConfig
	DormanConnector() connector.ConnectorI
}

type DoormanConfig struct {
	ServiceUrl string `fig:"service_url,required"`
}

func NewDoormanConfiger(getter kv.Getter) DoormanConfiger {
	return &doormanConfig{
		getter: getter,
	}
}

type doormanConfig struct {
	getter kv.Getter
	once   comfig.Once
}

func (c *doormanConfig) DoormanConfig() *DoormanConfig {
	return c.once.Do(func() interface{} {
		raw := kv.MustGetStringMap(c.getter, "doorman")
		config := DoormanConfig{}
		err := figure.Out(&config).From(raw).Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out"))
		}

		return &config
	}).(*DoormanConfig)
}
func (c *doormanConfig) DormanConnector() connector.ConnectorI {
	return connector.NewConnectorMockKyc(c.DoormanConfig().ServiceUrl) //TODO remove mock
}
