package jenkins_api

// getAllJobs
type JenkinsJobs struct {
	Class   string `json:"_class"`
	Name    string `json:"name"`
	JobsURL string `json:"url"`
}
type JenkinsJobsResponse struct {
	_class string        `json:"_class"`
	Jobs   []JenkinsJobs `json:"jobs"`
}

type JenkinsJobsRequest struct {
	JenkinsURL         string             `json:"jenkinsURL"`
	JenkinsBaseRequest JenkinsBaseRequest `json:"jenkinsBaseRequest"`
}

// getBuildHistory
type JenkinsBuildsRequest struct {
	JenkinsName string `json:"jenkinsName"`
	JobURL      string `json:"jobURL"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}

type JenkinsBuildsResponse struct {
	Class           string     `json:"_class"`
	Actions         []struct{} `json:"actions"`
	Description     string     `json:"description"`
	DisplayName     string     `json:"displayName"`
	FullDisplayName string     `json:"fullDisplayName"`
	FullName        string     `json:"fullName"`
	Name            string     `json:"name"`
	URL             string     `json:"url"`
	Buildable       bool       `json:"buildable"`
	Builds          []struct {
		BuildsClass string `json:"_class"`
		Number      int    `json:"number"`
		URL         string `json:"url"`
	} `json:"builds"`
	Color      string `json:"color"`
	FirstBuild struct {
		BuildsClass string `json:"_class"`
		Number      int    `json:"number"`
		URL         string `json:"url"`
	} `json:"firstBuild"`
	HealthReport []struct {
		Description   string `json:"description"`
		IconClassName string `json:"iconClassName"`
		IconURL       string `json:"iconUrl"`
		Score         int    `json:"score"`
	} `json:"healthReport"`
	InQueue          bool `json:"inQueue"`
	KeepDependencies bool `json:"keepDependencies"`
	LastBuild        struct {
		BuildsClass string `json:"_class"`
		Number      int    `json:"number"`
		URL         string `json:"url"`
	} `json:"lastBuild"`
	LastCompletedBuild struct {
		BuildsClass string `json:"_class"`
		Number      int    `json:"number"`
		URL         string `json:"url"`
	} `json:"lastCompletedBuild"`
	LastFailedBuild struct {
		BuildsClass string `json:"_class"`
		Number      int    `json:"number"`
		URL         string `json:"url"`
	} `json:"lastFailedBuild"`
	LastStableBuild struct {
		BuildsClass string `json:"_class"`
		Number      int    `json:"number"`
		URL         string `json:"url"`
	} `json:"lastStableBuild"`
	LastSuccessfulBuild struct {
		BuildsClass string `json:"_class"`
		Number      int    `json:"number"`
		URL         string `json:"url"`
	} `json:"lastSuccessfulBuild"`
	LastUnstableBuild     interface{} `json:"lastUnstableBuild"`
	LastUnsuccessfulBuild struct {
		BuildsClass string `json:"_class"`
		Number      int    `json:"number"`
		URL         string `json:"url"`
	} `json:"lastUnsuccessfulBuild"`
	NextBuildNumber int `json:"nextBuildNumber"`
	Property        []struct {
		Class string `json:"_class"`
	} `json:"property"`
	QueueItem       interface{} `json:"queueItem"`
	ConcurrentBuild bool        `json:"concurrentBuild"`
	ResumeBlocked   bool        `json:"resumeBlocked"`
}

type JenkinsBuildsHistoryMap struct {
	JenkinsName string  `json:"jenkinsName"`
	Builds      []Build `json:"builds"`
	JobName     string  `json:"jobName"`
	Id          int     `json:"id"`
}

type JenkinsBaseRequest struct {
	JenkinsName string `json:"jenkinsName"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}

type JenkinsJobUrlRequest struct {
	JobURL             string             `json:"jobURL"`
	JenkinsBaseRequest JenkinsBaseRequest `json:"jenkinsBaseRequest"`
}

// 第3部分
type Builds struct {
	BuildsClass string `json:"_class"`
	Number      int    `json:"number"`
	URL         string `json:"url"`
}

type Build struct {
	Number          int    `json:"number"`
	URL             string `json:"url"`
	ID              string `json:"number"`
	Building        bool   `json:"building"`
	FullDisplayName string `json:"fullDisplayName"`
	Result          string `json:"result"`
	Timestamp       int64  `json:"timestamp"`
}

type JenkinsTodayReleaseStatusRequest struct {
	Number    int    `json:"number"`
	BuildsURL string `json:"buildsURL"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type Cause struct {
	Class     string `json:"_class"`
	ShortDesc string `json:"shortDescription"`
	UserID    string `json:"userId"`
	UserName  string `json:"userName"`
}

type BuildsData struct {
	Class        string `json:"_class"`
	BuildNumber  int    `json:"buildNumber"`
	BuildResult  string `json:"buildResult"`
	MarkedSHA1   string `json:"marked>SHA1"`
	MarkedBranch struct {
		SHA1 string `json:"SHA1"`
		Name string `json:"name"`
	} `json:"marked>branch"`
	RevisionSHA1   string `json:"revision>SHA1"`
	RevisionBranch struct {
		SHA1 string `json:"SHA1"`
		Name string `json:"name"`
	} `json:"revision>branch"`
}

type BuildData struct {
	Class              string           `json:"_class"`
	BuildsByBranchName map[string]Build `json:"buildsByBranchName"`
	LastBuiltRevision  struct {
		SHA1   string `json:"SHA1"`
		Branch []struct {
			SHA1 string `json:"SHA1"`
			Name string `json:"name"`
		} `json:"branch"`
	} `json:"lastBuiltRevision"`
	RemoteUrls []string `json:"remoteUrls"`
	ScmName    string   `json:"scmName"`
}

type JenkinsRun struct {
	Class             string        `json:"_class"`
	Actions           []interface{} `json:"actions"`
	Artifacts         []interface{} `json:"artifacts"`
	Building          bool          `json:"building"`
	Description       string        `json:"description"`
	DisplayName       string        `json:"displayName"`
	Duration          int           `json:"duration"`
	EstimatedDuration int           `json:"estimatedDuration"`
	Executor          interface{}   `json:"executor"`
	FullDisplayName   string        `json:"fullDisplayName"`
	ID                string        `json:"id"`
	KeepLog           bool          `json:"keepLog"`
	Number            int           `json:"number"`
	QueueID           int           `json:"queueId"`
	Result            string        `json:"result"`
	Timestamp         int64         `json:"timestamp"`
	URL               string        `json:"url"`
	ChangeSets        []interface{} `json:"changeSets"`
	Culprits          []interface{} `json:"culprits"`
	NextBuild         struct {
		Number int    `json:"number"`
		URL    string `json:"url"`
	} `json:"nextBuild"`
	PreviousBuild interface{} `json:"previousBuild"`
}

type TodayReleaseStatus struct {
	FullDisplayName string `json:"fullDisplayName"`
	Number          string `json:"number"`
	Id              string `json:"id"`
	Result          string `json:"result"`
	Timestamp       int64  `json:"timestamp"`
	URL             string `json:"url"`
}

type JenkinsTodayReleaseStatusResponse struct {
	Todaytimestamp     int64                `json:"todaytimestamp"`
	TodayReleaseStatus []TodayReleaseStatus `json:"todayReleaseStatus"`
}
