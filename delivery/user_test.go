/*
 *
 * user.go
 * apis
 *
 * Created by lintao on 2019-01-29 16:17
 * Copyright © 2017-2019 PYL. All rights reserved.
 *
 */

package delivery

import (
	"go-template/delivery/server"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestUserController_getUser(t *testing.T) {
	c, _ := server.request(echo.GET, "/api/users", server.e)
	assert.Equal(t, http.StatusOK, c)
}

func TestUserController_createUser(t *testing.T) {
	c, _ := server.request(echo.POST, "/api/users", server.e)
	assert.Equal(t, http.StatusOK, c)
}

func TestUserController_updateUser(t *testing.T) {
	c, _ := server.request(echo.PUT, "/api/users/:id", server.e)
	assert.Equal(t, http.StatusOK, c)
}

func TestUserController_deleteUser(t *testing.T) {
	c, _ := server.request(echo.DELETE, "/api/users/:id", server.e)
	assert.Equal(t, http.StatusOK, c)
}

func TestUserController_getUserDetail(t *testing.T) {
	c, _ := server.request(echo.GET, "/api/users/:id", server.e)
	assert.Equal(t, http.StatusOK, c)
}
