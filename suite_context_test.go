package main

import (
	"fmt"
	"github.com/cucumber/godog"
	"github.com/go-openapi/strfmt"
	"github.com/pawmart/wp-atrd-task/api"
	"github.com/pawmart/wp-atrd-task/models"
	"io/ioutil"
	"mime"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	opts := godog.Options{
		Format:    "progress",
		Paths:     []string{"features"},
		Randomize: time.Now().UTC().UnixNano(),
	}

	status := godog.TestSuite{
		Name:                 "wp-atrd-task",
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenario,
		Options:              &opts,
	}.Run()

	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	scenario := NewScenario(ctx)
	ctx.Step(`^I send a "([^"]*)" request to "([^"]*)"$`, scenario.iSendARequestTo)
	ctx.Step(`^I send a "([^"]*)" request to "([^"]*)" with "([^"]*)"$`, scenario.iSendARequestToWith)
	ctx.Step(`^the JSON response should contain secret data$`, scenario.theJSONResponseShouldContainSecretData)
	ctx.Step(`^the response code should be (\d+)$`, scenario.theResponseCodeShouldBe)
}

type Scenario struct {
	api      api.Api
	response *http.Response
}

func NewScenario(ctx *godog.ScenarioContext) (ret *Scenario) {
	ret = &Scenario{}
	ret.api = api.NewApi()
	return ret
}

func (this *Scenario) iSendARequestTo(method, endpoint string) (err error) {
	var request *http.Request
	request, err = http.NewRequest(method, endpoint, nil)
	if err != nil {
		return err
	}
	recorder := httptest.NewRecorder()
	this.api.ServeHTTP(recorder, request)
	this.response = recorder.Result()
	return
}

func (this *Scenario) iSendARequestToWith(method, endpoint, parameters string) (err error) {
	request, err := http.NewRequest(method, endpoint, strings.NewReader(parameters))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Length", strconv.Itoa(len(parameters)))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Accept", "application/json; charset=utf-8")

	recorder := httptest.NewRecorder()
	this.api.ServeHTTP(recorder, request)
	this.response = recorder.Result()
	return err
}

func (this *Scenario) theJSONResponseShouldContainSecretData() (err error) {
	if !isJSONMediaType(this.response.Header) {
		return fmt.Errorf("Returned value is not of JSON media type, but %s", this.response.Header.Get("Content-Type"))
	}

	secret := models.Secret{}

	var buf []byte
	buf, err = ioutil.ReadAll(this.response.Body)
	if err != nil {
		return err
	}

	err = secret.UnmarshalBinary(buf)
	if err != nil {
		return err
	}

	err = secret.Validate(strfmt.Default)
	if err != nil {
		return fmt.Errorf("Invalid secret JSON '%s'; reason: %s", string(buf), err.Error())
	}

	return err
}

func (this *Scenario) theResponseCodeShouldBe(code int) error {
	if code != this.response.StatusCode {
		return fmt.Errorf("Expected response code %d, but received %d", code, this.response.StatusCode)
	}
	return nil
}

///

func isJSONMediaType(header http.Header) bool {
	contentType := header.Get("Content-Type")
	mediaType, _, _ := mime.ParseMediaType(contentType)
	m := strings.TrimPrefix(mediaType, "application/")
	if len(m) == len(mediaType) {
		return false
	}
	// Look for +json suffix. See https://tools.ietf.org/html/rfc6838#section-4.2.8
	// We recognize multiple suffixes too (e.g. application/something+json+other)
	// as that seems to be a possibility.
	for {
		i := strings.Index(m, "+")
		if i == -1 {
			return m == "json"
		}
		if m[0:i] == "json" {
			return true
		}
		m = m[i+1:]
	}
}
