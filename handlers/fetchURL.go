package handlers

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
		"webScraper/helpers"
	"github.com/chromedp/chromedp"
)

type FetchSuccessMsg string      
type ScreenshotSuccessMsg string 
type FetchUrListSuccessMsg string
type SaveSuccessMsg string       
type ErrorMsg error   

var htmlContent string
var list[] string

type FetchResultMsg struct {
	Err error
}

func CheckStatusCmd(url string) tea.Cmd {
    return func() tea.Msg {
        err := helpers.CheckHttpStatus(url)
        if err != nil {
            return err 
        }
        return StatusOkMsg{} 
    }
}

type StatusOkMsg struct{}

func SavingHandlers(url string)tea.Cmd{
	return func() tea.Msg {

	if err := SaveHtml(htmlContent,list, url); err != nil {
		log.Error("HTML verileri kaydedilemedi", err)
		return fmt.Errorf("failed to save HTML: %w", err)
	}
	log.Info("Veriler başarıyla kaydedildi")
	return SaveSuccessMsg("Veriler kaydedildi")

}



}

func ListURLsCmd() tea.Cmd {
	return func() tea.Msg {
		list=FetchUrlList(htmlContent)

		log.Info("Sayfadaki URLler listelendi")
		return FetchUrListSuccessMsg("Sayfadaki URLler listelendi")

		
	}
}

func FetchHTMLCmd(url string) tea.Cmd {
	
    return func() tea.Msg {
		

	ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithLogf(func(string, ...interface{}) {}))
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Evaluate(`navigator.userAgent`, nil),
		chromedp.Navigate(url),
		
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.OuterHTML("html", &htmlContent),
		
	)
	if err != nil {
		log.Error("HTML verileri alınamadı", err)
		return fmt.Errorf("failed to fetch page: %w", err)
	}
        log.Info("HTML içeriği başarıyla çekildi!")
        return FetchSuccessMsg("HTML içeriği başarıyla çekildi!")
    }
}



func TakeScreenshotCmd( url string) tea.Cmd {
    return func() tea.Msg {
        ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithLogf(func(string, ...interface{}) {}))
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Evaluate(`navigator.userAgent`, nil),
		chromedp.Navigate(url),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.OuterHTML("html", &htmlContent),
	)
	if err != nil {
		log.Error("HTML verileri alınamadı", err)
		return fmt.Errorf("Hata: %w", err)
	}
	if err := captureFullPageScreenshot(ctx, url); err != nil {
		log.Error("Ekran görüntüsü alma sırarsında bir hata oldu:", err)
		return fmt.Errorf("Hata: %w",err)
	}
		log.Info("Ekran görüntüsü kaydedildi!")
        return ScreenshotSuccessMsg("Ekran görüntüsü kaydedildi!")
    }
}
