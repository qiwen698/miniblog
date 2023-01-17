package version

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/gosuri/uitable"
)

var (
	//GitVersion 是语义化的版本
	GitVersion = "v0.0.0-master+$Format:%h$"
	//BuildDate 是ISO8601 格式的构建时间，$(date -u +'%Y-%m-%dT%H:%M:%SZ') 命令的输出
	BuildDate = "1970-01-01T00:00:00Z"
	//GitCommit 是Git的SHA1值，$(git rev-parse HEAD)  命令的输出
	GitCommit = "$Format:%H$"

	// GitTreeState

	GitTreeState = ""
)

type Info struct {
	GitVersion   string `json:"gitVersion"`
	GitCommit    string `json:"gitCommit"`
	GitTreeState string `json:"gitTreeState"`
	BuildDate    string `json:"buildDate"`
	GoVersion    string `json:"goVersion"`
	Compiler     string `json:"compiler"`
	Platform     string `json:"platform"`
}

func (info Info) String() string {
	if s, err := info.Text(); err == nil {
		return string(s)
	}
	return info.GitVersion

}

// ToJSON 以JSON格式返回版本信息

func (info Info) ToJSON() string {
	s, _ := json.Marshal(info)
	return string(s)
}

func (info Info) Text() ([]byte, error) {
	table := uitable.New()
	table.RightAlign(0)
	table.MaxColWidth = 80
	table.Separator = " "
	table.AddRow("gitVersion:", info.GitVersion)
	table.AddRow("gitCommit:", info.GitCommit)
	table.AddRow("gitTreeState:", info.GitTreeState)
	table.AddRow("buildDate:", info.BuildDate)
	table.AddRow("goVersion", info.GoVersion)
	table.AddRow("compiler:", info.Compiler)
	table.AddRow("platform:", info.Platform)
	return table.Bytes(), nil
}

// Get 返回详尽的代码版本信息，用来标明二进制文件由哪个版本的代码构建
func Get() Info {
	// 以下变量通常由 -ldflags 进行设置
	return Info{
		GitVersion:   GitVersion,
		GitCommit:    GitCommit,
		GitTreeState: GitTreeState,
		BuildDate:    BuildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
