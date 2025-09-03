/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/NSObjects/echo-admin/internal/configs"

	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := configs.InitConfig(cfgFile)
		if err != nil {
			panic(err)
		}
		GenMysql(configs.Mysql)
		fmt.Println("gen called")
	},
}

func init() {
	rootCmd.AddCommand(genCmd)
}

// Querier Dynamic SQL
type Querier interface {
	// GetById
	// SELECT * FROM @@table WHERE id = @id
	GetById(id int) (gen.T, error)

	// DeleteByID
	// DELETE FROM @@table WHERE id = @id
	DeleteByID(id int64) error
}

func GenMysql(cfg configs.MysqlConfig) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	// 数据库迁移需要根据实际业务模型来实现
	// err = db.AutoMigrate(&model.YourModel{})

	if err != nil {
		panic(err)
	}
	g := gen.NewGenerator(gen.Config{
		OutPath:      "./internal/api/data/query", // output path
		ModelPkgPath: "./internal/api/data/model", // model package path
		WithUnitTest: true,
		Mode:         gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode

	})
	g.UseDB(db)

	// Generate basic type-safe DAO API for business models
	// g.ApplyBasic(model.User{}, model.Product{})
	// Generate Type Safe API with Dynamic SQL defined on Querier interface
	// g.ApplyInterface(func(Querier) {}, model.User{}, model.Product{})

	// Generate the code
	g.Execute()
}
