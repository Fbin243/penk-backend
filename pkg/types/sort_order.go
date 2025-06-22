package types

type SortOrder string

const (
	ASC  SortOrder = "ASC"
	DESC SortOrder = "DESC"
)

func (s SortOrder) ToInt() int {
	switch s {
	case ASC:
		return 1
	case DESC:
		return -1
	default:
		return 0
	}
}
