package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gorm.io/gorm"

	"github.com/miladbarzideh/shortify/domain"
)

var cmdMigrate = func(log *logrus.Logger, postgresDb *gorm.DB) *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Migrate the database",
		Run: func(cmd *cobra.Command, args []string) {
			if err := postgresDb.AutoMigrate(&domain.URL{}); err != nil {
				log.Fatalf("failed to migrate database: %v", err)
			}
		},
	}
}
