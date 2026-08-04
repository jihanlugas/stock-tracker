package response

import (
	"math"
	"stock-tracker/request"
)

type Pagination struct {
	Page        int         `json:"page"`
	DataPerPage int         `json:"dataPerPage"`
	TotalData   int64       `json:"totalData"`
	TotalPage   int         `json:"totalPage"`
	List        interface{} `json:"list" swaggertype:"array,object"`
}

func PayloadPagination(req request.IPaging, list interface{}, totalData int64) *Pagination {
	dataPerPage := int(totalData)
	totalPage := 1
	if req.GetLimit() > 0 {
		dataPerPage = req.GetLimit()
		totalPage = int(math.Ceil(float64(totalData) / float64(req.GetLimit())))
	}
	pgn := Pagination{
		Page:        req.GetPage(),
		DataPerPage: dataPerPage,
		TotalData:   totalData,
		TotalPage:   totalPage,
		List:        list,
	}

	req.SetPage(0)

	return &pgn
}
