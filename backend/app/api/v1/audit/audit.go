package audit

import (
	"datauptwo/app/dto"
	"datauptwo/app/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type AuditApi struct {
	auditService *service.AuditService
}

func NewAuditApi() *AuditApi {
	return &AuditApi{
		auditService: service.NewAuditService(),
	}
}

// ListOperationLogs 获取操作日志列表
// @Summary 获取操作日志列表
// @Tags Audit
// @Security Bearer
// @Success 200 {object} dto.Response
// @Router /api/audit/operation-logs [get]
func (api *AuditApi) ListOperationLogs(c *gin.Context) {
	var req dto.OperationLogListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 20
	}

	result, err := api.auditService.ListOperationLogs(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: result})
}

// ExportOperationLogs 导出操作日志
// @Summary 导出操作日志
// @Tags Audit
// @Security Bearer
// @Success 200 {object} dto.Response
// @Router /api/audit/operation-logs/export [get]
func (api *AuditApi) ExportOperationLogs(c *gin.Context) {
	var req dto.OperationLogListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	logs, exportResult, err := api.auditService.ListOperationLogsForExport(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	// 使用 excelize 生成 Excel
	f := excelize.NewFile()
	defer f.Close()

	sheet := "操作日志"
	f.NewSheet(sheet)
	f.DeleteSheet("Sheet1")

	// 设置表头样式
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "#FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#005bbf"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "#e0e2ec", Style: 1},
			{Type: "right", Color: "#e0e2ec", Style: 1},
			{Type: "top", Color: "#e0e2ec", Style: 1},
			{Type: "bottom", Color: "#e0e2ec", Style: 1},
		},
	})

	// 标题行
	headers := []string{"ID", "用户名", "功能菜单", "动作", "资源类型", "资源名称", "IP地址", "详情", "时间"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	// 数据行
	for idx, log := range logs {
		row := idx + 2
		f.SetCellValue(sheet, "A"+strconv.Itoa(row), log.ID)
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), log.Username)
		f.SetCellValue(sheet, "C"+strconv.Itoa(row), log.MenuName)
		f.SetCellValue(sheet, "D"+strconv.Itoa(row), log.ActionText)
		f.SetCellValue(sheet, "E"+strconv.Itoa(row), log.ResourceType)
		f.SetCellValue(sheet, "F"+strconv.Itoa(row), log.ResourceName)
		f.SetCellValue(sheet, "G"+strconv.Itoa(row), log.IPAddress)
		f.SetCellValue(sheet, "H"+strconv.Itoa(row), log.Detail)
		f.SetCellValue(sheet, "I"+strconv.Itoa(row), log.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	// 设置数据行边框
	dataStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "#e0e2ec", Style: 1},
			{Type: "right", Color: "#e0e2ec", Style: 1},
			{Type: "top", Color: "#e0e2ec", Style: 1},
			{Type: "bottom", Color: "#e0e2ec", Style: 1},
		},
	})
	if len(logs) > 0 {
		startCell, _ := excelize.CoordinatesToCellName(1, 2)
		endCell, _ := excelize.CoordinatesToCellName(9, len(logs)+1)
		f.SetCellStyle(sheet, startCell, endCell, dataStyle)
	}

	// 设置列宽
	f.SetColWidth(sheet, "A", "A", 8)
	f.SetColWidth(sheet, "B", "B", 15)
	f.SetColWidth(sheet, "C", "C", 20)
	f.SetColWidth(sheet, "D", "D", 10)
	f.SetColWidth(sheet, "E", "E", 15)
	f.SetColWidth(sheet, "F", "F", 25)
	f.SetColWidth(sheet, "G", "G", 15)
	f.SetColWidth(sheet, "H", "H", 30)
	f.SetColWidth(sheet, "I", "I", 20)

	// 生成文件
	file, err := f.WriteToBuffer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: "生成Excel文件失败"})
		return
	}

	// 截断提示通过文件名反馈
	filename := "operation_logs.xlsx"
	if exportResult.Truncated {
		filename = "operation_logs_(已截断).xlsx"
	}
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", file.Bytes())
}

// ListLoginLogs 获取登录日志列表
// @Summary 获取登录日志列表
// @Tags Audit
// @Security Bearer
// @Success 200 {object} dto.Response
// @Router /api/audit/login-logs [get]
func (api *AuditApi) ListLoginLogs(c *gin.Context) {
	var req dto.LoginLogListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 20
	}

	result, err := api.auditService.ListLoginLogs(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: result})
}

// ExportLoginLogs 导出登录日志
// @Summary 导出登录日志
// @Tags Audit
// @Security Bearer
// @Success 200 {object} dto.Response
// @Router /api/audit/login-logs/export [get]
func (api *AuditApi) ExportLoginLogs(c *gin.Context) {
	var req dto.LoginLogListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	logs, exportResult, err := api.auditService.ListLoginLogsForExport(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	// 使用 excelize 生成 Excel
	f := excelize.NewFile()
	defer f.Close()

	sheet := "登录日志"
	f.NewSheet(sheet)
	f.DeleteSheet("Sheet1")

	// 设置表头样式
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "#FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#005bbf"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "#e0e2ec", Style: 1},
			{Type: "right", Color: "#e0e2ec", Style: 1},
			{Type: "top", Color: "#e0e2ec", Style: 1},
			{Type: "bottom", Color: "#e0e2ec", Style: 1},
		},
	})

	// 标题行
	headers := []string{"ID", "用户名", "状态", "IP地址", "失败原因", "是否MFA", "登录方式", "时间"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	// 数据行
	for idx, log := range logs {
		row := idx + 2
		f.SetCellValue(sheet, "A"+strconv.Itoa(row), log.ID)
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), log.Username)
		f.SetCellValue(sheet, "C"+strconv.Itoa(row), log.StatusText)
		f.SetCellValue(sheet, "D"+strconv.Itoa(row), log.IPAddress)
		f.SetCellValue(sheet, "E"+strconv.Itoa(row), log.FailReason)
		mfaText := "否"
		if log.MFAUsed {
			mfaText = "是"
		}
		f.SetCellValue(sheet, "F"+strconv.Itoa(row), mfaText)
		f.SetCellValue(sheet, "G"+strconv.Itoa(row), log.LoginMethod)
		f.SetCellValue(sheet, "H"+strconv.Itoa(row), log.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	// 设置数据行边框
	dataStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "#e0e2ec", Style: 1},
			{Type: "right", Color: "#e0e2ec", Style: 1},
			{Type: "top", Color: "#e0e2ec", Style: 1},
			{Type: "bottom", Color: "#e0e2ec", Style: 1},
		},
	})
	if len(logs) > 0 {
		startCell, _ := excelize.CoordinatesToCellName(1, 2)
		endCell, _ := excelize.CoordinatesToCellName(8, len(logs)+1)
		f.SetCellStyle(sheet, startCell, endCell, dataStyle)
	}

	// 设置列宽
	f.SetColWidth(sheet, "A", "A", 8)
	f.SetColWidth(sheet, "B", "B", 15)
	f.SetColWidth(sheet, "C", "C", 10)
	f.SetColWidth(sheet, "D", "D", 15)
	f.SetColWidth(sheet, "E", "E", 25)
	f.SetColWidth(sheet, "F", "F", 10)
	f.SetColWidth(sheet, "G", "G", 12)
	f.SetColWidth(sheet, "H", "H", 20)

	// 生成文件
	file, err := f.WriteToBuffer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: "生成Excel文件失败"})
		return
	}

	// 截断提示通过文件名反馈
	filename := "login_logs.xlsx"
	if exportResult.Truncated {
		filename = "login_logs_(已截断).xlsx"
	}
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", file.Bytes())
}
