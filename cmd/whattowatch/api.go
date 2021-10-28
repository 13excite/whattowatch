package cmd

import (
	cli "github.com/spf13/cobra"
	"go.uber.org/zap"

	"ww/internal/api"
	"ww/internal/conf"
	"ww/internal/store/postgres"
)

func init() {
	rootCmd.AddCommand(apiCmd)
}

var (
	apiCmd = &cli.Command{
		Use:   "api",
		Short: "Start API",
		Long:  `Start API`,
		Run: func(cmd *cli.Command, args []string) { // Initialize the database
			conf.C.Defaults()
			if &configFile != nil && configFile != "" {
				conf.C.ReadConfigFile(configFile)
			}
			conf.InitLogger(&conf.C)

			logger = zap.S().With("package", "cmd")

			// Database
			pg, err := postgres.New(&conf.C)
			if err != nil {
				logger.Fatalw("Database error", "error", err)
			}

			// Create the server
			s, err := api.New(&conf.C, pg)
			if err != nil {
				logger.Fatalw("Could not create server", "error", err)
			}

			s.Router().HandleFunc("/random", s.GetRandomFilm).Methods("GET")
			s.Router().HandleFunc("/version", conf.GetVersion()).Methods("GET")
			if err = s.ListenAndServe(&conf.C); err != nil {
				logger.Fatalw("Could not start server", "error", err)
			}

			conf.Stop.InitInterrupt()
			<-conf.Stop.Chan() // Wait until Stop
			conf.Stop.Wait()   // Wait until everyone cleans up
			_ = zap.L().Sync() // Flush the logger

		},
	}
)
