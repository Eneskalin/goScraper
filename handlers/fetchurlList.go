package handlers


import(
	"github.com/PuerkitoBio/goquery"
	
	"log"
	"strings"
)

func FetchUrlList(body string)(list[] string){
	var urls[] string

	doc,err:=goquery.NewDocumentFromReader(strings.NewReader(body))
	if err !=nil {
	log.Fatal("Failed to parse the HTML document", err)
	}

	doc.Find("a,link,area,base").Each(func (i int,s *goquery.Selection)  {
			if href, exists := s.Attr("href"); exists {
				if href != "" && !strings.HasPrefix(href, "javascript") {
				urls = append(urls, href)
				
			}
			}	
	})
	doc.Find("img, script, iframe, audio, video, source, embed").Each(func (i int ,s*goquery.Selection)  {
		if src, exists := s.Attr("src"); exists {
			if src != "" {
				urls = append(urls, src)
			}
		}
	})
	doc.Find("form").Each(func(i int, s *goquery.Selection) {
		if action, exists := s.Attr("action"); exists {
			if action != "" {
				urls = append(urls, action)
			}
		}
	})
	return urls

}