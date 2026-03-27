import service from '@/utils/request'

// @Tags Product
// @Summary 创建产品
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body business.Product true "创建产品"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /product [post]
export const createProduct = (data) => {
  return service({
    url: '/product',
    method: 'post',
    data
  })
}

// @Tags Product
// @Summary 更新产品
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body business.Product true "更新产品"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /product [put]
export const updateProduct = (data) => {
  return service({
    url: '/product',
    method: 'put',
    data
  })
}

// @Tags Product
// @Summary 删除产品
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.Product true "删除产品"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /product [delete]
export const deleteProduct = (data) => {
  return service({
    url: '/product',
    method: 'delete',
    data
  })
}

// @Tags Product
// @Summary 获取单一产品
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.Product true "获取单一产品"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /product/{id} [get]
export const getProduct = (params) => {
  return service({
    url: '/product/' + params.id,
    method: 'get',
    params
  })
}

// @Tags Product
// @Summary 获取产品列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body modelInterface.PageInfo true "获取产品列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /product/list [post]
export const getProductList = (params) => {
  return service({
    url: '/product/list',
    method: 'post',
    data: params
  })
}

// @Tags Product
// @Summary 获取公开产品列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body modelInterface.PageInfo true "获取公开产品列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /product/public/list [post]
export const getPublicProductList = (params) => {
  return service({
    url: '/product/public/list',
    method: 'post',
    params
  })
}

// @Tags Product
// @Summary 获取指定分类的公开产品
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body modelInterface.PageInfo true "获取指定分类产品"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /product/public/by-category [get]
export const getProductsByCategory = (params) => {
  return service({
    url: '/product/public/by-category',
    method: 'get',
    params
  })
}
