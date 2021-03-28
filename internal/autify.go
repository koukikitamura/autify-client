package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
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

type Scenario struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (a *Autify) FetchScenario(projectId, scenarioId int) (*Scenario, error) {
	path := fmt.Sprintf("projects/%d/scenarios/%d", projectId, scenarioId)

	request, err := http.NewRequest(http.MethodGet, strings.Join([]string{a.baseUrl, path}, "/"), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.accessToken))

	response, err := a.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var scenario Scenario
	if err := json.NewDecoder(response.Body).Decode(&scenario); err != nil {
		return nil, err
	}

	return &scenario, nil
}

type Result struct {
	Id         int       `json:"id"`
	Status     string    `json:"status"`
	Duration   int       `json:"duration"`
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	TestPlan   `json:"test_plan"`
}

type TestPlan struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const (
	ResultStatusWarning = "warning"
	ResultStatusRunning = "running"
	ResultStatusPassed  = "passed"
	ResultStatusFailed  = "failed"
)

func (a *Autify) FetchResult(projectId, resultId int) (*Result, error) {
	path := fmt.Sprintf("/projects/%d/results/%d", projectId, resultId)

	request, err := http.NewRequest(http.MethodGet, strings.Join([]string{a.baseUrl, path}, "/"), nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.accessToken))

	response, err := a.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var result Result
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func CheckAccessToken() bool {
	return len(GetAccessToken()) > 0
}

func GetAccessToken() string {
	return os.Getenv(AccessTokenEnvName)
}
