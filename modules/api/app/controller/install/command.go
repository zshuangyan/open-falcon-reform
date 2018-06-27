package install

import (
	"net/http"
	"github.com/gin-gonic/gin"
	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
)

func GetInstallCommand(c *gin.Context){
	command := "curl -SL   http://10.202.42.2:8089/install.sh | sh && cd /home/work/open-falcon && ./open-falcon start agent"
	h.JSONResponse(c, http.StatusOK, 0, "get install commands succeed", command)
	return
}