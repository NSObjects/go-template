package mysql

import (
	"testing"

	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func TestOpenMysqlConnectionDoesNotRegisterTemplateCreateCallback(t *testing.T) {
	db, err := openMysqlConnection(fakeDialector{})
	if err != nil {
		t.Fatalf("openMysqlConnection() error = %v", err)
	}

	if db.Callback().Create().Get("role:menu_after_create") != nil {
		t.Fatal("template create callback is registered, want no template callback on user mysql storage")
	}
}

type fakeDialector struct{}

func (fakeDialector) Name() string {
	return "fake"
}

func (fakeDialector) Initialize(db *gorm.DB) error {
	callbacks.RegisterDefaultCallbacks(db, &callbacks.Config{})
	return nil
}

func (fakeDialector) Migrator(*gorm.DB) gorm.Migrator {
	return nil
}

func (fakeDialector) DataTypeOf(*schema.Field) string {
	return ""
}

func (fakeDialector) DefaultValueOf(*schema.Field) clause.Expression {
	return clause.Expr{SQL: "DEFAULT"}
}

func (fakeDialector) BindVarTo(writer clause.Writer, _ *gorm.Statement, _ any) {
	writer.WriteByte('?')
}

func (fakeDialector) QuoteTo(writer clause.Writer, value string) {
	writer.WriteByte('`')
	writer.WriteString(value)
	writer.WriteByte('`')
}

func (fakeDialector) Explain(sql string, vars ...any) string {
	return logger.ExplainSQL(sql, nil, `"`, vars...)
}
