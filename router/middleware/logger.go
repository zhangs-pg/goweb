package middleware


// 日志初始化包  调用qmlog.QMLog.Info 记录日志 24小时切割 日志保存7天 可自行设置
import (
	"fmt"
	"m/utils"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"time"

	"bytes"
	"net/http/httputil"
	"strings"

	"github.com/gin-gonic/gin"
)

var QMLog = logrus.New()

//禁止logrus的输出
func InitLog() *logrus.Logger {
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}
	QMLog.Out = src
	QMLog.SetLevel(logrus.DebugLevel)
	if ok, _ := utils.PathExists("./log"); !ok {
		// Directory not exist
		fmt.Println("Create log.")
		_ = os.Mkdir("log", os.ModePerm)
	}
	apiLogPath := "./log/api.log"
	logWriter, err := rotatelogs.New(
		apiLogPath+".%Y-%m-%d-%H-%M.log",
		rotatelogs.WithLinkName(apiLogPath),       // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
	)
	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
	}
	lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{})
	QMLog.AddHook(lfHook)
	return QMLog
}


func Logger() gin.HandlerFunc {
	log := QMLog
	return func(c *gin.Context) {
		// request time
		start := time.Now()
		// request path
		path := c.Request.URL.Path
		logFlag := true
		if strings.Contains(path, "swagger") {
			logFlag = false
		}
		// request ip
		clientIP := c.ClientIP()
		// method
		method := c.Request.Method
		// copy request content
		req, _ := httputil.DumpRequest(c.Request, true)
		if logFlag {
			log.Infof(`| %s | %s | %s | %5s | %s\n`,
				`Request :`, method, clientIP, path, string(req))
		}
		// replace writer
		cusWriter := &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = cusWriter
		// handle request
		c.Next()
		// ending time
		end := time.Now()
		//execute time
		latency := end.Sub(start)
		statusCode := c.Writer.Status()
		if logFlag {
			log.Infof(`| %s | %3d | %13v | %s \n`,
				`Response:`,
				statusCode,
				latency,
				cusWriter.body.String())
		}
	}
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
