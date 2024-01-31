package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gocolly/colly"
)

type Product struct {
	Image, ProductName, Price string
}

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

func main() {
	defer timer("main")()
	c := colly.NewCollector(colly.Async(true))

	products := []Product{}

	c.OnHTML("li.next a", func(h *colly.HTMLElement) {
		c.Visit(h.Request.AbsoluteURL(h.Attr("href")))
	})

	c.OnHTML("article.product_pod", func(h *colly.HTMLElement) {
		iterations := Product{
			Image:       "https://books.toscrape.com" + h.ChildAttr("img", "src"),
			ProductName: h.ChildAttr("img", "alt"),
			Price:       h.ChildText("p.price_color"),
		}

		products = append(products, iterations)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting", r.URL)
	})

	c.Visit("https://books.toscrape.com/catalogue/page-1.html")
	c.Wait()

	fmt.Println(products)

	// opening the CSV file
	fileName := fmt.Sprintf("products-test.csv")
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalln("Failed to create output CSV file", err)
	}
	defer file.Close()

	// initializing a file writer
	writer := csv.NewWriter(file)

	// defining the CSV headers
	headers := []string{
		"Image",
		"Product Name",
		"Price",
	}
	// writing the column headers
	writer.Write(headers)

	// adding each  product to the CSV output file
	for _, productItem := range products {
		// converting a Product to an array of strings
		record := []string{
			productItem.Image,
			productItem.ProductName,
			productItem.Price,
		}

		// writing a new CSV record
		writer.Write(record)
	}
	defer writer.Flush()
}
