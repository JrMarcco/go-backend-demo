package models

type Account struct {
	Model
	AccountOwner string `json:"accountOwner"`
	Balance      int64  `json:"balance"`
	Currency     string `json:"currency"`
}

func (a Account) TableName() string {
	return "account"
}

func GetAccount(id uint64) (*Account, error) {
	var account Account

	if err := db.Where("id = ?", id).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

type AddAccountParam struct {
	AccountOwner string `json:"accountOwner"`
	Balance      int64  `json:"balance"`
	Currency     string `json:"currency"`
}

func AddAccount(param AddAccountParam) (*Account, error) {
	account := Account{
		AccountOwner: param.AccountOwner,
		Balance:      param.Balance,
		Currency:     param.Currency,
	}

	if err := db.Create(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

func DelAccount(id uint64) error {
	if err := db.Where("id = ?", id).Delete(Account{}).Error; err != nil {
		return err
	}
	return nil
}
