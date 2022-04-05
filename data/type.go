package data

const (
	TypeInMemory Type = "in_memory"
	TypePostgres Type = "postgres"
)

type Type string

func (t Type) String() string {
	return string(t)
}
