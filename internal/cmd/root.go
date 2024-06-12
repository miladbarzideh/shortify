package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	infra2 "github.com/miladbarzideh/shortify/internal/infra"
)

var rooCmd = &cobra.Command{
	Use:   "shortify",
	Short: "Simple URL shortener",
}

func Execute() {
	cfg, err := infra2.Load()
	if err != nil {
		logrus.Fatal(err)
	}

	log := infra2.InitLogger(cfg)
	postgresDb, err := infra2.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatal("database connection failed")
	}

	redis, err := infra2.NewRedisClient(cfg)
	if err != nil {
		log.Fatal("redis client failed")
	}

	cmdServe := cmdServer(cfg, log, postgresDb, redis)
	cmdServe.Flags().IntP("port", "p", 8080,
		"Optional port number.Default value will be read from the config file")
	rooCmd.AddCommand(cmdServe)
	rooCmd.AddCommand(cmdMigrate(log, postgresDb))
	if err = rooCmd.Execute(); err != nil {
		log.Fatalf("failed to execute root command %s", err)
	}
}
