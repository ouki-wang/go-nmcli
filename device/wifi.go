package device

import (
	//"regexp"
	"context"
	"fmt"
	"github.com/leberKleber/go-nmcli/utils"
	"strings"
)

type WiFiListOptions struct {
	IfName string
	BSSID  string
	Rescan WiFiListOptionsRescan
}

type WiFiListOptionsRescan string

var (
	WiFiListOptionsRescanAuto WiFiListOptionsRescan = "auto"
	WiFiListOptionsRescanYes  WiFiListOptionsRescan = "yes"
	WiFiListOptionsRescanNo   WiFiListOptionsRescan = "no"
)

func (a WiFiListOptions) rawArgs() []string {
	var args []string

	args = appendWhenNotEmpty(args, a.IfName, "ifname")
	args = appendWhenNotEmpty(args, a.BSSID, "bssid")
	args = appendWhenNotEmpty(args, string(a.Rescan), "--rescan")

	return args
}

type WiFi struct {
	Name      string
	SSID      string
	SSIDHEX   string
	BSSID     string
	Mode      string
	Chan      string
	Frequency string
	Rate      string
	Signal    string
	Bars      string
	Security  string
	WPAFlags  string
	RSNFlags  string
	Device    string
	Active    string
	InUse     string
	DBusPath  string
}

// WiFiList List available Wi-Fi access points.
// The IfName and BSSID options can be used to list APs for a particular interface, or with a specific BSSID.
// The Rescan flag tells whether a new Wi-Fi scan should be triggered.
func (m Manager) WiFiList(ctx context.Context, args WiFiListOptions) ([]WiFi, error) {
	fields := []string{"NAME", "SSID", "SSID-HEX", "BSSID", "MODE", "CHAN", "FREQ", "RATE", "SIGNAL", "BARS", "SECURITY", "WPA-FLAGS", "RSN-FLAGS", "DEVICE", "ACTIVE", "IN-USE", "DBUS-PATH"}

	cmdArgs := []string{"-g", strings.Join(fields, ",")}
	cmdArgs = append(cmdArgs, "device", "wifi", "list")
	cmdArgs = append(cmdArgs, args.rawArgs()...)

	output, err := m.CommandContext(ctx, nmcliCmd, cmdArgs...).Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute nmcli with args %+q: %w", cmdArgs, err)
	}

	parsedOutput, err := utils.ParseCmdOutput(output, len(fields))
	if err != nil {
		return nil, fmt.Errorf("failed to parse nmcli output: %w", err)
	}

	var wifis []WiFi
	for _, fields := range parsedOutput {
		wifis = append(wifis, WiFi{
			Name:      fields[0],
			SSID:      fields[1],
			SSIDHEX:   fields[2],
			BSSID:     fields[3],
			Mode:      fields[4],
			Chan:      fields[5],
			Frequency: fields[6],
			Rate:      fields[7],
			Signal:    fields[8],
			Bars:      fields[9],
			Security:  fields[10],
			WPAFlags:  fields[11],
			RSNFlags:  fields[12],
			Device:    fields[13],
			Active:    fields[14],
			InUse:     fields[15],
			DBusPath:  fields[16],
		})
	}

	return wifis, nil
}

type WiFiConnectOptions struct {
	Password   string
	WEPKeyType WiFiConnectOptionsWEPKeyType
	IfName     string
	BSSID      string
	Name       string
	Private    WiFiConnectOptionsPrivate
	Hidden     WiFiConnectOptionsHidden
}

type WiFiConnectOptionsWEPKeyType string
type WiFiConnectOptionsPrivate string
type WiFiConnectOptionsHidden string

var (
	WiFiConnectOptionsWEPKeyTypeKey    WiFiConnectOptionsWEPKeyType = "key"
	WiFiConnectOptionsWEPKeyTypePhrase WiFiConnectOptionsWEPKeyType = "phrase"
	WiFiConnectOptionsPrivateYes       WiFiConnectOptionsPrivate    = "yes"
	WiFiConnectOptionsPrivateNo        WiFiConnectOptionsPrivate    = "no"
	WiFiConnectOptionsHiddenYes        WiFiConnectOptionsHidden     = "yes"
	WiFiConnectOptionsHiddenNo         WiFiConnectOptionsHidden     = "no"
)

func (a WiFiConnectOptions) rawArgs() []string {
	var args []string

	args = appendWhenNotEmpty(args, a.Password, "password")
	args = appendWhenNotEmpty(args, string(a.WEPKeyType), "wep-key-type")
	args = appendWhenNotEmpty(args, a.IfName, "ifname")
	args = appendWhenNotEmpty(args, a.BSSID, "bssid")
	args = appendWhenNotEmpty(args, a.Name, "name")
	args = appendWhenNotEmpty(args, string(a.Private), "private")
	args = appendWhenNotEmpty(args, string(a.Hidden), "hidden")

	return args
}

// WiFiConnect Connect to a Wi-Fi network specified by BSSID which could also be a SSID.
// The command finds a matching connection or creates one and then activates it on a device.
// This is a command-line counterpart of clicking an SSID in a GUI client.
// If a connection for the network already exists, it is possible to bring up the existing profile as follows: nmcli con up id <name>.
// Note that only open, WEP and WPA-PSK networks are supported if no previous connection exists.
// It is also assumed that IP configuration is obtained via DHCP.
func (m Manager) WiFiConnect(ctx context.Context, BSSID string, args WiFiConnectOptions) (string, error) {
	cmdArgs := []string{"device", "wifi", "connect", BSSID}
	cmdArgs = append(cmdArgs, args.rawArgs()...)

	output, err := m.CommandContext(ctx, nmcliCmd, cmdArgs...).Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute nmcli with args %+q: %w", cmdArgs, err)
	}

	return string(output), nil
}

func (m Manager) WiFiDisConnect(ctx context.Context,iface string) (string, error) {
        cmdArgs := []string{"device", "disconnect", iface}

        output, err := m.CommandContext(ctx, nmcliCmd,cmdArgs...).Output()
        if err != nil {
                return "", fmt.Errorf("failed to execute nmcli with args %+q: %w", cmdArgs, err)
        }

        return string(output), nil
}

func (m Manager) IsConnected(ctx context.Context, iface string)([]string,error){
        cmdArgs := []string{"device", "status", iface}

        output, err := m.CommandContext(ctx, nmcliCmd,cmdArgs...).Output()
        if err != nil {
                return []string{""}, fmt.Errorf("failed to execute nmcli with args %+q: %w", cmdArgs, err)
        }
	out := string(string(output)[strings.Index(string(output),iface):])
	out = string(out[:strings.Index(out,"\n")])
	out1 := strings.Fields(out)
	//out1 := strings.Split(out," ")
	return out1[0:4], nil

}

func (m Manager) QueryClients(ctx context.Context, iface string)([]string,error){
        // iw dev p2p0 station dump | grep Station | wc -l
	cmdArgs := []string{"dev", iface,"station","dump"}

        output, err := m.CommandContext(ctx, "iw", cmdArgs...).Output()
        if err != nil {
                return []string{""}, fmt.Errorf("failed to execute nmcli with args %+q: %w", cmdArgs, err)
        }
	//re := regexp.MustCompile(`Station.*`)
	//matches := re.FindAllString(output, -1)
	//re := regexp.MustCompile(`^Station.*$`)
	//matches := re.FindAllString(string(output), -1)
	//return matches,nil
	start := 0
	l := 0
	clients := []string{}
	input := string(output)
	for{
		i := strings.Index(input[start:],"Station")
		if i==-1{
			break
		}
		start = start + i
		l = strings.Index(input[start:],"\n")
		//fmt.Println(input[start:start+l])
		clients = append(clients, input[start:start+l])
		start = start + l
	}
	return clients, nil
}

func appendWhenNotEmpty(slice []string, toCheck string, preAppend string) []string {
	if toCheck != "" {
		slice = append(slice, preAppend, toCheck)
	}

	return slice
}
