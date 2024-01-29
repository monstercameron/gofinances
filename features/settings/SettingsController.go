package settings

import (
	"net/http"
	"fmt"
)

func init() {
	fmt.Println("SettingsController.init(): \t\tInitializing SettingsController...")
}

func GetSettingsPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SettingsController.GetSettingsPage(): \t\tGetting Settings Page...")
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("Settings Page from SettingsController.GetSettingsPage()"))
}