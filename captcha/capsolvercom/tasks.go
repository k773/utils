package capsolvercom

type BaseTask struct {
	ClientKey string `json:"clientKey"`
	TaskId    string `json:"taskId,omitempty"`
	Task      *Task  `json:"task,omitempty"`
}

type Task interface {
	fillDefaults()
}

/*
	Common
*/

// TaskType - taskType
type TaskType string

type TaskTypeW struct {
	Type TaskType `json:"type"`
}

type EnterprisePayload struct {
	S string `json:"s"`
}

// Cookies - cookies
type Cookies []Cookie

type CookiesW struct {
	Cookies Cookies `json:"CookiesW"`
}

type Cookie struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// TaskProxyW - task proxy
type TaskProxyW struct {
	Proxy string `json:"proxy,omitempty"` // socks5:ip:port:user:pass
}

/*
	ReCaptchaV2
*/

type ReCaptchaV2Task struct {
	TaskTypeW
	TaskProxyW
	CookiesW
	WebsiteURL        string            `json:"websiteURL"`
	WebsiteKey        string            `json:"websiteKey"`
	EnterprisePayload EnterprisePayload `json:"enterprisePayload,omitempty"`
	IsInvisible       bool              `json:"isInvisible"`
	PageAction        string            `json:"pageAction,omitempty"`
	ApiDomain         string            `json:"apiDomain,omitempty"`
	UserAgent         string            `json:"userAgent,omitempty"`
	Anchor            string            `json:"anchor,omitempty"`
	Reload            string            `json:"reload,omitempty"`
}

// Recaptcha v2 tasks
const (
	TaskTypeReCaptchaV2Task                    TaskType = "ReCaptchaV2Task"
	TaskTypeReCaptchaV2EnterpriseTask          TaskType = "TaskTypeReCaptchaV2EnterpriseTask"
	TaskTypeReCaptchaV2TaskProxyLess           TaskType = "TaskTypeReCaptchaV2TaskProxyLess"
	TaskTypeReCaptchaV2EnterpriseTaskProxyLess TaskType = "TaskTypeReCaptchaV2EnterpriseTaskProxyLess"
)

/*
	Funcaptcha
*/

type FunCaptchaTask struct {
	TaskTypeW
	TaskProxyW
	WebsiteURL               string `json:"websiteURL"`
	WebsitePublicKey         string `json:"websitePublicKey"`
	FuncaptchaApiJSSubdomain string `json:"funcaptchaApiJSSubdomain,omitempty"`
	Data                     string `json:"data,omitempty"`
}
