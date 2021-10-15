package main

import (
	"bytes"
	"time"

	"github.com/tagirmukail/dionysus"
)

func main() {

	d := struct {
		Data data
	}{
		Data: data{
			Products: []product{
				{
					Id:    1,
					Name:  "test1",
					Price: "234.23",
					Count: 12,
				},
				{
					Id:    2,
					Name:  "test2",
					Price: "5454.234",
					Count: 30,
				},
			},
			Categories: []category{
				{
					Id:   1,
					Name: "testcat1",
				},
				{
					Id:   2,
					Name: "testcat2",
				},
			},
		},
	}

	tmpl := &dionysus.Template{}

	tmpl.ToOutputFileType(dionysus.XML)

	catalogNode := dionysus.Node().To("catalog").AddArg(dionysus.Arg().To("date").StaticVal(time.Date(2020, 12, 12, 12, 12, 12, 0, time.UTC)))

	productsNode := dionysus.Node().To("products").Bind("Data.Products")
	productsNode = productsNode.AddNode(dionysus.Node().To("product").AddArg(dionysus.Arg().To("id").From("Id")).
		AddNode(dionysus.Node().To("name").From("Name")).
		AddNode(dionysus.Node().To("id").From("Id")).
		AddNode(dionysus.Node().To("price").From("Price")).
		AddNode(dionysus.Node().To("amount").From("Count")))
	catalogNode = catalogNode.AddNode(productsNode)

	categoriesNode := dionysus.Node().To("categories").Bind("Data.Categories")
	categoriesNode = categoriesNode.AddNode(dionysus.Node().To("category").From("Name").AddArg(dionysus.Arg().To("cat_id").From("Id")))
	catalogNode = catalogNode.AddNode(categoriesNode)

	tmpl.AddNode(catalogNode)

	buf := &bytes.Buffer{}

	err := tmpl.NewEncoder(buf).Encode(d)
	if err != nil {
		panic(err)
	}

	gotXML := buf.String()
	println(gotXML)
}

type product struct {
	Id    int
	Name  string
	Price string
	Count int
}

type category struct {
	Id   int
	Name string
}

type data struct {
	Products   []product
	Categories []category
}
