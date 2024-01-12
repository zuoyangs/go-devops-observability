package impl

import "net/http"

func NewJenkinsServiceImpl(client *http.Client) *JenkinsServiceImpl {
	return &JenkinsServiceImpl{}
}

type JenkinsServiceImpl struct {
	client *http.Client
}
