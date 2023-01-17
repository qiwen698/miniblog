package miniblog

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/qiwen698/miniblog/internal/pkg/log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	recommendedHomeDir = ".miniblog"
	defaultConfigName  = "miniblog.yaml"
)

// initConfig 设置读取的配置文件名，环境变量，并读取配置文件内容到 viper中

func initConfig() {
	if cfgFile != "" {
		// 从命令行选项指定的配置文件读取
		// 用来设置viper需要读取的配置文件（该配置文件通过--config参数指定）
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		// 获取用户主目录
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".miniblog" (without extension).
		// 将 `$HOME/<recommendedHomeDir>` 目录加入到配置文件的搜索路径中
		viper.AddConfigPath(filepath.Join(home, recommendedHomeDir))
		// 把当前目录加入到配置文件的搜索路劲中
		viper.AddConfigPath(".")
		//viper.AddConfigPath("configs")
		// 设置配置文件格式
		viper.SetConfigType("yaml")
		//设置配置文件名（没有文件扩展名）
		viper.SetConfigName(defaultConfigName)
	}

	//通过 设置viper查找是否有跟配置文件相匹配的环境变量，如果有则将该环境变量的值设置为配置项的值
	viper.AutomaticEnv() // read in environment variables that match

	//读取环境变量的前缀为MINIBLOG,如果是miniblog，将自动转变为大写
	viper.SetEnvPrefix("MINIBLOG")

	//以下2行，将viper.Get(key) key 字符串中 '.' 和 '-' 替换为'_'
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// If a config file is found, read it in.
	// 读取设置的配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Errorw("Failed to read viper configuration file", "err", err)
	}
	//打印 viper当前使用的配置文件，方便Debug
	log.Infow("Using config file", "file", viper.ConfigFileUsed())
}

func logOptions() *log.Options {
	return &log.Options{
		DisableCaller:     viper.GetBool("log.disable-caller"),
		DisableStacktrace: viper.GetBool("log.disable-stacktrace"),
		Level:             viper.GetString("log.level"),
		Format:            viper.GetString("log.format"),
		OutputPaths:       viper.GetStringSlice("log.output-paths"),
	}
}
