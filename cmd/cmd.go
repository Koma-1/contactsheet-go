package cmd

import (
	"github.com/Koma-1/contactsheet-go/contactsheet"
)

func run(userConfig userConfig) error {
	if err := userConfig.validate(); err != nil {
		return err
	}

	config, err := convertToFixedTileConfig(userConfig)
	if err != nil {
		return err
	}
	generator, err := contactsheet.NewGenerator(config)
	if err != nil {
		return err
	}
	err = generator.GenerateFromDir()
	if err != nil {
		return err
	}
	return nil
}
