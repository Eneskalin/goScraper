package handlers

import (
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"
	"webScraper/helpers"
	"github.com/charmbracelet/log"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func captureFullPageScreenshot(ctx context.Context, url string) error {

	folderName := helpers.Normalized(url)

	docPath := filepath.Join("docs", folderName)

	if err := os.MkdirAll(docPath, 0755); err != nil {
		return fmt.Errorf("Klasör oluşturulamadı: %w", err)
	}

	screenshotPath := filepath.Join(docPath, "screenshot.png")
	var buf []byte

	
	if err := chromedp.Run(ctx,
		chromedp.Sleep(5*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {

			_, _, contentSize, _, _, _, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			width := int64(math.Ceil(contentSize.Width))
			height := int64(math.Ceil(contentSize.Height))

			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}
		
			buf, err = page.CaptureScreenshot().
				WithQuality(90).
				WithClip(&page.Viewport{
					X:      0,
					Y:      0,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).
				Do(ctx)

			return err
		}),
	); err != nil {
		return fmt.Errorf("screenshot çekilemedi: %w", err)
	}
	if err := os.WriteFile(screenshotPath, buf, 0644); err != nil {
        return fmt.Errorf("dosya kaydedilemedi: %w", err)
    }

	log.Info("Screenshot kaydedildi: ", screenshotPath)
	return nil
}


