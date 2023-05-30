package gdb

import (
	gcore "github.com/snail007/gmc/core"
	gconfig "github.com/snail007/gmc/module/config"
	gmap "github.com/snail007/gmc/util/map"
	"github.com/stretchr/testify/assert"
	"log"
	"reflect"
	"testing"
	"time"
)

func TestModel(t *testing.T) {
	db := createTestDB() // 创建测试用的数据库对象

	// 创建测试用的表和数据
	dropTestTable(db)
	createTestTable(db)
	insertTestData(db)
	time.Sleep(time.Second)
	// 创建Model对象并设置数据库连接
	model := Table("test_table", db)

	// 测试ExecSQL函数
	_, count, err := model.ExecSQL("UPDATE test_table SET column1 = ?", "update_value1")
	if err != nil {
		t.Errorf("ExecSQL failed: %v", err)
	}
	// 检查返回结果是否符合预期
	expectedCount := int64(3)
	if count != expectedCount {
		t.Errorf("Count returned unexpected result. Expected: %d, Got: %d", expectedCount, count)
	}
	model.ExecSQL("UPDATE test_table SET column1 = ?", "value1")

	// 测试Count函数
	count, err = model.Count(map[string]interface{}{"column1": "value1"})
	if err != nil {
		t.Errorf("Count failed: %v", err)
	}
	// 检查返回结果是否符合预期
	expectedCount = int64(3)
	if count != expectedCount {
		t.Errorf("Count returned unexpected result. Expected: %d, Got: %d", expectedCount, count)
	}

	// 测试GetByID函数
	ret, err := model.GetByID("1")
	if err != nil {
		t.Errorf("GetByID failed: %v", err)
	}
	// 检查返回结果是否符合预期
	expectedData := map[string]interface{}{"test_table_id": "1", "column1": "value1", "column2": "value2"}
	if !assert.Equal(t, expectedData, gmap.ToAny(ret)) {
		t.Errorf("GetByID returned unexpected result. Expected: %v, Got: %v", expectedData, ret)
	}

	// 测试QuerySQL函数
	rs, err := model.QuerySQL("select * from test_table where test_table_id=1")
	if err != nil {
		t.Errorf("QuerySQL failed: %v", err)
	}
	// 检查返回结果是否符合预期
	expectedData = map[string]interface{}{"test_table_id": "1", "column1": "value1", "column2": "value2"}
	if !assert.Equal(t, expectedData, gmap.ToAny(rs[0])) {
		t.Errorf("QuerySQL returned unexpected result. Expected: %v, Got: %v", expectedData, ret)
	}

	// 测试GetBy函数
	ret, err = model.GetBy(map[string]interface{}{"column1": "value1"})
	if err != nil {
		t.Errorf("GetBy failed: %v", err)
	}
	// 检查返回结果是否符合预期
	expectedData = map[string]interface{}{"test_table_id": "1", "column1": "value1", "column2": "value2"}
	if !reflect.DeepEqual(expectedData, gmap.ToAny(ret)) {
		t.Errorf("GetBy returned unexpected result. Expected: %v, Got: %v", expectedData, ret)
	}

	// 测试MGetByIDs函数
	a, err := model.MGetByIDs([]string{"1", "2", "3"})
	if err != nil {
		t.Errorf("MGetByIDs failed: %v", err)
	}
	a0 := []map[string]interface{}{}
	for _, v := range a {
		a0 = append(a0, gmap.ToAny(v))
	}
	// 检查返回结果是否符合预期
	b := []map[string]interface{}{
		{"test_table_id": "1", "column1": "value1", "column2": "value2"},
		{"test_table_id": "2", "column1": "value1", "column2": "value2"},
		{"test_table_id": "3", "column1": "value1", "column2": "value2"},
	}
	if !reflect.DeepEqual(b, a0) {
		t.Errorf("MGetByIDs returned unexpected result. Expected: %v, Got: %v", b, ret)
	}

	// 测试MGetByIDsRs函数
	a1, err := model.MGetByIDsRs([]string{"1", "2", "3"})
	if err != nil {
		t.Errorf("MGetByIDsRs failed: %v", err)
	}
	a2 := []map[string]interface{}{}
	for _, v := range a1.Rows() {
		a2 = append(a2, gmap.ToAny(v))
	}
	// 检查返回结果是否符合预期
	b2 := []map[string]interface{}{
		{"test_table_id": "1", "column1": "value1", "column2": "value2"},
		{"test_table_id": "2", "column1": "value1", "column2": "value2"},
		{"test_table_id": "3", "column1": "value1", "column2": "value2"},
	}
	if !reflect.DeepEqual(b2, a2) {
		t.Errorf("MGetByIDsRs returned unexpected result. Expected: %v, Got: %v", b, ret)
	}

	// 测试GetAll函数
	c, err := model.GetAll()
	if err != nil {
		t.Errorf("GetAll failed: %v", err)
	}
	c0 := []map[string]interface{}{}
	for _, v := range c {
		c0 = append(c0, gmap.ToAny(v))
	}
	// 检查返回结果是否符合预期
	d := []map[string]interface{}{
		{"test_table_id": "1", "column1": "value1", "column2": "value2"},
		{"test_table_id": "2", "column1": "value1", "column2": "value2"},
		{"test_table_id": "3", "column1": "value1", "column2": "value2"},
	}
	if !reflect.DeepEqual(d, c0) {
		t.Errorf("GetAll returned unexpected result. Expected: %v, Got: %v", d, c)
	}

	// 测试GetAllRs函数
	rs0, err := model.GetAllRs()
	if err != nil {
		t.Errorf("GetAllRs failed: %v", err)
	}
	c1 := []map[string]interface{}{}
	for _, v := range rs0.Rows() {
		c1 = append(c1, gmap.ToAny(v))
	}
	// 检查返回结果是否符合预期
	d1 := []map[string]interface{}{
		{"test_table_id": "1", "column1": "value1", "column2": "value2"},
		{"test_table_id": "2", "column1": "value1", "column2": "value2"},
		{"test_table_id": "3", "column1": "value1", "column2": "value2"},
	}
	if !reflect.DeepEqual(d1, c1) {
		t.Errorf("GetAllRs returned unexpected result. Expected: %v, Got: %v", d, c)
	}

	// 测试MGetBy函数
	e, err := model.MGetBy(map[string]interface{}{"column1": "value1"})
	if err != nil {
		t.Errorf("MGetBy failed: %v", err)
	}
	e0 := []map[string]interface{}{}
	for _, v := range e {
		e0 = append(e0, gmap.ToAny(v))
	}
	// 检查返回结果是否符合预期
	f := []map[string]interface{}{
		{"test_table_id": "1", "column1": "value1", "column2": "value2"},
		{"test_table_id": "2", "column1": "value1", "column2": "value2"},
		{"test_table_id": "3", "column1": "value1", "column2": "value2"},
	}
	if !reflect.DeepEqual(f, e0) {
		t.Errorf("MGetBy returned unexpected result. Expected: %v, Got: %v", f, e)
	}

	// 测试MGetByRs函数
	rs2, err := model.MGetByRs(map[string]interface{}{"column1": "value1"})
	if err != nil {
		t.Errorf("MGetByRs failed: %v", err)
	}
	e2 := []map[string]interface{}{}
	for _, v := range rs2.Rows() {
		e2 = append(e2, gmap.ToAny(v))
	}
	// 检查返回结果是否符合预期
	f2 := []map[string]interface{}{
		{"test_table_id": "1", "column1": "value1", "column2": "value2"},
		{"test_table_id": "2", "column1": "value1", "column2": "value2"},
		{"test_table_id": "3", "column1": "value1", "column2": "value2"},
	}
	if !reflect.DeepEqual(f2, e2) {
		t.Errorf("MGetByRs returned unexpected result. Expected: %v, Got: %v", f, e)
	}

	// 测试Insert函数
	lastInsertID, err := model.Insert(map[string]interface{}{"column1": "insert_value1"})
	if err != nil {
		t.Errorf("Insert failed: %v", err)
	}
	// 检查返回结果是否符合预期
	expectedID := int64(4)
	if lastInsertID != expectedID {
		t.Errorf("Insert returned unexpected result. Expected lastInsertID: %v, Got: %v", expectedID, lastInsertID)
	}

	// 测试DeleteBy函数
	cnt, err := model.DeleteBy(map[string]interface{}{"column1": "insert_value1"})
	if err != nil {
		t.Errorf("DeleteBy failed: %v", err)
	}
	// 检查返回结果是否符合预期
	expectedCount = 1
	if cnt != expectedCount {
		t.Errorf("DeleteBy returned unexpected result. Expected count: %d, Got: %d", expectedCount, cnt)
	}

	// 测试DeleteByIds函数
	model.Insert(map[string]interface{}{"column1": "DeleteByIds_value1"})
	row, _ := model.GetBy(map[string]interface{}{"column1": "DeleteByIds_value1"})
	cnt, err = model.DeleteByIDs([]string{row["test_table_id"]})
	if err != nil {
		t.Errorf("DeleteByIDs failed: %v", err)
	}
	// 检查返回结果是否符合预期
	expectedCount = 1
	if cnt != expectedCount {
		t.Errorf("DeleteByIDs returned unexpected result. Expected count: %d, Got: %d", expectedCount, cnt)
	}

	// 测试InsertBatch函数
	cnt, lastInsertID, err = model.InsertBatch([]map[string]interface{}{{"column1": "InsertBatch_value1"}, {"column1": "InsertBatch_value2"}})
	if err != nil {
		t.Errorf("InsertBatch failed: %v", err)
	}
	// 检查返回结果是否符合预期
	expectedCount = 2
	expectedID = int64(6)
	if cnt != expectedCount || lastInsertID != expectedID {
		t.Errorf("InsertBatch returned unexpected result. Expected count: %d, lastInsertID: %v, Got: %d, %v", expectedCount, expectedID, cnt, lastInsertID)
	}
	model.DeleteBy(map[string]interface{}{"column1": []string{"InsertBatch_value1", "InsertBatch_value2"}})

	// 测试UpdateByIDs函数
	cnt, err = model.UpdateByIDs([]string{"1", "2", "3"}, map[string]interface{}{"column1": "UpdateByIDs_value1"})
	if err != nil {
		t.Errorf("UpdateByIDs failed: %v", err)
	}
	// 检查返回结果是否符合预期
	expectedCount = 3
	if cnt != expectedCount {
		t.Errorf("UpdateByIDs returned unexpected result. Expected count: %d, Got: %d", expectedCount, cnt)
	}
	model.UpdateByIDs([]string{"1", "2", "3"}, map[string]interface{}{"column1": "value1"})

	// 测试UpdateBy函数
	cnt, err = model.UpdateBy(map[string]interface{}{"column1": "value1"}, map[string]interface{}{"column1": "UpdateBy_value1"})
	if err != nil {
		t.Errorf("UpdateBy failed: %v", err)
	}
	// 检查返回结果是否符合预期
	expectedCount = 3
	if cnt != expectedCount {
		t.Errorf("UpdateBy returned unexpected result. Expected count: %d, Got: %d", expectedCount, cnt)
	}
	model.UpdateBy(map[string]interface{}{"column1": "UpdateBy_value1"}, map[string]interface{}{"column1": "value1"})

	// 测试Page函数
	x, total, err := model.Page(map[string]interface{}{"column1": "value1"}, 0, 10)
	if err != nil {
		t.Errorf("Page failed: %v", err)
	}
	x0 := []map[string]interface{}{}
	for _, v := range x {
		x0 = append(x0, gmap.ToAny(v))
	}
	// 检查返回结果是否符合预期
	y := []map[string]interface{}{
		{"test_table_id": "1", "column1": "value1", "column2": "value2"},
		{"test_table_id": "2", "column1": "value1", "column2": "value2"},
		{"test_table_id": "3", "column1": "value1", "column2": "value2"},
	}
	expectedTotal := 3
	if !reflect.DeepEqual(y, x0) || total != expectedTotal {
		t.Errorf("Page returned unexpected result. Expected data: %v, total: %d, Got: %v, %d", y, expectedTotal, x, total)
	}

	// 测试List函数
	w, err := model.List(map[string]interface{}{"column1": "value1"}, 0, 10, "test_table_id", "desc")
	if err != nil {
		t.Errorf("List failed: %v", err)
	}
	w0 := []map[string]interface{}{}
	for _, v := range w {
		w0 = append(w0, gmap.ToAny(v))
	}
	// 检查返回结果是否符合预期
	z := []map[string]interface{}{
		{"test_table_id": "3", "column1": "value1", "column2": "value2"},
		{"test_table_id": "2", "column1": "value1", "column2": "value2"},
		{"test_table_id": "1", "column1": "value1", "column2": "value2"},
	}
	if !reflect.DeepEqual(z, w0) {
		t.Errorf("List returned unexpected result. Expected: %v, Got: %v", z, w)
	}

	// 清理测试用的表和数据
	dropTestTable(db)
}

func createTestDB() gcore.Database {
	// 返回您的测试用数据库连接对象
	cfg := gconfig.New()
	cfg.SetConfigFile("../app/app.toml")
	err := cfg.ReadInConfig()
	if err != nil {
		panic(err)
	}
	Init(cfg)
	db := DBMySQL()
	db.Config.TablePrefix = ""
	db.Config.TablePrefixSQLIdentifier = ""
	return db
}

func createTestTable(db gcore.Database) {
	// 创建测试表的 SQL 语句
	sqlStr := `
		CREATE TABLE test_table (
			test_table_id INT AUTO_INCREMENT,
			column1 VARCHAR(255),
			column2 VARCHAR(255),
			PRIMARY KEY (test_table_id)
		)
	`

	// 执行创建表的 SQL 语句
	_, err := db.ExecSQL(sqlStr)
	if err != nil {
		log.Fatalf("Failed to create test table: %v", err)
	}
}

func insertTestData(db gcore.Database) {
	// 插入测试数据的 SQL 语句
	sqlStr := `
		INSERT INTO test_table (column1, column2) VALUES
			('value1', 'value2'),
			('value1', 'value2'),
			('value1', 'value2')
	`

	// 执行插入数据的 SQL 语句
	_, err := db.ExecSQL(sqlStr)
	if err != nil {
		log.Fatalf("Failed to insert test data: %v", err)
	}
}

func dropTestTable(db gcore.Database) {
	// 执行删除测试表的 SQL 语句
	_, err := db.ExecSQL(`
		DROP TABLE IF EXISTS test_table
	`)
	if err != nil {
		panic(err)
	}
}
