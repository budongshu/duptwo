package api

import (
	"datauptwo/app/dto"
	"datauptwo/app/service"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type UploadRecordApi struct {
	uploadRecordService *service.UploadRecordService
}

func NewUploadRecordApi() *UploadRecordApi {
	return &UploadRecordApi{
		uploadRecordService: service.NewUploadRecordService(),
	}
}

// Create 创建上传记录
// @Summary 创建上传记录
// @Tags UploadRecord
// @Accept json
// @Param request body dto.UploadRecordCreateReq true "上传记录信息"
// @Success 200 {object} dto.Response
// @Router /api/v1/upload-records [post]
func (api *UploadRecordApi) Create(c *gin.Context) {
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

// List 记录列表
// @Summary 上传记录列表
// @Tags UploadRecord
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param dataType query string false "数据类型"
// @Param status query string false "状态"
// @Param uploader query string false "上传人"
// @Param startDate query string false "开始日期"
// @Param endDate query string false "结束日期"
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} dto.Response
// @Router /api/v1/upload-records [get]
func (api *UploadRecordApi) List(c *gin.Context) {
	var req dto.UploadRecordListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		req.Page = 1
		req.PageSize = 20
	}

	result, err := api.uploadRecordService.List(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: result})
}

// GetByID 获取记录详情
// @Summary 获取上传记录详情
// @Tags UploadRecord
// @Param id path int true "记录ID"
// @Success 200 {object} dto.Response
// @Router /api/v1/upload-records/{id} [get]
func (api *UploadRecordApi) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	detail, err := api.uploadRecordService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{Code: 404, Message: "记录不存在"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: detail})
}

// Update 更新记录
// @Summary 更新上传记录
// @Tags UploadRecord
// @Accept json
// @Param request body dto.UploadRecordUpdateReq true "记录信息"
// @Success 200 {object} dto.Response
// @Router /api/v1/upload-records [put]
func (api *UploadRecordApi) Update(c *gin.Context) {
	var req dto.UploadRecordUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	record, err := api.uploadRecordService.Update(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "更新成功", Data: record})
}

// Delete 删除记录
// @Summary 删除上传记录
// @Tags UploadRecord
// @Param id path int true "记录ID"
// @Success 200 {object} dto.Response
// @Router /api/v1/upload-records/{id} [delete]
func (api *UploadRecordApi) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := api.uploadRecordService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "删除成功"})
}

// BatchDelete 批量删除上传记录
// @Summary 批量删除上传记录
// @Tags UploadRecord
// @Security Bearer
// @Accept json
// @Param request body dto.BatchDeleteReq true "记录ID列表"
// @Success 200 {object} dto.Response
// @Router /api/upload-records/batch-delete [post]
func (api *UploadRecordApi) BatchDelete(c *gin.Context) {
	var req dto.BatchDeleteReq
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "请选择要删除的记录"})
		return
	}

	if err := api.uploadRecordService.BatchDelete(req.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: fmt.Sprintf("成功删除 %d 条记录", len(req.IDs))})
}

// 获取统计数据
// @Summary 获取上传记录统计
// @Tags UploadRecord
// @Param startDate query string false "开始日期"
// @Param endDate query string false "结束日期"
// @Param projectName query string false "项目名称"
// @Param dataType query string false "数据类型"
// @Param status query string false "状态"
// @Param uploader query string false "上传人"
// @Success 200 {object} dto.Response
// @Router /api/v1/upload-records/statistics [get]
func (api *UploadRecordApi) GetStatistics(c *gin.Context) {
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	projectName := c.Query("projectName")
	dataType := c.Query("dataType")
	status := c.Query("status")
	uploader := c.Query("uploader")

	stats, err := api.uploadRecordService.GetStatistics(startDate, endDate, projectName, dataType, status, uploader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: stats})
}

// GetRecent 获取最近上传记录
// @Summary 获取最近上传记录
// @Tags UploadRecord
// @Param limit query int false "数量限制"
// @Param projectName query string false "项目名称"
// @Param dataType query string false "数据类型"
// @Param status query string false "状态"
// @Param uploader query string false "上传人"
// @Success 200 {object} dto.Response
// @Router /api/v1/upload-records/recent [get]
func (api *UploadRecordApi) GetRecent(c *gin.Context) {
	limit := 10
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}
	projectName := c.Query("projectName")
	dataType := c.Query("dataType")
	status := c.Query("status")
	uploader := c.Query("uploader")

	records, err := api.uploadRecordService.GetRecent(limit, projectName, dataType, status, uploader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: records})
}

// GetUploaderList 获取所有上传者列表
// @Summary 获取上传者列表
// @Tags UploadRecord
// @Success 200 {object} dto.Response
// @Router /api/v1/upload-records/uploaders [get]
func (api *UploadRecordApi) GetUploaderList(c *gin.Context) {
	uploaders, err := api.uploadRecordService.GetUploaderList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: uploaders})
}

// GetTemplate 获取导入模板信息（字段定义）
// @Summary 获取导入模板字段定义
// @Tags UploadRecord
// @Security Bearer
// @Success 200 {object} dto.Response
// @Router /api/v1/upload-records/template [get]
func (api *UploadRecordApi) GetTemplate(c *gin.Context) {
	template := api.uploadRecordService.GetImportTemplate()
	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: template})
}

// Import 批量导入上传记录
// @Summary 批量导入上传记录
// @Tags UploadRecord
// @Security Bearer
// @Accept multipart/form-data
// @Param file formData file true "Excel文件"
// @Success 200 {object} dto.Response
// @Router /api/v1/upload-records/import [post]
func (api *UploadRecordApi) Import(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "请上传Excel文件"})
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: "无法读取上传文件"})
		return
	}
	defer src.Close()

	// 读取Excel
	f, err := excelize.OpenReader(src)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "无法解析Excel文件，请确保是标准xlsx格式"})
		return
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "Excel文件中没有工作表"})
		return
	}
	sheetName := sheets[0]

	rows, err := f.GetRows(sheetName)
	if err != nil || len(rows) < 2 {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "Excel文件数据为空或格式错误（需要表头+数据行）"})
		return
	}

	// 第一行是表头，建立列索引映射
	headerRow := rows[0]
	colIndex := make(map[string]int)
	for idx, colName := range headerRow {
		colIndex[colName] = idx
	}

	// 将数据行转换为 map[字段代码]值
	dataRows := []map[string]string{}
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		rowMap := make(map[string]string)
		for colName, idx := range colIndex {
			if idx < len(row) {
				rowMap[colName] = row[idx]
			}
		}
		// 跳过空行
		if rowMap["dataType"] == "" && rowMap["destPath"] == "" && rowMap["uploader"] == "" {
			continue
		}
		dataRows = append(dataRows, rowMap)
	}

	result := api.uploadRecordService.ImportRecords(dataRows)
	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: fmt.Sprintf("导入完成：成功 %d 行，失败 %d 行", result.Success, result.Failed), Data: result})
}

// Export 导出上传记录为 Excel
// @Summary 导出上传记录
// @Tags UploadRecord
// @Param dataType query string false "数据类型"
// @Param status query string false "状态"
// @Param uploader query string false "上传人"
// @Param startDate query string false "开始日期"
// @Param endDate query string false "结束日期"
// @Param keyword query string false "搜索关键词"
// @Success 200 {file} file "Excel文件"
// @Router /api/v1/upload-records/export [get]
func (api *UploadRecordApi) Export(c *gin.Context) {
	var req dto.UploadRecordListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		req = dto.UploadRecordListReq{}
	}

	data, err := api.uploadRecordService.Export(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: "导出失败: " + err.Error()})
		return
	}

	filename := fmt.Sprintf("上传记录_%s.xlsx", time.Now().Format("20060102150405"))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", filename))
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}
