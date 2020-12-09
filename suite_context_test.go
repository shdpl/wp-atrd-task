package main

import (
	"github.com/cucumber/godog"
	"os"
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
	ctx.Step(`^I send a "([^"]*)" request to "([^"]*)"$`, iSendARequestTo)
	ctx.Step(`^I send a "([^"]*)" request to "([^"]*)" with "([^"]*)"$`, iSendARequestToWith)
	ctx.Step(`^the JSON response should contain secret data$`, theJSONResponseShouldContainSecretData)
	ctx.Step(`^the response code should be (\d+)$`, theResponseCodeShouldBe)
}

func iSendARequestTo(arg1, arg2 string) error {
	return godog.ErrPending
}

func iSendARequestToWith(arg1, arg2, arg3 string) error {
	return godog.ErrPending
}

func theJSONResponseShouldContainSecretData() error {
	return godog.ErrPending
}

func theResponseCodeShouldBe(arg1 int) error {
	return godog.ErrPending
}
