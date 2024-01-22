package router

import (
	"errors"

	"github.com/zuoyangs/go-devops-observability/internal/jenkins_api"
)

func GetJenkinsConfig(jenkinsConfig interface{}) (*jenkins_api.JenkinsJobsRequest, error) {

	// 检查 jenkinsConfig 是否为 map
	jenkinsConfigMap, ok := jenkinsConfig.(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid jenkins config: not a map")
	}

	jenkinsURL, ok := jenkinsConfigMap["jenkinsurl"].(string)
	if !ok || jenkinsURL == "" {
		return nil, errors.New("invalid jenkins config: missing or invalid jenkinsurl")
	}

	jenkinsUsername, ok := jenkinsConfigMap["username"].(string)
	if !ok || jenkinsUsername == "" {
		return nil, errors.New("invalid jenkins config: missing or invalid username")
	}

	jenkinsPassword, ok := jenkinsConfigMap["password"].(string)
	if !ok || jenkinsPassword == "" {
		return nil, errors.New("invalid jenkins config: missing or invalid password")
	}

	baseRequest := jenkins_api.JenkinsBaseRequest{
		Username: jenkinsUsername,
		Password: jenkinsPassword,
	}

	return &jenkins_api.JenkinsJobsRequest{
		JenkinsURL:         jenkinsURL,
		JenkinsBaseRequest: baseRequest,
	}, nil
}
