package cmd

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
	"go.skymeyer.dev/app"
	"go.uber.org/zap"

	"go.reefassistant.com/apex-proxy/pkg/logger"
	"go.reefassistant.com/apex-proxy/pkg/proxy"
)

const (
	APEX_PROXY_DEV = "APEX_PROXY_DEV"
)

var (
	allowedIPs  = []string{"0.0.0.0/0"}
	allowedURLs = append(proxy.JSONEndpoints, proxy.XMLEndpoints...)
	apex        = "http://apex.local"
	bind        = ":8080"
	logLevel    = "info"
	logFile     = ""
	timeout     = 10 * time.Second
)

// New creates the root command.
func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "apex-proxy",
		Short:   fmt.Sprintf("Apex Proxy %s", app.Version),
		PreRunE: preRunApexProxy,
		RunE:    runApexProxy,
	}

	cmd.Flags().DurationVar(&timeout, "timeout", timeout, "Proxy read/write timeout")
	cmd.Flags().StringSliceVar(&allowedIPs, "allow-ip", allowedIPs, "List of IP ranges to allow access")
	cmd.Flags().StringSliceVar(&allowedURLs, "allow-url", allowedURLs, "List of Apex URLs to enable")
	cmd.Flags().StringVar(&apex, "apex", apex, "The Apex URL")
	cmd.Flags().StringVar(&bind, "bind", bind, "The socket to bind the server to (host:port)")
	cmd.Flags().StringVar(&logFile, "log-file", logFile, "Write logs to file")
	cmd.Flags().StringVar(&logLevel, "log-level", logLevel, fmt.Sprintf("Log level (%s|%s|%s|%s)",
		logger.LEVEL_ERROR, logger.LEVEL_WARNING, logger.LEVEL_INFO, logger.LEVEL_DEBUG))

	cmd.AddCommand(app.NewVersionCmd())

	return cmd
}

func preRunApexProxy(cmd *cobra.Command, args []string) error {
	var opts = []logger.LoggerOption{
		logger.WithLogLevel(logLevel),
		logger.WithLogFile(logFile),
	}
	// export APEX_PROXY_DEV=true for development
	if _, dev := os.LookupEnv(APEX_PROXY_DEV); dev {
		opts = append(opts, logger.WithDevelopment())
	}
	return logger.Initialize(opts...)
}

func runApexProxy(cmd *cobra.Command, args []string) error {

	mainLogger := zap.L().Named("apex-proxy")
	mainLogger.Sugar().Infow("startup",
		"timeout", timeout,
		"ips", allowedIPs,
		"urls", allowedURLs,
		"apex", apex,
		"bind", bind,
		"level", logLevel,
		"file", logFile,
	)

	var proxyOpts = []proxy.ProxyOption{
		proxy.WithAllowIPs(allowedIPs),
		proxy.WithAllowedURLs(allowedURLs),
	}

	apexProxy, err := proxy.New(apex, proxyOpts...)
	if err != nil {
		return err
	}

	loggerMiddleware := logger.ContextualHandler(
		logger.WithBaseLogger(mainLogger),
		logger.WithHeaders("user-agent"),
	)

	srv := &http.Server{
		Handler:      loggerMiddleware(apexProxy),
		Addr:         bind,
		WriteTimeout: timeout,
		ReadTimeout:  timeout,
	}

	return srv.ListenAndServe()
}
