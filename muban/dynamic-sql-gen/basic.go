package dynamicsql

import (
	"fmt"
	"strings"
	"sync/atomic"
	"unicode"

	"github.com/NSObjects/go-template/internal/configs"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// InterfaceBinder controls how query interfaces are attached to generated models.
type InterfaceBinder func(*gen.Generator, []interface{})

var interfaceBinder atomic.Value

func init() {
	interfaceBinder.Store(InterfaceBinder(defaultInterfaceBinder))
}

// SetInterfaceBinder overrides the default interface binding behavior used during
// dynamic SQL generation. Passing nil resets the default behavior which only
// registers basic CRUD methods.
func SetInterfaceBinder(binder InterfaceBinder) {
	if binder == nil {
		interfaceBinder.Store(InterfaceBinder(defaultInterfaceBinder))
		return
	}
	interfaceBinder.Store(binder)
}

func defaultInterfaceBinder(g *gen.Generator, models []interface{}) {
	if g == nil || len(models) == 0 {
		return
	}

	g.ApplyBasic(models...)
}

type Options struct {
	Config string
	Table  string

	OutPath  string
	ModelPkg string
	JSONTag  string

	WithContext  bool
	WithUnitTest bool

	FieldNullable     bool
	FieldCoverable    bool
	FieldSignable     bool
	FieldWithIndexTag bool
	FieldWithTypeTag  bool
}

func NewCommand() *cobra.Command {
	opts := Options{
		OutPath:           "./internal/api/data/query",
		JSONTag:           "snake",
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
	}

	cmd := &cobra.Command{
		Use:   "dynamicsql",
		Short: "Generate dynamic SQL helpers based on the current database schema",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(opts)
		},
	}

	cmd.Flags().StringVar(&opts.Config, "config", "configs/config.toml", "config file path")
	cmd.Flags().StringVarP(&opts.Table, "table", "t", "", "specify target table (comma separated for multiple)")
	cmd.Flags().StringVar(&opts.OutPath, "out", opts.OutPath, "output directory for generated query code")
	cmd.Flags().StringVar(&opts.ModelPkg, "model-pkg", opts.ModelPkg, "package import path for generated models")
	cmd.Flags().StringVar(&opts.JSONTag, "json-tag-style", opts.JSONTag, "json tag naming strategy: snake, camel, pascal, none")
	cmd.Flags().BoolVar(&opts.WithContext, "with-context", opts.WithContext, "generate query code with context support")
	cmd.Flags().BoolVar(&opts.WithUnitTest, "with-unit-test", opts.WithUnitTest, "generate unit tests for query helpers")
	cmd.Flags().BoolVar(&opts.FieldNullable, "field-nullable", opts.FieldNullable, "generate pointer for nullable database columns")
	cmd.Flags().BoolVar(&opts.FieldCoverable, "field-coverable", opts.FieldCoverable, "generate pointer for columns with default values")
	cmd.Flags().BoolVar(&opts.FieldSignable, "field-signable", opts.FieldSignable, "detect unsigned integer fields")
	cmd.Flags().BoolVar(&opts.FieldWithIndexTag, "field-index-tag", opts.FieldWithIndexTag, "include gorm index tags in generated models")
	cmd.Flags().BoolVar(&opts.FieldWithTypeTag, "field-type-tag", opts.FieldWithTypeTag, "include gorm column type tags in generated models")

	return cmd
}

func Run(opts Options) error {
	cfg := configs.NewCfg(opts.Config)

	db, cleanup, err := openDatabase(cfg)
	if err != nil {
		return err
	}
	defer func() {
		if cleanup == nil {
			return
		}
		if err := cleanup(); err != nil {
			fmt.Printf("warning: failed to close database connection: %v\n", err)
		}
	}()

	generator, err := buildGenerator(opts)
	if err != nil {
		return err
	}

	generator.UseDB(db)

	models, err := getTargetModels(generator, opts.Table)
	if err != nil {
		return err
	}
	if len(models) == 0 {
		return fmt.Errorf("no models generated, please verify database schema and table filters")
	}

	binder := interfaceBinder.Load().(InterfaceBinder)
	binder(generator, models)

	if err := executeGenerator(generator); err != nil {
		return err
	}

	fmt.Printf("Dynamic SQL generation completed for %d model(s) at %s\n", len(models), opts.OutPath)

	return nil
}

func buildGenerator(opts Options) (*gen.Generator, error) {
	cfg := gen.Config{
		OutPath:           opts.OutPath,
		ModelPkgPath:      opts.ModelPkg,
		WithUnitTest:      opts.WithUnitTest,
		FieldNullable:     opts.FieldNullable,
		FieldCoverable:    opts.FieldCoverable,
		FieldSignable:     opts.FieldSignable,
		FieldWithIndexTag: opts.FieldWithIndexTag,
		FieldWithTypeTag:  opts.FieldWithTypeTag,
		Mode:              gen.WithDefaultQuery | gen.WithQueryInterface,
	}
	if !opts.WithContext {
		cfg.Mode |= gen.WithoutContext
	}

	jsonStrategy, err := buildJSONTagStrategy(opts.JSONTag)
	if err != nil {
		return nil, err
	}
	cfg.WithJSONTagNameStrategy(jsonStrategy)

	generator := gen.NewGenerator(cfg)

	return generator, nil
}

func getTargetModels(g *gen.Generator, tableOpt string) ([]interface{}, error) {
	if tableOpt == "" {
		models := sanitizeModels(g.GenerateAllTable())
		if len(models) == 0 {
			return nil, fmt.Errorf("no tables found in database")
		}
		return models, nil
	}

	tableNames := parseTableNames(tableOpt)
	if len(tableNames) == 0 {
		return nil, fmt.Errorf("no valid table names provided")
	}

	models := make([]interface{}, 0, len(tableNames))
	skipped := make([]string, 0)

	for _, name := range tableNames {
		model, err := safeGenerateModel(g, name)
		if err != nil {
			skipped = append(skipped, fmt.Sprintf("%s (%v)", name, err))
			continue
		}
		models = append(models, model)
	}

	if len(models) == 0 {
		return nil, fmt.Errorf("failed to generate models for tables: %s", strings.Join(skipped, ", "))
	}

	if len(skipped) > 0 {
		fmt.Printf("warning: skipped tables during generation: %s\n", strings.Join(skipped, "; "))
	}

	return models, nil
}

func sanitizeModels(models []interface{}) []interface{} {
	filtered := make([]interface{}, 0, len(models))
	for _, m := range models {
		if m != nil {
			filtered = append(filtered, m)
		}
	}
	return filtered
}

func parseTableNames(tableOpt string) []string {
	if tableOpt == "" {
		return nil
	}

	raw := strings.Split(tableOpt, ",")
	names := make([]string, 0, len(raw))
	for _, item := range raw {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			names = append(names, trimmed)
		}
	}
	return names
}

func safeGenerateModel(g *gen.Generator, table string) (model interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("generate model panic: %v", r)
		}
	}()

	model = g.GenerateModel(table)
	if model == nil {
		return nil, fmt.Errorf("table not found or ignored")
	}
	return model, nil
}

func executeGenerator(g *gen.Generator) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("generator execution failed: %v", r)
		}
	}()

	g.Execute()
	return nil
}

func buildJSONTagStrategy(style string) (func(string) string, error) {
	normalized := strings.TrimSpace(strings.ToLower(style))
	naming := schema.NamingStrategy{}

	switch normalized {
	case "", "snake":
		return func(column string) string {
			return naming.ColumnName("", column)
		}, nil
	case "camel":
		return func(column string) string {
			return lowerFirstRune(naming.SchemaName(column))
		}, nil
	case "pascal":
		return func(column string) string {
			return naming.SchemaName(column)
		}, nil
	case "none":
		return func(column string) string {
			return column
		}, nil
	default:
		return nil, fmt.Errorf("unsupported json tag style: %s", style)
	}
}

func lowerFirstRune(s string) string {
	if s == "" {
		return s
	}

	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

func openDatabase(cfg configs.Config) (*gorm.DB, func() error, error) {
	mysqlCfg := cfg.Mysql
	if mysqlCfg.Database == "" {
		return nil, nil, fmt.Errorf("mysql database is not configured")
	}

	host := mysqlCfg.Host
	if cfg.System.Env == "docker" && mysqlCfg.DockerHost != "" {
		host = mysqlCfg.DockerHost
	}
	if host == "" {
		return nil, nil, fmt.Errorf("mysql host is not configured")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlCfg.User,
		mysqlCfg.Password,
		host,
		mysqlCfg.Port,
		mysqlCfg.Database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	if mysqlCfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(mysqlCfg.MaxOpenConns)
	}
	if mysqlCfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(mysqlCfg.MaxIdleConns)
	}

	cleanup := func() error {
		return sqlDB.Close()
	}

	return db, cleanup, nil
}
