package cmd

import (
	"fmt"
	devauth "github.com/rockset/device-authorization"
	"github.com/spf13/cobra"
	"log"
)

const Auth0ClientID = "0dJNiGWClbLjg7AdtXtAyPCeE0jKOFet"

func newAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Args:  cobra.NoArgs,
		Short: "authenticate",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			p := devauth.NewProvider("auth0")
			cfg := p.Config("rockset", Auth0ClientID)
			cfg.Audience = "https://rockset.sh/"
			cfg.OAuth2Config.Endpoint.AuthURL = "https://auth.rockset.com/oauth/device/code"
			cfg.OAuth2Config.Endpoint.TokenURL = "https://auth.rockset.com/oauth/token"
			cfg.OAuth2Config.Scopes = append(cfg.OAuth2Config.Scopes, "email")

			a := devauth.NewAuthorizer(cfg)

			code, err := a.RequestCode(ctx)
			if err != nil {
				log.Fatalf("failed to request a device code: %v", err)
			}

			fmt.Printf(`Attempting to automatically open the SSO authorization page in your default browser.
If the browser does not open or you wish to use a different device to authorize this request, open the following URL:

%s

Then enter the code:
%s
`, code.VerificationURIComplete, code.UserCode)

			token, err := a.WaitForAuthorization(ctx, code)
			if err != nil {
				log.Fatalf("failed to wait for authorization: %v", err)
			}

			fmt.Printf("Successfully logged in!\n")

			fmt.Printf("token:\n%s\n", token.AccessToken)

			return nil
		},
	}

	return cmd
}
