package format

import "github.com/rockset/rockset-go-client/openapi"

func userHeader() []string {
	return []string{
		"email",
		"first name",
		"last name",
		"organization",
		"state",
		"created at",
	}
}

func user(u openapi.User) []string {
	return []string{
		u.GetEmail(),
		u.GetFirstName(),
		u.GetLastName(),
		u.GetOrg(),
		u.GetState(),
		u.GetCreatedAt(),
	}
}
