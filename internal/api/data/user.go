/*
 * Created by lintao on 2023/7/27 下午2:36
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package data

import (
	"github.com/NSObjects/go-template/internal/api/data/db"
	"gorm.io/gorm"

	"github.com/NSObjects/go-template/internal/api/data/model"
	"github.com/NSObjects/go-template/internal/api/service/param"
	"github.com/pkg/errors"
)

type UserRepository interface {
	GetUserByID(id int64) (user model.User, err error)
	FindUser(param model.User, query param.APIQuery) (user []model.User, total int64, err error)
	DeleteUserByID(id int64) (err error)
	UpdateUser(param model.User, id int64) (err error)
	CreateUser(param model.User) (id int64, err error)
}

type UserDataSource struct {
	dataSource *db.DataSource
}

func NewUserDataSource(dataSource *db.DataSource) UserRepository {
	return &UserDataSource{dataSource: dataSource}
}

func (u *UserDataSource) CreateUser(param model.User) (id int64, err error) {
	err = u.dataSource.Mysql.Create(&param).Error
	if err != nil {
		return 0, errors.Wrap(err, "创建用户失败")
	}
	return
}

func (u *UserDataSource) GetUserByID(id int64) (user model.User, err error) {

	if err = u.dataSource.Mysql.First(&user, id).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return user, errors.Wrap(err, "查询用户失败")
		}
		return model.User{}, nil
	}
	return
}

func (u *UserDataSource) FindUser(param model.User, query param.APIQuery) (user []model.User, total int64, err error) {
	if err = u.dataSource.Mysql.Offset(query.Offset()).Limit(query.Limit()).Where(&param).Find(&user).Error; err != nil {
		return nil, 0, errors.Wrap(err, "")
	}
	if err = u.dataSource.Mysql.Model(&model.User{}).Where(&param).Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(err, "")
	}
	return
}

func (u *UserDataSource) DeleteUserByID(id int64) (err error) {
	if err = u.dataSource.Mysql.Delete(&model.User{}, id).Error; err != nil {
		return errors.Wrap(err, "删除用户失败")
	}
	return
}

func (u *UserDataSource) UpdateUser(param model.User, id int64) (err error) {
	param.ID = uint(id)
	if err = u.dataSource.Mysql.Select("name", "phone", "status", "password").Updates(&param).Error; err != nil {
		return errors.Wrap(err, "更新用户失败")
	}

	return
}
