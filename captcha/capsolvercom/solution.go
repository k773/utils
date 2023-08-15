package capsolvercom

import "github.com/k773/utils/fixedPoint"

type Solution struct {
	Object             []bool    `json:"objects,omitempty"`
	Box                []float64 `json:"box,omitempty"`
	ImageSizes         []int     `json:"imageSize,omitempty"`
	Text               string    `json:"text,omitempty"`
	UserAgent          string    `json:"userAgent,omitempty"`
	ExpireTime         int64     `json:"expireTime,omitempty"`
	GRecaptchaResponse string    `json:"gRecaptchaResponse,omitempty"`
	Challenge          string    `json:"challenge,omitempty"`
	Validate           string    `json:"validate,omitempty"`
	CaptchaId          string    `json:"captcha-id,omitempty"`
	CaptchaOutput      string    `json:"captcha-output,omitempty"`
	GenTime            string    `json:"gen_time,omitempty"`
	LogNumber          string    `json:"log_number,omitempty"`
	PassToken          string    `json:"pass_token,omitempty"`
	RiskType           string    `json:"risk_Type,omitempty"`
	Token              string    `json:"token,omitempty"`
	Cookie             string    `json:"cookie,omitempty"`
	Type               string    `json:"type,omitempty"`
}

type CapSolverResponseError struct {
	ErrorId          int    `json:"errorId"`
	ErrorCode        string `json:"errorCode"`
	ErrorDescription string `json:"errorDescription"`
}

func (c *CapSolverResponseError) Error() string {
	return c.ErrorDescription
}

type CapSolverResponse struct {
	CapSolverResponseError
	Status   string                 `json:"status,omitempty"`
	Solution Solution               `json:"solution,omitempty"`
	TaskId   string                 `json:"taskId,omitempty"`
	Balance  fixedPoint.IntScaledP6 `json:"balance,omitempty"`
	Packages []string               `json:"packages,omitempty"`
}

func (c *CapSolverResponse) isReady() bool {
	return c.Status == "ready"
}
