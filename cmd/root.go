package cmd

import (
	"fmt"
	"os"

	"crypto/rand"

	"encoding/hex"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/workflow-intelligence/gather/search"
	"github.com/workflow-intelligence/gather/server"
	"math"
)

var (
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "gather",
	Short: "Application to gather github workflow data",
	Long: `This application waits for a github workflow to contact it and then fetches workflow metrics
from github about it.`,
	Run: func(cmd *cobra.Command, args []string) {
		opensearch_client, err := search.New(viper.GetString("opensearch_user"), viper.GetString("opensearch_password"), []string{viper.GetString("opensearch_url")})
		if err != nil {
			fmt.Println("Could not connect to Opensearch")
			os.Exit(1)
		}
		server.Server(viper.GetBool("ssl"), viper.GetInt("port"), viper.GetString("logfile"), viper.GetString("loglevel"), viper.GetString("jwtsecretkey"), opensearch_client)
	},
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "gather.yaml", "Config file")
	rootCmd.PersistentFlags().BoolP("ssl", "s", true, "Whether to use SSL or not")
	rootCmd.PersistentFlags().IntP("port", "p", 1443, "Port to run the server on")
	rootCmd.PersistentFlags().StringP("logfile", "l", "", "Where to log, defaults to stdout")
	rootCmd.PersistentFlags().StringP("loglevel", "L", "info", "Loglevel can be one of panic, fatal, error, warn, info, debug, trace")
	rootCmd.PersistentFlags().StringP("jwtsecretkey", "j", "", "JWT secret key for token generation")
	rootCmd.PersistentFlags().StringP("opensearch_url", "u", "https://localhost:9200", "URL for the opensearch server")
	rootCmd.PersistentFlags().StringP("opensearch_user", "U", "gather", "Username for storing data in opensearch")
	rootCmd.PersistentFlags().StringP("opensearch_password", "P", "", "Password for the opensearch user")
	viper.BindPFlag("ssl", rootCmd.PersistentFlags().Lookup("ssl"))
	viper.SetDefault("ssl", true)
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.SetDefault("port", 1443)
	viper.BindPFlag("logfile", rootCmd.PersistentFlags().Lookup("logfile"))
	viper.SetDefault("logfile", "")
	viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))
	viper.SetDefault("loglevel", "info")
	viper.BindPFlag("jwtsecretkey", rootCmd.PersistentFlags().Lookup("jwtsecretkey"))
	viper.SetDefault("jwtsecretkey", randomBase16String(32))
	viper.BindPFlag("opensearch_url", rootCmd.PersistentFlags().Lookup("opensearch_url"))
	viper.SetDefault("opensearch_url", "https://localhost:9200")
	viper.BindPFlag("opensearch_user", rootCmd.PersistentFlags().Lookup("opensearch_user"))
	viper.SetDefault("opensearch_user", "gather")
	viper.BindPFlag("opensearch_password", rootCmd.PersistentFlags().Lookup("opensearch_password"))
	viper.SetDefault("opensearch_password", "")
}

func initConfig() {
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file ", viper.ConfigFileUsed())
		os.Exit(1)
	}
}

func randomBase16String(l int) string {
	buff := make([]byte, int(math.Ceil(float64(l)/2)))
	rand.Read(buff)
	str := hex.EncodeToString(buff)
	return str[:l] // strip 1 extra character we get from odd length results
}
