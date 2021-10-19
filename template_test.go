package gotemplconstr

import (
	"testing"
	"time"
)

func Test_TemplateMarshalJSON(t *testing.T) {
	want := `{"node":{"attrs":[{"from":"","staticVal":"2020-12-12T12:12:12Z","to":"date"}],"bind":"","from":"","nodes":[{"attrs":[],"bind":"Data.Products","from":"","nodes":[{"attrs":[{"from":"Id","staticVal":null,"to":"id"}],"bind":"","from":"","nodes":[{"attrs":[],"bind":"","from":"Name","nodes":[],"staticVal":null,"to":"name"},{"attrs":[],"bind":"","from":"Id","nodes":[],"staticVal":null,"to":"id"},{"attrs":[],"bind":"","from":"Price","nodes":[],"staticVal":null,"to":"price"},{"attrs":[],"bind":"","from":"Count","nodes":[],"staticVal":null,"to":"amount"}],"staticVal":null,"to":"product"}],"staticVal":null,"to":"products"},{"attrs":[],"bind":"Data.Categories","from":"","nodes":[{"attrs":[{"from":"Id","staticVal":null,"to":"cat_id"}],"bind":"","from":"Name","nodes":[],"staticVal":null,"to":"category"}],"staticVal":null,"to":"categories"}],"staticVal":null,"to":"catalog"},"outputType":1}`

	date := time.Date(2020, 12, 12, 12, 12, 12, 0, time.UTC)

	tmpl := NewTemplate().ToOutputFileType(XML)

	tmpl.AddNode(
		Node().To("catalog").AddAttr(Attr().To("date").StaticVal(date)).AddNode(
			Node().To("products").Bind("Data.Products").AddNode(
				Node().To("product").AddAttr(Attr().To("id").From("Id")).AddNode(
					Node().To("name").From("Name"),
					Node().To("id").From("Id"),
					Node().To("price").From("Price"),
					Node().To("amount").From("Count"),
				),
			),
			Node().To("categories").Bind("Data.Categories").AddNode(
				Node().To("category").From("Name").AddAttr(Attr().To("cat_id").From("Id")),
			),
		),
	)

	bts, err := tmpl.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	if want != string(bts) {
		t.Fatalf("want: %s,\ngot: %s", want, string(bts))
	}
}

func Test_TemplateUnmarshalJSON(t *testing.T) {
	jsonText := `{"node":{"attrs":[{"from":"","staticVal":"2020-12-12T12:12:12Z","to":"date"}],"bind":"","from":"","nodes":[{"attrs":[],"bind":"Data.Products","from":"","nodes":[{"attrs":[{"from":"Id","staticVal":null,"to":"id"}],"bind":"","from":"","nodes":[{"attrs":[],"bind":"","from":"Name","nodes":[],"staticVal":null,"to":"name"},{"attrs":[],"bind":"","from":"Id","nodes":[],"staticVal":null,"to":"id"},{"attrs":[],"bind":"","from":"Price","nodes":[],"staticVal":null,"to":"price"},{"attrs":[],"bind":"","from":"Count","nodes":[],"staticVal":null,"to":"amount"}],"staticVal":null,"to":"product"}],"staticVal":null,"to":"products"},{"attrs":[],"bind":"Data.Categories","from":"","nodes":[{"attrs":[{"from":"Id","staticVal":null,"to":"cat_id"}],"bind":"","from":"Name","nodes":[],"staticVal":null,"to":"category"}],"staticVal":null,"to":"categories"}],"staticVal":null,"to":"catalog"},"outputType":1}`

	templ := &Template{}

	err := templ.UnmarshalJSON([]byte(jsonText))
	if err != nil {
		t.Fatal(err)
	}

	if templ.node.to != "catalog" {
		t.Fatal("first node is not catalog")
	}

	if len(templ.node.nodes) != 2 {
		t.Fatal("catalog nodes must be 2")
	}
}
