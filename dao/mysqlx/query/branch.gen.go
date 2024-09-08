// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"github.com/tiny-sky/Tdtm/entity"
)

func newBranch(db *gorm.DB, opts ...gen.DOOption) branch {
	_branch := branch{}

	_branch.branchDo.UseDB(db, opts...)
	_branch.branchDo.UseModel(&entity.Branch{})

	tableName := _branch.branchDo.TableName()
	_branch.ALL = field.NewAsterisk(tableName)
	_branch.GID = field.NewString(tableName, "g_id")
	_branch.BranchId = field.NewString(tableName, "branch_id")
	_branch.Url = field.NewString(tableName, "url")
	_branch.ReqData = field.NewString(tableName, "req_data")
	_branch.ReqHeader = field.NewString(tableName, "req_header")
	_branch.TranType = field.NewString(tableName, "tran_type")
	_branch.Protocol = field.NewString(tableName, "protocol")
	_branch.Action = field.NewString(tableName, "action")
	_branch.State = field.NewString(tableName, "state")
	_branch.Level = field.NewUint32(tableName, "level")
	_branch.LastErrMsg = field.NewString(tableName, "last_err_msg")
	_branch.Timeout = field.NewInt64(tableName, "timeout")
	_branch.CreateTime = field.NewInt64(tableName, "create_time")
	_branch.UpdateTime = field.NewInt64(tableName, "update_time")

	_branch.fillFieldMap()

	return _branch
}

type branch struct {
	branchDo branchDo

	ALL        field.Asterisk
	GID        field.String
	BranchId   field.String
	Url        field.String
	ReqData    field.String
	ReqHeader  field.String
	TranType   field.String
	Protocol   field.String
	Action     field.String
	State      field.String
	Level      field.Uint32
	LastErrMsg field.String
	Timeout    field.Int64
	CreateTime field.Int64
	UpdateTime field.Int64

	fieldMap map[string]field.Expr
}

func (b branch) Table(newTableName string) *branch {
	b.branchDo.UseTable(newTableName)
	return b.updateTableName(newTableName)
}

func (b branch) As(alias string) *branch {
	b.branchDo.DO = *(b.branchDo.As(alias).(*gen.DO))
	return b.updateTableName(alias)
}

func (b *branch) updateTableName(table string) *branch {
	b.ALL = field.NewAsterisk(table)
	b.GID = field.NewString(table, "g_id")
	b.BranchId = field.NewString(table, "branch_id")
	b.Url = field.NewString(table, "url")
	b.ReqData = field.NewString(table, "req_data")
	b.ReqHeader = field.NewString(table, "req_header")
	b.TranType = field.NewString(table, "tran_type")
	b.Protocol = field.NewString(table, "protocol")
	b.Action = field.NewString(table, "action")
	b.State = field.NewString(table, "state")
	b.Level = field.NewUint32(table, "level")
	b.LastErrMsg = field.NewString(table, "last_err_msg")
	b.Timeout = field.NewInt64(table, "timeout")
	b.CreateTime = field.NewInt64(table, "create_time")
	b.UpdateTime = field.NewInt64(table, "update_time")

	b.fillFieldMap()

	return b
}

func (b *branch) WithContext(ctx context.Context) *branchDo { return b.branchDo.WithContext(ctx) }

func (b branch) TableName() string { return b.branchDo.TableName() }

func (b branch) Alias() string { return b.branchDo.Alias() }

func (b branch) Columns(cols ...field.Expr) gen.Columns { return b.branchDo.Columns(cols...) }

func (b *branch) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := b.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (b *branch) fillFieldMap() {
	b.fieldMap = make(map[string]field.Expr, 14)
	b.fieldMap["g_id"] = b.GID
	b.fieldMap["branch_id"] = b.BranchId
	b.fieldMap["url"] = b.Url
	b.fieldMap["req_data"] = b.ReqData
	b.fieldMap["req_header"] = b.ReqHeader
	b.fieldMap["tran_type"] = b.TranType
	b.fieldMap["protocol"] = b.Protocol
	b.fieldMap["action"] = b.Action
	b.fieldMap["state"] = b.State
	b.fieldMap["level"] = b.Level
	b.fieldMap["last_err_msg"] = b.LastErrMsg
	b.fieldMap["timeout"] = b.Timeout
	b.fieldMap["create_time"] = b.CreateTime
	b.fieldMap["update_time"] = b.UpdateTime
}

func (b branch) clone(db *gorm.DB) branch {
	b.branchDo.ReplaceConnPool(db.Statement.ConnPool)
	return b
}

func (b branch) replaceDB(db *gorm.DB) branch {
	b.branchDo.ReplaceDB(db)
	return b
}

type branchDo struct{ gen.DO }

func (b branchDo) Debug() *branchDo {
	return b.withDO(b.DO.Debug())
}

func (b branchDo) WithContext(ctx context.Context) *branchDo {
	return b.withDO(b.DO.WithContext(ctx))
}

func (b branchDo) ReadDB() *branchDo {
	return b.Clauses(dbresolver.Read)
}

func (b branchDo) WriteDB() *branchDo {
	return b.Clauses(dbresolver.Write)
}

func (b branchDo) Session(config *gorm.Session) *branchDo {
	return b.withDO(b.DO.Session(config))
}

func (b branchDo) Clauses(conds ...clause.Expression) *branchDo {
	return b.withDO(b.DO.Clauses(conds...))
}

func (b branchDo) Returning(value interface{}, columns ...string) *branchDo {
	return b.withDO(b.DO.Returning(value, columns...))
}

func (b branchDo) Not(conds ...gen.Condition) *branchDo {
	return b.withDO(b.DO.Not(conds...))
}

func (b branchDo) Or(conds ...gen.Condition) *branchDo {
	return b.withDO(b.DO.Or(conds...))
}

func (b branchDo) Select(conds ...field.Expr) *branchDo {
	return b.withDO(b.DO.Select(conds...))
}

func (b branchDo) Where(conds ...gen.Condition) *branchDo {
	return b.withDO(b.DO.Where(conds...))
}

func (b branchDo) Order(conds ...field.Expr) *branchDo {
	return b.withDO(b.DO.Order(conds...))
}

func (b branchDo) Distinct(cols ...field.Expr) *branchDo {
	return b.withDO(b.DO.Distinct(cols...))
}

func (b branchDo) Omit(cols ...field.Expr) *branchDo {
	return b.withDO(b.DO.Omit(cols...))
}

func (b branchDo) Join(table schema.Tabler, on ...field.Expr) *branchDo {
	return b.withDO(b.DO.Join(table, on...))
}

func (b branchDo) LeftJoin(table schema.Tabler, on ...field.Expr) *branchDo {
	return b.withDO(b.DO.LeftJoin(table, on...))
}

func (b branchDo) RightJoin(table schema.Tabler, on ...field.Expr) *branchDo {
	return b.withDO(b.DO.RightJoin(table, on...))
}

func (b branchDo) Group(cols ...field.Expr) *branchDo {
	return b.withDO(b.DO.Group(cols...))
}

func (b branchDo) Having(conds ...gen.Condition) *branchDo {
	return b.withDO(b.DO.Having(conds...))
}

func (b branchDo) Limit(limit int) *branchDo {
	return b.withDO(b.DO.Limit(limit))
}

func (b branchDo) Offset(offset int) *branchDo {
	return b.withDO(b.DO.Offset(offset))
}

func (b branchDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *branchDo {
	return b.withDO(b.DO.Scopes(funcs...))
}

func (b branchDo) Unscoped() *branchDo {
	return b.withDO(b.DO.Unscoped())
}

func (b branchDo) Create(values ...*entity.Branch) error {
	if len(values) == 0 {
		return nil
	}
	return b.DO.Create(values)
}

func (b branchDo) CreateInBatches(values []*entity.Branch, batchSize int) error {
	return b.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (b branchDo) Save(values ...*entity.Branch) error {
	if len(values) == 0 {
		return nil
	}
	return b.DO.Save(values)
}

func (b branchDo) First() (*entity.Branch, error) {
	if result, err := b.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*entity.Branch), nil
	}
}

func (b branchDo) Take() (*entity.Branch, error) {
	if result, err := b.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*entity.Branch), nil
	}
}

func (b branchDo) Last() (*entity.Branch, error) {
	if result, err := b.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*entity.Branch), nil
	}
}

func (b branchDo) Find() ([]*entity.Branch, error) {
	result, err := b.DO.Find()
	return result.([]*entity.Branch), err
}

func (b branchDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*entity.Branch, err error) {
	buf := make([]*entity.Branch, 0, batchSize)
	err = b.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (b branchDo) FindInBatches(result *[]*entity.Branch, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return b.DO.FindInBatches(result, batchSize, fc)
}

func (b branchDo) Attrs(attrs ...field.AssignExpr) *branchDo {
	return b.withDO(b.DO.Attrs(attrs...))
}

func (b branchDo) Assign(attrs ...field.AssignExpr) *branchDo {
	return b.withDO(b.DO.Assign(attrs...))
}

func (b branchDo) Joins(fields ...field.RelationField) *branchDo {
	for _, _f := range fields {
		b = *b.withDO(b.DO.Joins(_f))
	}
	return &b
}

func (b branchDo) Preload(fields ...field.RelationField) *branchDo {
	for _, _f := range fields {
		b = *b.withDO(b.DO.Preload(_f))
	}
	return &b
}

func (b branchDo) FirstOrInit() (*entity.Branch, error) {
	if result, err := b.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*entity.Branch), nil
	}
}

func (b branchDo) FirstOrCreate() (*entity.Branch, error) {
	if result, err := b.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*entity.Branch), nil
	}
}

func (b branchDo) FindByPage(offset int, limit int) (result []*entity.Branch, count int64, err error) {
	result, err = b.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = b.Offset(-1).Limit(-1).Count()
	return
}

func (b branchDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = b.Count()
	if err != nil {
		return
	}

	err = b.Offset(offset).Limit(limit).Scan(result)
	return
}

func (b branchDo) Scan(result interface{}) (err error) {
	return b.DO.Scan(result)
}

func (b branchDo) Delete(models ...*entity.Branch) (result gen.ResultInfo, err error) {
	return b.DO.Delete(models)
}

func (b *branchDo) withDO(do gen.Dao) *branchDo {
	b.DO = *do.(*gen.DO)
	return b
}