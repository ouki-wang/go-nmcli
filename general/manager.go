//go:generate go run github.com/golang/mock/mockgen -destination=utils_cmd_mock_test.go -package=general_test github.com/ouki-wang/go-nmcli/utils Cmd

package general

import (
	"context"

	"github.com/ouki-wang/go-nmcli/utils"
)

const nmcliCmd = "nmcli"

type Manager struct {
	CommandContext func(ctx context.Context, name string, args ...string) utils.Cmd
}
