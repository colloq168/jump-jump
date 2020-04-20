package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jwma/jump-jump/internal/app/config"
	"github.com/jwma/jump-jump/internal/app/db"
	"github.com/jwma/jump-jump/internal/app/models"
	"github.com/jwma/jump-jump/internal/app/repository"
	"log"
	"net/http"
)

func Redirect(c *gin.Context) {
	slRepo := repository.GetShortLinkRepo(db.GetRedisClient())
	s, err := slRepo.Get(c.Param("id"))

	if err != nil {
		log.Printf("查找短链接失败，error: %v\n", err)
		cc := config.GetConfig().GetStringStringMapValue("shortLinkNotFoundConfig",
			config.GetDefaultShortLinkNotFoundConfig())

		switch cc["mode"] {
		case config.ShortLinkNotFoundContentMode:
			c.String(http.StatusOK, cc["value"])
			break
		case config.ShortLinkNotFoundRedirectMode:
			c.Redirect(http.StatusTemporaryRedirect, cc["value"])
			break
		default:
			c.String(http.StatusOK, "你访问的页面不存在哦")
		}

		return
	}

	if !s.IsEnable {
		c.String(http.StatusOK, "你访问的页面不存在哦")
		return
	}

	// 保存短链接请求记录（IP、User-Agent）
	rhRepo := repository.GetRequestHistoryRepo(db.GetRedisClient())
	go rhRepo.Save(models.NewRequestHistory(s, c.ClientIP(), c.Request.UserAgent()))

	c.Redirect(http.StatusTemporaryRedirect, s.Url)
}
