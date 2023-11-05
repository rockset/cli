//go:build integration

package cmd_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/rockset/cli/cmd"
	"github.com/rockset/cli/internal/test"
)

type AliasTestSuite struct {
	suite.Suite
	name string
}

func TestAliasSuite(t *testing.T) {
	test.SkipUnlessIntegrationTest(t)

	s := AliasTestSuite{name: "dummy"} // TODO use random name
	suite.Run(t, &s)
}

func (s *AliasTestSuite) Test_0_Create() {
	c := cmd.NewRootCmd("test")
	out := test.WrapAndExecute(s.T(), c, "create", "alias", s.name, "commons._events")

	s.Equal(fmt.Sprintf("alias %s created\n", s.name), out.String())
}

func (s *AliasTestSuite) Test_1_Get() {
	c := cmd.NewRootCmd("test")
	out := test.WrapAndExecute(s.T(), c, "get", "alias", s.name)

	s.NotEmpty(out.String())
}

func (s *AliasTestSuite) Test_2_List() {
	c := cmd.NewRootCmd("test")
	out := test.WrapAndExecute(s.T(), c, "list", "aliases")

	s.NotEmpty(out.String())
}

func (s *AliasTestSuite) Test_3_Update() {
	c := cmd.NewRootCmd("test")
	out := test.WrapAndExecute(s.T(), c, "update", "alias", s.name, "commons._events")

	s.Equal(fmt.Sprintf("alias %s updated\n", s.name), out.String())
}

func (s *AliasTestSuite) Test_4_Delete() {
	c := cmd.NewRootCmd("test")
	out := test.WrapAndExecute(s.T(), c, "delete", "alias", s.name)

	s.Equal(fmt.Sprintf("alias %s deleted\n", s.name), out.String())
}
