package rflink

import "strconv"

func strToUint16(s string, base int) (uint16, error) {
	u, err := strconv.ParseUint(s, base, 16)
	if err != nil {
		return 0, err
	}
	return uint16(u), nil
}
