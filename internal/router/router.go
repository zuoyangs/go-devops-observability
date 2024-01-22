package router

import (
	"context"
	"math"
	"net/http"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/zuoyangs/go-devops-observability/internal/jenkins_api"
	"github.com/zuoyangs/go-devops-observability/internal/jenkins_api/impl"
)

type ResponseData struct {
	JenkinsInstanceName string `json:"jenkinsInstanceName"`
	JobName             string `json:"jobName"`
	SuccessCount        int    `json:"successCount"`
	FailureCount        int    `json:"failureCount"`
	Building            bool   `json:"building"`
	DisplayName         string `json:"displayName"`
	Duration            int    `json:"duration"`
	Number              int    `json:"number"`
	Result              string `json:"result"`
	Timestamp           int64  `json:"timestamp"`
	BuildsURL           string `json:"url"`
	JobURL              string `json:"jobUrl"`
}

type JenkinsJobStatsExtended struct {
	JenkinsInstanceName string  `json:"jenkinsInstanceName"`
	JobName             string  `json:"jobName"`
	SuccessCount        int     `json:"successCount"`
	FailureCount        int     `json:"failureCount"`
	TotalCount          int     `json:"totalCount"`
	SuccessRate         float64 `json:"successRate"`
	FailureRate         float64 `json:"failureRate"`
}

func SetupAPIRouters(r *gin.Engine) {
	r.GET("/metrics", getMetricsHandler)
}

func getMetricsHandler(c *gin.Context) {

	jenkinsConfigs := viper.AllSettings()
	buildsHistoryMap := make(map[string][]jenkins_api.JenkinsBuildsHistoryMap)
	var responseData []interface{}
	jenkinsService := &impl.JenkinsServiceImpl{}

	log.Println("正在从'etc/config.yaml'配置文件中获取Jenkins实例的用户名、密码以及访问地址...")

	for jenkinsInstanceName, jenkinsConfig := range jenkinsConfigs {

		//读取 Jenkins 配置信息
		jobs_config, err := GetJenkinsConfig(jenkinsConfig)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Jenkins配置格式错误"})
			return
		}

		log.Printf("正在从Jenkins实例[%s]中获取Jobs信息...", jenkinsInstanceName)
		// 调用 GetAllJobs 方法获取 Jenkins 中的 jobs 信息
		jobs, err := jenkinsService.GetAllJobs(context.Background(), jobs_config)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 检查Jobs列表是否为空或其第一个元素的JobsURL属性是否为空字符串
		if jobs.Jobs == nil || jobs.Jobs[0].JobsURL == "" {
			log.Warningf("对于%s实例, job.Jobs为nil或jobs.Jobs的第一个元素的JobsURL为空字符串", jenkinsInstanceName)
			continue
		}

		var wg sync.WaitGroup
		resultCh := make(chan gin.H, len(jobs.Jobs))

		//填充 jobsMap 逻辑
		for _, job := range jobs.Jobs {
			wg.Add(1)

			// 启动一个goroutine来处理每个job
			go func(job jenkins_api.JenkinsJobs) {
				defer wg.Done()

				jobURL := job.JobsURL
				builds_config := &jenkins_api.JenkinsBuildsRequest{
					JobURL:   jobURL,
					Username: jobs_config.JenkinsBaseRequest.Username,
					Password: jobs_config.JenkinsBaseRequest.Password,
				}
				buildsHistory, err := jenkinsService.GetBuildsHistory(context.Background(), builds_config)
				if err != nil {
					resultCh <- gin.H{"error": err.Error()}
					return
				}

				var buildData []JenkinsBuildsResponse
				for _, build := range buildsHistory.Builds {
					data := getBuildInfo(c, jenkinsInstanceName, job.Name, job.JobsURL, job, build, buildsHistoryMap, jobs_config.JenkinsBaseRequest.Username, jobs_config.JenkinsBaseRequest.Password)
					buildData = append(buildData, data)
				}
				log.Printf("goroutine: jenkins实例: %v,  builds: %v", jenkinsInstanceName, buildData)

				// 将结果发送到channel
				resultCh <- gin.H{"builds": buildData}
			}(job)
		}

		wg.Wait()

		close(resultCh)

		// 从channel中读取结果并处理它们
		for result := range resultCh {
			if errStr, ok := result["error"]; ok {
				// 处理错误
				c.JSON(http.StatusInternalServerError, gin.H{"error": errStr})
				return
			} else if errNull, ok := result["null"]; ok {
				// 处理错误
				c.JSON(http.StatusInternalServerError, gin.H{"error": errNull})
				continue
			} else {
				// 处理结果
				builds := result["builds"].([]JenkinsBuildsResponse)
				responseData = append(responseData, builds)
			}
		}
	}

	now := time.Now()
	millis := now.UnixNano() / int64(time.Millisecond)

	todayStart := time.Now().Truncate(24 * time.Hour)
	todayMillis := todayStart.UnixNano() / int64(time.Millisecond)

	resultMap := make(map[string]map[string]*ResponseData)

	// 遍历构建数据并计算成功和失败次数
	for _, buildItem := range responseData {

		if buildItem == nil {
			log.Printf("buildItem is null,可能未发布过")
			continue
		}

		// 遍历每个构建
		if build, ok := buildItem.([]JenkinsBuildsResponse); ok {
			for _, items := range build {
				for _, item := range items.Jobs {
					if item.Timestamp >= todayMillis && item.Timestamp <= millis {
						if _, ok := resultMap[item.JenkinsInstanceName]; !ok {
							resultMap[item.JenkinsInstanceName] = make(map[string]*ResponseData)
						}
						if _, ok := resultMap[item.JenkinsInstanceName][item.JobName]; !ok {
							resultMap[item.JenkinsInstanceName][item.JobName] = &ResponseData{
								JenkinsInstanceName: item.JenkinsInstanceName,
								JobName:             item.JobName,
								JobURL:              item.JobURL,
								SuccessCount:        0,
								FailureCount:        0,
							}
						}
						//更新成功或失败次数
						switch item.Result {
						case "SUCCESS":
							resultMap[item.JenkinsInstanceName][item.JobName].SuccessCount++
							log.Printf("SUCCESS:%v", resultMap[item.JenkinsInstanceName][item.JobName])
						case "FAILURE":
							resultMap[item.JenkinsInstanceName][item.JobName].FailureCount++
							log.Printf("FAILURE:%v", resultMap[item.JenkinsInstanceName][item.JobName])

						default:
							// 处理未知Result，例如记录日志或增加错误计数
							log.Printf("Unknown result for job %s in instance %s: %s", item.JobName, item.JenkinsInstanceName, item.Result)
						}
					}
				}
			}
		} else {
			// 处理buildItem不是JenkinsBuildsResponse的情况
			log.Printf("处理buildItem不是JenkinsBuildsResponse的情况: %v", buildItem)
		}
	}

	// 计算成功率和失败率，并将结果转换为新的结构体输出
	stats := make([]JenkinsJobStatsExtended, 0)

	for jenkinsInstanceName, jobs := range resultMap {

		for jobName, data := range jobs {
			successCount := data.SuccessCount
			failureCount := data.FailureCount
			totalCount := successCount + failureCount

			successRate := 0.0
			failureRate := 0.0
			if totalCount > 0 {
				successRate = math.Round(float64(successCount)/float64(totalCount)*100.0*1e2) / 1e2
				failureRate = math.Round(float64(failureCount)/float64(totalCount)*100.0*1e2) / 1e2
			}

			stats = append(stats, JenkinsJobStatsExtended{
				JenkinsInstanceName: jenkinsInstanceName,
				JobName:             jobName,
				SuccessCount:        successCount,
				FailureCount:        failureCount,
				SuccessRate:         successRate,
				FailureRate:         failureRate,
				TotalCount:          totalCount,
			})
		}

	}

	//如果stat是空的，则添置默认值
	if len(stats) == 0 {
		stats = append(stats, JenkinsJobStatsExtended{
			JenkinsInstanceName: "无",
			JobName:             "无",
			SuccessCount:        0,
			FailureCount:        0,
			SuccessRate:         0.0,
			FailureRate:         0.0,
			TotalCount:          0,
		})
	}
	log.Printf("stats:%v", stats)
	c.JSON(http.StatusOK, gin.H{"持续集成发布稳定性指标": stats})

}
