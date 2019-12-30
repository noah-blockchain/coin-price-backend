package models

import "time"

type Address struct {
	ID        uint64    `json:"id" sql:",pk"`
	Address   string    `json:"address"          sql:"varchar(40)"`
	Symbol    string    `json:"symbol"           sql:"type:varchar(20)"`
	Amount    string    `json:"amount"           sql:"type:numeric(70)"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
