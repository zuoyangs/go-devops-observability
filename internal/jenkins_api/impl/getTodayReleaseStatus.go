package impl

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/zuoyangs/go-devops-observability/internal/jenkins_api"
)

func (j *JenkinsServiceImpl) GetTodayReleaseStatus(c context.Context, config *jenkins_api.JenkinsTodayReleaseStatusRequest) (string, error) {

	req, err := http.NewRequest("GET", config.BuildsURL+"/api/json?pretty=true", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Accept", AcceptHeader)
	req.Header.Set("Accept-Language", AcceptLanguageHeader)
	req.Header.Set("Connection", ConnectionHeader)
	req.SetBasicAuth(config.Username, config.Password)

	j.client = &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := j.client.Do(req)
	if err != nil {
		return "resp:%s", err
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data jenkins_api.JenkinsRun
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		log.Println("Error parsing JSON:", err)
	}

	return "", nil
}
