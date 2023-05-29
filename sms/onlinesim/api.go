package onlinesim

import (
	"context"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/k773/utils"
	"github.com/k773/utils/fixedPoint"
	"strconv"
	"time"
)

type Provider struct {
	Ses *resty.Client
	// BusyPollingInterval is used by the methods that are continuously polling data from the api.
	// Default: 15sec.
	BusyPollingInterval time.Duration
}

func New(key string) *Provider {
	var ses = resty.New()
	ses.OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {
		request.SetQueryParam("apikey", key)
		return nil
	})
	return &Provider{Ses: ses, BusyPollingInterval: 15 * time.Second}
}

/*
	Get balance
*/

type GetBalanceResponse struct {
	// Balance: currency - usd.
	Balance fixedPoint.IntScaledP6 `json:"balance"`
	// ZBalance is the frozen amount of funds.
	// ZBalance: currency - usd.
	ZBalance fixedPoint.IntScaledP6 `json:"zbalance"`
}

func (p *Provider) ApiGetBalance(ctx context.Context) (resp GetBalanceResponse, e error) {
	resp, e = ApiRequest[GetBalanceResponse](ctx, p.Ses.R(), "getBalance.php", []string{"1"})
	return
}

/*
	Get states
*/

type GetStatesResponse []*State

func (g *GetStatesResponse) ByTzid(tzid int) *State {
	for _, state := range *g {
		if state.Tzid == tzid {
			return state
		}
	}
	return nil
}

type State struct {
	Tzid     int            `json:"tzid"`
	Response StringResponse `json:"response"`
	// Msg contains the code (or the entire message, depending on the showOnlyCode flag)
	Msg string `json:"msg"`

	Country int    `json:"country"`
	Service string `json:"service"`
	Number  string `json:"number"`
	// Sum represents how much the number costs
	Sum fixedPoint.IntScaledP6 `json:"sum"`

	// Time is the amount of time the number after which the activation will be deactivated.
	Time DurationSeconds `json:"time"`
	// Deadline is derived from Time, the field does not exist in the online-sim's api, it is here only for convenience purposes
	ActivationDeadline time.Time `json:"activation_deadline"`

	Form string `json:"form"`
}

func (p *Provider) ApiGetStates(ctx context.Context, showOnlyCode bool) (resp GetStatesResponse, e error) {
	return p.ApiGetStatesWithTzId(ctx, showOnlyCode, 0)
}

// ApiGetStatesWithTzId
// pass tzid <= 0 to not send the parameter
func (p *Provider) ApiGetStatesWithTzId(ctx context.Context, showOnlyCode bool, tzId int) (resp GetStatesResponse, e error) {
	request := p.Ses.R().
		SetQueryParam("msg_list", "0").
		SetQueryParam("message_to_code", utils.If(showOnlyCode, "1", "0"))
	if tzId > 0 {
		request.SetQueryParam("tzid", strconv.Itoa(tzId))
	}
	resp, e = ApiRequest[GetStatesResponse](ctx, request, "getState.php", nil)
	if e == nil {
		for _, state := range resp {
			if state.Time != 0 {
				state.ActivationDeadline = time.Now().Add(time.Duration(state.Time))
			}
		}
	}
	return
}

/*
	Get number
*/

// GetPhone returns an Activation object. For documentation see ApiGetPhone, Activation.
func (p *Provider) GetPhone(ctx context.Context, country int, service string) (phone *Activation, e error) {
	phoneResp, e := p.ApiGetPhone(ctx, country, service)
	if e != nil {
		return
	}
	phone = NewActivation(p, &phoneResp)
	return
}

type GetPhoneResponse struct {
	Response StringResponse `json:"response"`
	// Tzid is the activation id
	Tzid int `json:"tzid"`
}

// ApiGetPhone gets the phone for the given service.
// Arguments:
// country - phone's country, usually it's just country's prefix - 7, 380, ... (all available: https://onlinesim.io/docs/api/en/sms/getNumbersStats)
// service - service, usually just a plain service name - google, vkcom, ... (see all here: https://onlinesim.io/api/getNumbersStats.php)
func (p *Provider) ApiGetPhone(ctx context.Context, country int, service string) (resp GetPhoneResponse, e error) {
	resp, e = ApiRequest[GetPhoneResponse](ctx, p.Ses.R().SetQueryParam("country", strconv.Itoa(country)).SetQueryParam("service", service), "getNum.php", []string{"1"})
	return
}

/*
	Cancel activation
*/

type SetOperationOkResponse struct {
	Response StringResponse `json:"response"`
	// Tzid is the activation id
	Tzid int `json:"tzid"`
}

func (p *Provider) SetOperationOk(ctx context.Context, tzid int) (resp SetOperationOkResponse, e error) {
	resp, e = ApiRequest[SetOperationOkResponse](ctx, p.Ses.R().SetQueryParam("tzid", strconv.Itoa(tzid)), "setOperationOk.php", []string{"1"})
	return
}

/*
	Generic api
*/

// ApiRequest executes provided request in the provided context, and checks that the code returned by the server is acceptable by the caller.
// allowedResponses can be empty, in such case only responses with an empty (non-existent) response field are valid.
func ApiRequest[T any](ctx context.Context, req *resty.Request, dst string, allowedResponses []string) (res T, e error) {
	r, e := req.SetContext(ctx).Post("https://onlinesim.io/api/" + dst)
	if e != nil {
		return
	}

	// Checking for api errors
	e = checkApiResponse(r, allowedResponses)
	if e != nil {
		return
	}
	e = json.Unmarshal(r.Body(), &res)
	return
}

/*
	Some types wrappers (thanks to the api design)
*/

// DurationSeconds is a wrapper that converts seconds (retrieved from the api) back to time.Duration

type DurationSeconds time.Duration

func (d *DurationSeconds) UnmarshalJSON(data []byte) error {
	var raw time.Duration
	if e := json.Unmarshal(data, &raw); e != nil {
		return e
	}
	*d = DurationSeconds(raw * time.Second)
	return nil
}

// StringResponse is a wrapper that converts shitty api data back to string (in api, field Response sometimes can be an int, other times - a string)
type StringResponse string

func (s *StringResponse) UnmarshalJSON(data []byte) error {
	if len(data) > 2 && data[0] == '"' && data[len(data)-1] == '"' {
		*s = StringResponse(data[1 : len(data)-1])
	} else {
		*s = StringResponse(data)
	}
	return nil
}

// baseApiResponseRaw is a struct for internal usage only; used for checking api response
type baseApiResponseRaw struct {
	Response *StringResponse `json:"response"`
}

func checkApiResponse(r *resty.Response, allowedResponses []string) (e error) {
	// Service's api is a mess. That's all you need to know.
	var res baseApiResponseRaw
	if e = json.Unmarshal(r.Body(), &res); e == nil {
		if res.Response == nil {
			if len(allowedResponses) == 0 {
				return nil
			}
			e = ErrorUnexpectedResponse{Response: "*response field was not found*"}
			return
		}

		var allowed bool
		for _, v := range allowedResponses {
			if string(*res.Response) == v {
				allowed = true
				break
			}
		}
		if !allowed {
			e = ErrorUnexpectedResponse{Response: *res.Response}
		}
	} else {
		// If the error tells us that data type is an array, and we don't expect any response, ignore this error
		if err, ok := e.(*json.UnmarshalTypeError); ok {
			if err.Value == "array" {
				e = nil
			}
		}
	}
	return
}
