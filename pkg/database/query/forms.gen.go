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

	"beneburg/pkg/database/model"
)

func newForm(db *gorm.DB) form {
	_form := form{}

	_form.formDo.UseDB(db)
	_form.formDo.UseModel(&model.Form{})

	tableName := _form.formDo.TableName()
	_form.ALL = field.NewAsterisk(tableName)
	_form.ID = field.NewUint(tableName, "id")
	_form.CreatedAt = field.NewTime(tableName, "created_at")
	_form.UpdatedAt = field.NewTime(tableName, "updated_at")
	_form.DeletedAt = field.NewField(tableName, "deleted_at")
	_form.UserTelegramId = field.NewInt64(tableName, "user_telegram_id")
	_form.Name = field.NewString(tableName, "name")
	_form.Age = field.NewInt32(tableName, "age")
	_form.Sex = field.NewString(tableName, "sex")
	_form.About = field.NewString(tableName, "about")
	_form.Hobbies = field.NewString(tableName, "hobbies")
	_form.Work = field.NewString(tableName, "work")
	_form.Education = field.NewString(tableName, "education")
	_form.CoverLetter = field.NewString(tableName, "cover_letter")
	_form.Contacts = field.NewString(tableName, "contacts")
	_form.Status = field.NewString(tableName, "status")
	_form.User = formBelongsToUser{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("User", "model.User"),
	}

	_form.fillFieldMap()

	return _form
}

type form struct {
	formDo formDo

	ALL            field.Asterisk
	ID             field.Uint
	CreatedAt      field.Time
	UpdatedAt      field.Time
	DeletedAt      field.Field
	UserTelegramId field.Int64
	Name           field.String
	Age            field.Int32
	Sex            field.String
	About          field.String
	Hobbies        field.String
	Work           field.String
	Education      field.String
	CoverLetter    field.String
	Contacts       field.String
	Status         field.String
	User           formBelongsToUser

	fieldMap map[string]field.Expr
}

func (f form) Table(newTableName string) *form {
	f.formDo.UseTable(newTableName)
	return f.updateTableName(newTableName)
}

func (f form) As(alias string) *form {
	f.formDo.DO = *(f.formDo.As(alias).(*gen.DO))
	return f.updateTableName(alias)
}

func (f *form) updateTableName(table string) *form {
	f.ALL = field.NewAsterisk(table)
	f.ID = field.NewUint(table, "id")
	f.CreatedAt = field.NewTime(table, "created_at")
	f.UpdatedAt = field.NewTime(table, "updated_at")
	f.DeletedAt = field.NewField(table, "deleted_at")
	f.UserTelegramId = field.NewInt64(table, "user_telegram_id")
	f.Name = field.NewString(table, "name")
	f.Age = field.NewInt32(table, "age")
	f.Sex = field.NewString(table, "sex")
	f.About = field.NewString(table, "about")
	f.Hobbies = field.NewString(table, "hobbies")
	f.Work = field.NewString(table, "work")
	f.Education = field.NewString(table, "education")
	f.CoverLetter = field.NewString(table, "cover_letter")
	f.Contacts = field.NewString(table, "contacts")
	f.Status = field.NewString(table, "status")

	f.fillFieldMap()

	return f
}

func (f *form) WithContext(ctx context.Context) *formDo { return f.formDo.WithContext(ctx) }

func (f form) TableName() string { return f.formDo.TableName() }

func (f form) Alias() string { return f.formDo.Alias() }

func (f *form) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := f.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (f *form) fillFieldMap() {
	f.fieldMap = make(map[string]field.Expr, 16)
	f.fieldMap["id"] = f.ID
	f.fieldMap["created_at"] = f.CreatedAt
	f.fieldMap["updated_at"] = f.UpdatedAt
	f.fieldMap["deleted_at"] = f.DeletedAt
	f.fieldMap["user_telegram_id"] = f.UserTelegramId
	f.fieldMap["name"] = f.Name
	f.fieldMap["age"] = f.Age
	f.fieldMap["sex"] = f.Sex
	f.fieldMap["about"] = f.About
	f.fieldMap["hobbies"] = f.Hobbies
	f.fieldMap["work"] = f.Work
	f.fieldMap["education"] = f.Education
	f.fieldMap["cover_letter"] = f.CoverLetter
	f.fieldMap["contacts"] = f.Contacts
	f.fieldMap["status"] = f.Status

}

func (f form) clone(db *gorm.DB) form {
	f.formDo.ReplaceDB(db)
	return f
}

type formBelongsToUser struct {
	db *gorm.DB

	field.RelationField
}

func (a formBelongsToUser) Where(conds ...field.Expr) *formBelongsToUser {
	if len(conds) == 0 {
		return &a
	}

	exprs := make([]clause.Expression, 0, len(conds))
	for _, cond := range conds {
		exprs = append(exprs, cond.BeCond().(clause.Expression))
	}
	a.db = a.db.Clauses(clause.Where{Exprs: exprs})
	return &a
}

func (a formBelongsToUser) WithContext(ctx context.Context) *formBelongsToUser {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a formBelongsToUser) Model(m *model.Form) *formBelongsToUserTx {
	return &formBelongsToUserTx{a.db.Model(m).Association(a.Name())}
}

type formBelongsToUserTx struct{ tx *gorm.Association }

func (a formBelongsToUserTx) Find() (result *model.User, err error) {
	return result, a.tx.Find(&result)
}

func (a formBelongsToUserTx) Append(values ...*model.User) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a formBelongsToUserTx) Replace(values ...*model.User) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a formBelongsToUserTx) Delete(values ...*model.User) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a formBelongsToUserTx) Clear() error {
	return a.tx.Clear()
}

func (a formBelongsToUserTx) Count() int64 {
	return a.tx.Count()
}

type formDo struct{ gen.DO }

func (f formDo) Debug() *formDo {
	return f.withDO(f.DO.Debug())
}

func (f formDo) WithContext(ctx context.Context) *formDo {
	return f.withDO(f.DO.WithContext(ctx))
}

func (f formDo) ReadDB() *formDo {
	return f.Clauses(dbresolver.Read)
}

func (f formDo) WriteDB() *formDo {
	return f.Clauses(dbresolver.Write)
}

func (f formDo) Clauses(conds ...clause.Expression) *formDo {
	return f.withDO(f.DO.Clauses(conds...))
}

func (f formDo) Returning(value interface{}, columns ...string) *formDo {
	return f.withDO(f.DO.Returning(value, columns...))
}

func (f formDo) Not(conds ...gen.Condition) *formDo {
	return f.withDO(f.DO.Not(conds...))
}

func (f formDo) Or(conds ...gen.Condition) *formDo {
	return f.withDO(f.DO.Or(conds...))
}

func (f formDo) Select(conds ...field.Expr) *formDo {
	return f.withDO(f.DO.Select(conds...))
}

func (f formDo) Where(conds ...gen.Condition) *formDo {
	return f.withDO(f.DO.Where(conds...))
}

func (f formDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *formDo {
	return f.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (f formDo) Order(conds ...field.Expr) *formDo {
	return f.withDO(f.DO.Order(conds...))
}

func (f formDo) Distinct(cols ...field.Expr) *formDo {
	return f.withDO(f.DO.Distinct(cols...))
}

func (f formDo) Omit(cols ...field.Expr) *formDo {
	return f.withDO(f.DO.Omit(cols...))
}

func (f formDo) Join(table schema.Tabler, on ...field.Expr) *formDo {
	return f.withDO(f.DO.Join(table, on...))
}

func (f formDo) LeftJoin(table schema.Tabler, on ...field.Expr) *formDo {
	return f.withDO(f.DO.LeftJoin(table, on...))
}

func (f formDo) RightJoin(table schema.Tabler, on ...field.Expr) *formDo {
	return f.withDO(f.DO.RightJoin(table, on...))
}

func (f formDo) Group(cols ...field.Expr) *formDo {
	return f.withDO(f.DO.Group(cols...))
}

func (f formDo) Having(conds ...gen.Condition) *formDo {
	return f.withDO(f.DO.Having(conds...))
}

func (f formDo) Limit(limit int) *formDo {
	return f.withDO(f.DO.Limit(limit))
}

func (f formDo) Offset(offset int) *formDo {
	return f.withDO(f.DO.Offset(offset))
}

func (f formDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *formDo {
	return f.withDO(f.DO.Scopes(funcs...))
}

func (f formDo) Unscoped() *formDo {
	return f.withDO(f.DO.Unscoped())
}

func (f formDo) Create(values ...*model.Form) error {
	if len(values) == 0 {
		return nil
	}
	return f.DO.Create(values)
}

func (f formDo) CreateInBatches(values []*model.Form, batchSize int) error {
	return f.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (f formDo) Save(values ...*model.Form) error {
	if len(values) == 0 {
		return nil
	}
	return f.DO.Save(values)
}

func (f formDo) First() (*model.Form, error) {
	if result, err := f.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Form), nil
	}
}

func (f formDo) Take() (*model.Form, error) {
	if result, err := f.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Form), nil
	}
}

func (f formDo) Last() (*model.Form, error) {
	if result, err := f.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Form), nil
	}
}

func (f formDo) Find() ([]*model.Form, error) {
	result, err := f.DO.Find()
	return result.([]*model.Form), err
}

func (f formDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Form, err error) {
	buf := make([]*model.Form, 0, batchSize)
	err = f.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (f formDo) FindInBatches(result *[]*model.Form, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return f.DO.FindInBatches(result, batchSize, fc)
}

func (f formDo) Attrs(attrs ...field.AssignExpr) *formDo {
	return f.withDO(f.DO.Attrs(attrs...))
}

func (f formDo) Assign(attrs ...field.AssignExpr) *formDo {
	return f.withDO(f.DO.Assign(attrs...))
}

func (f formDo) Joins(fields ...field.RelationField) *formDo {
	for _, _f := range fields {
		f = *f.withDO(f.DO.Joins(_f))
	}
	return &f
}

func (f formDo) Preload(fields ...field.RelationField) *formDo {
	for _, _f := range fields {
		f = *f.withDO(f.DO.Preload(_f))
	}
	return &f
}

func (f formDo) FirstOrInit() (*model.Form, error) {
	if result, err := f.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Form), nil
	}
}

func (f formDo) FirstOrCreate() (*model.Form, error) {
	if result, err := f.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Form), nil
	}
}

func (f formDo) FindByPage(offset int, limit int) (result []*model.Form, count int64, err error) {
	result, err = f.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = f.Offset(-1).Limit(-1).Count()
	return
}

func (f formDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = f.Count()
	if err != nil {
		return
	}

	err = f.Offset(offset).Limit(limit).Scan(result)
	return
}

func (f formDo) Scan(result interface{}) (err error) {
	return f.DO.Scan(result)
}

func (f formDo) Delete(models ...*model.Form) (result gen.ResultInfo, err error) {
	return f.DO.Delete(models)
}

func (f *formDo) withDO(do gen.Dao) *formDo {
	f.DO = *do.(*gen.DO)
	return f
}