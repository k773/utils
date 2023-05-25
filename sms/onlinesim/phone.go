package onlinesim

import (
	"context"
	"github.com/k773/utils"
	"strconv"
)

// Activation is a wrapper around activation-connected data that simplifies interactions with an overcomplicated api
type Activation struct {
	api *Provider

	GetPhoneResponse *GetPhoneResponse
	State            *State

	number string
}

func NewActivation(api *Provider, getPhoneResponse *GetPhoneResponse) *Activation {
	return &Activation{
		api:              api,
		GetPhoneResponse: getPhoneResponse,
	}
}

func (p *Activation) GetMessage(ctx context.Context) (code string, e error) {
	e = p.BusyStatePolling(ctx, func(state *State) (next bool, e error) {
		if state.Msg != "" {
			code = state.Msg
		}
		return code == "", e
	})
	return
}

func (p *Activation) GetNumber(ctx context.Context) (number string, e error) {
	if p.number != "" {
		number = p.number
		return
	}

	e = p.BusyStatePolling(ctx, func(state *State) (next bool, e error) {
		if state.Number != "" {
			number = state.Number
		}
		return number == "", e
	})
	if e == nil {
		p.number = number
	}
	return
}

// BusyStatePolling polls states from the api until provided func returns false or a non-nil error.
func (p *Activation) BusyStatePolling(ctx context.Context, do func(*State) (next bool, e error)) (e error) {
	var next = true
	for e == nil && next {
		if e = ctx.Err(); e != nil {
			continue
		}
		var state *State
		if state, e = p.GetState(ctx); e != nil {
			continue
		}
		next, e = do(state)
		if next && e == nil {
			e = utils.SleepWithContext(ctx, p.api.BusyPollingInterval)
		}
	}
	return
}

// GetState tries to find a state with this activation's tzid
func (p *Activation) GetState(ctx context.Context) (state *State, e error) {
	var states GetStatesResponse
	if states, e = p.api.ApiGetStates(ctx, true, p.GetPhoneResponse.Tzid); e != nil {
		return
	}
	if state = states.ByTzid(p.GetPhoneResponse.Tzid); state == nil {
		e = ErrorUnexpectedResponse{Response: "no activation found with tzid: " + strconv.Itoa(p.GetPhoneResponse.Tzid)}
	}
	p.State = state
	return
}
