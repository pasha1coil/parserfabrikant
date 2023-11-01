package worker

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/geziyor/geziyor/export"
	"strconv"
	"strings"
	"testingParser/internal/models"
)

// функция для получения url со всех доступных страниц
func GetUrl(value string) {
	numPages := getNumPages(value)
	if numPages == 0 {
		numPages = 1
	}
	for i := 1; i <= numPages; i++ {
		geziyor.NewGeziyor(&geziyor.Options{
			StartURLs: []string{"https://www.fabrikant.ru/trades/procedure/search/?query=" + value + "&page=" + strconv.Itoa(i)},
			ParseFunc: parseData,
			Exporters: []export.Exporter{&export.JSON{}},
		}).Start()
	}
}

// функция парсинга данных с url
func parseData(g *geziyor.Geziyor, r *client.Response) {
	r.HTMLDoc.Find("div.marketplace-unit").Each(func(i int, s *goquery.Selection) {
		data := models.DataFrame{}

		data.Number = strings.TrimSpace(strings.ReplaceAll(s.Find("div.marketplace-unit__info__name span").Text(), "\n", " "))
		data.Name = strings.TrimSpace(s.Find("h4.marketplace-unit__title a").Text())
		data.URL, _ = s.Find("h4.marketplace-unit__title a").Attr("href")
		data.Price = strings.TrimSpace(strings.ReplaceAll(s.Find("div.marketplace-unit__price strong").Text(), "\n", " "))

		organizerText := s.Find("div.marketplace-unit__organizer span").Text()
		data.Organizer = strings.TrimSpace(strings.TrimPrefix(organizerText, "Организатор:"))

		stateText := s.Find("div.marketplace-unit__state").Text()
		stateText = strings.ReplaceAll(stateText, "Прием заявок", "")
		stateText = strings.ReplaceAll(stateText, "Дата и время начала приема заявок:", "")
		stateText = strings.ReplaceAll(stateText, "Дата и время окончания приема заявок:", "")
		stateText = strings.TrimSpace(stateText)
		data.Date = stateText

		g.Exports <- data
	})
}

// функция для определения количества страниц
func getNumPages(value string) int {
	lastPage := 0
	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{"https://www.fabrikant.ru/trades/procedure/search/?query=" + value},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			r.HTMLDoc.Find("div.col-xs-8").Each(func(i int, s *goquery.Selection) {
				Page := s.Find("li.pagination__lt__el span").Last()
				lastPage, _ = strconv.Atoi(Page.Text())
			})
		},
	}).Start()

	return lastPage
}
