# xorm 使用说明

[TOC]

http://gobook.io/read/github.com/go-xorm/manual-zh-CN/

## 特性

- 支持Struct和数据库表之间的灵活映射，并支持自动同步表结构
- 事务支持
- 支持原始SQL语句和ORM操作的混合执行
- 使用连写来简化调用
- 支持使用Id, In, Where, Limit, Join, Having, Table, Sql, Cols等函数和结构体等方式作为条件
- 支持级联加载Struct
- 支持LRU缓存(支持memory, memcache, leveldb, redis缓存Store) 和 Redis缓存
- 支持反转，即根据数据库自动生成xorm的结构体
- 支持事件
- 支持created, updated, deleted和version记录版本（即乐观锁）

## 驱动支持

- Mysql: github.com/go-sql-driver/mysql
- MyMysql: github.com/ziutek/mymysql/godrv
- Postgres: github.com/lib/pq
- Tidb: github.com/pingcap/tidb
- SQLite: github.com/mattn/go-sqlite3
- MsSql: github.com/denisenkom/go-mssqldb
- MsSql: github.com/lunny/godbc
- Oracle: github.com/mattn/go-oci8 (试验性支持)
- ql: github.com/cznic/ql (试验性支持)

## 说明

### 使用Table和Tag改变名称映射

- 如果结构体拥有TableName() string的成员方法，那么此方法的返回值即是该结构体对应的数据库表名。
- 通过engine.Table()方法可以改变struct对应的数据库表的名称，通过sturct中field对应的Tag中使用xorm:"'column_name'"
  可以使该field对应的Column名称为指定名称。这里使用两个单引号将Column名称括起来是为了防止名称冲突，因为我们在Tag中还可以对这个Column进行更多的定义。如果名称不冲突的情况，单引号也可以不使用。

- 表名的优先级顺序如下：
    - engine.Table() 指定的临时表名优先级最高
    - TableName() string 其次
    - Mapper 自动映射的表名优先级最后
- 字段名的优先级顺序如下：
    - 结构体tag指定的字段名优先级较高
    - Mapper 自动映射的表名优先级较低

### Column属性定义

```
type User struct {
    Id   int64
    Name string  `xorm:"varchar(25) notnull unique 'usr_name'"`
}
```

具体的Tag规则如下，另Tag中的关键字均不区分大小写，但字段名根据不同的数据库是区分大小写：

| 字段  | 说明  |
|---|---|
|name    |当前field对应的字段的名称，可选，如不写，则自动根据field名字和转换规则命名，如与其它关键字冲突，请使用单引号括起来。|
|pk    |是否是Primary Key，如果在一个struct中有多个字段都使用了此标记，则这多个字段构成了复合主键，单主键当前支持int32,int,int64,uint32,uint,uint64,string这7种Go的数据类型，复合主键支持这7种Go的数据类型的组合。|
|当前支持30多种字段类型，详情参见本文最后一个表格    |字段类型|
|autoincr    |是否是自增|
|[not ]null 或 notnull    |是否可以为空|
|unique或unique(uniquename)    |是否是唯一，如不加括号则该字段不允许重复；如加上括号，则括号中为联合唯一索引的名字，此时如果有另外一个或多个字段和本unique的uniquename相同，则这些uniquename相同的字段组成联合唯一索引|
|index或index(indexname)    |是否是索引，如不加括号则该字段自身为索引，如加上括号，则括号中为联合索引的名字，此时如果有另外一个或多个字段和本index的indexname相同，则这些indexname相同的字段组成联合索引|
|extends    |应用于一个匿名成员结构体或者非匿名成员结构体之上，表示此结构体的所有成员也映射到数据库中，extends可加载无限级|
|-    |这个Field将不进行字段映射|
|->    |这个Field将只写入到数据库而不从数据库读取|
|<-    |这个Field将只从数据库读取，而不写入到数据库|
|created    |这个Field将在Insert时自动赋值为当前时间|
|updated    |这个Field将在Insert或Update时自动赋值为当前时间|
|deleted    |这个Field将在Delete时设置为当前时间，并且当前记录不删除|
|version    |这个Field将会在insert时默认为1，每次更新自动加1|
|default 0或default(0)    |设置默认值，紧跟的内容如果是Varchar等需要加上单引号|
|json    |表示内容将先转成Json格式，然后存储到数据库中，数据库中的字段类型可以为Text或者二进制|

自动映射的规则：

- 1.如果field名称为Id而且类型为int64并且没有定义tag，则会被xorm视为主键，并且拥有自增属性。如果想用Id以外的名字或非int64类型做为主键名，必须在对应的Tag上加上xorm:"pk"来定义主键，加上xorm:"
  autoincr"作为自增。这里需要注意的是，有些数据库并不允许非主键的自增属性。
- 2.string类型默认映射为varchar(255)，如果需要不同的定义，可以在tag中自定义，如：varchar(1024)
- 3.支持type MyString string等自定义的field，支持Slice,
  Map等field成员，这些成员默认存储为Text类型，并且默认将使用Json格式来序列化和反序列化。也支持数据库字段类型为Blob类型。如果是Blob类型，则先使用Json格式序列化再转成[]byte格式。如果是[]byte或者[]
  uint8，则不做转换二十直接以二进制方式存储。具体参见 Go与字段类型对应表
- 4.实现了Conversion接口的类型或者结构体，将根据接口的转换方式在类型和数据库记录之间进行相互转换，这个接口的优先级是最高的。

    ```
    type Conversion interface {
    FromDB([]byte) error
    ToDB() ([]byte, error)
    }
    ```

- 5.如果一个结构体包含一个Conversion的接口类型，那么在获取数据时，必须要预先设置一个实现此接口的struct或者struct的指针。此时可以在此struct中实现BeforeSet(name string, cell
  xorm.Cell)方法来进行预先给Conversion赋值。

### Go与字段类型对应表

| go type's kind  | value method  |xorm type|
|---| --- |---| 
|implemented Conversion    |Conversion.ToDB / Conversion.FromDB    |Text|
|int, int8, int16, int32, uint, uint8, uint16, uint32|		|Int|
|int64, uint64    |	|BigInt|
|float32    |	|Float|
|float64    |	|Double|
|complex64, complex128    |json.Marshal / json.UnMarshal    |Varchar(64)|
|[]uint8    |	|Blob|
|array, slice, map except []uint8    |json.Marshal / json.UnMarshal    |Text|
|bool    |1 or 0    |Bool|
|string    |	|Varchar(255)|
|time.Time|		|DateTime|
|cascade |struct	primary key field value    |BigInt|
|struct    |json.Marshal / json.UnMarshal    |Text|
|Others    |	|Text|

## 操作

### 表结构操作

#### 获取数据库信息

- DBMetas()

  xorm支持获取表结构信息，通过调用engine.DBMetas()可以获取到数据库中所有的表，字段，索引的信息。

- TableInfo()

  根据传入的结构体指针及其对应的Tag，提取出模型对应的表结构信息。这里不是数据库当前的表结构信息，而是我们通过struct建模时希望数据库的表的结构信息

#### 表操作

- CreateTables()

  创建表使用engine.CreateTables()，参数为一个或多个空的对应Struct的指针。同时可用的方法有Charset()和StoreEngine()
  ，如果对应的数据库支持，这两个方法可以在创建表时指定表的字符编码和使用的引擎。Charset()和StoreEngine()当前仅支持Mysql数据库。

- IsTableEmpty()

  判断表是否为空，参数和CreateTables相同

- IsTableExist()

  判断表是否存在

- DropTables()

  删除表使用engine.DropTables()，参数为一个或多个空的对应Struct的指针或者表的名字。如果为string传入，则只删除对应的表，如果传入的为Struct，则删除表的同时还会删除对应的索引。

#### 创建索引和唯一索引

- CreateIndexes

  根据struct中的tag来创建索引

- CreateUniques

  根据struct中的tag来创建唯一索引

#### 同步数据库结构

- Sync

Sync将进行如下的同步操作：

* 自动检测和创建表，这个检测是根据表的名字
* 自动检测和新增表中的字段，这个检测是根据字段名
* 自动检测和创建索引和唯一索引，这个检测是根据索引的一个或多个字段名，而不根据索引名称 调用方法如下：

```
err := engine.Sync(new(User), new(Group))
```

- Sync2 Sync2对Sync进行了改进，目前推荐使用Sync2。Sync2函数将进行如下的同步操作：

* 自动检测和创建表，这个检测是根据表的名字
* 自动检测和新增表中的字段，这个检测是根据字段名，同时对表中多余的字段给出警告信息
* 自动检测，创建和删除索引和唯一索引，这个检测是根据索引的一个或多个字段名，而不根据索引名称。因此这里需要注意，如果在一个有大量数据的表中引入新的索引，数据库可能需要一定的时间来建立索引。
* 自动转换varchar字段类型到text字段类型，自动警告其它字段类型在模型和数据库之间不一致的情况。
* 自动警告字段的默认值，是否为空信息在模型和数据库之间不匹配的情况

以上这些警告信息需要将`engine.ShowWarn` 设置为 `true` 才会显示。 调用方法和Sync一样：

```
err := engine.Sync2(new(User), new(Group))
```

#### Dump数据库结构和数据

如果需要在程序中Dump数据库的结构和数据可以调用

```
engine.DumpAll(w io.Writer)
```

和

```
engine.DumpAllFile(fpath string)
```

DumpAll方法接收一个io.Writer接口来保存Dump出的数据库结构和数据的SQL语句，这个方法导出的SQL语句并不能通用。只针对当前engine所对应的数据库支持的SQL。

Import 执行数据库SQL脚本 如果你需要将保存在文件或者其它存储设施中的SQL脚本执行，那么可以调用
`
engine.Import(r io.Reader)`

和
`
engine.ImportFile(fpath string)`

同样，这里需要对应的数据库的SQL语法支持。

### 插入数据

插入数据使用Insert方法，Insert方法的参数可以是一个或多个Struct的指针，一个或多个Struct的Slice的指针。

如果传入的是Slice并且当数据库支持批量插入时，Insert会使用批量插入的方式进行插入。

- 插入一条数据，此时可以用Insert或者InsertOne

```
user := new(User)
user.Name = "myname"
affected, err := engine.Insert(user)
// INSERT INTO user (name) values (?)
```

在插入单条数据成功后，如果该结构体有自增字段(设置为autoincr)，则自增字段会被自动赋值为数据库中的id。这里需要注意的是，如果插入的结构体中，自增字段已经赋值，则该字段会被作为非自增字段插入。

```
fmt.Println(user.Id)
```

- 插入同一个表的多条数据，此时如果数据库支持批量插入，那么会进行批量插入，但是这样每条记录就无法被自动赋予id值。如果数据库不支持批量插入，那么就会一条一条插入。

```
users := make([]User, 1)
users[0].Name = "name0"
...
affected, err := engine.Insert(&users)
```

- 使用指针Slice插入多条记录，同上

```
users := make([]*User, 1)
users[0] = new(User)
users[0].Name = "name0"
...
affected, err := engine.Insert(&users)
```

- 插入多条记录并且不使用批量插入，此时实际生成多条插入语句，每条记录均会自动赋予Id值。

```
users := make([]*User, 1)
users[0] = new(User)
users[0].Name = "name0"
...
affected, err := engine.Insert(users...)
```

- 插入不同表的一条记录

```
user := new(User)
user.Name = "myname"
question := new(Question)
question.Content = "whywhywhwy?"
affected, err := engine.Insert(user, question)
```

- 插入不同表的多条记录

```
users := make([]User, 1)
users[0].Name = "name0"
...
questions := make([]Question, 1)
questions[0].Content = "whywhywhwy?"
affected, err := engine.Insert(&users, &questions)
```

- 插入不同表的一条或多条记录

```
user := new(User)
user.Name = "myname"
...
questions := make([]Question, 1)
questions[0].Content = "whywhywhwy?"
affected, err := engine.Insert(user, &questions)
```

这里需要注意以下几点：

- 这里虽然支持同时插入，但这些插入并没有事务关系。因此有可能在中间插入出错后，后面的插入将不会继续。此时前面的插入已经成功，如果需要回滚，请开启事务。
- 批量插入会自动生成Insert into table values (),(),()
  的语句，因此各个数据库对SQL语句有长度限制，因此这样的语句有一个最大的记录数，根据经验测算在150条左右。大于150条后，生成的sql语句将太长可能导致执行失败。因此在插入大量数据时，目前需要自行分割成每150条插入一次。

### 查询和统计数据

所有的查询条件不区分调用顺序，但必须在调用Get，Exist, Sum, Find，Count, Iterate,
Rows这几个函数之前调用。同时需要注意的一点是，在调用的参数中，如果采用默认的SnakeMapper所有的字符字段名均为映射后的数据库的字段名，而不是field的名字。

#### 查询条件方法

查询和统计主要使用Get, Find, Count, Rows, Iterate这几个方法，同时大部分函数在调用Update, Delete时也是可用的。在进行查询时可以使用多个方法来形成查询条件，条件函数如下：

- Alias(string)

给Table设定一个别名

```
engine.Alias("o").Where("o.name = ?", name).Get(&order)
```

- And(string, …interface{})

和Where函数中的条件基本相同，作为条件

```
engine.Where(...).And(...).Get(&order)
```

- Asc(…string)

指定字段名正序排序，可以组合

```
engine.Asc("id").Find(&orders)
```

- Desc(…string)

指定字段名逆序排序，可以组合

```
engine.Asc("id").Desc("time").Find(&orders)
```

- ID(interface{})

传入一个主键字段的值，作为查询条件，如

```
var user User
engine.ID(1).Get(&user)
// SELECT * FROM user Where id = 1
```

如果是复合主键，则可以

```
engine.ID(core.PK{1, "name"}).Get(&user)
// SELECT * FROM user Where id =1 AND name= 'name'
```

传入的两个参数按照struct中pk标记字段出现的顺序赋值。

- Or(interface{}, …interface{})

和Where函数中的条件基本相同，作为条件

- OrderBy(string)

按照指定的顺序进行排序

- Select(string)

指定select语句的字段部分内容，例如：

```
engine.Select("a.*, (select name from b limit 1) as name").Find(&beans)

engine.Select("a.*, (select name from b limit 1) as name").Get(&bean)
```

- SQL(string, …interface{})

执行指定的Sql语句，并把结果映射到结构体。有时，当选择内容或者条件比较复杂时，可以直接使用Sql，例如：

```
engine.SQL("select * from table").Find(&beans)
```

- Where(string, …interface{})

和SQL中Where语句中的条件基本相同，作为条件

```
engine.Where("a = ? AND b = ?", 1, 2).Find(&beans)

engine.Where(builder.Eq{"a":1, "b": 2}).Find(&beans)

engine.Where(builder.Eq{"a":1}.Or(builder.Eq{"b": 2})).Find(&beans)
```

- In(string, …interface{})

某字段在一些值中，这里需要注意必须是[]interface{}才可以展开，由于Go语言的限制，[]int64等不可以直接展开，而是通过传递一个slice。第二个参数也可以是一个*builder.Builder 指针。示例代码如下：

```
// select from table where column in (1,2,3)
engine.In("cloumn", 1, 2, 3).Find()

// select from table where column in (1,2,3)
engine.In("column", []int{1, 2, 3}).Find()

// select from table where column in (select column from table2 where a = 1)
engine.In("column", builder.Select("column").From("table2").Where(builder.Eq{"a":1})).Find()
```

- Cols(…string)

只查询或更新某些指定的字段，默认是查询所有映射的字段或者根据Update的第一个参数来判断更新的字段。例如：

```
engine.Cols("age", "name").Get(&usr)
// SELECT age, name FROM user limit 1
engine.Cols("age", "name").Find(&users)
// SELECT age, name FROM user
engine.Cols("age", "name").Update(&user)
// UPDATE user SET age=? AND name=?
```

- AllCols()

查询或更新所有字段，一般与Update配合使用，因为默认Update只更新非0，非”“，非bool的字段。

```
engine.AllCols().Id(1).Update(&user)
// UPDATE user SET name = ?, age =?, gender =? WHERE id = 1
```

- MustCols(…string)
  某些字段必须更新，一般与Update配合使用。

- Omit(…string)

和cols相反，此函数指定排除某些指定的字段。注意：此方法和Cols方法不可同时使用。

```
// 例1：
engine.Omit("age", "gender").Update(&user)
// UPDATE user SET name = ? AND department = ?
// 例2：
engine.Omit("age, gender").Insert(&user)
// INSERT INTO user (name) values (?) // 这样的话age和gender会给默认值
// 例3：
engine.Omit("age", "gender").Find(&users)
// SELECT name FROM user //只select除age和gender字段的其它字段
```

- Distinct(…string)

按照参数中指定的字段归类结果。

```
engine.Distinct("age", "department").Find(&users)
// SELECT DISTINCT age, department FROM user
```

注意：当开启了缓存时，此方法的调用将在当前查询中禁用缓存。因为缓存系统当前依赖Id，而此时无法获得Id

- Table(nameOrStructPtr interface{})

传入表名称或者结构体指针，如果传入的是结构体指针，则按照IMapper的规则提取出表名

- Limit(int, …int)

限制获取的数目，第一个参数为条数，第二个参数表示开始位置，如果不传则为0

- Top(int)

相当于Limit(int, 0)

- Join(string,interface{},string)

第一个参数为连接类型，当前支持INNER, LEFT OUTER, CROSS中的一个值， 第二个参数为string类型的表名，表对应的结构体指针或者为两个值的[]string，表示表名和别名， 第三个参数为连接条件

- GroupBy(string)

Groupby的参数字符串

- Having(string)

Having的参数字符串

#### Get方法

查询单条数据使用Get方法，在调用Get方法时需要传入一个对应结构体的指针，同时结构体中的非空field自动成为查询的条件和前面的方法条件组合在一起查询。

如：

1) 根据Id来获得单条数据:

```
user := new(User)
has, err := engine.Id(id).Get(user)
// 复合主键的获取方法
// has, errr := engine.Id(xorm.PK{1,2}).Get(user)
```

2) 根据Where来获得单条数据：

```
user := new(User)
has, err := engine.Where("name=?", "xlw").Get(user)
```

3) 根据user结构体中已有的非空数据来获得单条数据：

```
user := &User{Id:1}
has, err := engine.Get(user)
```

或者其它条件

```
user := &User{Name:"xlw"}
has, err := engine.Get(user)
```

返回的结果为两个参数，一个has为该条记录是否存在，第二个参数err为是否有错误。不管err是否为nil，has都有可能为true或者false。

#### Exist系列方法

判断某个记录是否存在可以使用Exist, 相比Get，Exist性能更好。

#### Find方法

查询多条数据使用Find方法，Find方法的第一个参数为slice的指针或Map指针，即为查询后返回的结果，第二个参数可选，为查询的条件struct的指针。

1) 传入Slice用于返回数据

```
everyone := make([]Userinfo, 0)
err := engine.Find(&everyone)

pEveryOne := make([]*Userinfo, 0)
err := engine.Find(&pEveryOne)
```

2) 传入Map用户返回数据，map必须为map[int64]Userinfo的形式，map的key为id，因此对于复合主键无法使用这种方式。

```
users := make(map[int64]Userinfo)
err := engine.Find(&users)

pUsers := make(map[int64]*Userinfo)
err := engine.Find(&pUsers)
```

3) 也可以加入各种条件

```
users := make([]Userinfo, 0)
err := engine.Where("age > ? or name = ?", 30, "xlw").Limit(20, 10).Find(&users)
```

4) 如果只选择单个字段，也可使用非结构体的Slice

```
var ints []int64
err := engine.Table("user").Cols("id").Find(&ints)
```

#### Join的使用

#### Iterate方法

Iterate方法提供逐条执行查询到的记录的方法，他所能使用的条件和Find方法完全相同

#### Count方法

统计数据使用Count方法，Count方法的参数为struct的指针并且成为查询条件。

#### Rows方法

Rows方法和Iterate方法类似，提供逐条执行查询到的记录的方法，不过Rows更加灵活好用。

```
user := new(User)
rows, err := engine.Where("id >?", 1).Rows(user)
if err != nil {
}
defer rows.Close()
for rows.Next() {
    err = rows.Scan(user)
    //...
}
```

#### Sum系列方法

求和数据可以使用Sum, SumInt, Sums 和 SumsInt 四个方法，Sums系列方法的参数为struct的指针并且成为查询条件。

### 更新数据

更新数据使用Update方法，Update方法的第一个参数为需要更新的内容，可以为一个结构体指针或者一个Map[string]
interface{}类型。当传入的为结构体指针时，只有非空和0的field才会被作为更新的字段。当传入的为Map类型时，key为数据库Column的名字，value为要更新的内容。

Update方法将返回两个参数，第一个为 更新的记录数，需要注意的是 SQLITE 数据库返回的是根据更新条件查询的记录数而不是真正受更新的记录数。

```
user := new(User)
user.Name = "myname"
affected, err := engine.Id(id).Update(user)
```

这里需要注意，Update会自动从user结构体中提取非0和非nil得值作为需要更新的内容，因此，如果需要更新一个值为0，则此种方法将无法实现，因此有两种选择：

- 1.通过添加Cols函数指定需要更新结构体中的哪些值，未指定的将不更新，指定了的即使为0也会更新。

```
affected, err := engine.Id(id).Cols("age").Update(&user)
```

- 2.通过传入map[string]interface{}来进行更新，但这时需要额外指定更新到哪个表，因为通过map是无法自动检测更新哪个表的。

```
affected, err := engine.Table(new(User)).Id(id).Update(map[string]interface{}{"age":0})
```

#### 乐观锁Version

要使用乐观锁，需要使用version标记

```
type User struct {
    Id int64
    Name string
    Version int `xorm:"version"`
}
```

在Insert时，version标记的字段将会被设置为1，在Update时，Update的内容必须包含version原来的值。

```
var user User
engine.Id(1).Get(&user)
// SELECT * FROM user WHERE id = ?
engine.Id(1).Update(&user)
// UPDATE user SET ..., version = version + 1 WHERE id = ? AND version = ?
```

#### 更新时间Updated

Updated可以让您在记录插入或每次记录更新时自动更新数据库中的标记字段为当前时间，需要在xorm标记中使用updated标记，如下所示进行标记，对应的字段可以为time.Time或者自定义的time.Time或者int,int64等int类型。

```
type User struct {
    Id int64
    Name string
    UpdatedAt time.Time `xorm:"updated"`
}
```

在Insert(), InsertOne(), Update()方法被调用时，updated标记的字段将会被自动更新为当前时间，如下所示：

```
var user User
engine.Id(1).Get(&user)
// SELECT * FROM user WHERE id = ?
engine.Id(1).Update(&user)
// UPDATE user SET ..., updaetd_at = ? WHERE id = ?
```

如果你希望临时不自动插入时间，则可以组合NoAutoTime()方法：

```
engine.NoAutoTime().Insert(&user)
```

这个在从一张表拷贝字段到另一张表时比较有用。

### 删除数据

删除数据Delete方法，参数为struct的指针并且成为查询条件。

```
user := new(User)
affected, err := engine.Id(id).Delete(user)
```

Delete的返回值第一个参数为删除的记录数，第二个参数为错误。

注意：当删除时，如果user中包含有bool,float64或者float32类型，有可能会使删除失败。

#### 软删除Deleted

Deleted可以让您不真正的删除数据，而是标记一个删除时间。使用此特性需要在xorm标记中使用deleted标记，如下所示进行标记，对应的字段必须为time.Time类型。

```
type User struct {
    Id int64
    Name string
    DeletedAt time.Time `xorm:"deleted"`
}
```

在Delete()时，deleted标记的字段将会被自动更新为当前时间而不是去删除该条记录，如下所示：

```
var user User
engine.Id(1).Get(&user)
// SELECT * FROM user WHERE id = ?
engine.Id(1).Delete(&user)
// UPDATE user SET ..., deleted_at = ? WHERE id = ?
engine.Id(1).Get(&user)
// 再次调用Get，此时将返回false, nil，即记录不存在
engine.Id(1).Delete(&user)
// 再次调用删除会返回0, nil，即记录不存在
```

那么如果记录已经被标记为删除后，要真正的获得该条记录或者真正的删除该条记录，需要启用Unscoped，如下所示：

```
var user User
engine.Id(1).Unscoped().Get(&user)
// 此时将可以获得记录
engine.Id(1).Unscoped().Delete(&user)
// 此时将可以真正的删除记录
```

### 执行SQL查询

也可以直接执行一个SQL查询，即Select命令。在Postgres中支持原始SQL语句中使用 ` 和 ? 符号。

```
sql := "select * from userinfo"
results, err := engine.Query(sql)
```

当调用Query时，第一个返回值results为[]map[string][]byte的形式。

### 执行SQL命令

也可以直接执行一个SQL命令，即执行Insert， Update， Delete 等操作。此时不管数据库是何种类型，都可以使用 ` 和 ? 符号。

```
sql = "update `userinfo` set username=? where id=?"
res, err := engine.Exec(sql, "xiaolun", 1) 
```

### 事务处理

当使用事务处理时，需要创建Session对象。在进行事物处理时，可以混用ORM方法和RAW方法，如下代码所示：

```
session := engine.NewSession()
defer session.Close()
// add Begin() before any action
err := session.Begin()
user1 := Userinfo{Username: "xiaoxiao", Departname: "dev", Alias: "lunny", Created: time.Now()}
_, err = session.Insert(&user1)
if err != nil {
    session.Rollback()
    return
}
user2 := Userinfo{Username: "yyy"}
_, err = session.Where("id = ?", 2).Update(&user2)
if err != nil {
    session.Rollback()
    return
}

_, err = session.Exec("delete from userinfo where username = ?", user2.Username)
if err != nil {
    session.Rollback()
    return
}

// add Commit() after all actions
err = session.Commit()
if err != nil {
    return
}
```

注意如果您使用的是mysql，数据库引擎为innodb事务才有效，myisam引擎是不支持事务的。

### 缓存

xorm内置了一致性缓存支持，不过默认并没有开启。要开启缓存，需要在engine创建完后进行配置，如： 启用一个全局的内存缓存

```
cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 1000)
engine.SetDefaultCacher(cacher)
```

上述代码采用了LRU算法的一个缓存，缓存方式是存放到内存中，缓存struct的记录数为1000条，缓存针对的范围是所有具有主键的表，没有主键的表中的数据将不会被缓存。 如果只想针对部分表，则：

```
cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 1000)
engine.MapCacher(&user, cacher)
```

如果要禁用某个表的缓存，则：

```
engine.MapCacher(&user, nil)
```

设置完之后，其它代码基本上就不需要改动了，缓存系统已经在后台运行。

当前实现了内存存储的CacheStore接口MemoryStore，如果需要采用其它设备存储，可以实现CacheStore接口。

不过需要特别注意不适用缓存或者需要手动编码的地方：

- 当使用了Distinct,Having,GroupBy方法将不会使用缓存

- 在Get或者Find时使用了Cols,Omit方法，则在开启缓存后此方法无效，系统仍旧会取出这个表中的所有字段。

- 在使用Exec方法执行了方法之后，可能会导致缓存与数据库不一致的地方。因此如果启用缓存，尽量避免使用Exec。如果必须使用，则需要在使用了Exec之后调用ClearCache手动做缓存清除的工作。比如：

```
engine.Exec("update user set name = ? where id = ?", "xlw", 1)
engine.ClearCache(new(User))
```

### 事件

xorm支持两种方式的事件，一种是在Struct中的特定方法来作为事件的方法，一种是在执行语句的过程中执行事件。

在Struct中作为成员方法的事件如下：

- BeforeInsert()

在将此struct插入到数据库之前执行

- BeforeUpdate()

在将此struct更新到数据库之前执行

- BeforeDelete()

在将此struct对应的条件数据从数据库删除之前执行

- func BeforeSet(name string, cell xorm.Cell)

在 Get 或 Find 方法中，当数据已经从数据库查询出来，而在设置到结构体之前调用，name为数据库字段名称，cell为数据库中的字段值。

- func AfterSet(name string, cell xorm.Cell)

在 Get 或 Find 方法中，当数据已经从数据库查询出来，而在设置到结构体之后调用，name为数据库字段名称，cell为数据库中的字段值。

- AfterInsert()

在将此struct成功插入到数据库之后执行

- AfterUpdate()

在将此struct成功更新到数据库之后执行

- AfterDelete()

在将此struct对应的条件数据成功从数据库删除之后执行

在语句执行过程中的事件方法为：

- Before(beforeFunc interface{})

临时执行某个方法之前执行

```
before := func(bean interface{}){
    fmt.Println("before", bean)
}
engine.Before(before).Insert(&obj)
```

- After(afterFunc interface{})

临时执行某个方法之后执行

```
after := func(bean interface{}){
    fmt.Println("after", bean)
}
engine.After(after).Insert(&obj)
```

其中beforeFunc和afterFunc的原型为func(bean interface{}).

