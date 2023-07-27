/*
 * Created by lintao on 2023/7/27 下午1:44
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package service

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestUserController_getUser(t *testing.T) {
	c, _ := Request(echo.GET, "/api/users")
	assert.Equal(t, http.StatusOK, c)
}

func TestUserController_createUser(t *testing.T) {
	c, _ := Request(echo.POST, "/api/users")
	assert.Equal(t, http.StatusOK, c)
}

func TestUserController_updateUser(t *testing.T) {
	c, _ := Request(echo.PUT, "/api/users/:id")
	assert.Equal(t, http.StatusOK, c)
}

func TestUserController_deleteUser(t *testing.T) {
	c, _ := Request(echo.DELETE, "/api/users/:id")
	assert.Equal(t, http.StatusOK, c)
}

func TestUserController_getUserDetail(t *testing.T) {
	c, _ := Request(echo.GET, "/api/users/:id")
	assert.Equal(t, http.StatusOK, c)
}
