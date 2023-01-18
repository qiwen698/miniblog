/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package miniblog

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/qiwen698/miniblog/pkg/version/verflag"

	"github.com/qiwen698/miniblog/internal/pkg/log"

	"github.com/spf13/viper"

	mw "github.com/qiwen698/miniblog/internal/pkg/middleware"
	"github.com/spf13/cobra"
)

var cfgFile string

func NewMiniBlogCommand() *cobra.Command {
	// cmd represents the base command when called without any subcommands
	cmd := &cobra.Command{
		Use:   "miniblog",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
		// 命令出错时，不打印帮助信息。不需要打印帮助信息，设置为 true 可以保持命令出错时一眼就能看到错误信息
		SilenceUsage: true,
		// 指定调用cmd.Execute()时，执行的Run 函数，函数执行失败会返回错误信息
		RunE: func(cmd *cobra.Command, args []string) error {
			//如果 `--version=true`，则打印版本并退出
			verflag.PrintAndExitIfRequested()
			//初始化日志
			log.Init(logOptions())
			defer log.Sync() // Sync 将缓存中的日志刷新到磁盘文件中
			return run()
		},
		//这里设置命令运行时，不需要指定命令行参数
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments,got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
	}

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "The path to the mini blog configuration file. Empty string for no configuration file.")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// 添加 --version 标志
	verflag.AddFlags(cmd.PersistentFlags())
	return cmd
}

// run 函数是实际的业务代码入口函数
func run() error {

	// 初始化 store 层
	if err := initStore(); err != nil {
		return err
	}

	//设置Gin模式
	gin.SetMode(viper.GetString("runmode"))
	//创建Gin引擎
	g := gin.New()
	mws := []gin.HandlerFunc{gin.Recovery(), mw.NoCache, mw.Cors, mw.Secure, mw.RequestID()}
	g.Use(mws...)
	if err := installRouters(g); err != nil {
		return err
	}
	// 创建 HTTP Server 实例
	httpsrv := &http.Server{
		Addr:    viper.GetString("addr"),
		Handler: g,
	}
	// 运行HTTP服务
	// 打印一条日志，用来提示HTTP服务已经起来，方便排障
	log.Infow("Start to listening the incoming requests on http address", "addr", viper.GetString("addr"))
	go func() {
		if err := httpsrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalw(err.Error())

		}
	}()
	// 等待中断信号优雅地关闭服务器 （10 秒超时）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) //此处不会阻塞
	<-quit                                               //阻塞在此，当接到上述两种信号时才会往下执行
	log.Infow("Shutting down server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpsrv.Shutdown(ctx); err != nil {
		log.Errorw("Insecure Server forced to shutdown", "err", err)
		return err
	}
	log.Infow("Server exiting")
	return nil
}
