//go:build integration

package cmd_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/rockset/cli/cmd"
	"github.com/rockset/cli/internal/test"
)

type QueryLambdaSuite struct {
	suite.Suite
	name    string
	version string
}

func TestQueryLambda(t *testing.T) {
	s := QueryLambdaSuite{name: "dummy"}
	suite.Run(t, &s)
}

func (s *QueryLambdaSuite) Test_0_Create() {
	c := cmd.NewRootCmd("test")
	out := test.WrapAndExecute(s.T(), c, "create", "ql", "--sql", "testdata/query_lambda.sql", "--wait", s.name)

	re := regexp.MustCompile(`created query lambda commons\.\w+:(\w+)`)
	m := re.FindStringSubmatch(out.String())
	s.NotNil(m)
	s.version = m[1]
}

func (s *QueryLambdaSuite) Test_1_Get() {
	c := cmd.NewRootCmd("test")
	out := test.WrapAndExecute(s.T(), c, "get", "ql", "--version", s.version, s.name)

	s.NotEmpty(out.String())
}

func (s *QueryLambdaSuite) Test_2_List() {
	c := cmd.NewRootCmd("test")
	out := test.WrapAndExecute(s.T(), c, "list", "ql")

	s.NotEmpty(out.String())
}

func (s *QueryLambdaSuite) Test_3_Execute() {
	c := cmd.NewRootCmd("test")
	out := test.WrapAndExecute(s.T(), c, "execute", "ql", "--version", s.version, s.name)

	s.NotEmpty(out.String())
}

func (s *QueryLambdaSuite) Test_4_Update() {
	c := cmd.NewRootCmd("test")
	out := test.WrapAndExecute(s.T(), c, "update", "ql", "--sql", "testdata/query_lambda_updated.sql", "--wait", s.name)

	re := regexp.MustCompile(`updated query lambda commons\.\w+:(\w+)`)
	m := re.FindStringSubmatch(out.String())
	s.NotNil(m)
	s.version = m[1]
}

func (s *QueryLambdaSuite) Test_5_Execute() {
	c := cmd.NewRootCmd("test")
	out := test.WrapAndExecute(s.T(), c, "execute", "ql", "--version", s.version, s.name)

	s.NotEmpty(out.String())
}

// TODO test invoking using a tag too

func (s *QueryLambdaSuite) Test_6_Delete() {
	c := cmd.NewRootCmd("test")
	out := test.WrapAndExecute(s.T(), c, "delete", "ql", s.name)

	s.Equal(fmt.Sprintf("deleted query lambda %s in commons\n", s.name), out.String())
}
