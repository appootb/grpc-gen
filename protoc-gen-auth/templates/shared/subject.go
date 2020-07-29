package shared

import "github.com/appootb/protobuf/go/permission"

type Subjects []permission.Subject

func (s Subjects) Len() int {
	return len(s)
}

func (s Subjects) Less(i int, j int) bool {
	return s[i] < s[j]
}

func (s Subjects) Swap(i int, j int) {
	s[i], s[j] = s[j], s[i]
}
