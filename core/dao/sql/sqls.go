package sql

var (
	branchSql string
	global    string
)

func Sql() []string {
	return []string{branchSql, global}
}
