/*
 * Copyright 2023-present by Damon All Rights Reserved
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package db

import (
	"errors"
	"gorm.io/gorm"
)

/*
* @author Damon
* @date   2023/5/8 16:46
 */

// 开始一个新的事务
func Begin() *gorm.DB {
	return sqlDB.Begin()
}

func Transaction_Handler(model interface{})  {
	sqlDB.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		if err := tx.Create(model).Error; err != nil {
			// 返回任何错误都会回滚事务
			return err
		}

		if err := tx.Create(nil).Error; err != nil {
			return err
		}

		// 返回 nil 提交事务
		return nil
	})
}


func More_Transaction_Handler(model1 interface{}, model2 interface{}, model3 interface{})  {
	sqlDB.Transaction(func(tx *gorm.DB) error {
		tx.Create(&model1)

		tx.Transaction(func(tx2 *gorm.DB) error {
			tx2.Create(&model2)
			return errors.New("rollback model2") // Rollback model2
		})

		tx.Transaction(func(tx2 *gorm.DB) error {
			tx2.Create(&model3)
			return nil
		})

		return nil
	})

	// Commit model1, model3
}

//TODO 特殊case
/*func CreateAnimals(db *gorm.DB) error {
	// 再唠叨一下，事务一旦开始，你就应该使用 tx 处理数据
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(&Animal{Name: "Giraffe"}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Create(&Animal{Name: "Lion"}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}*/


// 创建记录，如果出错则回滚事务，手动事务
func Create(model interface{}) error {
	tx := sqlDB.Begin()
	//此处一定要用tx来执行，不能用DB
	if err := tx.Create(model).Error; err != nil {
		// 返回任何错误都会回滚事务
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// 查询单条记录
func FindOne(model interface{}, query interface{}, args ...interface{}) error {
	if err := sqlDB.Where(query, args...).First(model).Error; err != nil {
		return err
	}
	return nil
}

// 查询多条记录
func FindAll(models interface{}, query interface{}, args ...interface{}) error {
	if err := sqlDB.Where(query, args...).Find(models).Error; err != nil {
		return err
	}
	return nil
}

// 更新记录，如果出错则回滚事务
func Update(model interface{}) error {
	tx := sqlDB.Begin()
	if err := tx.Save(model).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// 删除记录，如果出错则回滚事务
func Delete(tx *gorm.DB, model interface{}) error {
	if err := tx.Delete(model).Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// 提交事务
func Commit(tx *gorm.DB) error {
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
