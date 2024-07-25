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
	pagesToScrape := make(map[string]string)
	getPagesToScrape(pagesToScrape)
	pagesToScrape["https://www.scrapingcourse.com/ecommerce/page/7/"] = "https://www.scrapingcourse.com/ecommerce/page/7/"
	fmt.Println(len(pagesToScrape))
	wg.Add(len(pagesToScrape))
	go writeRecords(productRecords)
	i := 0
	for p := range pagesToScrape {
		if len(p) > 0 {
			i++
			go scrapeProducts(p, productRecords, i == len(pagesToScrape)-1)

		}
	}
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

func scrapeProducts(page string, productRecords chan Product, closeChan bool) {
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
		// fmt.Println(product.name, product.price)
	})
	c.OnScraped(func(r *colly.Response) {
		fmt.Println(page, "Is scraped.")
		if closeChan == true {
			close(productRecords)
		}
	})
	c.Visit(page)
}

func getPagesToScrape(pagesToScrape map[string]string) {
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
		pagesToScrape[page] = page
	})

	c.Visit("https://www.scrapingcourse.com/ecommerce/page/7/")
}
