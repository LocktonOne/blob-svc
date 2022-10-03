package types

import "regexp"

var AddressRegexp = regexp.MustCompile("^(0x)?[0-9a-fA-F]{40}$")
