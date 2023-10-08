package cmd

import (
	"fmt"
	"github.com/rockset/cli/config"
	"log/slog"
	"time"

	"github.com/pkg/browser"
	devauth "github.com/rockset/device-authorization"
	"github.com/spf13/cobra"
)

const Auth0ClientID = "0dJNiGWClbLjg7AdtXtAyPCeE0jKOFet"

func newAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth NAME CLUSTER ORGANIZATION",
		Args:  cobra.ExactArgs(3),
		Short: "authenticate",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			name := args[0]
			cluster := args[1]
			org := args[2]

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

			fmt.Printf("Successfully logged in!\n")

			fmt.Printf("token:\n%s\n", token.IDToken)

			// TODO save token & expiration in config
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			if err = cfg.AddToken(name, config.Token{
				Token:      token.IDToken,
				Org:        org,
				Server:     cluster,
				Expiration: time.Now().Add(time.Duration(token.ExpiresIn) * time.Second),
			}); err != nil {
				return err
			}

			if err = config.Store(cfg); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
