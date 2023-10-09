package domain

type MonoResponse struct {
	Message string `json:"message"`
}
type Alive struct {
	Version        string `json:"version"`
	ArtifactId     string `json:"artifactId"`
	GroupId        string `json:"groupId"`
	Echo           string `json:"echo"`
	BuildTimestamp string `json:"buildTimestamp"`
}

type RequestInput struct {
	UniqueID  string `json:"profileId"`
	Timestamp string `json:"referenceId"`
}
type RequestDetailsInput struct {
	UniqueID string `json:"profileId"`
}

type RequestReferenceIdInput struct {
	UniqueID    string `json:"profileId"`
	ReferenceId string `json:"referenceId"`
}

type RequestTimestampInput struct {
	UniqueID  string `json:"profileId"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}
type SaveSignInInfo struct {
	UniqueId    string `dynamodbav:"uniqueId" json:"uniqueId,omitempty"`
	TimeStamp   string `dynamodbav:"timestamp" json:"timeStamp,omitempty"`
	CalledId    string `dynamodbav:"calledId" json:"calledId,omitempty"`
	IpAddress   string `dynamodbav:"ipAddress" json:"ipAddress,omitempty"`
	UserAgent   string `dynamodbav:"userAgent" json:"userAgent,omitempty"`
	SourceId    string `dynamodbav:"sourceId" json:"sourceId,omitempty"`
	Region      string `dynamodbav:"region" json:"region,omitempty"`
	ReferenceId string `dynamodbav:"referenceId" json:"referenceId,omitempty"`
	SsoOrgId    string `dynamodbav:"ssoOrgId" json:"ssoOrgId,omitempty"`
}

type SignInInfo struct {
	UniqueId    string `dynamodbav:"uniqueId" json:"uniqueId,omitempty"`
	TimeStamp   string `dynamodbav:"timestamp" json:"timeStamp,omitempty"`
	CalledId    string `dynamodbav:"calledId" json:"calledId,omitempty"`
	IpAddress   string `dynamodbav:"ipAddress" json:"ipAddress,omitempty"`
	UserAgent   string `dynamodbav:"userAgent" json:"userAgent,omitempty"`
	SourceId    string `dynamodbav:"sourceId" json:"sourceId,omitempty"`
	Region      string `dynamodbav:"region" json:"region,omitempty"`
	ReferenceId string `dynamodbav:"referenceId" json:"referenceId,omitempty"`
}

type ErrorResponse struct {
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	Error        string `json:"error"`
}
