/**
* @Author: mo tan
* @Description:
* @Date 2021/1/1 21:58
 */
package logger

import (
	"fmt"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/maotan/go-truffle/yaml_config"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"path"
	"time"
)

var (
	logPath = "./"
	serverName = ""
	fileSuffix = ".log"
	rotationTime = time.Hour * 24
	rotationCount uint = 8
)

func init()  {
	if yaml_config.YamlConf.LogConf.LogPath != ""{
		logPath = yaml_config.YamlConf.LogConf.LogPath
	}
	if yaml_config.YamlConf.ServerConf.Name != ""{
		serverName = yaml_config.YamlConf.ServerConf.Name
	}
	ConfigLocalFileLogger()
}

func writer(logPath string, level string) *rotatelogs.RotateLogs {
	appLevel := fmt.Sprintf("%s-%s", serverName, level)
	logFullPath := path.Join(logPath, appLevel)
	var cstSh, _ = time.LoadLocation("Asia/Shanghai") //上海
	fileSuffix := time.Now().In(cstSh).Format("2006-01-02") + fileSuffix

	logier, err := rotatelogs.New(
		logFullPath+"-"+fileSuffix,
		//rotatelogs.WithLinkName(logFullPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithRotationCount(rotationCount), // 文件最大保存份数
		rotatelogs.WithRotationTime(rotationTime),   // 日志切割时间间隔
	)

	if err != nil {
		panic(err)
	}
	return logier
}

func LogerMiddleware() gin.HandlerFunc {
	/*lfHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer(logPath, "debug"),
		log.InfoLevel:  writer(logPath, "info"),
		log.WarnLevel:  writer(logPath, "warn"),
		log.ErrorLevel: writer(logPath, "error"),
		log.FatalLevel: writer(logPath, "fatal"),
		log.PanicLevel: writer(logPath, "panic"),
	}, &MineFormatter{})
	// &MineFormatter{}  &log.TextFormatter{DisableColors: true}
	log.AddHook(lfHook)*/

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

		if reqUrl!="/actuator/health"{
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
}

//切割日志和清理过期日志
func ConfigLocalFileLogger() {
	/*writer := writer(logPath, "info")
	log.SetOutput(writer)*/
	log.SetLevel(log.InfoLevel)
	log.SetReportCaller(true)
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer(logPath, "debug"), // 为不同级别设置不同的输出目的
		log.InfoLevel:  writer(logPath, "info"),
		log.WarnLevel:  writer(logPath, "warn"),
		log.ErrorLevel: writer(logPath, "error"),
		log.FatalLevel: writer(logPath, "fatal"),
		log.PanicLevel: writer(logPath, "panic"),
	}, &log.TextFormatter{DisableColors: true})
	log.AddHook(lfHook)
}