package cmd

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/rockset/rockset-go-client/option"
	"log/slog"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pkg/browser"
	devauth "github.com/rockset/device-authorization"
	"github.com/spf13/cobra"

	"github.com/rockset/cli/config"
	"github.com/rockset/cli/tui"
)

const Auth0ClientID = "0dJNiGWClbLjg7AdtXtAyPCeE0jKOFet"

func newAuthLoginCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "login [NAME CLUSTER ORGANIZATION]",
		Args:  cobra.RangeArgs(0, 3),
		Short: "authenticate using the Rockset console",
		Long:  "authenticate using the Rockset console and save a bearer token on local disk",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if len(args) > 0 && len(args) < 3 {
				return fmt.Errorf("you must provide either zero or three arguments")
			}

			p := devauth.NewProvider("auth0")
			authCfg := p.Config("rockset", Auth0ClientID)
			authCfg.Audience = "https://rockset.sh/"
			authCfg.OAuth2Config.Endpoint.AuthURL = "https://auth.rockset.com/oauth/device/code"
			authCfg.OAuth2Config.Endpoint.TokenURL = "https://auth.rockset.com/oauth/token"
			authCfg.OAuth2Config.Scopes = append(authCfg.OAuth2Config.Scopes, "email")

			a := devauth.NewAuthorizer(authCfg)

			code, err := a.RequestCode(ctx)
			if err != nil {
				slog.Error("failed to request a device code", "err", err)
			}

			fmt.Printf(`Attempting to automatically open the SSO authorization page in your default browser.
If the browser does not open or you wish to use a different device to authorize this request, open the following URL:

%s

Then enter the code:
%s
`, code.VerificationURIComplete, code.UserCode)

			if err := browser.OpenURL(code.VerificationURIComplete); err != nil {
				slog.Warn("could not open", "url", code.VerificationURIComplete, "err", err)
			}

			token, err := a.WaitForAuthorization(ctx, code)
			if err != nil {
				slog.Error("failed to wait for authorization", "err", err)
			}

			fmt.Printf("Successfully logged in!\n\n")

			var name, cluster, org string
			if len(args) > 0 {
				name = args[0]
			}
			if len(args) > 1 {
				cluster = args[1]
			}
			if len(args) > 2 {
				org = args[2]
			}

			model := tui.NewInput("Enter authentication context information", []tui.InputConfig{
				{Placeholder: name, Prompt: "Name: "},
				{Placeholder: cluster, Prompt: "Cluster: ", Validate: func(s string) error {
					for _, c := range config.Clusters {
						if strings.HasPrefix(c, s) {
							return nil
						}
					}
					return fmt.Errorf("cluster must be one of: %s", strings.Join(config.Clusters, ", "))
				}},
				{Placeholder: org, Prompt: "Organization: "},
			})
			input := tea.NewProgram(model)

			if _, err = input.Run(); err != nil {
				return err
			}

			cfg, err := config.Load()
			if err != nil {
				return err
			}

			exp := time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
			if err = cfg.AddToken(model.Fields[0], config.Token{
				Token:      token.IDToken,
				Org:        model.Fields[2],
				Server:     fmt.Sprintf("https://api.%s.rockset.com", model.Fields[1]),
				Expiration: exp,
			}); err != nil {
				return err
			}

			if err = config.Store(cfg); err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "authentication context %s saved (expires in %s)\n",
				model.Fields[0], humanize.Time(exp))
			// should we select the new context, or tell the user how to do it?

			return nil
		},
	}

	return &cmd
}

func newAuthKeyCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "key NAME ROLE",
		Args:  cobra.ExactArgs(2),
		Short: "create an apikey",
		Long:  "create an apikey using the current auth context and save it in the configuration file",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			name := args[0]
			role := args[1]

			rs, err := rockClient(cmd)
			if err != nil {
				return err
			}

			key, err := rs.CreateAPIKey(ctx, "name", option.WithRole(role))
			if err != nil {
				return err
			}

			cfg, err := config.Load()
			if err != nil {
				return err
			}

			cfg.Keys[name] = config.APIKey{
				Key:    key.Key,
				Server: rs.APIServer,
			}
			if err = config.Store(cfg); err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "apikey %s created and saved as authentication context %s\n", name, name)
			// should we select the new context, or tell the user how to do it?

			return nil
		},
	}

	return &cmd
}
