package models

import "time"

type Coin struct {
	ID             uint64    `json:"id" sql:",pk"`
	Volume         string    `json:"volume"          sql:"type:numeric(70)"`
	ReserveBalance string    `json:"reserve_balance" db:"reserve_balance" sql:"type:numeric(70)"`
	Price          string    `json:"price"           sql:"type:numeric(100)"`
	Capitalization string    `json:"capitalization"  sql:"type:numeric(100)"`
	Symbol         string    `json:"symbol"          sql:"type:varchar(20)"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}
