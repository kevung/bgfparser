// Copyright (c) 2024 LeLuxNet
// Licensed under the MIT License
// Original source: https://gitlab.com/LeLuxNet/X/-/blob/c09411c26dfb/encoding/smile/share.go

package smile

type shared []string

func (sPtr *shared) add(val string) {
	s := *sPtr
	if len(s) >= 1024 {
		s = s[:0]
	}
	s = append(s, val)
	*sPtr = s
}
