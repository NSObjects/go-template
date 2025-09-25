package dynamicsql

import (
	"fmt"
	"time"

	"github.com/NSObjects/go-template/internal/configs"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"

	"github.com/spf13/cobra"
)

// 通用查询接口 - 适用于所有模型的基础CRUD操作
type ICommonQuery interface {
	// GetByID
	// SELECT * FROM @@table WHERE id = @id
	GetByID(id uint) (gen.T, error)

	// GetByIDs
	// SELECT * FROM @@table WHERE id IN @ids
	GetByIDs(ids []uint) ([]gen.T, error)

	// CountRecords
	// SELECT COUNT(*) FROM @@table
	CountRecords() (int64, error)

	// Exists
	// SELECT 1 FROM @@table WHERE id = @id LIMIT 1
	Exists(id uint) (bool, error)

	// DeleteByID
	// DELETE FROM @@table WHERE id = @id
	DeleteByID(id uint) error

	// DeleteByIDs
	// DELETE FROM @@table WHERE id IN @ids
	DeleteByIDs(ids []uint) error
}

// 分页查询接口 - 适用于需要分页的模型
type IPaginationQuery interface {
	// GetPage
	// SELECT * FROM @@table ORDER BY @orderBy LIMIT @limit OFFSET @offset
	GetPage(offset, limit int, orderBy string) ([]gen.T, error)

	// GetPageWithCondition
	// SELECT * FROM @@table WHERE @condition ORDER BY @orderBy LIMIT @limit OFFSET @offset
	GetPageWithCondition(condition string, offset, limit int, orderBy string) ([]gen.T, error)
}

// 搜索查询接口 - 适用于需要搜索功能的模型
type ISearchQuery interface {
	// Search
	// SELECT * FROM @@table WHERE @field LIKE @keyword
	Search(field, keyword string) ([]gen.T, error)

	// SearchMultiple
	// SELECT * FROM @@table WHERE @field1 LIKE @keyword OR @field2 LIKE @keyword
	SearchMultiple(field1, field2, keyword string) ([]gen.T, error)
}

// 状态查询接口 - 适用于有状态字段的模型
type IStatusQuery interface {
	// GetByStatus
	// SELECT * FROM @@table WHERE status = @status
	GetByStatus(status int) ([]gen.T, error)

	// UpdateStatus
	// UPDATE @@table SET status = @status WHERE id = @id
	UpdateStatus(id uint, status int) error

	// GetActive
	// SELECT * FROM @@table WHERE status = 1
	GetActive() ([]gen.T, error)

	// GetInactive
	// SELECT * FROM @@table WHERE status = 0
	GetInactive() ([]gen.T, error)
}

// 高级查询接口 - 支持模板表达式
type IAdvancedQuery interface {
	// FilterWithCondition - 使用where模板表达式
	// SELECT * FROM @@table
	// {{where}}
	//   {{if condition != ""}}
	//     @condition
	//   {{end}}
	// {{end}}
	FilterWithCondition(condition string) ([]gen.T, error)

	// FilterWithTime - 使用if/else模板表达式
	// SELECT * FROM @@table
	// {{if !start.IsZero()}}
	//   WHERE created_at > @start
	// {{end}}
	// {{if !end.IsZero()}}
	//   AND created_at < @end
	// {{end}}
	FilterWithTime(start, end time.Time) ([]gen.T, error)

	// UpdateWithSet - 使用set模板表达式
	// UPDATE @@table
	// {{set}}
	//   {{if name != ""}} name=@name, {{end}}
	//   {{if age > 0}} age=@age, {{end}}
	//   updated_at=NOW()
	// {{end}}
	// WHERE id=@id
	UpdateWithSet(name string, age int, id uint) error
}

// 业务特定查询接口 - 适用于特定业务场景
type IBusinessQuery interface {
	// GetByField - 通用字段查询
	// SELECT * FROM @@table WHERE @@field = @value
	GetByField(field, value string) (gen.T, error)

	// GetByFields - 多字段查询
	// SELECT * FROM @@table WHERE @@field1 = @value1 AND @@field2 = @value2
	GetByFields(field1, value1, field2, value2 string) ([]gen.T, error)

	// BatchUpdate - 批量更新
	// UPDATE @@table SET @@field = @value WHERE id IN @ids
	BatchUpdate(field, value string, ids []uint) error
}

type Options struct {
	Config string
}

func NewCommand() *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:   "dynamic-sql",
		Short: "Generate dynamic SQL helpers based on the current database schema",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(opts)
		},
	}

	cmd.Flags().StringVar(&opts.Config, "config", "configs/config.toml", "config file path")

	return cmd
}

func Run(opts Options) error {
	cfg := configs.NewCfg(opts.Config)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Mysql.User,
		cfg.Mysql.Password,
		cfg.Mysql.Host,
		cfg.Mysql.Port,
		cfg.Mysql.Database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:           "./internal/api/data/query",
		Mode:              gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
	})

	g.UseDB(db)

	g.ApplyBasic(g.GenerateAllTable()...)

	g.ApplyInterface(func(ICommonQuery) {}, g.GenerateAllTable()...)
	g.ApplyInterface(func(IPaginationQuery) {}, g.GenerateAllTable()...)
	g.ApplyInterface(func(ISearchQuery) {}, g.GenerateAllTable()...)
	g.ApplyInterface(func(IStatusQuery) {}, g.GenerateAllTable()...)
	g.ApplyInterface(func(IAdvancedQuery) {}, g.GenerateAllTable()...)
	g.ApplyInterface(func(IBusinessQuery) {}, g.GenerateAllTable()...)

	g.Execute()

	fmt.Println("Basic Dynamic SQL generation completed!")

	return nil
}
