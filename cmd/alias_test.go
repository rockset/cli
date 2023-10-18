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
	c := cmd.NewCreateAliasCmd()
	out := test.Wrapper(c, []string{s.name, "commons._events"})
	err := c.Execute()
	s.Require().NoError(err)
	s.Equal(fmt.Sprintf("alias %s created", s.name), out.String())
}

func (s *AliasTestSuite) Test_1_Get() {
	c := cmd.NewGetAliasCmd()
	out := test.Wrapper(c, []string{s.name})

	err := c.Execute()
	s.Require().NoError(err)
	s.NotEmpty(out.String())
}

func (s *AliasTestSuite) Test_2_List() {
	c := cmd.NewListAliasesCmd()
	out := test.Wrapper(c, []string{})

	err := c.Execute()
	s.Require().NoError(err)
	s.NotEmpty(out.String())
}

func (s *AliasTestSuite) Test_3_Update() {
	c := cmd.NewUpdateAliasCmd()
	out := test.Wrapper(c, []string{s.name, "commons._events"})
	err := c.Execute()
	s.Require().NoError(err)
	s.NotEmpty(fmt.Sprintf("alias %s deleted", s.name), out.String())
}

func (s *AliasTestSuite) Test_4_Delete() {
	c := cmd.NewDeleteAliasCmd()
	out := test.Wrapper(c, []string{s.name})
	err := c.Execute()
	s.Require().NoError(err)
	s.NotEmpty(fmt.Sprintf("alias %s deleted", s.name), out.String())
}
