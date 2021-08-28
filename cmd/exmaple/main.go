package main

import (
	"github.com/tagirmukail/dionysus"
	"log"
	"os"
	"time"
)

func main() {

	tmpl := &dionysus.Template{}

	tmpl.ToOutputFileType(dionysus.XML).
		AddNode(
			dionysus.Node{}.To("catalog").AddArg(dionysus.Arg{}.To("date").StaticVal(time.Now())).
				AddNode(
					dionysus.Node{}.To("products").
						AddNode(dionysus.Node{}.To("product").AddArg(dionysus.Arg{}.To("id").BindTo("data.id")).
							AddNode(dionysus.Node{}.To("name").BindTo("data.products.name")).
							AddNode(dionysus.Node{}.To("id").BindTo("data.products.id")).
							AddNode(dionysus.Node{}.To("price").BindTo("data.products.price")),
						),
				).
				AddNode(
					dionysus.Node{}.To("categories").
						AddNode(dionysus.Node{}.To("category").BindTo("data.categories.name").
							AddArg(dionysus.Arg{}.To("id").BindTo("data.categories.id")),
						),
				),
		)

	f, _ := os.Create("output." + tmpl.FileType().String())
	defer f.Close()

	err := dionysus.NewEncoder(f).Encode(tmpl)
	if err != nil {
		log.Fatalln(err)
	}
}
