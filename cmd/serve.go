package cmd

import (
	"fmt"
	"net"
	"net/http"

	"github.com/scryner/wdav/dir"
	"github.com/spf13/cobra"
)

const (
	defaultListenPort = 8080
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serving webdav server for given directories",
	RunE: func(cmd *cobra.Command, args []string) error {
		// get port
		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			port = defaultListenPort
		}

		// get directories
		dirs, err := cmd.Flags().GetStringArray("dir")
		if err != nil {
			return fmt.Errorf("failed to get directories from command-line arguments: %v", err)
		}

		if len(dirs) < 1 {
			return fmt.Errorf("directories are must specified")
		}

		handler, err := dir.NewHandler("webdav", dirs)
		if err != nil {
			return fmt.Errorf("failed to make handler: %v", err)
		}

		// listen
		l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
		if err != nil {
			return fmt.Errorf("failed to listen at '%d': %v", port, err)
		}

		// serve http
		fmt.Printf("Starting to listen to serve WebDAV at '0.0.0.0:%d'\n", port)
		return http.Serve(l, handler)
	},
}

func init() {
	serveCmd.Flags().IntP("port", "p", defaultListenPort, "Specify listen port")
	serveCmd.Flags().StringArrayP("dir", "d", nil, "Specify directories")
}
