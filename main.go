package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
)

const (
	package_name = "user"
	fields       = "ID:int64,Name:string,Email:string,CreatedAt:time.Time,UpdateAt:time.Time"

	modelTypeEntity modelType = "entity"
	modelTypeModel  modelType = "model"
	modelTypeView   modelType = "view"

	nameTypeSep = ":"
	fieldSep    = ","
)

type modelType string

func (t modelType) String() string {
	return string(t)
}

func main() {
	if err := generateModelFile("User", "sample/entity/user.go", modelTypeEntity); err != nil {
		log.Fatalf("failed to save: %v", err)
	}
	if err := generateModelFile("User", "sample/model/user.go", modelTypeModel); err != nil {
		log.Fatalf("failed to save: %v", err)
	}
	if err := generateModelFile("User", "sample/view/user.go", modelTypeView); err != nil {
		log.Fatalf("failed to save: %v", err)
	}
}

func generateModelFile(name, path string, t modelType) error {
	f := jen.NewFile(package_name)

	f.ImportName("time", "time")

	fs := NewFields(fields)

	var jenFields []jen.Code
	for _, f := range fs {
		jenFields = append(jenFields, f.ToJen(t))
	}

	f.Type().Id(name).Struct(jenFields...)

	return f.Save(path)
}

type Field struct {
	Name string
	Type string
}

func (f Field) ToJen(t modelType) *jen.Statement {
	s := jen.Id(f.Name)
	switch f.Type {
	case "time.Time":
		s = s.Qual("time", "Time")
	default:
		s = s.Op(f.Type)
	}

	switch t {
	case modelTypeEntity:
		return s.Tag(gormTag(f.Name))
	case modelTypeModel:
		return s
	case modelTypeView:
		return s.Tag(jsonTag(f.Name))
	default:
		return s
	}
}

func NewField(nameAndType string) Field {
	s := strings.Split(nameAndType, nameTypeSep)
	return Field{
		Name: s[0],
		Type: s[1],
	}
}

func NewFields(nameAndTypes string) (fs []Field) {
	nats := strings.Split(nameAndTypes, fieldSep)
	for _, nat := range nats {
		fs = append(fs, NewField(nat))
	}
	return fs
}

func jsonTag(field string) map[string]string {
	return map[string]string{
		"json": strcase.ToSnake(field),
	}
}

func gormTag(field string) map[string]string {
	return map[string]string{
		"gorm": fmt.Sprintf("column:%s", strcase.ToSnake(field)),
	}
}
