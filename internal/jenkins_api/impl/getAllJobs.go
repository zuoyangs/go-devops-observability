package impl

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/zuoyangs/go-jenkins-api/internal/jenkins_api"
)

// 获取 Jenkins 中的 jobs 信息
func (j *JenkinsServiceImpl) GetAllJobs(c context.Context, config *jenkins_api.JenkinsJobsRequest) (*jenkins_api.JenkinsJobsResponse, error) {

	if config == nil {
		return nil, errors.New("config is nil")
	}

	req, err := http.NewRequest("GET", config.JenkinsURL+"/api/json?pretty=true", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Accept", AcceptHeader)
	req.Header.Set("Accept-Language", AcceptLanguageHeader)
	req.Header.Set("Connection", ConnectionHeader)
	req.SetBasicAuth(config.JenkinsBaseRequest.Username, config.JenkinsBaseRequest.Password)

	j.client = &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := j.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var JenkinsJobsResponse jenkins_api.JenkinsJobsResponse
	err = json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&JenkinsJobsResponse)
	if err != nil {
		return nil, err
	}

	return &JenkinsJobsResponse, nil
}
