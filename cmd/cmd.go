package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "Hugo is a very fast static site generator",
	Long: `A Fast and Flexible Static Site Generator built with
	love by spf13 and friends in Go.
	Complete documentation is available at http://hugo.spf13.com`,
	Run: func(cmd *cobra.Command, args []string) {
		runServer()
	},
}

// dbCmd is the parent command for db operations
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database operations",
	Long:  "Commands for managing the database (up, down, seed).",
}

// upCmd - db up
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Migrate the database up",
	Long:  "Run database migrations to upgrade to the latest version.",
	Run: func(cmd *cobra.Command, args []string) {
		dbUp()
	},
}

// upTableCmd - db up view
var upTableCmd = &cobra.Command{
	Use:   "table",
	Short: "Migrate the database up table",
	Long:  "Run database migrations to upgrade to the latest version.",
	Run: func(cmd *cobra.Command, args []string) {
		dbUpTable()
	},
}

// upViewCmd - db up view
var upViewCmd = &cobra.Command{
	Use:   "view",
	Short: "Migrate the database up view",
	Long:  "Run database migrations to upgrade to the latest version.",
	Run: func(cmd *cobra.Command, args []string) {
		dbUpView()
	},
}

// downCmd - db down
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Migrate the database down",
	Long:  "Revert the database to the previous version by undoing migrations.",
	Run: func(cmd *cobra.Command, args []string) {
		dbDown()
	},
}

// seedCmd - db seed
var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database",
	Long:  "Seed the database with initial data for testing or development.",
	Run: func(cmd *cobra.Command, args []string) {
		dbSeed()
	},
}

// seedCmd - db reset
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset the database",
	Long:  "Reset the database run down then up and seed.",
	Run: func(cmd *cobra.Command, args []string) {
		dbReset()
	},
}

func init() {
	rootCmd.AddCommand(dbCmd)

	dbCmd.AddCommand(upCmd)
	dbCmd.AddCommand(downCmd)
	dbCmd.AddCommand(seedCmd)
	dbCmd.AddCommand(resetCmd)

	upCmd.AddCommand(upViewCmd)
	upCmd.AddCommand(upTableCmd)

}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
