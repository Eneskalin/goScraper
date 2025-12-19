package helpers

import (
	"fmt"
	"net/http"
)

func CheckHttpStatus(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close() 

	if resp.StatusCode != http.StatusOK {
		
		return fmt.Errorf("siteye ulaşılamadı, durum kodu: %d", resp.StatusCode)
	}

	return nil
}
