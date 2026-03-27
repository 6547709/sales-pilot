import service from '@/utils/request'

// ===================== Layer APIs =====================
// @Tags Topology
// @Summary 创建拓扑层级
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body business.TopologyLayer true "创建拓扑层级"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /topology/layer [post]
export const createLayer = (data) => {
  return service({
    url: '/topology/layer',
    method: 'post',
    data
  })
}

// @Tags Topology
// @Summary 更新拓扑层级
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body business.TopologyLayer true "更新拓扑层级"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /topology/layer [put]
export const updateLayer = (data) => {
  return service({
    url: '/topology/layer',
    method: 'put',
    data
  })
}

// @Tags Topology
// @Summary 删除拓扑层级
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.TopologyLayer true "删除拓扑层级"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /topology/layer [delete]
export const deleteLayer = (data) => {
  return service({
    url: '/topology/layer',
    method: 'delete',
    data
  })
}

// @Tags Topology
// @Summary 获取所有拓扑层级
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /topology/layer/all [get]
export const getLayers = () => {
  return service({
    url: '/topology/layer/all',
    method: 'get'
  })
}

// ===================== Category APIs =====================
// @Tags Topology
// @Summary 创建拓扑分类
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body business.TopologyCategory true "创建拓扑分类"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /topology/category [post]
export const createCategory = (data) => {
  return service({
    url: '/topology/category',
    method: 'post',
    data
  })
}

// @Tags Topology
// @Summary 更新拓扑分类
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body business.TopologyCategory true "更新拓扑分类"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /topology/category [put]
export const updateCategory = (data) => {
  return service({
    url: '/topology/category',
    method: 'put',
    data
  })
}

// @Tags Topology
// @Summary 删除拓扑分类
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.TopologyCategory true "删除拓扑分类"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /topology/category [delete]
export const deleteCategory = (data) => {
  return service({
    url: '/topology/category',
    method: 'delete',
    data
  })
}

// @Tags Topology
// @Summary 按层级获取拓扑分类
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query layerId true "层级ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /topology/category/by-layer [get]
export const getCategoriesByLayer = (params) => {
  return service({
    url: '/topology/category/by-layer',
    method: 'get',
    params
  })
}

// ===================== Vendor APIs =====================
// @Tags Topology
// @Summary 创建拓扑厂商
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body business.TopologyVendor true "创建拓扑厂商"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /topology/vendor [post]
export const createVendor = (data) => {
  return service({
    url: '/topology/vendor',
    method: 'post',
    data
  })
}

// @Tags Topology
// @Summary 更新拓扑厂商
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body business.TopologyVendor true "更新拓扑厂商"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /topology/vendor [put]
export const updateVendor = (data) => {
  return service({
    url: '/topology/vendor',
    method: 'put',
    data
  })
}

// @Tags Topology
// @Summary 删除拓扑厂商
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.TopologyVendor true "删除拓扑厂商"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /topology/vendor [delete]
export const deleteVendor = (data) => {
  return service({
    url: '/topology/vendor',
    method: 'delete',
    data
  })
}

// @Tags Topology
// @Summary 按分类获取拓扑厂商
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query categoryId true "分类ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /topology/vendor/by-category [get]
export const getVendorsByCategory = (params) => {
  return service({
    url: '/topology/vendor/by-category',
    method: 'get',
    params
  })
}

// ===================== Public APIs (No Auth Required) =====================
// @Tags Topology
// @Summary 获取完整拓扑结构（供首页全景图使用）
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /topology/full [get]
export const getFullTopology = () => {
  return service({
    url: '/topology/full',
    method: 'get'
  })
}
