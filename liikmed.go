package main

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type liige struct {
	nimi           string
	url            string
	enesetutvustus string
}

var (
	liikmed []liige
)

func main() {
	// Kogu liikmete lehtede lingid
	c1 := colly.NewCollector(
		// Külasta ainult domeene:
		colly.AllowedDomains("www.riigikogu.ee"),
	)

	// Lehe igal lingil:
	c1.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Filtreeri: link sisaldab '/saadik/...' ja
		// lingi tekst ei sisalda '<img'

		r1, _ := regexp.Compile("/saadik/.")

		// Link -> []byte
		linkbytes := []byte(link)

		if (r1.Find(linkbytes) != nil) &&
			!strings.Contains(e.Text, "<img") {

			// Eralda liikme nimi: link lõpust, enne viimast /
			r2 := regexp.MustCompile(`/`)
			t := r2.Split(link, -1)
			nimi := t[len(t)-1]

			liikmed = append(liikmed, liige{nimi: nimi, url: link})
		}
	})

	// Alusta korjet
	c1.Visit("https://www.riigikogu.ee/riigikogu/koosseis/riigikogu-liikmed/")

	// Prindi liikmete arv
	fmt.Println("Liikmeid: ", len(liikmed))

	// Koguja liikmete enesetutvustuste kogumiseks
	c2 := colly.NewCollector(
		colly.AllowedDomains("www.riigikogu.ee"),
	)

	// Eralda enesetutvustused
	// CSS klass profile-desc all 2. p-element
	c2.OnHTML(".profile-desc p:nth-child(3)", func(e *colly.HTMLElement) {
		enesetutvustus := e.Text
		r := e.Request
		ctx := r.Ctx
		liikme_nr_string := ctx.Get(`liikme nr`)
		liikme_nr, err := strconv.Atoi(liikme_nr_string)
		if err != nil {
			log.Fatal(err)
		}
		liikmed[liikme_nr].enesetutvustus = enesetutvustus
	})

	// Kogu enesetutvustused
	// for i, _ := range liikmed {
	for i := 10; i <= 100; i++ {
		/* if i == 102 {
			break

		} */

		// Kontekst enesetutvustuste kogumisele
		ctx2 := colly.NewContext()
		ctx2.Put(`liikme nr`, strconv.Itoa(i))

		// Kui ei soovi konteksti edastada
		// c2.Visit(liikmed[i].url)

		// Edenemise näitaja
		fmt.Print(".")

		// Konteksti edasiandmiseks
		c2.Request("GET", liikmed[i].url, nil, ctx2, nil)

		// Viivitus päringute vahel
		viivitus := 3 * time.Second
		juhuViivitus := time.Duration(rand.Int63n(4))
		time.Sleep(viivitus + juhuViivitus)
	}

	// Edenemise lõpp
	fmt.Println()

	// Prindi enesetutvustused
	for i, l := range liikmed {
		if l.enesetutvustus != "" {
			fmt.Printf("%v : %s\n%s\n\n",
				i+1, l.nimi,
				l.enesetutvustus)
		}
	}

}
