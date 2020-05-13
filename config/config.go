package config

import (
	"flag"
	"errors"

	"github.com/BurntSushi/toml"
)

type ProducerConfig struct {
	fs                         *flag.FlagSet
	Topic                      string      `toml:"topic" json:"topic"`
	Group                      string      `toml:"group" json:"group"`
	BusinessAppkey             string      `toml:"business-appkey" json:"business-appkey"`
	CastleAppkey               string      `toml:"castle-appkey" json:"castle-appkey"`
	CastleName                 string      `toml:"castle-name" json:"castle-name"`
	StallThreshold             int         `toml:"stall-threshold" json:"stall-threshold"`
	WaitThreshold              int64       `toml:"wait-threshold" json:"wait-threshold"`
	MaxRetryTimes              int         `toml:"max-retry-times" json:"max-retry-times"`
	MaxAsyncBufferChanSize     int64       `toml:"max-chan-size" json:"max-chan-size"`
	SafeMode                   int         `toml:"safe-mode" json:"safe-mode"`

	LocalHost                  string
}

func NewProducerConfig() *ProducerConfig {
	cfg := &ProducerConfig{}
	fs := flag.NewFlagSet("syncer", flag.ContinueOnError)
	fs.StringVar(&cfg.Topic, "topic", "", "mafka 's topic")
	fs.StringVar(&cfg.Group, "group", "", "mafka 's group")
	fs.StringVar(&cfg.BusinessAppkey, "business-appkey", "", "business appkey")
	fs.StringVar(&cfg.CastleAppkey, "castle-appkey", "", "castle appkey")
	fs.StringVar(&cfg.CastleName, "castle-name", "", "castle name")
	fs.IntVar(&cfg.StallThreshold, "stall-threshold", 90 * 1024 * 1024, "stall threshold")
	fs.Int64Var(&cfg.WaitThreshold, "wait-threshold", 30000, "wait threshold (ms)")
	fs.IntVar(&cfg.MaxRetryTimes, "max-retry-times", 10000, "max retry times")
	fs.Int64Var(&cfg.MaxAsyncBufferChanSize, "max-chan-size", 1 << 30, "max async buffer chan size")
	fs.IntVar(&cfg.SafeMode, "safe-mode", 0, "enable safe mode to make syncer reentrant")

	cfg.fs = fs
	return cfg
}

func (c *ProducerConfig) StrictDecodeFile(path string) error {
	metaData, err := toml.DecodeFile(path, c)
	if err != nil {
		return err
	}

	if undecoded := metaData.Undecoded(); len(undecoded) > 0 {
		var undecodedItems []string
		for _, item := range undecoded {
			undecodedItems = append(undecodedItems, item.String())
		}
		err = errors.New("contained unknown configuration options")
	}

	return err
}

func (c *ProducerConfig) Parse(filename string) error {
	// load config file if specified
	if filename == "" {
		return errors.New("config file name is nil")
	}
	if err := c.StrictDecodeFile(filename); err != nil {
		return err
	}
	return nil
}
