package router

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zuoyangs/go-devops-observability/internal/jenkins_api"
)

type JenkinsJobsResponse struct {
	_class      string `json:"_class"`
	JenkinsJobs struct {
		Class string `json:"_class"`
		Name  string `json:"name"`
		URL   string `json:"url"`
	} `json:"jobs"`
}

type BuildInfo struct {
	JenkinsInstanceName string `json:"jenkinsInstanceName"`
	JobName             string `json:"jobName"`
	Building            bool   `json:"building"`
	DisplayName         string `json:"displayName"`
	Duration            int    `json:"duration"`
	Number              int    `json:"number"`
	Result              string `json:"result"`
	Timestamp           int64  `json:"timestamp"`
	BuildsURL           string `json:"url"`
	JobURL              string `json:"jobUrl"`
}

type JenkinsBuildsResponse struct {
	JenkinsInstanceName string      `json:"jenkinsInstanceName"`
	Jobs                []BuildInfo `json:"jobs"`
}

func getBuildInfo(c *gin.Context, jenkinsInstanceName, jobName, jobURL string, job jenkins_api.JenkinsJobs, builds jenkins_api.Builds, buildsHistoryMap map[string][]jenkins_api.JenkinsBuildsHistoryMap, jenkinsUsername string, jenkinsPassword string) JenkinsBuildsResponse {
	buildDetailsURL := builds.URL + "api/json?pretty=true"
	req, err := http.NewRequest("GET", buildDetailsURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.SetBasicAuth(jenkinsUsername, jenkinsPassword)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	//解析json数据
	var buildData map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &buildData); err != nil {
		log.Fatal(err)
	}

	result, ok := buildData["result"].(string)
	if !ok || result == "" {
		result = "Building"
	}

	buildInfo := BuildInfo{
		JenkinsInstanceName: jenkinsInstanceName,
		JobName:             jobName,
		Building:            buildData["building"].(bool),
		DisplayName:         buildData["displayName"].(string),
		Duration:            int(buildData["duration"].(float64)),
		Result:              result, //buildData["result"].(string),
		Timestamp:           int64(buildData["timestamp"].(float64)),
		JobURL:              jobURL,
		BuildsURL:           builds.URL,
		Number:              int(buildData["number"].(float64)),
	}

	// 创建 JenkinsBuildsResponse 实例并填充数据
	jenkinsBuildsResponse := JenkinsBuildsResponse{
		JenkinsInstanceName: jenkinsInstanceName,
		Jobs: []BuildInfo{
			buildInfo,
		},
	}

	return jenkinsBuildsResponse
}
