package main

import (
	"Dionysus"
	"log"
	"os"
	"time"
)

func main() {

	tmpl := &Dionysus.Template{}

	tmpl.ToOutputFileType(Dionysus.XML).
		AddNode(
			Dionysus.Node{}.To("catalog").AddArg(Dionysus.Arg{}.To("date").StaticVal(time.Now())).
				AddNode(
					Dionysus.Node{}.To("products").
						AddNode(Dionysus.Node{}.To("product").AddArg(Dionysus.Arg{}.To("id").BindTo("data.id")).
							AddNode(Dionysus.Node{}.To("name").BindTo("data.products.name")).
							AddNode(Dionysus.Node{}.To("id").BindTo("data.products.id")).
							AddNode(Dionysus.Node{}.To("price").BindTo("data.products.price")),
						),
				).
				AddNode(
					Dionysus.Node{}.To("categories").
						AddNode(Dionysus.Node{}.To("category").BindTo("data.categories.name").
							AddArg(Dionysus.Arg{}.To("id").BindTo("data.categories.id")),
						),
				),
		)

	f, _ := os.Create("output." + tmpl.FileType().String())
	defer f.Close()

	err := Dionysus.NewEncoder(f).Encode(tmpl)
	if err != nil {
		log.Fatalln(err)
	}
}
