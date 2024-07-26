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

var productRecords = []Product{}

func main() {
	start := time.Now()
	pagesToScrape := make(map[string]string)
	pagesToScrape["https://www.scrapingcourse.com/ecommerce/page/7/"] = "https://www.scrapingcourse.com/ecommerce/page/7/"
	getPagesToScrape("https://www.scrapingcourse.com/ecommerce/page/7/", pagesToScrape)
	wg.Add(len(pagesToScrape))
	for _, p := range pagesToScrape {
		go scrapeProducts(p)
	}
	wg.Wait()
	fmt.Println(len(productRecords))
	writeRecords()
	fmt.Printf("scraping took: %v\n", time.Since(start))
}

func writeRecords() {
	// defer wg.Done()
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

	for _, p := range productRecords {
		record := []string{p.name, p.price, p.url, p.img}
		writer.Write(record)
	}
	defer writer.Flush()
}

func scrapeProducts(page string) {
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
		productRecords = append(productRecords, product)
	})
	c.OnScraped(func(r *colly.Response) {
		fmt.Println(page, "Is scraped.")
	})
	c.Visit(page)
}

func getPagesToScrape(page string, pagesToScrape map[string]string) {
	c := colly.NewCollector()
	c.OnError(func(r *colly.Response, err error) {
		log.Fatalln(err)
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("getting pages to scrape...")
	})
	c.OnScraped(func(r *colly.Response) {
		fmt.Println("All pages to scraped have been fetched.")
	})
	c.OnHTML("ul.page-numbers li", func(h *colly.HTMLElement) {
		page := h.ChildAttr("a", "href")
		if len(page) > 0 {

			pagesToScrape[page] = page
		}
	})

	c.Visit(page)
}
