package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
)

type Product struct {
	url, name, price, img string
}

var wg sync.WaitGroup

func main() {
	start := time.Now()
	productRecords := make(chan Product)
	// pagesToScrape := make(chan string)
	wg.Add(2)
	go writeRecords(productRecords)
	go scrapeProducts("https://www.scrapingcourse.com/ecommerce/page/7/", productRecords)
	wg.Wait()
	fmt.Printf("scraping took: %v\n", time.Since(start))
}

func writeRecords(productRecords chan Product) {
	defer wg.Done()
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

	for p := range productRecords {
		record := []string{p.name, p.price, p.url, p.img}
		writer.Write(record)
	}

	defer writer.Flush()
}

func scrapeProducts(page string, productRecords chan Product) {
	defer wg.Done()
	c := colly.NewCollector()

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
		productRecords <- product
		fmt.Println(product.name, product.price)
	})
	c.OnScraped(func(r *colly.Response) {
		fmt.Println(page, "Is scraped.")
		close(productRecords)
	})
	c.Visit(page)
}
