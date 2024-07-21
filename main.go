package main

import (
	"fmt"
	"log"

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
		for _, p := range products {
			fmt.Println(p.name, ":", p.url, p.price)
		}
	})
	c.Visit("https://www.scrapingcourse.com/ecommerce/")

}
