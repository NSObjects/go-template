/*
 * Created by lintao on 2023/7/18 下午3:59
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package data

import (
	"github.com/NSObjects/go-template/internal/api/data/model"
	"github.com/NSObjects/go-template/tools/db"

	"github.com/pkg/errors"
)

type UserRepository interface {
	GetUserById(id int64) (user model.User, err error)
	FindUser(param model.UserParam) (user []model.User, total int64, err error)
	DeleteUserById(id int64) (err error)
	UpdateUser(param model.UserParam, id int64) (err error)
	CreateUser(param model.UserParam) (id int64, err error)
}

type UserDataSource struct {
	dataSource *db.DataSource
}

func NewUserDataSource(dataSource *db.DataSource) *UserDataSource {
	return &UserDataSource{dataSource: dataSource}
}

func (this *UserDataSource) CreateUser(param model.UserParam) (id int64, err error) {
	var user model.User
	user.Account = param.Account
	user.Password = param.Password
	user.Phone = param.Phone
	user.Status = param.Status
	user.Name = param.Name

	_, err = this.dataSource.Engine.Insert(&user)
	if err != nil {
		return 0, errors.Wrap(err, "创建用户失败")
	}

	return
}

func (this *UserDataSource) GetUserById(id int64) (user model.User, err error) {

	if _, err = this.dataSource.Engine.ID(id).Get(&user); err != nil {
		return user, errors.Wrap(err, "查询用户失败")
	}

	return

}

func (this *UserDataSource) FindUser(param model.UserParam) (user []model.User, total int64, err error) {

	session := this.dataSource.Engine.Table(new(model.User))

	if param.Phone != "" {
		session = session.Where("phone = ?", param.Phone)
	}

	if param.Name != "" {
		session = session.Where("name = ?", param.Phone)
	}

	if param.Status != 0 {
		session = session.Where("status = ?", param.Phone)
	}

	if param.Account != "" {
		session = session.Where("account = ?", param.Phone)
	}

	if total, err = session.Count(); err != nil {
		return nil, 0, errors.Wrap(err, "查询用户数量出错")
	}

	if err = session.Limit(param.Limit()).Find(&user); err != nil {
		return nil, 0, errors.Wrap(err, "查询用户出错")
	}

	return
}

func (this *UserDataSource) DeleteUserById(id int64) (err error) {

	if _, err = this.dataSource.Engine.ID(id).Delete(new(model.User)); err != nil {
		return errors.Wrap(err, "删除用户失败")
	}

	return
}

func (u *UserDataSource) UpdateUser(param model.UserParam, id int64) (err error) {

	if _, err = u.dataSource.Engine.ID(id).Update(&param.User); err != nil {
		return errors.Wrap(err, "更新用户失败")
	}
	//error = 更新用户失败: ExecQuery
	// 'UPDATE `user` SET `name` = ? WHERE `id`=?',
	// 'UPDATE `user` SET `name` = ? WHERE `id`=?', wantErr false
	return
}
