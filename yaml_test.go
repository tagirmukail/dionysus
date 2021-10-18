package gotemplconstr

import (
	"bytes"
	"testing"
)

func TestTemplate_encodeYAML(t *testing.T) {
	want := `first_service:
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
`

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

	tmpl := &Template{}

	tmpl.ToOutputFileType(YAML)

	tmpl = tmpl.AddNode(Node().To("first_service").
		AddNode(
			Node().To("debug").StaticVal(true),
			Node().To("db").Bind("FooService.Database").AddNode(
				Node().To("host").From("Hostname"),
				Node().To("user").From("Username"),
				Node().To("port").From("Port"),
				Node().To("passwd").From("Password"),
				Node().To("db_name").From("Name"),
			),
			Node().To("broker").Bind("FooService.Broker").AddNode(
				Node().To("address").From("Addr"),
				Node().To("user").From("User"),
				Node().To("Password").From("Password"),
			),
			Node().To("bar_items").Bind("FooService.Bars").AddNode(
				Node().To("").AddNode(
					Node().To("field_1").From("Field1"),
					Node().To("field_2").From("Field2"),
				),
			),
		),
	)

	buf := &bytes.Buffer{}

	err := tmpl.NewEncoder(buf).Encode(cfg)
	if err != nil {
		t.Fatal(err)
	}

	if want != buf.String() {
		t.Fatalf("want:%s\ngot:%s", want, buf.String())
	}
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
