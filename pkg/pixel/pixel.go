package pixel

import (
	"context"
)

type Pixel struct {
	CookieID           string `json:"cid" query:"cid"`
	Country            string `json:"c" query:"c"`
	Hotel              string `json:"h" query:"h"`
	ConfirmationNumber string `json:"cf" query:"cf"`
	ExtraField         string `json:"ex"`
}

type Repository interface {
	Store(context.Context, Pixel) error
}
