package cmd

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
)

type AliasTestSuite struct {
	suite.Suite
	name string
}

func TestAliasSuite(t *testing.T) {
	s := AliasTestSuite{name: "dummy"}
	suite.Run(t, &s)
}

func (s *AliasTestSuite) Test_0_Create() {
	cmd := newCreateAliasCmd()
	out := testWrapper(cmd, []string{s.name, "commons._events"})
	err := cmd.Execute()
	s.Require().NoError(err)
	s.Equal(fmt.Sprintf("alias %s created", s.name), out.String())
}

func (s *AliasTestSuite) Test_1_Get() {
	cmd := newGetAliasCmd()
	out := testWrapper(cmd, []string{s.name})

	err := cmd.Execute()
	s.Require().NoError(err)
	s.NotEmpty(out.String())
}

func (s *AliasTestSuite) Test_2_List() {
	cmd := newListAliasesCmd()
	out := testWrapper(cmd, []string{})

	err := cmd.Execute()
	s.Require().NoError(err)
	s.NotEmpty(out.String())
}

func (s *AliasTestSuite) Test_3_Update() {
	cmd := newUpdateAliasCmd()
	out := testWrapper(cmd, []string{s.name, "commons._events"})
	err := cmd.Execute()
	s.Require().NoError(err)
	s.NotEmpty(fmt.Sprintf("alias %s deleted", s.name), out.String())
}

func (s *AliasTestSuite) Test_4_Delete() {
	cmd := newDeleteAliasCmd()
	out := testWrapper(cmd, []string{s.name})
	err := cmd.Execute()
	s.Require().NoError(err)
	s.NotEmpty(fmt.Sprintf("alias %s deleted", s.name), out.String())
}
