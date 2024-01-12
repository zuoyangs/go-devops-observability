package jenkins_api

import (
	"context"
)

type JenkinsAPIInterface interface {

	//查询 Jenkins 所有 job 列表
	GetAllJobs(context.Context, *JenkinsJobsRequest) (*JenkinsJobsResponse, error)

	//通过 Jenkins 所有 job 列表，获取指定 job 的构建历史列表
	GetBuildsHistory(context.Context, *JenkinsBuildsRequest) (*JenkinsBuildsResponse, error)

	//通过 job 的构建历史列表，获取指定 job 的当天发布状态
	GetTodayReleaseStatus(context.Context, *JenkinsTodayReleaseStatusRequest) (*JenkinsTodayReleaseStatusResponse, error)

	//通过获取指定 job 的当天发布状态，计算当天该 job 的发版成功率
	CalculateSuccessRate(context.Context, *JenkinsJobsRequest) (float64, error)

	//通过获取指定 job 的当天发布状态，计算当天该 job 的发版失败率
	CalculateFailureRate(context.Context, *JenkinsJobsRequest) (float64, error)

	//获取按天同比、环比的发布成功率和失败率
	CompareDailySuccessFailureRates(context.Context, *JenkinsJobsRequest) (float64, float64, error)

	//获取按星期同比、环比的发布成功率和失败率
	CompareWeeklySuccessFailureRates(context.Context, *JenkinsJobsRequest) (float64, float64, error)
}
