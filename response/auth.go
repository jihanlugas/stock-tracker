package response

import "stock-tracker/model"

type Init struct {
	User model.UserView `json:"user,omitempty"`
}
