package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const AccessTokenEnvName = "AUTIFY_PERSONAL_ACCESS_TOKEN"
const defaultBaseUrl = "https://app.autify.com/api/v1"

type Autify struct {
	httpClient  *http.Client
	accessToken string
	baseUrl     string
}

type AutifyOption func(a *Autify)

func AutifyOptionBaseUrl(baseUrl string) AutifyOption {
	return func(autify *Autify) {
		autify.baseUrl = baseUrl
	}
}

func AutifyOptionHTTPClient(c *http.Client) AutifyOption {
	return func(autify *Autify) {
		autify.httpClient = c
	}
}

func NewAutfiy(accessToken string, options ...AutifyOption) *Autify {
	autify := &Autify{
		accessToken: accessToken,
		httpClient:  http.DefaultClient,
		baseUrl:     defaultBaseUrl,
	}

	for _, option := range options {
		option(autify)
	}

	return autify
}

type RunResult struct {
	TestPlanResultId string `json:"id"`
	Type             string `json:"type"`
	Attributes       struct {
		Id int `json:"id"`
	} `json:"attributes"`
}

type Scenario struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TestPlan struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TestPlanResult struct {
	Id         int       `json:"id"`
	Status     string    `json:"status"`
	Duration   int       `json:"duration"`
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	TestPlan   `json:"test_plan"`
}

const (
	TestPlanStatusWarning = "warning"
	TestPlanStatusQueuing = "queuing"
	TestPlanStatuWaiting  = "waiting"
	TestPlanStatusRunning = "running"
	TestPlanStatusPassed  = "passed"
	TestPlanStatusFailed  = "failed"
)

type RuntTestPlanResponse struct {
	Data RunResult `json:"data"`
}

func (a *Autify) RunTestPlan(planId int) (*RunResult, error) {
	path := fmt.Sprintf("/schedules/%d", planId)

	request, err := http.NewRequest(http.MethodPost, strings.Join([]string{a.baseUrl, path}, "/"), nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.accessToken))

	logrus.WithFields(logrus.Fields{
		"method": request.Method,
		"url":    request.URL,
	}).Debug("Request to autify")

	response, err := a.httpClient.Do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	logrus.WithFields(logrus.Fields{
		"method": response.Request.Method,
		"url":    response.Request.URL,
		"status": response.Status,
	}).Debug("Respond from autify")
	logrus.Debug("Body is ", string(body))

	if response.StatusCode == http.StatusUnauthorized {
		return nil, errors.New("Unauthorized: Bad credentials")
	}

	var result RuntTestPlanResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, errors.WithStack(err)
	}

	return &result.Data, nil
}

func (a *Autify) FetchScenario(projectId, scenarioId int) (*Scenario, error) {
	path := fmt.Sprintf("projects/%d/scenarios/%d", projectId, scenarioId)

	request, err := http.NewRequest(http.MethodGet, strings.Join([]string{a.baseUrl, path}, "/"), nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.accessToken))

	logrus.WithFields(logrus.Fields{
		"method": request.Method,
		"url":    request.URL,
	}).Debug("Request to autify")

	response, err := a.httpClient.Do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	logrus.WithFields(logrus.Fields{
		"method": response.Request.Method,
		"url":    response.Request.URL,
		"status": response.Status,
	}).Debug("Respond from autify")
	logrus.Debug("Body is ", string(body))

	if response.StatusCode == http.StatusUnauthorized {
		return nil, errors.New("Unauthorized: Bad credentials")
	}

	var scenario Scenario
	if err := json.Unmarshal(body, &scenario); err != nil {
		return nil, errors.WithStack(err)
	}

	return &scenario, nil
}

func (a *Autify) FetchResult(projectId, resultId int) (*TestPlanResult, error) {
	path := fmt.Sprintf("/projects/%d/results/%d", projectId, resultId)

	request, err := http.NewRequest(http.MethodGet, strings.Join([]string{a.baseUrl, path}, "/"), nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.accessToken))

	logrus.WithFields(logrus.Fields{
		"method": request.Method,
		"url":    request.URL,
	}).Debug("Request to autify")

	response, err := a.httpClient.Do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	logrus.WithFields(logrus.Fields{
		"method": response.Request.Method,
		"url":    response.Request.URL,
		"status": response.Status,
	}).Debug("Respond from autify")
	logrus.Debug("Body is ", string(body))

	if response.StatusCode == http.StatusUnauthorized {
		return nil, errors.New("Unauthorized: Bad credentials")
	}

	var result TestPlanResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, errors.WithStack(err)
	}

	return &result, nil
}

func CheckAccessToken() bool {
	return len(GetAccessToken()) > 0
}

func GetAccessToken() string {
	return os.Getenv(AccessTokenEnvName)
}
