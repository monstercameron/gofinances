package assets

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/monstercameron/gofinances/database"
	"github.com/monstercameron/gofinances/helpers"
	"log"
	"net/http"
	"strconv"
)

type Asset struct {
	ID              int     `json:"id"`
	AssetName       string  `json:"assetName"`
	AssetOwner      string  `json:"assetOwner"`
	InsertDate      string  `json:"insertDate"`
	AssetValue      float64 `json:"assetValue"`
	AssetGrowthRate float64 `json:"assetGrowthRate"`
	Notes           string  `json:"notes"`
}

func (a *Asset) Save() error {
	tx, err := database.DB.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	if a.ID == 0 {
		query := `INSERT INTO assets (asset_name, asset_owner, insert_date, asset_value, asset_growth_rate, notes) 
				  VALUES (?, ?, ?, ?, ?, ?)`
		result, err := tx.Exec(query, a.AssetName, a.AssetOwner, a.InsertDate, a.AssetValue, a.AssetGrowthRate, a.Notes)
		if err != nil {
			return fmt.Errorf("error inserting asset: %v", err)
		}
		lastInsertID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("error getting last insert ID: %v", err)
		}
		a.ID = int(lastInsertID)
	} else {
		query := `UPDATE assets SET asset_name=?, asset_owner=?, insert_date=?, asset_value=?, asset_growth_rate=?, notes=? 
				  WHERE id=?`
		_, err := tx.Exec(query, a.AssetName, a.AssetOwner, a.InsertDate, a.AssetValue, a.AssetGrowthRate, a.Notes, a.ID)
		if err != nil {
			return fmt.Errorf("error updating asset: %v", err)
		}
	}

	return tx.Commit()
}

func (a *Asset) Delete() error {
	_, err := database.DB.Exec("DELETE FROM assets WHERE id=?", a.ID)
	return err
}

func GetAsset(id int) (Asset, error) {
	var a Asset
	query := `SELECT id, asset_name, asset_owner, insert_date, asset_value, asset_growth_rate, notes 
			  FROM assets WHERE id=?`
	err := database.DB.QueryRow(query, id).Scan(
		&a.ID, &a.AssetName, &a.AssetOwner, &a.InsertDate, &a.AssetValue, &a.AssetGrowthRate, &a.Notes)
	if err != nil {
		if err == sql.ErrNoRows {
			return Asset{}, fmt.Errorf("asset with id %d not found", id)
		}
		return Asset{}, fmt.Errorf("error querying asset: %v", err)
	}
	return a, nil
}

func GetAllAssetsWithoutError() []Asset {
	var assets []Asset
	query := `SELECT id, asset_name, asset_owner, insert_date, asset_value, asset_growth_rate, notes FROM assets`
	rows, err := database.DB.Query(query)
	if err != nil {
		log.Printf("Error querying assets: %v", err)
		return []Asset{}
	}
	defer rows.Close()

	for rows.Next() {
		var a Asset
		err := rows.Scan(&a.ID, &a.AssetName, &a.AssetOwner, &a.InsertDate, &a.AssetValue, &a.AssetGrowthRate, &a.Notes)
		if err != nil {
			log.Printf("Error scanning asset: %v", err)
			continue
		}
		assets = append(assets, a)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error after scanning assets: %v", err)
	}

	return assets
}

// Handler functions

func GetAssetsIndex(w http.ResponseWriter, r *http.Request) {
	//assets := GetAllAssetsWithoutError()
	component := AssetsIndex()
	w.Header().Set("Content-Type", "text/html")
	component.Render(r.Context(), w)
}

func GetOneAsset(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ExtractSegmentFromPath(r.URL.Path, 2)
	if err != nil {
		http.Error(w, "Invalid asset ID", http.StatusBadRequest)
		return
	}

	assetID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid asset ID", http.StatusBadRequest)
		return
	}

	asset, err := GetAsset(assetID)
	if err != nil {
		http.Error(w, "Asset not found", http.StatusNotFound)
		return
	}

	component := AssetListItem(asset)
	w.Header().Set("Content-Type", "text/html")
	component.Render(r.Context(), w)
}

func AddAsset(w http.ResponseWriter, r *http.Request) {
	var asset Asset
	err := json.NewDecoder(r.Body).Decode(&asset)
	if err != nil {
		http.Error(w, "Invalid asset data", http.StatusBadRequest)
		return
	}

	err = asset.Save()
	if err != nil {
		http.Error(w, "Failed to save asset", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(asset)
}

func GetEditAssetForm(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ExtractSegmentFromPath(r.URL.Path, 3)
	if err != nil {
		http.Error(w, "Invalid asset ID", http.StatusBadRequest)
		return
	}

	assetID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid asset ID", http.StatusBadRequest)
		return
	}

	asset, err := GetAsset(assetID)
	if err != nil {
		http.Error(w, "Asset not found", http.StatusNotFound)
		return
	}

	component := EditAssetForm(asset)
	w.Header().Set("Content-Type", "text/html")
	component.Render(r.Context(), w)
}

func UpdateAsset(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ExtractSegmentFromPath(r.URL.Path, 2)
	if err != nil {
		http.Error(w, "Invalid asset ID", http.StatusBadRequest)
		return
	}

	assetID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid asset ID", http.StatusBadRequest)
		return
	}

	var asset Asset
	err = json.NewDecoder(r.Body).Decode(&asset)
	if err != nil {
		http.Error(w, "Invalid asset data", http.StatusBadRequest)
		return
	}

	asset.ID = assetID
	err = asset.Save()
	if err != nil {
		http.Error(w, "Failed to update asset", http.StatusInternalServerError)
		return
	}

	component := AssetListItem(asset)
	w.Header().Set("Content-Type", "text/html")
	component.Render(r.Context(), w)
}

func DeleteAsset(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ExtractSegmentFromPath(r.URL.Path, 2)
	if err != nil {
		http.Error(w, "Invalid asset ID", http.StatusBadRequest)
		return
	}

	assetID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid asset ID", http.StatusBadRequest)
		return
	}

	asset := Asset{ID: assetID}
	err = asset.Delete()
	if err != nil {
		http.Error(w, "Failed to delete asset", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetAssetsTotalValue(w http.ResponseWriter, r *http.Request) {
	assets := GetAllAssetsWithoutError()
	var total float64
	for _, asset := range assets {
		total += asset.AssetValue
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "%.2f", total)
}

func GetNewAssetForm(w http.ResponseWriter, r *http.Request) {
	newAsset := Asset{} // Empty asset for the new form
	component := EditAssetForm(newAsset)
	w.Header().Set("Content-Type", "text/html")
	component.Render(r.Context(), w)
}
