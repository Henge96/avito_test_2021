package app

import "time"

type Wallet struct {
	Id         int     `json:"Id"`         // айди кошелька пользователя (примару кей)
	UserId     int     `json:"UserId"`     // айди юзера (поступает извне)
	Balance    float64 `json:"Balance"`    // баланс кошелька
	Currency   string  `json:"Currency"`   // валюта операции
	ReceiverId int     `json:"ReceiverId"` // айди получателя
	Transaction
}

type Transaction struct {
	TransactionId    int       `json:"TransactionId"`    // айди транзакции (будет примару кей)
	WalletId         int       `json:"wallet_id"`        // откуда был перевод
	ReceiverWalletId int       `json:"ReceiverWalletId"` // куда перевод
	Money            float64   `json:"Money"`            // сумма перевода
	Date             time.Time `json:"Date"`             // дата операции
}
