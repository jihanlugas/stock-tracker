package response

import "stock-tracker/model"

type Dashboard struct {
	Items []model.ItemView `json:"items"`
}
