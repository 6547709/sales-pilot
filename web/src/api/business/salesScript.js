import service from '@/utils/request'

// @Tags SalesScript
// @Summary 创建销售话术
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body business.SalesScript true "创建销售话术"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /script [post]
export const createSalesScript = (data) => {
  return service({
    url: '/script',
    method: 'post',
    data
  })
}

// @Tags SalesScript
// @Summary 更新销售话术
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body business.SalesScript true "更新销售话术"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /script [put]
export const updateSalesScript = (data) => {
  return service({
    url: '/script',
    method: 'put',
    data
  })
}

// @Tags SalesScript
// @Summary 删除销售话术
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.SalesScript true "删除销售话术"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /script [delete]
export const deleteSalesScript = (data) => {
  return service({
    url: '/script',
    method: 'delete',
    data
  })
}

// @Tags SalesScript
// @Summary 获取单一销售话术
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.SalesScript true "获取单一销售话术"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /script/{id} [get]
export const getSalesScript = (params) => {
  return service({
    url: '/script/' + params.id,
    method: 'get',
    params
  })
}

// @Tags SalesScript
// @Summary 获取产品关联销售话术列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.PageInfo true "获取产品关联销售话术列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /script/product [get]
export const getSalesScriptsByProductID = (params) => {
  return service({
    url: '/script/product',
    method: 'get',
    params
  })
}
