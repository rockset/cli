package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"
	"sort"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/rockset/rockset-go-client"
	"gopkg.in/yaml.v3"
)

const (
	FileName        = "config.yaml"
	HistoryFileName = "cli.hist"
)

var Clusters = []string{"usw2a1", "use1a1", "euc1a1", "apt2a1"}

func init() {
	sort.Strings(Clusters)
}

var File string

func init() {
	file, err := configPath()
	if err != nil {
		panic(fmt.Sprintf("unable to locate config file %s: %v", FileName, err))
	}
	File = file
}

type Config struct {
	Current string            `yaml:"current"`
	Keys    map[string]APIKey `yaml:"keys"`
	Tokens  map[string]Token  `yaml:"tokens"`
}

func New() Config {
	return Config{
		Keys:   make(map[string]APIKey),
		Tokens: make(map[string]Token),
	}
}

func (c *Config) AsOptions(override string) ([]rockset.RockOption, error) {
	name := override
	if c.Current == "" && override == "" {
		return nil, NoSelectionErr
	} else if c.Current != "" {
		name = c.Current
	}

	if key, found := c.Keys[name]; found {
		return []rockset.RockOption{
			rockset.WithAPIServer(key.Server),
			rockset.WithAPIKey(key.Key),
		}, nil
	}

	if token, found := c.Tokens[name]; found {
		if time.Now().After(token.Expiration) {
			return nil, TokenExpiredErr
		}

		return []rockset.RockOption{
			rockset.WithAPIServer(token.Server),
			rockset.WithBearerToken(token.Token, token.Org),
		}, nil
	}

	return nil, fmt.Errorf("%w", NotFoundErr)
}

func (c *Config) DeleteContext(name string) error {
	if _, found := c.Tokens[name]; found {
		delete(c.Tokens, name)
		return nil
	}

	if _, found := c.Keys[name]; found {
		delete(c.Keys, name)
		return nil
	}

	return fmt.Errorf("context %s: %w", name, NotFoundErr)
}

func (c *Config) AddToken(name string, token Token) error {
	if _, found := c.Tokens[name]; found {
		return ContextAlreadyExistErr
	}

	c.Tokens[name] = token
	return nil
}

func (c *Config) AddKey(name string, key APIKey) error {
	if _, found := c.Keys[name]; found {
		return ContextAlreadyExistErr
	}

	c.Keys[name] = key
	return nil
}

func (c *Config) Use(name string) error {
	if _, found := c.Keys[name]; found {
		c.Current = name
		return nil
	}

	if _, found := c.Tokens[name]; found {
		c.Current = name
		return nil
	}

	return NotFoundErr
}

var (
	NoSelectionErr         = errors.New("no context selected")
	NotFoundErr            = errors.New("context not found")
	ContextAlreadyExistErr = errors.New("context already exists")
	TokenExpiredErr        = errors.New("token expired")
)

type APIKey struct {
	Key    string `yaml:"apikey"`
	Server string `yaml:"apiserver"`
}

func (a APIKey) APIServer() string { return a.Server }

type Token struct {
	Token      string    `yaml:"token"`
	Org        string    `yaml:"org"`
	Server     string    `yaml:"apiserver"`
	Expiration time.Time `yaml:"expiration"`
}

func (t Token) APIServer() string { return t.Server }

func configPath() (string, error) {
	return rocksetConfigDir(FileName)
}

func HistoryFile() (string, error) {
	return rocksetConfigDir(HistoryFileName)
}

func rocksetConfigDir(name string) (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	return path.Join(home, ".config", "rockset", name), nil
}

func Store(cfg Config) error {
	cfgPath, err := configPath()
	if err != nil {
		return err
	}

	return StoreFile(cfg, cfgPath)
}

func StoreFile(cfg Config, cfgPath string) error {
	dir := path.Dir(cfgPath)
	err := os.MkdirAll(dir, 0700)

	f, err := os.OpenFile(cfgPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	enc := yaml.NewEncoder(f)
	if err = enc.Encode(cfg); err != nil {
		return err
	}

	if err = enc.Close(); err != nil {
		slog.Error("failed to close config", "err", err)
	}

	return nil
}

// Load loads the CLI configuration, and if the config doesn't exist, it returns an empty config.
func Load() (Config, error) {
	cfg := New()

	cfgPath, err := configPath()
	if err != nil {
		return cfg, err
	}

	return LoadFile(cfgPath)
}

func LoadFile(cfgPath string) (Config, error) {
	cfg := New()

	f, err := os.Open(cfgPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("failed to read apikey config file: %w", err)
	}

	dec := yaml.NewDecoder(f)
	err = dec.Decode(&cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
