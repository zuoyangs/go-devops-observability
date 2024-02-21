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
	JenkinsInstanceName string `json:"jenkinsInstanceName"`
	JobName             string `json:"jobName"`

	// 当天的统计信息
	TodaySuccessCount int     `json:"todaySuccessCount"`
	TodayFailureCount int     `json:"todayFailureCount"`
	TodayTotalCount   int     `json:"todayTotalCount"`
	TodaySuccessRate  float64 `json:"todaySuccessRate"`
	TodayFailureRate  float64 `json:"todayFailureRate"`

	// 当前月份的统计信息
	CurrentMonthSuccessCount int     `json:"currentMonthSuccessCount"`
	CurrentMonthFailureCount int     `json:"currentMonthFailureCount"`
	CurrentMonthTotalCount   int     `json:"currentMonthTotalCount"`
	CurrentMonthSuccessRate  float64 `json:"currentMonthSuccessRate"`
	CurrentMonthFailureRate  float64 `json:"currentMonthFailureRate"`
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
				//log.Printf("goroutine: jenkins实例: %v,  builds: %v", jenkinsInstanceName, buildData)

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

	// 计算当前月份的第一天
	firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	firstDayOfMonthMillis := firstDayOfMonth.UnixNano() / int64(time.Millisecond)

	resultMap := make(map[string]map[string]*JenkinsJobStatsExtended)

	// 遍历构建数据并计算成功和失败次数
	for _, buildItem := range responseData {

		if buildItem == nil {
			log.Printf("构建项 (buildItem) 为空, 可能尚未发布过任何内容。")
			continue
		}

		// 遍历每个构建
		if build, ok := buildItem.([]JenkinsBuildsResponse); ok {
			for _, items := range build {
				for _, item := range items.Jobs {

					//计算当日的发布次数
					if item.Timestamp >= todayMillis && item.Timestamp <= millis {
						if _, ok := resultMap[item.JenkinsInstanceName]; !ok {
							resultMap[item.JenkinsInstanceName] = make(map[string]*JenkinsJobStatsExtended)
						}
						if _, ok := resultMap[item.JenkinsInstanceName][item.JobName]; !ok {
							resultMap[item.JenkinsInstanceName][item.JobName] = &JenkinsJobStatsExtended{
								JenkinsInstanceName: item.JenkinsInstanceName,
								JobName:             item.JobName,
								TodaySuccessCount:   0,
								TodayFailureCount:   0,
								TodaySuccessRate:    0.0,
								TodayFailureRate:    0.0,
								TodayTotalCount:     0,
							}
						}
						//更新成功或失败次数
						switch item.Result {
						case "SUCCESS":
							resultMap[item.JenkinsInstanceName][item.JobName].TodaySuccessCount++
						case "FAILURE":
							resultMap[item.JenkinsInstanceName][item.JobName].TodayFailureCount++
						default:
							// 处理未知Result，例如记录日志或增加错误计数
							log.Printf("Unknown result for job %s in instance %s: %s", item.JobName, item.JenkinsInstanceName, item.Result)
						}
						resultMap[item.JenkinsInstanceName][item.JobName].TodayTotalCount = resultMap[item.JenkinsInstanceName][item.JobName].TodaySuccessCount + resultMap[item.JenkinsInstanceName][item.JobName].TodayFailureCount

						if resultMap[item.JenkinsInstanceName][item.JobName].TodayTotalCount > 0 {
							resultMap[item.JenkinsInstanceName][item.JobName].TodaySuccessRate = math.Round(float64(resultMap[item.JenkinsInstanceName][item.JobName].TodaySuccessCount)/float64(resultMap[item.JenkinsInstanceName][item.JobName].TodayTotalCount)*100.0*1e2) / 1e2
							resultMap[item.JenkinsInstanceName][item.JobName].TodayFailureRate = math.Round(float64(resultMap[item.JenkinsInstanceName][item.JobName].TodayFailureCount)/float64(resultMap[item.JenkinsInstanceName][item.JobName].TodayTotalCount)*100.0*1e2) / 1e2
						}
					}

					//计算当前月的发布次数
					if item.Timestamp >= firstDayOfMonthMillis && item.Timestamp <= millis {
						if _, ok := resultMap[item.JenkinsInstanceName]; !ok {
							resultMap[item.JenkinsInstanceName] = make(map[string]*JenkinsJobStatsExtended)
						}
						if _, ok := resultMap[item.JenkinsInstanceName][item.JobName]; !ok {
							resultMap[item.JenkinsInstanceName][item.JobName] = &JenkinsJobStatsExtended{
								JenkinsInstanceName:      item.JenkinsInstanceName,
								JobName:                  item.JobName,
								CurrentMonthSuccessCount: 0,
								CurrentMonthFailureCount: 0,
								CurrentMonthTotalCount:   0,
								CurrentMonthSuccessRate:  0.0,
								CurrentMonthFailureRate:  0.0,
							}
						}
						//更新成功或失败次数
						switch item.Result {
						case "SUCCESS":
							resultMap[item.JenkinsInstanceName][item.JobName].CurrentMonthSuccessCount++
						case "FAILURE":
							resultMap[item.JenkinsInstanceName][item.JobName].CurrentMonthFailureCount++
						default:
							// 处理未知Result，例如记录日志或增加错误计数
							log.Printf("Unknown result for job %s in instance %s: %s", item.JobName, item.JenkinsInstanceName, item.Result)
						}
						resultMap[item.JenkinsInstanceName][item.JobName].CurrentMonthTotalCount = resultMap[item.JenkinsInstanceName][item.JobName].CurrentMonthSuccessCount + resultMap[item.JenkinsInstanceName][item.JobName].CurrentMonthFailureCount

						if resultMap[item.JenkinsInstanceName][item.JobName].CurrentMonthTotalCount > 0 {
							resultMap[item.JenkinsInstanceName][item.JobName].CurrentMonthSuccessRate = math.Round(float64(resultMap[item.JenkinsInstanceName][item.JobName].CurrentMonthSuccessCount)/float64(resultMap[item.JenkinsInstanceName][item.JobName].CurrentMonthTotalCount)*100.0*1e2) / 1e2
							resultMap[item.JenkinsInstanceName][item.JobName].CurrentMonthFailureRate = math.Round(float64(resultMap[item.JenkinsInstanceName][item.JobName].CurrentMonthFailureCount)/float64(resultMap[item.JenkinsInstanceName][item.JobName].CurrentMonthTotalCount)*100.0*1e2) / 1e2
						}
					}
				}
			}
		} else {
			// 处理buildItem不是JenkinsBuildsResponse的情况
			log.Printf("处理buildItem不是JenkinsBuildsResponse的情况: %v", buildItem)
		}
	}

	//遍历 stats，将 jobName 下的层级append到切片字符串中
	var stats []*JenkinsJobStatsExtended
	for _, jobName := range resultMap {
		log.Printf("jobName: %v\n", jobName)
		for _, job := range jobName {
			log.Printf("job: %v\n", job)
			stats = append(stats, job)
		}
	}

	c.JSON(http.StatusOK, gin.H{"持续集成发布稳定性指标": stats})
}
