package business

import (
	"github.com/gin-gonic/gin"
)

type StatisticsRouter struct{}

func (s *StatisticsRouter) InitStatisticsRouter(Router *gin.RouterGroup) {
	statisticsRouter := Router.Group("statistics")
	{
		statisticsRouter.GET("/overview", statisticsApi.GetOverview)
	}
}
