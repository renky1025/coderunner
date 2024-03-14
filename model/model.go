package model

import (
	log "coderunner/logger"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type RunTpl struct {
	Image   string  `json:"image"`            // docker iamge 名字
	File    string  `json:"file"`             // 代码要保存的文件路径
	PrevCmd *string `json:"precmd,omitempty"` // 写入之前执行的命令，主要用于设置一些变量，给cmd中的命令使用
	Cmd     string  `json:"cmd"`              // 保存代码之后要执行的命令
	Timeout int     `json:"timeout"`          // 容器执行超时时间
	Memory  string  `json:"memory"`           // 允许容器使用的内存,例如:20MB
	Cpuset  string  `json:"cpuset"`           // 使用的cpu核心
}

type CodeInfoDTO struct {
	Lang  string  `json:"lang"`            // 选择的语言
	Code  string  `json:"code"`            // 输入代码源码
	Input *string `json:"input,omitempty"` // 其他输入参数
}

// NoRoute 无路由的响应
func NoRoute(c *gin.Context) {
	Fail(c, Request404Error)
}

// NoMethod 无方法的响应
func NoMethod(c *gin.Context) {
	Fail(c, Request405Error)
}

// RespType 响应类型
// swagger:model RespType
type RespType struct {
	code int
	msg  string
	data interface{}
}

// Response 响应格式结构
// swagger:model Response
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

var (
	Success = RespType{code: 200, msg: "成功"}
	Failed  = RespType{code: 500, msg: "失败"}

	ParamsValidError    = RespType{code: 310, msg: "参数校验错误"}
	ParamsTypeError     = RespType{code: 311, msg: "参数类型错误"}
	RequestMethodError  = RespType{code: 312, msg: "请求方法错误"}
	AssertArgumentError = RespType{code: 313, msg: "断言参数错误"}

	LoginAccountError = RespType{code: 330, msg: "登录账号或密码错误"}
	LoginDisableError = RespType{code: 331, msg: "登录账号已被禁用了"}
	TokenEmpty        = RespType{code: 332, msg: "Could not find bearer token in Authorization header"}
	TokenInvalid      = RespType{code: 333, msg: "token参数无效"}

	NoPermission    = RespType{code: 403, msg: "无相关权限"}
	Request404Error = RespType{code: 404, msg: "请求接口不存在"}
	Request405Error = RespType{code: 405, msg: "请求方法不允许"}
	SystemAuthError = RespType{code: 424, msg: "token 失效"}

	SystemError        = RespType{code: 500, msg: "系统错误"}
	SystemHTTPError    = RespType{code: 500, msg: "远程调用服务错误"}
	SystemDataNotFound = RespType{code: 500, msg: "数据不存在"}
)

// Make 以响应类型生成信息
func (rt RespType) Make(msg string) RespType {
	rt.msg = msg
	return rt
}

// MakeData 以响应类型生成数据
func (rt RespType) MakeData(data interface{}) RespType {
	rt.data = data
	return rt
}

// Code 获取code
func (rt RespType) Code() int {
	return rt.code
}

// Msg 获取msg
func (rt RespType) Msg() string {
	return rt.msg
}

// Data 获取data
func (rt RespType) Data() interface{} {
	return rt.data
}

// Result 统一响应
func Result(c *gin.Context, resp RespType, data interface{}) {
	if data == nil {
		data = resp.data
	}
	c.JSON(http.StatusOK, Response{
		Code: resp.code,
		Msg:  resp.msg,
		Data: data,
	})
}

// Copy 拷贝结构体
func Copy(toValue interface{}, fromValue interface{}) interface{} {
	if err := copier.Copy(toValue, fromValue); err != nil {
		log.Logger.Errorf("Copy err: err=[%+v]", err)
		panic(SystemError)
	}
	return toValue
}

// Ok 正常响应
func Ok(c *gin.Context) {
	Result(c, Success, []string{})
}

// OkWithMsg 正常响应附带msg
func OkWithMsg(c *gin.Context, msg string) {
	resp := Success
	resp.msg = msg
	Result(c, resp, []string{})
}

// OkWithData 正常响应附带data
func OkWithData(c *gin.Context, data interface{}) {
	Result(c, Success, data)
}

// respLogger 打印日志
func respLogger(resp RespType, template string, args ...interface{}) {
	loggerFunc := log.Logger.Warnf
	if resp.code >= 500 {
		loggerFunc = log.Logger.Errorf
	}
	loggerFunc(template, args...)
}

// Fail 错误响应
func Fail(c *gin.Context, resp RespType) {
	respLogger(resp, "Request Fail: url=[%s], resp=[%+v]", c.Request.URL.Path, resp)
	Result(c, resp, []string{})
}

// FailWithMsg 错误响应附带msg
func FailWithMsg(c *gin.Context, resp RespType, msg string) {
	resp.msg = msg
	respLogger(resp, "Request FailWithMsg: url=[%s], resp=[%+v]", c.Request.URL.Path, resp)
	Result(c, resp, []string{})
}

// FailWithData 错误响应附带data
func FailWithData(c *gin.Context, resp RespType, data interface{}) {
	respLogger(resp, "Request FailWithData: url=[%s], resp=[%+v], data=[%+v]", c.Request.URL.Path, resp, data)
	Result(c, resp, data)
}

func FileDownload(c *gin.Context, content, filename string) {
	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Accept-Length", fmt.Sprintf("%d", len(content)))
	c.Writer.Write([]byte(content))
}
