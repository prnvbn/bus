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
	Use:          "bus",
	SilenceUsage: true,
	Short:        "Track arrival information for bus routes at a given stop",
	Long: `Track arrival information for bus routes at a given stop by running 'bus add'.
See the arrival times of the tracked buses by running 'bus arrivals'`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cfg = bus.Config{}
		yamlBytes, err := os.ReadFile(cfgPath)
		if os.IsNotExist(err) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("read config: %w", err)
		}
		if err := yaml.Unmarshal(yamlBytes, &cfg); err != nil {
			return fmt.Errorf("parse config %s: %w", cfgPath, err)
		}
		return cfg.ValidateSections()
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		yamlBytes, err := yaml.Marshal(cfg)
		if err != nil {
			return fmt.Errorf("marshal config: %w", err)
		}
		if err := os.MkdirAll(filepath.Dir(cfgPath), 0755); err != nil {
			return fmt.Errorf("create config directory: %w", err)
		}
		if err := os.WriteFile(cfgPath, yamlBytes, 0644); err != nil {
			return fmt.Errorf("write config: %w", err)
		}
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
