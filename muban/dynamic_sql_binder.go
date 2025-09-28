package main

import (
	"time"

	dynamicsql "github.com/NSObjects/go-template/muban/dynamic-sql-gen"
	"gorm.io/gen"
)

func init() {
	dynamicsql.SetInterfaceBinder(applyDynamicSQLInterfaces)
}

func applyDynamicSQLInterfaces(g *gen.Generator, models []interface{}) {
	if g == nil || len(models) == 0 {
		return
	}

	g.ApplyBasic(models...)

	interfaces := []interface{}{
		func(dynamicSQLCommonQuery) {},
		func(dynamicSQLPaginationQuery) {},
		func(dynamicSQLSearchQuery) {},
		func(dynamicSQLStatusQuery) {},
		func(dynamicSQLAdvancedQuery) {},
		func(dynamicSQLBusinessQuery) {},
	}

	for _, iface := range interfaces {
		g.ApplyInterface(iface, models...)
	}
}

type dynamicSQLCommonQuery interface {
	// SELECT * FROM @@table WHERE id = @id
	GetByID(id uint) (gen.T, error)

	// SELECT * FROM @@table WHERE id IN @ids
	GetByIDs(ids []uint) ([]gen.T, error)

	// SELECT COUNT(*) FROM @@table
	CountRecords() (int64, error)

	// SELECT 1 FROM @@table WHERE id = @id LIMIT 1
	Exists(id uint) (bool, error)

	// DELETE FROM @@table WHERE id = @id
	DeleteByID(id uint) error

	// DELETE FROM @@table WHERE id IN @ids
	DeleteByIDs(ids []uint) error
}

type dynamicSQLPaginationQuery interface {
	// SELECT * FROM @@table ORDER BY @orderBy LIMIT @limit OFFSET @offset
	GetPage(offset, limit int, orderBy string) ([]gen.T, error)

	// SELECT * FROM @@table WHERE @condition ORDER BY @orderBy LIMIT @limit OFFSET @offset
	GetPageWithCondition(condition string, offset, limit int, orderBy string) ([]gen.T, error)
}

type dynamicSQLSearchQuery interface {
	// SELECT * FROM @@table WHERE @field LIKE @keyword
	Search(field, keyword string) ([]gen.T, error)

	// SELECT * FROM @@table WHERE @field1 LIKE @keyword OR @field2 LIKE @keyword
	SearchMultiple(field1, field2, keyword string) ([]gen.T, error)
}

type dynamicSQLStatusQuery interface {
	// SELECT * FROM @@table WHERE status = @status
	GetByStatus(status int) ([]gen.T, error)

	// UPDATE @@table SET status = @status WHERE id = @id
	UpdateStatus(id uint, status int) error

	// SELECT * FROM @@table WHERE status = 1
	GetActive() ([]gen.T, error)

	// SELECT * FROM @@table WHERE status = 0
	GetInactive() ([]gen.T, error)
}

type dynamicSQLAdvancedQuery interface {
	// SELECT * FROM @@table
	// {{where}}
	//   {{if condition != ""}}
	//     @condition
	//   {{end}}
	// {{end}}
	FilterWithCondition(condition string) ([]gen.T, error)

	// SELECT * FROM @@table
	// {{if !start.IsZero()}}
	//   WHERE created_at > @start
	// {{end}}
	// {{if !end.IsZero()}}
	//   AND created_at < @end
	// {{end}}
	FilterWithTime(start, end time.Time) ([]gen.T, error)

	// UPDATE @@table
	// {{set}}
	//   {{if name != ""}} name=@name, {{end}}
	//   {{if age > 0}} age=@age, {{end}}
	//   updated_at=NOW()
	// {{end}}
	// WHERE id=@id
	UpdateWithSet(name string, age int, id uint) error
}

type dynamicSQLBusinessQuery interface {
	// SELECT * FROM @@table WHERE @@field = @value
	GetByField(field, value string) (gen.T, error)

	// SELECT * FROM @@table WHERE @@field1 = @value1 AND @@field2 = @value2
	GetByFields(field1, value1, field2, value2 string) ([]gen.T, error)

	// UPDATE @@table SET @@field = @value WHERE id IN @ids
	BatchUpdate(field, value string, ids []uint) error
}
