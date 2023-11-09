package lookup

import (
	"context"
	"errors"
	"fmt"
	"github.com/rockset/rockset-go-client"
	"regexp"
)

func VirtualInstanceNameOrIDtoID(ctx context.Context, rs *rockset.RockClient, nameOrID string) (string, error) {
	if !isUUID(nameOrID) {
		id, err := viNameToID(ctx, rs, nameOrID)
		if err != nil {
			return "", fmt.Errorf("failed to get virtual instance id for %s: %v", nameOrID, err)
		}

		return id, nil
	}

	return nameOrID, nil
}

var uuidRe = regexp.MustCompile(`[[:xdigit:]]{8}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{4-[[:xdigit:]]{12`)

func isUUID(id string) bool {
	return uuidRe.MatchString(id)
}

func viNameToID(ctx context.Context, rs *rockset.RockClient, name string) (string, error) {
	vis, err := rs.ListVirtualInstances(ctx)
	if err != nil {
		return "", err
	}

	for _, vi := range vis {
		if vi.GetName() == name {
			return vi.GetId(), nil
		}
	}

	return "", VINotFoundErr
}

var VINotFoundErr = errors.New("virtual instance not found")
