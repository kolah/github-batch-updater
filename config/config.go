package config

import "os"

type Config struct {
	GithubPersonalAccessToken string
}

func Load() Config {
	c := Config{
		GithubPersonalAccessToken: os.Getenv("GITHUB_PERSONAL_ACCESS_TOKEN"),
	}

	if c.GithubPersonalAccessToken == "" {
		panic("GITHUB_PERSONAL_ACCESS_TOKEN is required")
	}

	return c
}
