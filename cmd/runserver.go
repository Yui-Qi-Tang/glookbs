package cmd

import (
	"glookbs.github.com/api/httphandler"
	"glookbs.github.com/httpserver"
	"glookbs.github.com/storage"
	_ "glookbs.github.com/storage/drivers/skiplists"

	"github.com/spf13/cobra"
)

type tlsfile struct {
	key, cert string
}

func (t tlsfile) Key() string {
	return t.key
}

func (t tlsfile) Cert() string {
	return t.cert
}

func runserver() *cobra.Command {
	var (
		addr        string
		apiMode     string
		pathTLSKey  string
		pathTLSCert string
	)

	cmd := &cobra.Command{
		Use:   "runserver",
		Short: "Run http server",
		Args:  cobra.MaximumNArgs(4),
		Run: func(c *cobra.Command, args []string) {
			srv := httpserver.New(
				httpserver.WithAddr(addr),
				httpserver.WithHandler(httphandler.New(apiMode, &storage.DataStorage)),
			)

			if len(pathTLSCert) > 0 && len(pathTLSKey) > 0 {
				tls := tlsfile{
					key:  pathTLSKey,
					cert: pathTLSCert,
				}
				if err := srv.Run(tls); err != nil {
					panic(err)
				}
			} else {
				if err := srv.Run(nil); err != nil {
					panic(err)
				}
			}
		},
	}

	cmd.Flags().StringVarP(&addr, "addr", "a", ":8080", "address of httpserver")
	cmd.Flags().StringVarP(&apiMode, "mode", "m", "debug", "mode of api")
	cmd.Flags().StringVarP(&pathTLSKey, "tls-key", "k", "", "path of tls key")
	cmd.Flags().StringVarP(&pathTLSCert, "tls-cert", "c", "", "path of tls cert")

	return cmd
}
