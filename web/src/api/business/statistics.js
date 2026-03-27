import service from '@/utils/request'

// @Tags Statistics
// @Summary 获取业务统计数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /statistics/overview [get]
export const getStatisticsOverview = () => {
  return service({
    url: '/statistics/overview',
    method: 'get'
  })
}
