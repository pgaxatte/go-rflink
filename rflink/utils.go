package rflink

import "strconv"

// strToUint16 parses a string directly into an uint16  with the specified base
func strToUint16(s string, base int) (uint16, error) {
	u, err := strconv.ParseUint(s, base, 16)
	if err != nil {
		return 0, err
	}
	return uint16(u), nil
}
