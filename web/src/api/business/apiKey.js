import service from '@/utils/request'

// @Tags APIKey
// @Summary 创建API密钥
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body business.APIKey true "创建API密钥"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /api-key [post]
export const createAPIKey = (data) => {
  return service({
    url: '/api-key',
    method: 'post',
    data
  })
}

// @Tags APIKey
// @Summary 获取API密钥列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /api-key/list [get]
export const getAPIKeys = () => {
  return service({
    url: '/api-key/list',
    method: 'get'
  })
}

// @Tags APIKey
// @Summary 吊销API密钥
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body business.APIKey true "吊销API密钥"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /api-key/revoke [post]
export const revokeAPIKey = (data) => {
  return service({
    url: '/api-key/revoke',
    method: 'post',
    data
  })
}

// @Tags APIKey
// @Summary 删除API密钥
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.APIKey true "删除API密钥"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /api-key [delete]
export const deleteAPIKey = (data) => {
  return service({
    url: '/api-key',
    method: 'delete',
    data
  })
}

// @Tags APIKey
// @Summary 验证API密钥
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body object{key string} true "验证API密钥"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /api-key/validate [post]
export const validateAPIKey = (data) => {
  return service({
    url: '/api-key/validate',
    method: 'post',
    data
  })
}
