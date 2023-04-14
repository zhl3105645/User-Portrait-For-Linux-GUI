package main

import (
	dal2 "backend/cmd/dal"
	"gorm.io/gen"
)

func main() {
	dal2.Init()

	g := gen.NewGenerator(gen.Config{
		OutPath:           "./cmd/dal/query",
		Mode:              gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable:     true,
		FieldCoverable:    false,
		FieldSignable:     false,
		FieldWithIndexTag: false,
		FieldWithTypeTag:  true,
	})

	g.UseDB(dal2.DB)

	// 自定义字段的数据类型
	// 统一数字类型为int64,兼容protobuf
	dataMap := map[string]func(detailType string) (dataType string){
		"tinyint":   func(detailType string) (dataType string) { return "int64" },
		"smallint":  func(detailType string) (dataType string) { return "int64" },
		"mediumint": func(detailType string) (dataType string) { return "int64" },
		"bigint":    func(detailType string) (dataType string) { return "int64" },
		"int":       func(detailType string) (dataType string) { return "int64" },
	}
	g.WithDataTypeMap(dataMap)

	// 创建模型的结构体,生成文件在 model 目录; 先创建的结果会被后面创建的覆盖
	// 这里创建个别模型仅仅是为了拿到`*generate.QueryStructMeta`类型对象用于后面的模型关联操作中
	Test := g.GenerateModel("test")
	App := g.GenerateModel("app")
	Account := g.GenerateModel("account")
	Component := g.GenerateModel("component")
	DataSource := g.GenerateModel("data_source")
	Label := g.GenerateModel("label")
	LabelData := g.GenerateModel("label_data")
	DataModel := g.GenerateModel("data_model")
	ModelData := g.GenerateModel("model_data")
	Rule := g.GenerateModel("rule")
	RuleElement := g.GenerateModel("rule_element")
	User := g.GenerateModel("user")

	// 创建有关联关系的模型文件
	// 可以用于指定外键
	//Score := g.GenerateModel("score",
	//	append(
	//		fieldOpts,
	//		// user 一对多 address 关联, 外键`uid`在 address 表中
	//		gen.FieldRelate(field.HasMany, "user", User, &field.RelateConfig{GORMTag: "foreignKey:UID"}),
	//	)...,
	//)
	g.ApplyBasic(Test, App, Account, Component, DataSource, Label, LabelData, DataModel, ModelData, Rule, RuleElement, User)

	g.ApplyInterface(func(dal2.QueryAll) {}, Test)

	g.Execute()
}
