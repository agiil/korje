package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

type liige struct {
	nimi           string
	url            string
	enesetutvustus string
}

var (
	liikmeid int
	liikmed  []liige
)

func main() {
	// Esimene koguja : Kogub liikmete lehtede lingid
	c1 := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("www.riigikogu.ee"),
	)

	// On every a element which has href attribute call callback
	c1.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Filtreeri: link sisaldab '/saadik/...' ja
		// lingi tekst ei sisalda '<img'

		r1, _ := regexp.Compile("/saadik/.")

		// Link to []byte
		linkbytes := []byte(link)

		if (r1.Find(linkbytes) != nil) &&
			!strings.Contains(e.Text, "<img") {

			// Eralda liikme nimi: link lõpust, enne viimast /
			r2 := regexp.MustCompile(`/`)
			t := r2.Split(link, -1)
			nimi := t[len(t)-1]

			// Print link
			// fmt.Printf("%s %s\n", nimi, link)
			// fmt.Printf("Link found: %q -> %s\n", e.Text, link)
			// Visit link found on page
			// Only those links are visited which are in AllowedDomains
			// c.Visit(e.Request.AbsoluteURL(link))
			liikmeid++

			liikmed = append(liikmed, liige{nimi: nimi, url: link})
		}
	})

	// Before making a request print "Visiting ..."
	/* c1.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	}) */

	// Start scraping
	c1.Visit("https://www.riigikogu.ee/riigikogu/koosseis/riigikogu-liikmed/")

	// Prindi liikmete arv
	fmt.Println("Liikmeid: ", len(liikmed))

	// Prindi liikmed
	/* for _, l := range liikmed {
		fmt.Printf("%s %s\n", l.nimi, l.url)
	} */

	// Teine koguja : Kogub liikmete enesetutvustused
	c2 := colly.NewCollector(
		colly.AllowedDomains("www.riigikogu.ee"),
	)

	// Eralda enesetutvustused
	// CSS klass profile-desc all 2. p-element
	c2.OnHTML(".profile-desc p:nth-child(3)", func(e *colly.HTMLElement) {
		enesetutvustus := e.Text
		fmt.Println("*****: ", enesetutvustus)
	})

	// Käivita enesetutvustuste koguja
	// for i, _ := range liikmed {
	for i := 52; i <= 60; i++ {
		/* if i == 102 {
			break

		} */

		fmt.Println(i)

		// Kontekst enesetutvustuste kogumisele
		ctx2 := colly.NewContext()
		ctx2.Put(`liikme nr`, i)

		c2.Visit(liikmed[i].url)
		// Konteksti edasiandmiseks kasutan
		// c2.Request("GET", liikmed[i].url, nil, ctx2, nil)

		// time.Sleep(8 * time.Second)
	}

}
