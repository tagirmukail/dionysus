package main

import (
	"bytes"

	"github.com/tagirmukail/gotemplconstr"
)

/* RESULT
first_service:
  debug: true
  db:
    host:
      "localhost"
    user:
      "postgres"
    port:
      "5432"
    passwd:
      "password"
    db_name:
      "postgres"
  broker:
    address:
      "localhost:8888"
    user:
      "admin"
    Password:
      "admin"
  bar_items:
      - field_1: "value1"
        field_2: "value2"
      - field_1: "value11"
        field_2: "value22"

*/
func main() {
	cfg := configuration{
		FooService: serviceConfiguration{
			Database: Database{
				Hostname: "localhost",
				Port:     "5432",
				Username: "postgres",
				Password: "password",
				Name:     "postgres",
			},
			Broker: Broker{
				Addr:     "localhost:8888",
				User:     "admin",
				Password: "admin",
			},
			Bars: []Bar{
				{
					Field1: "value1",
					Field2: "value2",
				},
				{
					Field1: "value11",
					Field2: "value22",
				},
			},
		},
	}

	tmpl := &gotemplconstr.Template{}

	tmpl.ToOutputFileType(gotemplconstr.YAML)

	tmpl = tmpl.AddNode(gotemplconstr.Node().To("first_service").
		AddNode(
			gotemplconstr.Node().To("db").Bind("FooService.Database").AddNode(
				gotemplconstr.Node().To("host").From("Hostname"),
				gotemplconstr.Node().To("user").From("Username"),
				gotemplconstr.Node().To("port").From("Port"),
				gotemplconstr.Node().To("passwd").From("Password"),
				gotemplconstr.Node().To("db_name").From("Name"),
			),
			gotemplconstr.Node().To("broker").Bind("FooService.Broker").AddNode(
				gotemplconstr.Node().To("address").From("Addr"),
				gotemplconstr.Node().To("user").From("User"),
				gotemplconstr.Node().To("Password").From("Password"),
			),
			gotemplconstr.Node().To("bar_items").Bind("FooService.Bars").AddNode(
				gotemplconstr.Node().To("").AddNode(
					gotemplconstr.Node().To("field_1").From("Field1"),
					gotemplconstr.Node().To("field_2").From("Field2"),
				),
			),
		),
	)

	buf := &bytes.Buffer{}

	err := tmpl.NewEncoder(buf).Encode(cfg)
	if err != nil {
		panic(err)
	}

	print(buf.String())
}

type configuration struct {
	FooService serviceConfiguration
}

type serviceConfiguration struct {
	Database Database
	Broker   Broker
	Bars     []Bar
}

type Database struct {
	Hostname string
	Port     string
	Username string
	Password string
	Name     string
}

type Broker struct {
	Addr     string
	User     string
	Password string
}

type Bar struct {
	Field1 string
	Field2 string
}
