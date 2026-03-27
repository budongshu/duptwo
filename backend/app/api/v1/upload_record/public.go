package api

import (
	"datauptwo/app/dto"
	"datauptwo/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PublicUploadRecordApi struct {
	uploadRecordService *service.UploadRecordService
}

func NewPublicUploadRecordApi() *PublicUploadRecordApi {
	return &PublicUploadRecordApi{
		uploadRecordService: service.NewUploadRecordService(),
	}
}

// Create 公开创建上传记录（无需登录）
// @Summary 公开创建上传记录
// @Tags PublicUploadRecord
// @Accept json
// @Param request body dto.UploadRecordCreateReq true "上传记录信息"
// @Success 200 {object} dto.Response
// @Router /public/upload-records [post]
func (api *PublicUploadRecordApi) Create(c *gin.Context) {
	var req dto.UploadRecordCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	record, err := api.uploadRecordService.Create(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "创建成功", Data: record})
}

// GetBySerialNo 根据流水号查询上传记录（无需登录）
// @Summary 根据流水号查询上传记录
// @Tags PublicUploadRecord
// @Param serialNo path string true "流水号"
// @Success 200 {object} dto.Response
// @Router /public/upload-records/{serialNo} [get]
func (api *PublicUploadRecordApi) GetBySerialNo(c *gin.Context) {
	serialNo := c.Param("serialNo")
	if serialNo == "" {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "流水号不能为空"})
		return
	}

	record, err := api.uploadRecordService.GetBySerialNo(serialNo)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{Code: 404, Message: "记录不存在"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: record})
}
