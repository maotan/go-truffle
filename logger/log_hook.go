package logger

import (
	_ "fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"path"
	"time"
)
//(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration)
var (
	logPath = "E:/go-log"
	logFileName = "system.log"
	FileSuffix = ".log"
	RotationTime = time.Hour * 24
	RotationCount uint = 8
)
func writer(logPath string, level string) *rotatelogs.RotateLogs {
	logFullPath := path.Join(logPath, level)
	var cstSh, _ = time.LoadLocation("Asia/Shanghai") //上海
	fileSuffix := time.Now().In(cstSh).Format("2006-01-02") + FileSuffix

	logier, err := rotatelogs.New(
		logFullPath+"-"+fileSuffix,
		//rotatelogs.WithLinkName(logFullPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithRotationCount(RotationCount),   // 文件最大保存份数
		rotatelogs.WithRotationTime(RotationTime), // 日志切割时间间隔
	)

	if err != nil {
		panic(err)
	}
	return logier
}

func LogerMiddleware() gin.HandlerFunc {
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer(logPath, "debug"),
		log.InfoLevel:   writer(logPath, "info"),
		log.WarnLevel:   writer(logPath, "warn"),
		log.ErrorLevel:  writer(logPath, "error"),
		log.FatalLevel: writer(logPath, "fatal"),
		log.PanicLevel: writer(logPath, "panic"),
	},&log.TextFormatter{DisableColors: true})
	log.AddHook(lfHook)

	return func(c *gin.Context) {
		//开始时间
		startTime := time.Now()
		//处理请求
		c.Next()
		//结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		//请求方式
		reqMethod := c.Request.Method
		//请求路由
		reqUrl := c.Request.RequestURI
		//状态码
		statusCode := c.Writer.Status()
		//请求ip
		clientIP := c.ClientIP()

		// 日志格式
		log.WithFields(log.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"req_method":   reqMethod,
			"req_uri":      reqUrl,
		}).Info()
	}
}

func ConfigLocalFilesystemLogger(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration) {
	baseLogPaht := path.Join(logPath, logFileName)
	writer, err := rotatelogs.New(
		baseLogPaht+".%Y%m%d%H%M",
		//rotatelogs.WithLinkName(baseLogPaht), // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge), // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err != nil {
		log.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer, // 为不同级别设置不同的输出目的
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.FatalLevel: writer,
		log.PanicLevel: writer,
	},&log.TextFormatter{DisableColors: true})
	log.AddHook(lfHook)
}


//切割日志和清理过期日志
func ConfigLocalFilesystemLogger1(filePath string) {
	writer, err := rotatelogs.New(
		filePath+".%Y%m%d%H%M",
		//rotatelogs.WithLinkName(filePath),           // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(time.Second*60*3),     // 文件最大保存时间
		rotatelogs.WithRotationTime(time.Second*60), // 日志切割时间间隔
	)
	if err != nil {
		log.Fatal("Init log failed, err:", err)
	}
	log.SetOutput(writer)
	log.SetLevel(log.InfoLevel)
}