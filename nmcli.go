package go_nmcli

import (
	"context"
	"github.com/ouki-wang/go-nmcli/device"
	"os/exec"

	"github.com/ouki-wang/go-nmcli/general"
	"github.com/ouki-wang/go-nmcli/utils"
	"github.com/sirupsen/logrus"
)

type General interface {
	Hostname(ctx context.Context, args general.HostnameArgs) (string, error)
	Permissions(ctx context.Context) ([]general.Permission, error)
}

type Device interface {
	WiFiList(ctx context.Context, args device.WiFiListOptions) ([]device.WiFi, error)
	WiFiConnect(ctx context.Context, BSSID string, args device.WiFiConnectOptions) (string, error)
	WiFiDisConnect(ctx context.Context,iface string) (string, error)
	IsConnected(ctx context.Context,iface string) ([]string, error)
	QueryClients(ctx context.Context,iface string)([]string,error)
}

type NMCli struct {
	// should be used to exec custom nmcli commands
	CommandContext func(ctx context.Context, name string, args ...string) utils.Cmd
	logDebug       func(fmt string, args ...interface{})
	General        General
	Device         Device
}

type Option = func(cli *NMCli)

func NewNMCli(opts ...Option) NMCli {
	cli := NMCli{
		logDebug: logrus.Debugf,
		CommandContext: func(ctx context.Context, name string, args ...string) utils.Cmd {
			return exec.CommandContext(ctx, name, args...)
		},
	}
	for i := range opts {
		opts[i](&cli)
	}

	cli.General = general.Manager{CommandContext: cli.CommandContext}
	cli.Device = device.Manager{CommandContext: cli.CommandContext}

	return cli
}
