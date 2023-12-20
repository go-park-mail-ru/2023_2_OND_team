package share

import "errors"

var ErrUnknownRole = errors.New("unknown role")

func getContributorRoleFromString(role string) (int, bool) {
	switch role {
	case "read-write":
		return 1, true
	case "read-only":
		return 2, true
	default:
		return 0, false
	}
}

func getStringContributorRoleFromInt(role int) string {
	switch role {
	case 1:
		return "read-write"
	case 2:
		return "read-only"
	default:
		return ""
	}
}
