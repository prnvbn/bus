package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/prnvbn/bus/internal/bus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	cfgPath string
	cfg     bus.Config
)

func init() {
	var ok bool
	cfgPath, ok = os.LookupEnv("BUS_CONFIG_PATH")
	if !ok {
		cfgPath = xdg.ConfigHome + "/bus/config.yaml"
	}
}

var rootCmd = &cobra.Command{
	Use: "bq",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		yamlBytes, err := os.ReadFile(cfgPath)
		if err != nil {
			if os.IsNotExist(err) {
				err = os.MkdirAll(filepath.Dir(cfgPath), 0755)
				fatal(err, "error creating config directory at %s", filepath.Dir(cfgPath))

				f, err := os.Create(cfgPath)
				fatal(err, "error creating config file at %s", cfgPath)
				defer f.Close()

				yamlBytes, err := yaml.Marshal(cfg)
				fatal(err, "error marshaling config")

				_, err = f.Write(yamlBytes)
				fatal(err, "error writing config file at %s", cfgPath)
			}
			return err
		}
		fatal(err, "Config file %s does not exist!", cfgPath)

		err = yaml.Unmarshal(yamlBytes, &cfg)
		fatal(err, "YAML Config file at %s is malformed", cfgPath)

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		yamlBytes, err := yaml.Marshal(cfg)
		fatal(err, "error marshaling config")

		err = os.WriteFile(cfgPath, yamlBytes, 0644)
		fatal(err, "error writing config file at %s", cfgPath)

		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

}

// TODO: use something better than panic
func fatal(err error, message string, args ...interface{}) {
	if err != nil {
		panic(fmt.Errorf(message, args...))
	}
}
