package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/cucumber/godog"
	resty "github.com/go-resty/resty/v2"
	"github.com/ory/dockertest"
	"github.com/zakiafada32/shipping-go/config"
	"github.com/zakiafada32/shipping-go/handlers/rest"
)

type apiFeature struct {
	client   *resty.Client
	server   *httptest.Server
	word     string
	language string
}

var (
	pool     *dockertest.Pool
	database *dockertest.Resource
)

type contextKey string

const wordKey contextKey = "word"
const langKey contextKey = "lang"

func (api *apiFeature) iTranslateItTo(ctx context.Context, lang string) (context.Context, error) {
	return context.WithValue(ctx, langKey, lang), nil
}

func (api *apiFeature) theWord(ctx context.Context, word string) (context.Context, error) {
	return context.WithValue(ctx, wordKey, word), nil
}

func (api *apiFeature) theResponseShouldBe(ctx context.Context, result string) error {
	var ok bool
	api.word, ok = ctx.Value(wordKey).(string)
	if !ok {
		return fmt.Errorf("unable to get language from context")
	}

	api.language, ok = ctx.Value(langKey).(string)
	if !ok {
		return fmt.Errorf("unable to get language from context")
	}

	url := fmt.Sprintf("%s/translate/hello?word=%s", api.server.URL, api.word)

	resp, err := api.client.R().
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"language": api.language,
		}).
		SetResult(&rest.Resp{}).
		Get(url)

	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusNotFound {
		return fmt.Errorf("unable to call api: %d %s", resp.StatusCode(), url)
	}

	res := resp.Result().(*rest.Resp)
	if res.Translation != result {
		return fmt.Errorf("translation should be set to %s", result)
	}

	return nil
}

func InitializeTestSuite(sc *godog.TestSuiteContext) {
	var err error

	sc.BeforeSuite(func() {
		pool, err = dockertest.NewPool("")
		if err != nil {
			panic(fmt.Sprintf("unable to create connection pool %s", err))
		}

		wd, err := os.Getwd()
		if err != nil {
			panic(fmt.Sprintf("unable to get working directory %s", err))
		}

		mount := fmt.Sprintf("%s/data/:/data/", filepath.Dir(wd))
		fmt.Println(mount)
		redis, err := pool.RunWithOptions(&dockertest.RunOptions{
			Repository: "redis",
			Mounts:     []string{mount},
		})
		if err != nil {
			panic(fmt.Sprintf("unable to create container: %s", err))
		}
		if err := redis.Expire(600); err != nil {
			panic("unable to set expiration on container")
		}
		database = redis
	})

	sc.AfterSuite(func() {
		database.Close()
	})
}

func InitializeScenario(ctx *godog.ScenarioContext) {

	client := resty.New()
	api := &apiFeature{
		client: client,
	}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		cfg := config.Configuration{}
		cfg.LoadFromEnv()
		cfg.DatabasePort = database.GetPort("6379/tcp")
		cfg.DatabaseURL = "localhost"
		fmt.Printf("%+v\n", cfg)
		mux := API(cfg)
		server := httptest.NewServer(mux)
		api.server = server
		return ctx, nil
	})
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		api.server.Close()
		return ctx, nil
	})

	ctx.Step(`^I translate it to "([^"]*)"$`, api.iTranslateItTo)
	ctx.Step(`^the word "([^"]*)"$`, api.theWord)
	ctx.Step(`^the response should be "([^"]*)"$`, api.theResponseShouldBe)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer:  InitializeScenario,
		TestSuiteInitializer: InitializeTestSuite,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
