package di

import (
	"context"
	"net/http"
	"os"
	"sync"
	"sync/atomic"

	"github.com/charmbracelet/log"
	"github.com/google/go-github/v60/github"
	"github.com/kolah/github-batch-updater/internal/application"
	"github.com/kolah/github-batch-updater/internal/batch"
	"github.com/kolah/github-batch-updater/internal/config"
	"github.com/kolah/github-batch-updater/internal/domain/transformer"
	"github.com/kolah/github-batch-updater/internal/infrastructure"
	transformerImplementation "github.com/kolah/github-batch-updater/internal/infrastructure/transformer"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v3"
)

type safeLazyValue[T any] struct {
	once  sync.Once
	value atomic.Value
}

func (lv *safeLazyValue[T]) loadOrCreate(create func() T) T {
	if value := lv.value.Load(); value != nil {
		return value.(T)
	}

	lv.once.Do(func() {
		lv.value.Store(create())
	})

	return lv.value.Load().(T)
}

type Container struct {
	config                     safeLazyValue[config.Config]
	logger                     safeLazyValue[*log.Logger]
	githubAuthorizedHTTPClient safeLazyValue[*http.Client]
	githubRESTClient           safeLazyValue[*github.Client]
	githubRepositoryService    safeLazyValue[*infrastructure.GithubRepositoryService]
	githubContentService       safeLazyValue[*infrastructure.GithubContentService]
	transformerService         safeLazyValue[*transformer.Default]
	githubPullRequestService   safeLazyValue[*infrastructure.GithubPullRequestService]
	batchProcessingService     safeLazyValue[*application.BatchProcessingService]
}

func NewContainer() *Container {
	return &Container{}
}

func (c *Container) Config() config.Config {
	return c.config.loadOrCreate(config.Load)
}

func (c *Container) Logger() *log.Logger {
	return c.logger.loadOrCreate(func() *log.Logger {
		return log.Default()
	})
}

func (c *Container) GithubAuthorizedHTTPClient() *http.Client {
	return c.githubAuthorizedHTTPClient.loadOrCreate(func() *http.Client {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{
				AccessToken: c.Config().GithubPersonalAccessToken,
			},
		)

		return oauth2.NewClient(ctx, ts)
	})
}

func (c *Container) GithubRESTClient() *github.Client {
	return c.githubRESTClient.loadOrCreate(func() *github.Client {
		return github.NewClient(c.GithubAuthorizedHTTPClient())
	})
}

func (c *Container) GithubRepositoryService() *infrastructure.GithubRepositoryService {
	return c.githubRepositoryService.loadOrCreate(func() *infrastructure.GithubRepositoryService {
		return infrastructure.NewGithubRepositoryService(
			c.GithubRESTClient(),
		)
	})
}

func (c *Container) GithubContentService() *infrastructure.GithubContentService {
	return c.githubContentService.loadOrCreate(func() *infrastructure.GithubContentService {
		return infrastructure.NewGithubContentService(
			c.GithubRESTClient(),
		)
	})
}

func (c *Container) TransformerService() *transformer.Default {
	return c.transformerService.loadOrCreate(func() *transformer.Default {
		t := transformer.NewDefault()
		t.RegisterProcessor(transformerImplementation.NewReplaceProcessor())
		t.RegisterProcessor(transformerImplementation.NewYQProcessor())

		return t
	})
}

func (c *Container) GithubPullRequestService() *infrastructure.GithubPullRequestService {
	return c.githubPullRequestService.loadOrCreate(func() *infrastructure.GithubPullRequestService {
		return infrastructure.NewGithubPullRequestService(
			c.GithubRESTClient(),
		)
	})
}

func (c *Container) BatchProcessingService() *application.BatchProcessingService {
	return c.batchProcessingService.loadOrCreate(func() *application.BatchProcessingService {
		return application.NewBatchProcessingService(
			c.GithubRepositoryService(),
			c.GithubContentService(),
			c.TransformerService(),
			c.GithubPullRequestService(),
			c.Logger(),
		)
	})
}

func (c *Container) LoadBatch(filename string) (*batch.Config, error) {
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	batchConfig := &batch.Config{}

	err = yaml.Unmarshal(yamlFile, batchConfig)
	if err != nil {
		return nil, err
	}

	return batchConfig, nil
}
