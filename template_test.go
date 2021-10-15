package dionysus

import (
	"testing"
	"time"
)

func Test_TemplateMarshalJSON(t *testing.T) {
	want := `{"node":{"args":[{"from":"","staticVal":"2020-12-12T12:12:12Z","to":"date"}],"from":"","nodes":[{"args":[],"from":"","nodes":[{"args":[{"from":"Id","staticVal":null,"to":"id"}],"from":"","nodes":[{"args":[],"from":"Name","nodes":[],"staticVal":null,"to":"name"},{"args":[],"from":"Id","nodes":[],"staticVal":null,"to":"id"},{"args":[],"from":"Price","nodes":[],"staticVal":null,"to":"price"},{"args":[],"from":"Count","nodes":[],"staticVal":null,"to":"amount"}],"staticVal":null,"to":"product"}],"staticVal":null,"to":"products"},{"args":[],"from":"","nodes":[{"args":[{"from":"Id","staticVal":null,"to":"cat_id"}],"from":"Name","nodes":[],"staticVal":null,"to":"category"}],"staticVal":null,"to":"categories"}],"staticVal":null,"to":"catalog"},"outputType":1}`

	tmpl := &Template{}

	tmpl.ToOutputFileType(XML)

	catalogNode := Node().To("catalog").AddArg(Arg().To("date").StaticVal(time.Date(2020, 12, 12, 12, 12, 12, 0, time.UTC)))

	productsNode := Node().To("products")
	productsNode = productsNode.AddNode(Node().To("product").AddArg(Arg().To("id").From("Id")).
		AddNode(Node().To("name").From("Name")).
		AddNode(Node().To("id").From("Id")).
		AddNode(Node().To("price").From("Price")).
		AddNode(Node().To("amount").From("Count")))
	catalogNode = catalogNode.AddNode(productsNode)

	categoriesNode := Node().To("categories")
	categoriesNode = categoriesNode.AddNode(Node().To("category").From("Name").AddArg(Arg().To("cat_id").From("Id")))
	catalogNode = catalogNode.AddNode(categoriesNode)

	tmpl.AddNode(catalogNode)

	bts, err := tmpl.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	if want != string(bts) {
		t.Fatalf("want: %s,\ngot: %s", want, string(bts))
	}
}

func Test_TemplateUnmarshalJSON(t *testing.T) {
	jsonText := `{"node":{"args":[{"from":"","staticVal":"2020-12-12T12:12:12Z","to":"date"}],"from":"","nodes":[{"args":[],"from":"","nodes":[{"args":[{"from":"Id","staticVal":null,"to":"id"}],"from":"","nodes":[{"args":[],"from":"Name","nodes":[],"staticVal":null,"to":"name"},{"args":[],"from":"Id","nodes":[],"staticVal":null,"to":"id"},{"args":[],"from":"Price","nodes":[],"staticVal":null,"to":"price"},{"args":[],"from":"Count","nodes":[],"staticVal":null,"to":"amount"}],"staticVal":null,"to":"product"}],"staticVal":null,"to":"products"},{"args":[],"from":"","nodes":[{"args":[{"from":"Id","staticVal":null,"to":"cat_id"}],"from":"Name","nodes":[],"staticVal":null,"to":"category"}],"staticVal":null,"to":"categories"}],"staticVal":null,"to":"catalog"},"outputType":1}`

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
