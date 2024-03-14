package main

import (
	"coderunner/constants"
	"coderunner/middleware"
	"coderunner/model"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"coderunner/gocommand"
	log "coderunner/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func FuncIndex(ctx *gin.Context) {
	model.OkWithMsg(ctx, "index is ok")
}

func RunCodes(ctx *gin.Context) {
	var payload *model.CodeInfoDTO
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		model.FailWithMsg(ctx, model.ParamsValidError, err.Error())
		return
	}
	if len(payload.Code) > 1024*400 {
		model.FailWithMsg(ctx, model.ParamsValidError, "提交的代码太长，最多允许400KB")
		return
	}
	log.Logger.Info("输入源代码：==》" + payload.Code)
	val, ok := constants.TemplateData[payload.Lang]
	if !ok {
		model.FailWithMsg(ctx, model.ParamsValidError, "暂时不支持该语言:"+payload.Lang)
		return
	}
	temp := val
	if strings.HasPrefix(temp["file"], "regex::") {
		// java
		regexp, err := regexp.Compile(strings.Replace(temp["file"], "regex::", "", 1))
		match := regexp.FindStringSubmatch(payload.Code)
		if err != nil {
			log.Logger.Error(err.Error())
		}
		fmt.Println(match)
		temp["file"] = "/tmp/" + match[1] + ".java"
	}
	//语言选择镜像：payload.Lang
	//源码 后续保存成文件，payload.Code
	u4 := uuid.New() // a0d99f20-1dd1-459b-b516-dfeca4005203
	eof := u4.String()
	filename := temp["file"]
	cmd := temp["cmd"]
	finalcmd := fmt.Sprintf("\n/bin/cat>%s<<\\%s\n%s\n%s\n%s", filename, eof, payload.Code, eof, cmd)
	log.Logger.Info(finalcmd)
	container_name := "runing_container"
	memory := temp["memory"]
	cpuset := temp["cpuset"]
	image := temp["image"]
	lmd := fmt.Sprintf("docker run --name=%s --rm --network=none --cpus=1 --memory=%s --memory-swap=-1 --cpuset-cpus=%s -i %s /bin/bash %s",
		container_name, memory, cpuset, image, finalcmd)
	fmt.Println(lmd)
	mctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	out := make(chan string)
	gocommand.NewCommand().ExecAsync(mctx, out, lmd)
	t, _ := strconv.ParseInt(temp["timeout"], 10, 64)
	for {
		// 其实这段去掉程序也会正常运行，只是我们就不知道到底什么时候Command被停止了，而且如果我们需要实时给web端展示输出的话，这里可以作为依据 取消展示
		select {
		// 检测到ctx.Done()之后停止读取
		case <-mctx.Done():
			if mctx.Err() != nil {
				fmt.Printf("程序出现错误: %q", mctx.Err())
			} else {
				fmt.Println("程序被终止")
			}
			model.FailWithMsg(ctx, model.ParamsValidError, "程序被终止")
			return
		case <-time.After(time.Duration(t) * time.Second):
			fmt.Printf("timeout %d seconds", t)
			model.FailWithMsg(ctx, model.ParamsValidError, "运行超时 30s")
			return
		case output := <-out:
			fmt.Println(output)
			model.OkWithData(ctx, output)
			return
		}
	}

}
func main() {
	router := gin.New()
	router.NoRoute(model.NoRoute)
	router.NoMethod(model.NoMethod)
	router.Use(middleware.Cors(), middleware.ErrorRecover())

	router.GET("/", FuncIndex)
	router.POST("/run", RunCodes)
	serverPort := 8188
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", serverPort),
		Handler:      router,
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
	}

	go func() {
		// 服务连接
		log.Logger.Infof("Server runing at port: %d ", serverPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Logger.Errorf("listen: %s\n", err)
		}
	}()
	// // 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)
	<-quit
	log.Logger.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Logger.Error("Server Shutdown:", err)
	}
	log.Logger.Info("Server exiting")
}
