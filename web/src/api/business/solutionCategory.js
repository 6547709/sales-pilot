import service from '@/utils/request'

// @Tags SolutionCategory
// @Summary 创建解决方案分类
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body business.SolutionCategory true "创建解决方案分类"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /solution-category [post]
export const createSolutionCategory = (data) => {
  return service({
    url: '/solution-category',
    method: 'post',
    data
  })
}

// @Tags SolutionCategory
// @Summary 更新解决方案分类
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body business.SolutionCategory true "更新解决方案分类"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /solution-category [put]
export const updateSolutionCategory = (data) => {
  return service({
    url: '/solution-category',
    method: 'put',
    data
  })
}

// @Tags SolutionCategory
// @Summary 删除解决方案分类
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.SolutionCategory true "删除解决方案分类"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /solution-category [delete]
export const deleteSolutionCategory = (data) => {
  return service({
    url: '/solution-category',
    method: 'delete',
    data
  })
}

// @Tags SolutionCategory
// @Summary 获取单一解决方案分类
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query id true "分类ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /solution-category/{id} [get]
export const getSolutionCategory = (params) => {
  return service({
    url: '/solution-category/' + params.id,
    method: 'get',
    params
  })
}

// @Tags SolutionCategory
// @Summary 获取所有解决方案分类
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /solution-category/all [get]
export const getAllSolutionCategories = () => {
  return service({
    url: '/solution-category/all',
    method: 'get'
  })
}
