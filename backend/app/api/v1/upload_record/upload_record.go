package api

import (
	"datauptwo/app/dto"
	"datauptwo/app/service"
	"datauptwo/global"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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
// @Param diskLabel query string false "磁盘标签"
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

// BatchUpdateStatus 批量更新上传记录状态
// @Summary 批量更新上传记录状态
// @Tags UploadRecord
// @Security Bearer
// @Accept json
// @Param request body dto.BatchUpdateStatusReq true "记录ID列表和新状态"
// @Success 200 {object} dto.Response
// @Router /api/upload-records/batch-update-status [post]
func (api *UploadRecordApi) BatchUpdateStatus(c *gin.Context) {
	var req dto.BatchUpdateStatusReq
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "请选择要更新的记录"})
		return
	}

	// 验证状态值
	validStatuses := map[string]bool{"pending": true, "processing": true, "completed": true, "failed": true}
	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "状态值非法，支持值: pending/processing/completed/failed"})
		return
	}

	if err := api.uploadRecordService.BatchUpdateStatus(req.IDs, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: fmt.Sprintf("成功更新 %d 条记录的状态为「%s」", len(req.IDs), req.Status)})
}

// 获取统计数据
// @Summary 获取上传记录统计
// @Tags UploadRecord
// @Param startDate query string false "开始日期"
// @Param endDate query string false "结束日期"
// @Param projectName query string false "项目名称"
// @Param diskLabel query string false "磁盘标签"
// @Param status query string false "状态"
// @Param uploader query string false "上传人"
// @Success 200 {object} dto.Response
// @Router /api/v1/upload-records/statistics [get]
func (api *UploadRecordApi) GetStatistics(c *gin.Context) {
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	projectName := c.Query("projectName")
	diskLabel := c.Query("diskLabel")
	status := c.Query("status")
	uploader := c.Query("uploader")

	stats, err := api.uploadRecordService.GetStatistics(startDate, endDate, projectName, diskLabel, status, uploader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: stats})
}

// GetDiskLabelStatuses 获取所有磁盘标签及其综合状态
// @Summary 获取所有磁盘标签及其综合状态
// @Tags UploadRecord
// @Description 返回所有磁盘标签及其上传状态：completed全部完成(绿色)，failed全部失败(红色)，mixed部分失败(橙色)，pending处理中(灰色)
// @Security Bearer
// @Param startDate query string false "开始日期"
// @Param endDate query string false "结束日期"
// @Success 200 {object} dto.Response
// @Router /api/v1/upload-records/disk-labels [get]
func (api *UploadRecordApi) GetDiskLabelStatuses(c *gin.Context) {
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	labels, err := api.uploadRecordService.GetDiskLabelStatuses(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: labels})
}

// GetRecent 获取最近上传记录
// @Summary 获取最近上传记录
// @Tags UploadRecord
// @Param limit query int false "数量限制"
// @Param projectName query string false "项目名称"
// @Param diskLabel query string false "磁盘标签"
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
	diskLabel := c.Query("diskLabel")
	status := c.Query("status")
	uploader := c.Query("uploader")

	records, err := api.uploadRecordService.GetRecent(limit, projectName, diskLabel, status, uploader)
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

// Preview 预览上传文件（解析Excel，返回行数）
// @Summary 预览Excel文件行数
// @Tags UploadRecord
// @Security Bearer
// @Accept multipart/form-data
// @Param file formData file true "Excel文件"
// @Success 200 {object} dto.Response
// @Router /api/v1/upload-records/preview [post]
func (api *UploadRecordApi) Preview(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "请上传Excel文件"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: "无法读取上传文件"})
		return
	}
	defer src.Close()

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

	rows, err := f.GetRows(sheets[0])
	if err != nil || len(rows) < 2 {
		c.JSON(http.StatusOK, dto.Response{Code: 200, Data: map[string]interface{}{"totalRows": 0, "dataRows": 0, "sheetName": sheets[0]}})
		return
	}

	colNameToCode := map[string]string{
		"磁盘标签":         "diskLabel",
		"项目名称":        "projectName",
		"目标路径":        "destPath",
		"文件大小(字节)":   "fileSize",
		"上传人":          "uploader",
		"上传状态":        "status",
		"备注":           "remark",
		"创建时间":        "createdAt",
	}

	// 查找表头行（包含"磁盘标签"的行）
	headerRowIndex := -1
	var headerRow []string
	for i, row := range rows {
		for _, cell := range row {
			cellClean := strings.TrimSuffix(strings.TrimSpace(cell), " *")
			if cellClean == "磁盘标签" {
				headerRowIndex = i
				headerRow = row
				break
			}
		}
		if headerRowIndex >= 0 {
			break
		}
	}

	if headerRowIndex < 0 {
		c.JSON(http.StatusOK, dto.Response{Code: 200, Data: map[string]interface{}{"totalRows": 0, "dataRows": 0, "sheetName": sheets[0], "error": "未找到有效表头"}})
		return
	}

	// 建立列索引映射
	colIndex := make(map[string]int)
	for idx, colName := range headerRow {
		colNameClean := strings.TrimSuffix(strings.TrimSpace(colName), " *")
		colIndex[colNameClean] = idx
	}

	dataRows := 0
	for i := headerRowIndex + 1; i < len(rows); i++ {
		row := rows[i]
		rowMap := make(map[string]string)
		for colName, idx := range colIndex {
			if idx < len(row) {
				code := colNameToCode[colName]
				if code == "" {
					code = colName
				}
				rowMap[code] = row[idx]
			}
		}
		// 跳过空行
		if rowMap["diskLabel"] == "" && rowMap["destPath"] == "" && rowMap["uploader"] == "" {
			continue
		}
		dataRows++
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: map[string]interface{}{
		"totalRows":  len(rows) - headerRowIndex - 1,
		"dataRows":   dataRows,
		"sheetName":  sheets[0],
		"headers":    headerRow,
	}})
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

	// 中文列名 → 字段代码映射（与模板生成顺序一致）
	colNameToCode := map[string]string{
		"磁盘标签":         "diskLabel",
		"项目名称":        "projectName",
		"目标路径":        "destPath",
		"文件大小(字节)":   "fileSize",
		"上传人":          "uploader",
		"上传状态":        "status",
		"备注":           "remark",
		"创建时间":        "createdAt",
	}

	// 查找表头行（包含"磁盘标签"的行，可能是"磁盘标签"或"磁盘标签 *"）
	headerRowIndex := -1
	var headerRow []string
	for i, row := range rows {
		for _, cell := range row {
			// 去掉可能的 " *" 后缀
			cellClean := strings.TrimSuffix(strings.TrimSpace(cell), " *")
			if cellClean == "磁盘标签" {
				headerRowIndex = i
				headerRow = row
				break
			}
		}
		if headerRowIndex >= 0 {
			break
		}
	}

	if headerRowIndex < 0 {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "未找到有效的表头行（需包含「磁盘标签」列）"})
		return
	}

	global.AppLogger.Info("导入Excel表头行索引: %d, 列名: %v", headerRowIndex, headerRow)

	// 建立列索引映射
	colIndex := make(map[string]int)
	for idx, colName := range headerRow {
		colNameClean := strings.TrimSuffix(strings.TrimSpace(colName), " *")
		colIndex[colNameClean] = idx
	}

	// 创建时间列索引（用于处理Excel日期单元格）
	dateColIdx, dateColExists := colIndex["创建时间"]

	dataRows := []map[string]string{}
	for i := headerRowIndex + 1; i < len(rows); i++ {
		row := rows[i]
		rowMap := make(map[string]string)
		for colName, idx := range colIndex {
			// 特殊处理"创建时间"列：必须优先处理并跳过，避免被 general mapping 覆盖
			if dateColExists && idx == dateColIdx {
				// 日期列：即使该行没有这列数据，也跳过 general mapping
				if idx < len(row) {
					dateVal := row[idx]
					if dateVal != "" {
						// 优先：当 Excel 数字日期序列号解析
						if serialNum, err := strconv.ParseFloat(dateVal, 64); err == nil {
							if t, err := excelize.ExcelDateToTime(serialNum, false); err == nil {
								rowMap["createdAt"] = t.Format("2006-01-02 15:04:05")
							} else {
								rowMap["createdAt"] = dateVal // 服务层会再次解析
							}
						} else {
							// 非数字 → 当格式化日期字符串直接用
							rowMap["createdAt"] = dateVal
						}
					}
				}
				continue // ← 关键：日期列无论数据是否存在，都必须跳过 general mapping
			}
			// general mapping
			if idx < len(row) {
				code := colNameToCode[colName]
				if code == "" {
					code = colName
				}
				rowMap[code] = row[idx]
			}
		}
		// 跳过空行
		if rowMap["diskLabel"] == "" && rowMap["destPath"] == "" && rowMap["uploader"] == "" {
			continue
		}
		dataRows = append(dataRows, rowMap)
	}

	global.AppLogger.Info("导入Excel识别到 %d 行数据行（总行数=%d）", len(dataRows), len(rows))

	result := api.uploadRecordService.ImportRecords(dataRows)
	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: fmt.Sprintf("导入完成：成功 %d 行，失败 %d 行", result.Success, result.Failed), Data: result})
}

// Export 导出上传记录为 Excel
// @Summary 导出上传记录
// @Tags UploadRecord
// @Param diskLabel query string false "磁盘标签"
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
