//go:build integration

package cmd_test

import (
	"fmt"
	"github.com/rockset/cli/cmd"
	"github.com/rockset/cli/internal/test"
	"github.com/rockset/rockset-go-client"
	"github.com/rockset/rockset-go-client/option"
	"github.com/stretchr/testify/suite"
	"testing"
)

type APIKeySuite struct {
	suite.Suite
	rc   *rockset.RockClient
	name string
}

func TestAPIKeySuite(t *testing.T) {
	test.SkipUnlessIntegrationTest(t)

	s := APIKeySuite{
		name: "random",
	}
	suite.Run(t, &s)
}

func (s *APIKeySuite) TearDownSuite() {
	// try to remove the apikey in case it is lingering
	c := cmd.NewRootCmd("test")
	_ = test.Wrapper(s.T(), c, "delete", "apikey", s.name)
	_ = c.Execute()
}

func (s *APIKeySuite) Test_0_Create() {
	c := cmd.NewRootCmd("test")
	out := test.WrapAndExecute(s.T(), c, "create", "apikey", s.name)

	s.Equal(fmt.Sprintf("apikey %s created\n", s.name), out.String())
}

func (s *APIKeySuite) Test_1_Get() {
	c := cmd.NewRootCmd("test")
	out := test.WrapAndExecute(s.T(), c, "get", "apikey", s.name)

	s.NotEmpty(out.String())
}

func (s *APIKeySuite) Test_2_List() {
	c := cmd.NewRootCmd("test")
	out := test.WrapAndExecute(s.T(), c, "list", "apikey")

	s.NotEmpty(out.String())
}

func (s *APIKeySuite) Test_3_Update() {
	c := cmd.NewRootCmd("test")
	out := test.WrapAndExecute(s.T(), c, "update", "apikey", s.name, "--state", string(option.KeySuspended))

	s.Equal(fmt.Sprintf("apikey %s updated to %s\n", s.name, option.KeySuspended), out.String())
}

func (s *APIKeySuite) Test_8_Delete() {
	c := cmd.NewRootCmd("test")
	out := test.WrapAndExecute(s.T(), c, "delete", "apikey", s.name)

	s.Equal(fmt.Sprintf("apikey %s deleted\n", s.name), out.String())
}
