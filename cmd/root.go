package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

var (
	app struct {
		name string
		cli  *cli.App
		run  bool
	}
	server appServer
	config configuration
	log    logger
)

func Execute() {
	if err := app.cli.Run(os.Args); err != nil {
		fmt.Printf("%s: error: %s\n", app.name, err)
		fmt.Printf("Type %s --help to see a list of all options.", app.name)
		os.Exit(1)
	}

	if app.run {
		server.start()
	}
}

func init() {
	if len(os.Args) > 0 {
		app.name = filepath.Base(os.Args[0])
	}

	flags := []cli.Flag{
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "host",
			Value:       "",
			Usage:       "Server `host`",
			Destination: &config.server.host,
		}),

		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        "enable-http",
			Usage:       "Enable HTTP server",
			Value:       false,
			Destination: &config.http.enabled,
		}),
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:        "port-http",
			Value:       80,
			Usage:       "HTTP `port`",
			Destination: &config.http.port,
		}),

		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        "enable-https",
			Usage:       "Enable HTTPS server",
			Value:       false,
			Destination: &config.https.enabled,
		}),
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:        "port-https",
			Value:       443,
			Usage:       "HTTPS `port`",
			Destination: &config.https.port,
		}),

		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "crt-file",
			Usage:       "Location of the SSL certificate `file`",
			Value:       "",
			Destination: &config.https.cert,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "key-file",
			Usage:       "Location of the RSA private key `file`",
			Value:       "",
			Destination: &config.https.key,
		}),

		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "content",
			Value:       "ok",
			Usage:       "Response body",
			Destination: &config.content.content,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "content-file",
			Value:       "",
			Usage:       "Response body from `file`",
			Destination: &config.content.file,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "content-type",
			Value:       "text/plain; charset=UTF-8",
			Usage:       "Content-Type header",
			Destination: &config.content.contentType,
		}),

		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        "enable-tcp",
			Usage:       "Enable TCP echo server",
			Value:       false,
			Destination: &config.tcp.enabled,
		}),
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:        "port-tcp",
			Usage:       "TCP echo `port`",
			Value:       0,
			Destination: &config.tcp.port,
			DefaultText: "random",
		}),

		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        "enable-udp",
			Usage:       "Enable UDP echo server",
			Value:       false,
			Destination: &config.udp.enabled,
		}),
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:        "port-udp",
			Usage:       "UDP echo `port`",
			Value:       0,
			Destination: &config.udp.port,
			DefaultText: "random",
		}),

		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        "enable-log",
			Usage:       "Enable file logging",
			Value:       false,
			Destination: &config.log.enabled,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "log-dir",
			Value:       "log",
			Usage:       "Location of the log directory",
			Destination: &config.log.dir,
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        "log-requests",
			Usage:       "Log HTTP(S) requests",
			Value:       true,
			Destination: &config.log.requests,
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        "log-connections",
			Usage:       "Log TCP connections",
			Value:       true,
			Destination: &config.log.connections,
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        "log-packets",
			Usage:       "Log TCP/UDP echo packets",
			Value:       true,
			Destination: &config.log.packets,
		}),

		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        "quiet",
			Aliases:     []string{"q"},
			Usage:       "Activate quiet mode",
			Value:       false,
			Destination: &config.quiet,
		}),

		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Value:   "",
			Usage:   "Location of the configuration `file` in .yml format",
		},
	}

	cli.HelpFlag = &quietBoolFlag{
		BoolFlag: cli.BoolFlag{
			Name:    "help",
			Aliases: []string{"h"},
			Usage:   "Print this help text and exit",
		},
	}

	cli.VersionFlag = &quietBoolFlag{
		BoolFlag: cli.BoolFlag{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "Print program version and exit",
		},
	}

	app.cli = &cli.App{
		Name:                  "Echo Server",
		Version:               "v1.0.0",
		Compiled:              time.Now(),
		UsageText:             fmt.Sprintf("%s [global options]", app.name),
		HelpName:              app.name,
		HideHelpCommand:       true,
		Before:                altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("config")),
		Flags:                 flags,
		CustomAppHelpTemplate: helpTemplate,
		Action: func(cCtx *cli.Context) error {
			if !cCtx.Bool("help") && !cCtx.Bool("version") {
				if err := config.init(); err != nil {
					return err
				}

				if err := log.init(); err != nil {
					return err
				}

				app.run = true
			}

			return nil
		},
	}
}

type quietBoolFlag struct {
	cli.BoolFlag
}

func (q *quietBoolFlag) String() string {
	return cli.FlagStringer(q)
}

func (q *quietBoolFlag) GetDefaultText() string {
	return ""
}
