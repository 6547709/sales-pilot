package business

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/business"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
)

type APIKeyApi struct{}

func (b *APIKeyApi) Create(c *gin.Context) {
	var req struct {
		Name      string `json:"name" binding:"required"`
		ExpiresAt int64  `json:"expires_at"`
	}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userID := utils.GetUserID(c)
	var expiresAt *time.Time
	if req.ExpiresAt > 0 {
		t := time.Unix(req.ExpiresAt, 0)
		expiresAt = &t
	}

	apiKey, rawKey, err := APIKeyServiceApp.GenerateAPIKey(req.Name, userID, expiresAt)
	if err != nil {
		response.FailWithMessage("创建失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(gin.H{
		"api_key": apiKey,
		"raw_key": rawKey,
	}, "创建成功，请妥善保存 raw_key", c)
}

func (b *APIKeyApi) GetAPIKeys(c *gin.Context) {
	apiKeys, err := APIKeyServiceApp.GetAPIKeys()
	if err != nil {
		response.FailWithMessage("获取失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(apiKeys, "获取成功", c)
}

func (b *APIKeyApi) Revoke(c *gin.Context) {
	var reqId request.GetById
	err := c.ShouldBindJSON(&reqId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = APIKeyServiceApp.RevokeAPIKey(reqId.Uint())
	if err != nil {
		response.FailWithMessage("撤销失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("撤销成功", c)
}

func (b *APIKeyApi) Delete(c *gin.Context) {
	var reqId request.GetById
	err := c.ShouldBindJSON(&reqId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = APIKeyServiceApp.DeleteAPIKey(reqId.Uint())
	if err != nil {
		response.FailWithMessage("删除失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

func (b *APIKeyApi) Validate(c *gin.Context) {
	key := c.GetHeader("X-API-Key")
	if key == "" {
		response.FailWithMessage("缺少 API Key", c)
		return
	}

	apiKey, err := APIKeyServiceApp.ValidateAPIKey(key)
	if err != nil || apiKey == nil {
		response.FailWithMessage("无效的 API Key", c)
		return
	}

	APIKeyServiceApp.UpdateLastUsed(apiKey.ID)

	response.OkWithDetailed(apiKey, "验证成功", c)
}
