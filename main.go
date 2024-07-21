package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly/v2"
)

type Product struct {
	url, name, price, img string
}

func main() {
	fmt.Println("Hello")
	c := colly.NewCollector()
	var products []Product

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting: ", r.URL)
	})
	c.OnError(func(_ *colly.Response, err error) {
		log.Fatalln("Something went wrong: ", err)
	})
	c.OnHTML("li.product", func(h *colly.HTMLElement) {
		product := Product{}
		product.img = h.ChildAttr("img", "src")
		product.url = h.ChildAttr("a", "href")
		product.name = h.ChildText("h2")
		product.price = h.ChildText(".price")
		products = append(products, product)
	})
	c.OnScraped(func(r *colly.Response) {

		file, err := os.Create("products.csv")
		if err != nil {
			log.Fatalln(err)
		}
		defer file.Close()

		writer := csv.NewWriter(file)

		headers := []string{
			"Name",
			"Price",
			"URL",
			"Image",
		}
		writer.Write(headers)

		for _, p := range products {
			fmt.Println(p.name, ":", p.url, p.price)
			record := []string{p.name, p.price, p.url, p.img}
			writer.Write(record)
		}

		defer writer.Flush()
	})
	c.Visit("https://www.scrapingcourse.com/ecommerce/")

}
