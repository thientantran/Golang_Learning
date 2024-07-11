package middleware

import (
	"Food-delivery/common"
	"Food-delivery/component/appctx"
	"errors"
	"github.com/gin-gonic/gin"
)

func RoleRequired(appCtx appctx.AppContext, allowRoles ...string) func(c *gin.Context) {
	return func(c *gin.Context) {
		u := c.MustGet(common.CurrentUser).(common.Requester)

		hasFound := false
		for _, item := range allowRoles {
			if item == u.GetRole() {
				hasFound = true
				break
			}
		}
		if !hasFound {
			panic(common.ErrNoPermission(errors.New("you don't have permission")))
		}
		c.Next()
	}
}
