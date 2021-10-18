# Gotemplconstr

Is a template constructor. The package provides APIs for creating flexible templates, which can later be stored in json format. To do this, the methods `MarshalJSON` and `UnmarshalJSON` are presented.

### Get started
```go
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

	date := time.Date(2020, 12, 12, 12, 12, 12, 0, time.UTC)

	tmpl := gotemplconstr.NewTemplate().ToOutputFileType(gotemplconstr.XML)

	tmpl.AddNode(
		gotemplconstr.Node().To("catalog").AddAttr(gotemplconstr.Attr().To("date").StaticVal(date)).AddNode(
			gotemplconstr.Node().To("products").Bind("Data.Products").AddNode(
				gotemplconstr.Node().To("product").AddAttr(gotemplconstr.Attr().To("id").From("Id")).AddNode(
					gotemplconstr.Node().To("name").From("Name"),
					gotemplconstr.Node().To("id").From("Id"),
					gotemplconstr.Node().To("price").From("Price"),
					gotemplconstr.Node().To("amount").From("Count"),
				),
			),
			gotemplconstr.Node().To("categories").Bind("Data.Categories").AddNode(
				gotemplconstr.Node().To("category").From("Name").AddAttr(gotemplconstr.Attr().To("cat_id").From("Id")),
			),
		),
	)

	buf := &bytes.Buffer{}

	err := tmpl.NewEncoder(buf).Encode(d)
	if err != nil {
		panic(err)
	}

	gotXML := buf.String()
	println(gotXML)
}
```

Encoding result:
```xml
<?xml version="1.0" encoding="UTF-8"?>
<catalog date="2020-12-12 12:12:12 +0000 UTC">
    <products>
        <product id="1">
            <name>test1</name>
            <id>1</id>
            <price>234.23</price>
            <amount>12</amount>
        </product>
        <product id="2">
            <name>test2</name>
            <id>2</id>
            <price>5454.234</price>
            <amount>30</amount>
        </product>
    </products>
    <categories>
        <category cat_id="1">testcat1</category>
        <category cat_id="2">testcat2</category>
    </categories>
</catalog>
```

More examples in [cmd](./cmd)

### Supported output formats

- xml
- yaml

#### Coming soon:

- `json`
- `csv`
