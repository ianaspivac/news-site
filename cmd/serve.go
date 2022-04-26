package cmd

import (
	"context"
	"fmt"
	"github.com/ianaspivac/news-site-go/internal/server"
	"github.com/ianaspivac/news-site-go/internal/service"
	"github.com/ianaspivac/news-site-go/internal/store"
	"github.com/ianaspivac/news-site-go/internal/util"

	"github.com/spf13/cobra"
)

var port string

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long:  `bruh`,
	Run: func(cmd *cobra.Command, args []string) {
		if port == "" {
			fmt.Println("Cannot use empty port")
			return
		}

		_, cancel := context.WithCancel(context.Background())

		db := store.CreateDB()

		storeServer := server.New(
			util.NewJWT(),
			service.New(db),
		)

		if err := storeServer.Run(fmt.Sprintf(":%s", port)); err != nil {
			fmt.Printf("Unexpected error while running server: %v", err)
		}

		cancel()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.PersistentFlags().StringVarP(&port, "port", "p", "8080", "Defines port for memstore server. Default value is 8080")

}
