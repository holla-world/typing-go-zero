package new

import (
	_ "embed"
	"errors"
	"path/filepath"
	"strings"

	"github.com/holla-world/typing-go-zero/tools/goctl/api/gogen"
	conf "github.com/holla-world/typing-go-zero/tools/goctl/config"
	"github.com/holla-world/typing-go-zero/tools/goctl/util"
	"github.com/holla-world/typing-go-zero/tools/goctl/util/pathx"
	"github.com/spf13/cobra"
)

//go:embed api.tpl
var apiTemplate string

var (
	// VarStringHome describes the goctl home.
	VarStringHome string
	// VarStringRemote describes the remote git repository.
	VarStringRemote string
	// VarStringBranch describes the git branch.
	VarStringBranch string
	// VarStringStyle describes the style of output files.
	VarStringStyle string
)

// CreateServiceCommand fast create service
func CreateServiceCommand(_ *cobra.Command, args []string) error {
	dirName := args[0]
	if len(VarStringStyle) == 0 {
		VarStringStyle = conf.DefaultFormat
	}
	if strings.Contains(dirName, "-") {
		return errors.New("api new command service name not support strikethrough, because this will used by function name")
	}

	if len(VarStringRemote) > 0 {
		repo, _ := util.CloneIntoGitHome(VarStringRemote, VarStringBranch)
		if len(repo) > 0 {
			VarStringHome = repo
		}
	}

	if len(VarStringHome) > 0 {
		pathx.RegisterGoctlHome(VarStringHome)
	}

	abs, err := filepath.Abs(dirName)
	if err != nil {
		return err
	}

	err = pathx.MkdirIfNotExist(abs)
	if err != nil {
		return err
	}

	err = gogen.DoGenProject("zero", abs, VarStringStyle)
	return err
}
