package handlers

import (
	"os"
	"path/filepath"
	"webScraper/helpers"
	"encoding/json"
	"github.com/charmbracelet/log"
)

func SaveHtml(body string, list []string, url string) error {
	folderName := helpers.Normalized(url)

	docPath := filepath.Join("docs", folderName)

	err := os.MkdirAll(docPath, 0755)
	if err != nil {
		log.Error("Klasör oluşturulamadı: %v\n", err)
		return err
	}

	filePath := filepath.Join(docPath, "index.html")

	err = os.WriteFile(filePath, []byte(body), 0644)
	if err != nil {
		log.Error("Dosya yazılamadı: %v\n", err)
		return err
	}

	jsonData,err:=json.MarshalIndent(list,""," ")
	if err!=nil {
		log.Error("Url listesi json a cevrilmedi")
		return err
	}
	jsonPath := filepath.Join(docPath, "urls.json")
	err=os.WriteFile(jsonPath,jsonData,0644)
	if err != nil {
		log.Error("JSON dosyası yazılamadı", "error", err)
		return err
	}


	log.Info("Veriler başarıyla kaydedildi:", filePath)
	return nil
}
