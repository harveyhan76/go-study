package main

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Account struct {
	ID      int     `gorm:"primaryKey"`
	Balance float64 `gorm:"not null;defalt:0"`
}

type Transaction struct {
	ID            int       `gorm:"primaryKey"`
	FromAccountID uint      `gorm:"not null"`
	ToAccountID   uint      `gorm:"not null"`
	Amount        float64   `gorm:"not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}

// TransferMoney 转账函数
func TransferMoney(db *gorm.DB, fromAccountID, toAccountID uint, amount float64) error {
	// 检查金额是否有效
	if amount <= 0 {
		return errors.New("转账金额必须大于0")
	}

	// 检查是否同一账户
	if fromAccountID == toAccountID {
		return errors.New("不能向同一账户转账")
	}

	// 开始事务
	return db.Transaction(func(tx *gorm.DB) error {
		// 1. 查询转出账户
		var fromAccount Account
		if err := tx.First(&fromAccount, fromAccountID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.New("转出账户不存在")
			}
			return err
		}

		// 2. 检查余额是否足够
		if fromAccount.Balance < amount {
			return errors.New("账户余额不足")
		}

		// 3. 查询转入账户
		var toAccount Account
		if err := tx.First(&toAccount, toAccountID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.New("转入账户不存在")
			}
			return err
		}

		// 4. 执行转账：扣除转出账户余额
		if err := tx.Model(&fromAccount).Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
			return err
		}

		// 5. 执行转账：增加转入账户余额
		if err := tx.Model(&toAccount).Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
			return err
		}

		// 6. 记录转账交易
		transaction := Transaction{
			FromAccountID: fromAccountID,
			ToAccountID:   toAccountID,
			Amount:        amount,
			CreatedAt:     time.Now(),
		}

		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		return nil
	})
}

func main() {
	fmt.Println("Hello, World!")
	db, err := gorm.Open(mysql.Open("root:demo12!@@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic(err)
	}
	//建表
	db.AutoMigrate(&Account{}, &Transaction{})

	// accounts := []Account{
	// 	{Balance: 1000},
	// 	{Balance: 2000},
	// }
	// result := db.Create(&accounts)
	// if result.Error != nil {
	// 	panic(result.Error)
	// }
	var account1 Account
	// db.Debug().First(&account, "ID = ?", 1)
	db.First(&account1, "ID = ?", 1)
	fmt.Println(account1)
	var account2 Account
	db.First(&account2, "ID = ?", 2)
	fmt.Println(account2)

	// db.Model(&Account{}).Where("ID = ?", 1).Update("Balance", 100)
	// db.Model(&Account{}).Where("ID = ?", 2).Update("Balance", 1000)

	err = TransferMoney(db, 1, 2, 100.0)
	if err != nil {
		fmt.Printf("转账失败: %v\n", err)
	} else {
		fmt.Println("转账成功")
	}

}
