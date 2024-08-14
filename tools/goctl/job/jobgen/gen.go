package jobgen

import (
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/golang"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

const (
	defaultCronFile   = "job/etc/cron.yaml"
	defaultDaemonFile = "job/etc/daemon.yaml"
	defaultDir        = "job"
	handlerShuffix    = "Handler"
	internal          = "internal/"
	handlerDir        = internal + "handler"
	contextDir        = internal + "svc"
)

const (
	routesFilename = "routes"
)

const (
	Zero = int32(0)
)

const (
	JobCron   = "cron"   // 定时任务
	JobDaemon = "daemon" // 常驻任务
)

var (
	CronFile   string
	DaemonFile string
)

type jobConfig struct {
	JobType      string
	Action       string `yaml:"Action"`       // 执行命令方法
	Schedule     string `yaml:"Schedule"`     // cron定时参数
	Replicas     *int32 `yaml:"Replicas"`     // 副本数
	TestReplicas *int32 `yaml:"TestReplicas"` // 测试环境副本数-不设置取Replicas
	Desc         string `yaml:"Desc"`         // 命令描述
	Handler      string `yaml:"Handler"`      // Handler方法名
}

// fileExists checks if a file exists and is not a directory
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func GenJob(_ *cobra.Command, _ []string) error {
	if len(CronFile) == 0 {
		CronFile = defaultCronFile
	}
	if len(DaemonFile) == 0 {
		DaemonFile = defaultDaemonFile
	}
	if !fileExists(CronFile) {
		panic(fmt.Sprintf("cronfile:%s not found", CronFile))
	}
	if !fileExists(DaemonFile) {
		panic(fmt.Sprintf("daemonfile:%s not found", CronFile))
	}
	err := genJobCode(CronFile, DaemonFile)
	if err != nil {
		panic(err)
	}

	return nil
}

func genJobCode(cronFile string, daemonFile string) error {
	// 解析yaml配置文件
	var cfgCron map[string][]jobConfig
	var cfgDaemon map[string][]jobConfig

	// 读取配置文件
	cron, err := os.ReadFile(cronFile)
	if err != nil {
		panic(fmt.Sprintf("Failed to read YAML file: %v", err))
	}
	daemon, err := os.ReadFile(daemonFile)
	if err != nil {
		panic(fmt.Sprintf("Failed to read YAML file: %v", err))
	}

	// 解析 YAML 文件
	err = yaml.Unmarshal(cron, &cfgCron)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse YAML file: %v", err))
	}
	err = yaml.Unmarshal(daemon, &cfgDaemon)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse YAML file: %v", err))
	}

	// 合并配置
	mgrCfg := make(map[string][]jobConfig)
	for group, routes := range cfgCron {
		for i, route := range routes {
			route.JobType = JobCron
			routes[i] = route
		}
		cfgCron[group] = routes
		mgrCfg[group] = routes
	}
	for group, routes := range cfgDaemon {
		for i, route := range routes {
			route.JobType = JobDaemon
			if route.Replicas == nil {
				zero := Zero
				route.Replicas = &zero
			}
			if route.TestReplicas == nil {
				route.TestReplicas = route.Replicas
			}
			routes[i] = route
		}
		cfgDaemon[group] = routes
		if _, ok := mgrCfg[group]; ok {
			mgrCfg[group] = append(mgrCfg[group], routes...)
		} else {
			mgrCfg[group] = routes
		}
	}

	actionMap := make(map[string]struct{})
	for group, routes := range mgrCfg {
		// handler和action重复性检查
		handlerMap := make(map[string]struct{})
		for _, route := range routes {
			if _, ok := handlerMap[route.Handler]; !ok {
				handlerMap[route.Handler] = struct{}{}
			} else {
				panic(fmt.Sprintf("%s 中存在重复的handler命名: [%s], 请修正", group, route.Handler))

			}
			if _, ok := actionMap[route.Action]; !ok {
				actionMap[route.Action] = struct{}{}
			} else {
				panic(fmt.Sprintf("%s 中存在重复的action: [%s], 请修正", group, route.Action))
			}

		}
	}
	dir := defaultDir
	rootPkg, err := golang.GetParentPackage(dir)
	if err != nil {
		return err
	}
	// 按模板生成代码
	genHandlers(dir, rootPkg, mgrCfg)
	genRoutes(dir, rootPkg, cfgCron, cfgDaemon)

	return nil
}
func genHandlers(dir, rootPkg string, routesCfg map[string][]jobConfig) {
	for group, routes := range routesCfg {
		for _, route := range routes {
			genHandler(dir, rootPkg, group, route)
		}
	}
}

func genHandler(dir, rootPkg, group string, route jobConfig) error {
	handler := getHandlerName(route)
	handlerPath := getHandlerPath(group)
	pkgName := handlerPath[strings.LastIndex(handlerPath, "/")+1:]
	filename, err := format.FileNamingFormat(config.DefaultFormat, handler)
	if err != nil {
		return err
	}
	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          handlerPath,
		filename:        filename + ".go",
		templateName:    "handlerTemplate",
		category:        category,
		templateFile:    handlerTemplateFile,
		builtinTemplate: handlerTemplate,
		data: map[string]any{
			"PkgName":        pkgName,
			"ImportPackages": genHandlerImports(rootPkg),
			"HandlerName":    handler,
		},
	})
	return nil
}

func getHandlerName(route jobConfig) string {
	return route.Handler + handlerShuffix
}

func getHandlerPath(group string) string {
	return path.Join(handlerDir, strings.ToLower(group))
}

func genHandlerImports(parentPkg string) string {
	imports := []string{
		fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, contextDir)),
	}

	return strings.Join(imports, "\n\t")
}

func getAllGroups(cronRoutesCfg map[string][]jobConfig, daemonRouteCfg map[string][]jobConfig) []string {
	groupMaps := make(map[string]struct{})
	for group, _ := range cronRoutesCfg {
		groupMaps[group] = struct{}{}
	}
	for group, _ := range daemonRouteCfg {
		groupMaps[group] = struct{}{}
	}
	groups := make([]string, 0, len(groupMaps))
	for group, _ := range groupMaps {
		groups = append(groups, group)
	}
	return groups
}

func genRoutes(dir, rootPkg string, cronRoutesCfg map[string][]jobConfig, daemonRouteCfg map[string][]jobConfig) error {
	templateText, err := pathx.LoadTemplate(category, routesAdditionTemplateFile, routesAdditionTemplate)
	if err != nil {
		return err
	}
	gt := template.Must(template.New("groupTemplate").Parse(templateText))
	var builder strings.Builder
	groups := getAllGroups(cronRoutesCfg, daemonRouteCfg)

	cronRoutesString, err := genGroupRoutes(gt, cronRoutesCfg, JobCron)
	if err != nil {
		fmt.Printf("执行失败：%v", err)
		os.Exit(1)
	}
	daemonRoutesString, err := genGroupRoutes(gt, daemonRouteCfg, JobDaemon)
	if err != nil {
		fmt.Printf("执行失败：%v", err)
		os.Exit(1)
		return err
	}
	builder.WriteString(cronRoutesString)
	builder.WriteString("\n")
	builder.WriteString(daemonRoutesString)
	routeFilename, err := format.FileNamingFormat(config.DefaultFormat, routesFilename)
	if err != nil {
		return err
	}
	routeFilename = routeFilename + ".go"
	filename := path.Join(dir, handlerDir, routeFilename)
	os.Remove(filename)

	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          handlerDir,
		filename:        routeFilename,
		templateName:    "routesTemplate",
		category:        category,
		templateFile:    routesTemplateFile,
		builtinTemplate: routesTemplate,
		data: map[string]any{
			"rootPkg":         rootPkg,
			"importPackages":  genRouteImports(rootPkg, groups),
			"routesAdditions": strings.TrimSpace(builder.String()),
		},
	})
}

func genRouteImports(parentPkg string, groups []string) string {
	imports := []string{
		fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, contextDir)),
	}
	for _, group := range groups {
		group = strings.ToLower(group)
		imports = append(imports, fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, handlerDir, group)))
	}

	return strings.Join(imports, "\n\t")
}

func genGroupRoutes(gt *template.Template, routeConfigs map[string][]jobConfig, jobType string) (string, error) {
	var builder strings.Builder
	for group, routes := range routeConfigs {
		group = strings.ToLower(group)
		var gbuilder strings.Builder
		gbuilder.WriteString("[]engine.Route{")
		for _, route := range routes {
			var routeString string
			if jobType == JobDaemon {
				routeString = fmt.Sprintf(`
					{
						JobType: "%s",
						Action:   "%s",
						Desc:     "%s",
						Replicas: %d,
						TestReplicas: %d,
						Handler:  %s.%s%s(svc),
					},`, route.JobType, route.Action, route.Desc, *route.Replicas, *route.TestReplicas, group, route.Handler, handlerShuffix)
			} else {
				// 检查定时任务时间格式合法性
				var cronTxt string
				if len(route.Schedule) > 0 {
					if err := ValidCronExpression(route.Schedule); err != nil {
						fmt.Printf("定时任务[%s]的定时任务schedule格式错误，请检查[%s]", route.Action, route.Schedule)
						os.Exit(1)
					}
					cronTxt, _ = CronExpressionTxt(route.Schedule)
				}
				routeString = fmt.Sprintf(`
					{
						JobType:  "%s",
						Action:   "%s",
						Desc:     "%s",
						Schedule: "%s", // %s
						Handler:  %s.%s%s(svc),
					},`, route.JobType, route.Action, route.Desc, route.Schedule, cronTxt, group, route.Handler, handlerShuffix)
			}
			fmt.Fprint(&gbuilder, routeString)

		}
		var routesString string
		gbuilder.WriteString("\n},")
		routesString = strings.TrimSpace(gbuilder.String())
		if err := gt.Execute(&builder, map[string]string{
			"desc":   fmt.Sprintf("%s %s", group, jobType),
			"routes": routesString,
		}); err != nil {
			return "", err
		}
	}
	return strings.TrimSpace(builder.String()), nil
}
