package config

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/planetscale/cli/psapi"
)

const (
	defaultConfigPath = "~/.config/psctl"
)

type Config struct {
	AccessToken string
	BaseURL     string
}

func New() *Config {
	var accessToken []byte
	_, err := os.Stat(AccessTokenPath())
	if err != nil {
		if !os.IsNotExist(err) {
			// TODO(iheanyi): Is fatal the right move here?
			log.Fatal(err)
		}
	} else {
		accessToken, err = ioutil.ReadFile(AccessTokenPath())
		if err != nil {
			log.Fatal(err)
		}
	}

	return &Config{
		AccessToken: string(accessToken),
		BaseURL:     psapi.DefaultBaseURL,
	}
}

func (c *Config) IsAuthenticated() bool {
	return c.AccessToken != ""
}

// ConfigDir is the directory for psctl config.
func ConfigDir() string {
	dir, _ := homedir.Expand(defaultConfigPath)
	return dir
}

// AccessTokenPath is the path for the access token file
func AccessTokenPath() string {
	return path.Join(ConfigDir(), "access-token")
}

// NewClientFromConfig creates a PlaentScale API client from our configuration
func (c *Config) NewClientFromConfig(opts ...psapi.ClientOption) (*psapi.Client, error) {
	args := []psapi.ClientOption{psapi.SetBaseURL(c.BaseURL)}
	args = append(args, opts...)
	return psapi.NewClientFromToken(c.AccessToken, args...)
}