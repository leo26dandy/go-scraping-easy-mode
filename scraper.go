package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gocolly/colly"
)

// defining a data structure to store the scraped data
type ProductList struct {
	category, productName, description, price string
}

// it verifies if a string is present in a slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func main() {
	// initializing the slice of structs that will contain the scraped data
	var productItems []ProductList

	// initializing the list of pages to scrape with an empty slice
	var pagesToScrape []string

	// the first pagination URL to scrape
	pageToScrape := "https://sushi-diest.be/winkel/"

	// initializing the list of pages discovered with a pageToScrape
	pagesDiscovered := []string{pageToScrape}

	// current iteration
	i := 1
	// max pages to scrape
	limit := 5

	// initializing a Colly instance
	c := colly.NewCollector()
	// setting a valid User-Agent header
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"

	// iterating over the list of pagination links to implement the crawling logic
	c.OnHTML("a.page-numbers", func(e *colly.HTMLElement) {
		// discovering a new page
		newPaginationLink := e.Attr("href")

		// if the page discovered is new
		if !contains(pagesToScrape, newPaginationLink) {
			// if the page discovered should be scraped
			if !contains(pagesDiscovered, newPaginationLink) {
				pagesToScrape = append(pagesToScrape, newPaginationLink)
			}
			pagesDiscovered = append(pagesDiscovered, newPaginationLink)
		}
	})

	// scraping the product data
	c.OnHTML("li.cat-item", func(h *colly.HTMLElement) {
		productItem := ProductList{
			category: h.ChildText("h2"),
		}

		productItems = append(productItems, productItem)
	})

	c.OnHTML("li.prod-type-simple", func(e *colly.HTMLElement) {
		productItem := ProductList{
			productName: e.ChildText("h3"),
			description: e.ChildText("p"),
			price:       e.ChildText(".price"),
		}

		productItems = append(productItems, productItem)
	})

	c.OnScraped(func(response *colly.Response) {
		// until there is still a page to scrape
		if len(pagesToScrape) != 0 && i < limit {
			// getting the current page to scrape and removing it from the list
			pageToScrape = pagesToScrape[0]
			pagesToScrape = pagesToScrape[1:]

			// incrementing the iteration counter
			i++

			// visiting a new page
			c.Visit(pageToScrape)
		}
	})

	// visiting the first page
	c.Visit(pageToScrape)

	// opening the CSV file
	rand.Seed(time.Now().UnixNano())
	randomInt := rand.Intn(1000)
	fileName := fmt.Sprintf("products-test-%d.csv", randomInt)
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalln("Failed to create output CSV file", err)
	}
	defer file.Close()

	// initializing a file writer
	writer := csv.NewWriter(file)

	// defining the CSV headers
	headers := []string{
		"Category",
		"Product Name",
		"Description",
		"Price",
	}
	// writing the column headers
	writer.Write(headers)

	// adding each  product to the CSV output file
	for _, productItem := range productItems {
		// converting a Product to an array of strings
		record := []string{
			productItem.category,
			productItem.productName,
			productItem.description,
			productItem.price,
		}

		// writing a new CSV record
		writer.Write(record)
	}
	defer writer.Flush()
}
