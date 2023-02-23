package db_xorm

import (
	"github.com/chain5j/chain5j-pkg/util/reflectutil"
	"github.com/xwc1125/xwc1125-pkg/database"
)

// Insert 插入数据
func Insert(config database.MysqlConfig, value interface{}) error {
	_, err := MasterEngine(config).Insert(value)
	return err
}

// InsertOne 插入单条数据
func InsertOne(config database.MysqlConfig, value interface{}) error {
	_, err := MasterEngine(config).InsertOne(value)
	return err
}

// UpdateById 根据ID进行更新
func UpdateById(config database.MysqlConfig, id interface{}, value interface{}) error {
	_, err := MasterEngine(config).ID(id).Update(value)
	return err
}

// UpdateByIdAllCols 根据ID更新所有字段
func UpdateByIdAllCols(config database.MysqlConfig, id interface{}, value interface{}) error {
	_, err := MasterEngine(config).ID(id).AllCols().Update(value)
	return err
}

// Updates 根据条件更新
func Updates(config database.MysqlConfig, where interface{}, value interface{}) error {
	_, err := MasterEngine(config).Where(where).Update(value)
	return err
}

// DeleteByModel 根据对象进行删除
func DeleteByModel(config database.MysqlConfig, model interface{}) (count int64, err error) {
	return MasterEngine(config).Delete(model)
}

// DeleteByWhere 根据条件删除
func DeleteByWhere(config database.MysqlConfig, model, where interface{}) (count int64, err error) {
	return MasterEngine(config).Where(where).Delete(model)
}

// DeleteByID 根据ID进行删除
func DeleteByID(config database.MysqlConfig, model interface{}, id int64) (count int64, err error) {
	return MasterEngine(config).ID(id).Delete(model)
}

// Delete DeleteByIDS 根据ids进行删除
func DeleteByIDS(config database.MysqlConfig, model interface{}, ids []int64) (count int64, err error) {
	return MasterEngine(config).In("id", ids).Delete(model)
}

// TableFirst 获取第一条记录，按主键排序
// SELECT * FROM users ORDER BY id LIMIT 1;
func TableFirst(config database.MysqlConfig, out interface{}) (err error) {
	_, err = MasterEngine(config).Asc("id").Get(out)
	return
}

// TableLast 获取最后一条记录，按主键排序
// SELECT * FROM users ORDER BY id DESC LIMIT 1;
func TableLast(config database.MysqlConfig, out interface{}) (err error) {
	_, err = MasterEngine(config).Desc("id").Get(out)
	return
}

// FirstByID 使用主键获取记录
// SELECT * FROM users WHERE id = 10;
func FirstByID(config database.MysqlConfig, out interface{}, id int64) (err error) {
	_, err = MasterEngine(config).ID(id).Get(out)
	return
}

// First
// 获取第一个匹配记录
// db.Where("name = ?", "jinzhu").First(&user)
// SELECT * FROM users WHERE name = 'jinzhu' limit 1;
func First(config database.MysqlConfig, out interface{}, query interface{}, args ...interface{}) (err error) {
	// _, err = MasterEngine(config).Get(where)
	_, err = MasterEngine(config).Where(query, args).Get(out)
	return
}

// FirstByModel 根据model查询对象
func FirstByModel(config database.MysqlConfig, out interface{}) (err error) {
	_, err = MasterEngine(config).Get(out)
	return
}

// type UserDetail struct {
//    User `xorm:"extends"`
//    Detail `xorm:"extends"`
// }
//
// var users []UserDetail
// err := engine.Table("user").Select("user.*, detail.*").
//    Join("INNER", "detail", "detail.user_id = user.id").
//    Where("user.name = ?", name).Limit(10, 0).
//    Find(&users)

func Related(config database.MysqlConfig, model interface{}, out interface{}, foreignKeys string) (err error) {
	val := reflectutil.GetValueByFieldName(model, foreignKeys)
	_, err = MasterEngine(config).ID(val).Get(out)
	return
}

func Preload(config database.MysqlConfig, model interface{}, out interface{}, foreignKeys string) (err error) {
	val := reflectutil.GetValueByFieldName(model, foreignKeys)
	_, err = MasterEngine(config).In("id", val).Get(out)
	return
}

// Find 查找
func Find(config database.MysqlConfig, where interface{}, out interface{}, orders ...string) error {
	db := MasterEngine(config).Where(where)
	if len(orders) > 0 {
		for _, order := range orders {
			db = db.OrderBy(order)
		}
	}
	return db.Find(out)
}
