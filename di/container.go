package di

import (
	"context"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/charmbracelet/log"
	"github.com/google/go-github/v60/github"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/kolah/github-batch-updater/batch"
	"github.com/kolah/github-batch-updater/config"
	"golang.org/x/oauth2"
)

// safeLazyValue is a thread safe lazy loader.
// It is used to create components only when they are needed.
type safeLazyValue[T any] struct {
	once  sync.Once    //nolint: structcheck
	value atomic.Value //nolint: structcheck
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

func (c *Container) LoadBatch(filename string) (*batch.Config, error) {
	k := koanf.New(".")
	if err := k.Load(file.Provider(filename), yaml.Parser()); err != nil {
		return nil, err
	}

	batchConfig := &batch.Config{}
	if err := k.Unmarshal("batch", batchConfig); err != nil {
		return nil, err
	}

	return batchConfig, nil
}
