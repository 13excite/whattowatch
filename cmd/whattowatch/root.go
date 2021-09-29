package cmd

import (
	"fmt"
	"os"

	_ "net/http/pprof" // Import for pprof

	cli "github.com/spf13/cobra"
	"go.uber.org/zap"

	"ww/internal/conf"
)

var (

	// Config and global logger
	pidFile string
	logger  *zap.SugaredLogger

	// The Root Cli Handler
	rootCmd = &cli.Command{
		Version: conf.GitVersion,
		Use:     conf.Executable,
		PersistentPreRunE: func(cmd *cli.Command, args []string) error {
			// Create Pid File
			pidFile = conf.C.PidFile
			if pidFile != "" {
				file, err := os.OpenFile(pidFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
				if err != nil {
					return fmt.Errorf("Could not create pid file: %s Error:%v", pidFile, err)
				}
				defer file.Close()
				_, err = fmt.Fprintf(file, "%d\n", os.Getpid())
				if err != nil {
					return fmt.Errorf("Could not create pid file: %s Error:%v", pidFile, err)
				}
			}
			return nil
		},
		PersistentPostRun: func(cmd *cli.Command, args []string) {
			// Remove Pid file
			if pidFile != "" {
				os.Remove(pidFile)
			}
		},
	}
)

// Execute starts the program
func Execute() {

	// Load configuration
	conf.C.Defaults()
	configFile := rootCmd.PersistentFlags().StringP("config", "c", "", "config file")
	if configFile != nil && *configFile != "" {
		conf.C.ReadConfigFile(*configFile)
		fmt.Println("AAAAA")
	}
	conf.C.ReadConfigFile("./config.yaml")
	fmt.Println(rootCmd.PersistentFlags())
	fmt.Println(apiCmd.PersistentFlags().FlagUsages())
	fmt.Println(*configFile)
	fmt.Println("BBBBB")
	conf.InitLogger(&conf.C)

	logger = zap.S().With("package", "cmd")

	// Run the program
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}
