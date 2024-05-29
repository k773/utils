package capsolvercom

import (
	"context"
	"encoding/json"
	"github.com/k773/utils"
	"github.com/k773/utils/fixedPoint"
	"time"
)

const apiEndpoint = "https://api.capsolver.com"

func (p *Provider) BalanceUsd(ctx context.Context) (balance fixedPoint.FP, e error) {
	res, e := p.makeRequest(ctx, "/getBalance", BaseTask{})
	if e != nil {
		return
	}
	return res.Balance, nil
}

func (p *Provider) Solve(ctx context.Context, task any) (solution Solution, e error) {
	res, e := p.makeRequest(ctx, "/createTask", BaseTask{Task: &task})
	if e != nil {
		return
	}
	if res.isReady() {
		solution = res.Solution
		return
	}
	return p.wait(ctx, res.TaskId)
}

func (p *Provider) wait(ctx context.Context, taskId string) (solution Solution, e error) {
	const iterEvery = 5 * time.Second
	var res CapSolverResponse
	for i := 0; e == nil && res.Status != "ready"; i++ {
		if e = utils.SleepWithContext(ctx, iterEvery); e != nil {
			continue
		}
		if res, e = p.makeRequest(ctx, "/getTaskResult", BaseTask{TaskId: taskId}); e != nil {
			continue
		}
		if res.isReady() {
			break
		}
	}
	if e == nil {
		solution = res.Solution
	}
	return
}

func (p *Provider) makeRequest(ctx context.Context, path string, baseTask BaseTask) (res CapSolverResponse, e error) {
	baseTask.ClientKey = p.ApiKey

	r, e := p.S.R().
		SetContext(ctx).
		SetBody(baseTask).
		Post(apiEndpoint + path)
	if e != nil {
		return
	}
	if e = json.Unmarshal(r.Body(), &res); e != nil {
		return
	}
	if res.ErrorId == 1 {
		e = &res.CapSolverResponseError
		return
	}
	return
}
