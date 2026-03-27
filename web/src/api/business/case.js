import service from '@/utils/request'

// @Tags Case
// @Summary 创建客户案例
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body business.Case true "创建客户案例"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /case [post]
export const createCase = (data) => {
  return service({
    url: '/case',
    method: 'post',
    data
  })
}

// @Tags Case
// @Summary 更新客户案例
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body business.Case true "更新客户案例"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /case [put]
export const updateCase = (data) => {
  return service({
    url: '/case',
    method: 'put',
    data
  })
}

// @Tags Case
// @Summary 删除客户案例
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.Case true "删除客户案例"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /case [delete]
export const deleteCase = (data) => {
  return service({
    url: '/case',
    method: 'delete',
    data
  })
}

// @Tags Case
// @Summary 获取单一客户案例
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.Case true "获取单一客户案例"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /case/{id} [get]
export const getCase = (params) => {
  return service({
    url: '/case/' + params.id,
    method: 'get',
    params
  })
}

// @Tags Case
// @Summary 获取产品关联客户案例列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.PageInfo true "获取产品关联客户案例列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /case/product [get]
export const getCasesByProductID = (params) => {
  return service({
    url: '/case/product',
    method: 'get',
    params
  })
}
