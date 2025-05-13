package docgen

import (
	"fmt"
	"github.com/spf13/cobra"
)

const (
	DefaultApifoxToken = "APS-2tzNyIWDZYJuaFI2gGGd8HIkmKTlk2iL"
)

var (
	ApiFolder         string
	AdminFolder       string
	InternalFolder    string
	ApiProjectId      string
	AdminProjectId    string
	InternalProjectId string
	ApifoxToken       string
)

func GenDoc(_ *cobra.Command, _ []string) error {
	if len(ApifoxToken) == 0 {
		ApifoxToken = DefaultApifoxToken
	}
	// api配置
	configs := make([]APIDocConfig, 0)
	if len(ApiFolder) > 0 {
		configs = append(configs, APIDocConfig{
			FileName:        "api",
			APIFolder:       ApiFolder,
			APIPath:         fmt.Sprintf("%s/api.api", ApiFolder),
			ApifoxProjectID: ApiProjectId,
		})
	}
	if len(AdminFolder) > 0 {
		configs = append(configs, APIDocConfig{
			FileName:        "admin",
			APIFolder:       AdminFolder,
			APIPath:         fmt.Sprintf("%s/admin.api", AdminFolder),
			ApifoxProjectID: AdminProjectId,
		})
	}

	if len(InternalFolder) > 0 {
		configs = append(configs, APIDocConfig{
			FileName:        "internal",
			APIFolder:       InternalFolder,
			APIPath:         fmt.Sprintf("%s/internal.api", InternalFolder),
			ApifoxProjectID: InternalProjectId,
		})
	}

	err := GenerateAndImportSwagger(configs, ApifoxToken)
	if err != nil {
		fmt.Println("Error:", err)
	}
	return nil
}
