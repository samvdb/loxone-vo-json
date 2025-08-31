/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/samvdb/loxone-vo-json/proxy"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var (
	flagServer string
	flagPort   int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	RunE: func(cmd *cobra.Command, args []string) error {
		if flagServer == "" {
			return fmt.Errorf("--server is required (destination EMS host or URL)")
		}
		if flagPort <= 0 || flagPort > 65535 {
			return fmt.Errorf("--port must be a valid TCP port")
		}

		target, err := proxy.ParseTarget(flagServer)
		if err != nil {
			return fmt.Errorf("invalid --server value: %w", err)
		}

		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
		slog.SetDefault(logger)

		p := proxy.NewProxy(target)
		handler := proxy.LoggingMiddleware(p)

		addr := fmt.Sprintf(":%d", flagPort)
		slog.Info("starting proxy", "listen", addr, "upstream", target.String())

		srv := &http.Server{
			Addr:              addr,
			Handler:           handler,
			ReadHeaderTimeout: 15 * time.Second,
		}
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("server error: %w", err)
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.loxone-vo-json.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringVar(&flagServer, "server", "", "Proxy destination address (host or URL)")
	rootCmd.Flags().IntVar(&flagPort, "port", 8080, "HTTP Proxy listen port")
	// Read from environment if flag not set
	if envServer := os.Getenv("SERVER"); envServer != "" {
		flagServer = envServer
	}
	if envPort := os.Getenv("PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			flagPort = p
		}
	}
	_ = rootCmd.MarkFlagRequired("server")
}
