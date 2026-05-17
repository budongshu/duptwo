package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"datauptwo/app/dto"
	"datauptwo/app/service"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type PersonnelApi struct {
	personnelService *service.PersonnelService
	auditService     *service.AuditService
}

func NewPersonnelApi() *PersonnelApi {
	return &PersonnelApi{
		personnelService: service.NewPersonnelService(),
		auditService:     service.NewAuditService(),
	}
}

func (api *PersonnelApi) getUserID(c *gin.Context) uint {
	if id, exists := c.Get("userId"); exists {
		return id.(uint)
	}
	return 0
}

func (api *PersonnelApi) getUsername(c *gin.Context) string {
	if username, exists := c.Get("username"); exists {
		return username.(string)
	}
	return ""
}

// Create 创建人员
// @Summary 创建人员
// @Tags Personnel
// @Accept json
// @Param request body dto.PersonnelCreateReq true "人员信息"
// @Success 200 {object} dto.Response
// @Router /api/personnels [post]
func (api *PersonnelApi) Create(c *gin.Context) {
	var req dto.PersonnelCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	personnel, err := api.personnelService.Create(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	api.auditService.LogOperation(
		api.getUserID(c), api.getUsername(c),
		"人员管理", "create", "Personnel", personnel.ID,
		personnel.Name, c.ClientIP(), c.GetHeader("User-Agent"), nil,
	)

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "创建成功", Data: personnel})
}

// List 人员列表
// @Summary 人员列表
// @Tags Personnel
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param keyword query string false "关键词"
// @Param status query string false "状态"
// @Param onProject query string false "在项状态"
// @Success 200 {object} dto.Response
// @Router /api/personnels [get]
func (api *PersonnelApi) List(c *gin.Context) {
	var req dto.PersonnelListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		req.Page = 1
		req.PageSize = 20
	}

	result, err := api.personnelService.List(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: result})
}

// GetByID 获取人员详情
// @Summary 获取人员详情
// @Tags Personnel
// @Param id path int true "人员ID"
// @Success 200 {object} dto.Response
// @Router /api/personnels/{id} [get]
func (api *PersonnelApi) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	personnel, err := api.personnelService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{Code: 404, Message: "人员不存在"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: personnel})
}

// Update 更新人员
// @Summary 更新人员
// @Tags Personnel
// @Accept json
// @Param request body dto.PersonnelUpdateReq true "人员信息"
// @Success 200 {object} dto.Response
// @Router /api/personnels [put]
func (api *PersonnelApi) Update(c *gin.Context) {
	var req dto.PersonnelUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	// 先获取人员信息用于日志
	personnel, _ := api.personnelService.GetByID(req.ID)
	personnelName := ""
	if personnel != nil {
		personnelName = personnel.Name
	}

	if err := api.personnelService.Update(req); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	api.auditService.LogOperation(
		api.getUserID(c), api.getUsername(c),
		"人员管理", "update", "Personnel", req.ID,
		personnelName, c.ClientIP(), c.GetHeader("User-Agent"), nil,
	)

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "更新成功"})
}

// Delete 删除人员
// @Summary 删除人员
// @Tags Personnel
// @Param id path int true "人员ID"
// @Success 200 {object} dto.Response
// @Router /api/personnels/{id} [delete]
func (api *PersonnelApi) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	// 先获取人员信息用于日志
	personnel, _ := api.personnelService.GetByID(uint(id))
	personnelName := ""
	if personnel != nil {
		personnelName = personnel.Name
	}

	if err := api.personnelService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	api.auditService.LogOperation(
		api.getUserID(c), api.getUsername(c),
		"人员管理", "delete", "Personnel", uint(id),
		personnelName, c.ClientIP(), c.GetHeader("User-Agent"), nil,
	)

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "删除成功"})
}

// BatchDelete 批量删除人员
// @Summary 批量删除人员
// @Tags Personnel
// @Accept json
// @Param request body dto.BatchDeleteReq true "人员ID列表"
// @Success 200 {object} dto.Response
// @Router /api/personnels/batch-delete [post]
func (api *PersonnelApi) BatchDelete(c *gin.Context) {
	var req dto.BatchDeleteReq
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "请选择要删除的人员"})
		return
	}

	if err := api.personnelService.BatchDelete(req.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	api.auditService.LogOperation(
		api.getUserID(c), api.getUsername(c),
		"人员管理", "batch_delete", "Personnel", 0,
		fmt.Sprintf("批量删除 %d 个人员", len(req.IDs)), c.ClientIP(), c.GetHeader("User-Agent"),
		map[string]interface{}{"ids": req.IDs},
	)

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "成功删除"})
}

// ListAll 获取所有人员（用于下拉选择）
// @Summary 获取所有人员
// @Tags Personnel
// @Success 200 {object} dto.Response
// @Router /api/personnels/all [get]
func (api *PersonnelApi) ListAll(c *gin.Context) {
	personnels, err := api.personnelService.ListAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: personnels})
}

// Statistics 获取人员统计
// @Summary 人员统计
// @Tags Personnel
// @Param keyword query string false "关键词"
// @Param status query string false "状态"
// @Param onProject query string false "在项状态"
// @Success 200 {object} dto.Response
// @Router /api/personnels/statistics [get]
func (api *PersonnelApi) Statistics(c *gin.Context) {
	var req dto.PersonnelListReq
	c.ShouldBindQuery(&req)

	stats, err := api.personnelService.Statistics(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: stats})
}

// Export 导出人员Excel
// @Summary 导出人员Excel
// @Tags Personnel
// @Param keyword query string false "关键词"
// @Param status query string false "状态"
// @Param onProject query string false "在项状态"
// @Success 200 {file} file
// @Router /api/personnels/export [get]
func (api *PersonnelApi) Export(c *gin.Context) {
	var req dto.PersonnelListReq
	c.ShouldBindQuery(&req)

	list, err := api.personnelService.ListForExport(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	f := excelize.NewFile()
	defer f.Close()

	sheet := "人员列表"
	f.NewSheet(sheet)
	f.DeleteSheet("Sheet1")

	// 表头样式
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "#FFFFFF", Size: 11},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#1a5276"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "#d5dbe0", Style: 1},
			{Type: "right", Color: "#d5dbe0", Style: 1},
			{Type: "top", Color: "#d5dbe0", Style: 1},
			{Type: "bottom", Color: "#d5dbe0", Style: 1},
		},
	})
	// 数据行样式（隔行浅蓝）
	dataStyle1, _ := f.NewStyle(&excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#f8f9fc"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "#d5dbe0", Style: 1},
			{Type: "right", Color: "#d5dbe0", Style: 1},
			{Type: "top", Color: "#d5dbe0", Style: 1},
			{Type: "bottom", Color: "#d5dbe0", Style: 1},
		},
	})
	dataStyle2, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "#d5dbe0", Style: 1},
			{Type: "right", Color: "#d5dbe0", Style: 1},
			{Type: "top", Color: "#d5dbe0", Style: 1},
			{Type: "bottom", Color: "#d5dbe0", Style: 1},
		},
	})

	headers := []string{"序号", "姓名", "手机号", "邮箱", "所属公司", "职位", "工作经验", "入项时间", "立项时间", "在项状态", "薪资", "驻场地点", "状态", "备注", "排序", "创建时间"}
	headerWidths := map[string]float64{"A": 6, "B": 10, "C": 14, "D": 22, "E": 16, "F": 12, "G": 10, "H": 12, "I": 12, "J": 10, "K": 12, "L": 16, "M": 8, "N": 18, "O": 6, "P": 18}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}
	for col, w := range headerWidths {
		f.SetColWidth(sheet, col, col, w)
	}

	statusMap := map[string]string{"active": "启用", "inactive": "禁用", "": "启用"}

	for idx, p := range list {
		row := idx + 2
		rowData := []interface{}{
			row,
			p.Name,
			p.Phone,
			p.Email,
			p.Company,
			p.Position,
			p.WorkExperience,
			p.EntryDate,
			p.ProjectStartDate,
			p.OnProjectStatus,
			p.Salary,
			p.Location,
			statusMap[p.Status],
			p.Remark,
			p.Sort,
			p.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		for col, val := range rowData {
			cell, _ := excelize.CoordinatesToCellName(col+1, row)
			f.SetCellValue(sheet, cell, val)
			if idx%2 == 0 {
				f.SetCellStyle(sheet, cell, cell, dataStyle1)
			} else {
				f.SetCellStyle(sheet, cell, cell, dataStyle2)
			}
		}
	}

	f.SetActiveSheet(0)
	fileName := "人员列表_" + time.Now().Format("20060102150405") + ".xlsx"
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("FileName", fileName)
	c.Header("Access-Control-Expose-Headers", "Content-Disposition")

	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: "导出失败"})
	}
}

// GetTemplate 获取导入模板字段定义
// @Summary 获取导入模板字段定义
// @Tags Personnel
// @Security Bearer
// @Success 200 {object} dto.Response
// @Router /api/personnels/template [get]
func (api *PersonnelApi) GetTemplate(c *gin.Context) {
	fields := []dto.ImportTemplateField{
		{Field: "姓名", Code: "name", Required: true, Type: "text", MaxLength: 64, Example: "张三"},
		{Field: "手机号", Code: "phone", Required: false, Type: "text", MaxLength: 32, Example: "13800138000"},
		{Field: "邮箱", Code: "email", Required: false, Type: "text", MaxLength: 128, Example: "zhangsan@example.com"},
		{Field: "所属公司", Code: "company", Required: false, Type: "text", MaxLength: 128, Example: "XX科技有限公司"},
		{Field: "职位", Code: "position", Required: false, Type: "select", Options: "测试工程师,前端工程师,算法工程师,DBA数据库,网络工程师,安全工程师,开发工程师,运维工程师,运营人员,合规专家,解决方案,商务人员,成本人员,驻场人员,驻场人员-ODC,项目管理,合规负责人,产品人员,其他人员", Example: "测试工程师"},
		{Field: "工作经验", Code: "workExperience", Required: false, Type: "text", MaxLength: 64, Example: "3年"},
		{Field: "入项时间", Code: "entryDate", Required: false, Type: "date", Example: "2024-01-01"},
		{Field: "立项时间", Code: "projectStartDate", Required: false, Type: "date", Example: "2024-06-01"},
		{Field: "在项状态", Code: "onProjectStatus", Required: false, Type: "select", Options: "在项,离项", Example: "在项"},
		{Field: "薪资", Code: "salary", Required: false, Type: "text", MaxLength: 32, Example: "15000"},
		{Field: "驻场地点", Code: "location", Required: false, Type: "text", MaxLength: 128, Example: "北京市朝阳区"},
		{Field: "状态", Code: "status", Required: false, Type: "select", Options: "active,inactive", Example: "active"},
		{Field: "备注", Code: "remark", Required: false, Type: "text", MaxLength: 256, Example: "无"},
		{Field: "排序", Code: "sort", Required: false, Type: "number", MaxLength: 0, Example: "0"},
	}
	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: dto.ImportTemplateResp{Fields: fields, SheetName: "人员列表", Title: "人员导入模板"}})
}

// Preview 预览上传文件（解析Excel，返回行数）
// @Summary 预览Excel文件行数
// @Tags Personnel
// @Security Bearer
// @Accept multipart/form-data
// @Param file formData file true "Excel文件"
// @Success 200 {object} dto.Response
// @Router /api/personnels/preview [post]
func (api *PersonnelApi) Preview(c *gin.Context) {
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

	// 查找表头行（包含"姓名"的行）
	headerRowIndex := -1
	var headerRow []string
	for i, row := range rows {
		for _, cell := range row {
			cellClean := strings.TrimSpace(cell)
			if cellClean == "姓名" || strings.TrimSuffix(cellClean, " *") == "姓名" {
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

	colIndex := make(map[string]int)
	for idx, colName := range headerRow {
		colNameClean := strings.TrimSuffix(strings.TrimSpace(colName), " *")
		colIndex[colNameClean] = idx
	}

	dataRows := 0
	for i := headerRowIndex + 1; i < len(rows); i++ {
		row := rows[i]
		// 跳过空行
		hasData := false
		for _, cell := range row {
			if strings.TrimSpace(cell) != "" {
				hasData = true
				break
			}
		}
		if hasData {
			dataRows++
		}
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: map[string]interface{}{
		"totalRows":  len(rows) - headerRowIndex - 1,
		"dataRows":   dataRows,
		"sheetName":  sheets[0],
		"headers":    headerRow,
	}})
}

// Import 批量导入人员
// @Summary 批量导入人员
// @Tags Personnel
// @Security Bearer
// @Accept multipart/form-data
// @Param file formData file true "Excel文件"
// @Success 200 {object} dto.Response
// @Router /api/personnels/import [post]
func (api *PersonnelApi) Import(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "Excel文件数据为空或格式错误（需要表头+数据行）"})
		return
	}

	// 中文列名 → 字段代码映射
	colNameToCode := map[string]string{
		"姓名":      "name",
		"手机号":     "phone",
		"邮箱":      "email",
		"所属公司":    "company",
		"职位":      "position",
		"工作经验":    "workExperience",
		"入项时间":    "entryDate",
		"立项时间":    "projectStartDate",
		"在项状态":    "onProjectStatus",
		"薪资":      "salary",
		"驻场地点":    "location",
		"状态":      "status",
		"备注":      "remark",
		"排序":      "sort",
	}

	// 查找表头行
	headerRowIndex := -1
	var headerRow []string
	for i, row := range rows {
		for _, cell := range row {
			cellClean := strings.TrimSpace(cell)
			if cellClean == "姓名" || strings.TrimSuffix(cellClean, " *") == "姓名" {
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
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "未找到有效的表头行（需包含「姓名」列）"})
		return
	}

	// 建立列索引映射
	colIndex := make(map[string]int)
	for idx, colName := range headerRow {
		colNameClean := strings.TrimSuffix(strings.TrimSpace(colName), " *")
		colIndex[colNameClean] = idx
	}

	// 日期列索引
	entryDateColIdx, entryDateColExists := colIndex["入项时间"]
	projectDateColIdx, projectDateColExists := colIndex["立项时间"]

	dataRows := []map[string]string{}
	for i := headerRowIndex + 1; i < len(rows); i++ {
		row := rows[i]
		rowMap := make(map[string]string)
		for colName, idx := range colIndex {
			// 跳过日期列，统一在下面处理
			if (entryDateColExists && idx == entryDateColIdx) || (projectDateColExists && idx == projectDateColIdx) {
				continue
			}
			if idx < len(row) {
				code := colNameToCode[colName]
				if code == "" {
					code = colName
				}
				rowMap[code] = strings.TrimSpace(row[idx])
			}
		}

		// 处理日期列
		if entryDateColExists && entryDateColIdx < len(row) {
			rowMap["entryDate"] = parseExcelDate(row[entryDateColIdx])
		}
		if projectDateColExists && projectDateColIdx < len(row) {
			rowMap["projectStartDate"] = parseExcelDate(row[projectDateColIdx])
		}

		// 跳过空行
		if rowMap["name"] == "" {
			continue
		}
		dataRows = append(dataRows, rowMap)
	}

	// 执行导入
	success := 0
	failCount := 0
	var failRows []dto.ImportFailRow

	for idx, rowMap := range dataRows {
		req := dto.PersonnelCreateReq{
			Name:             rowMap["name"],
			Phone:            rowMap["phone"],
			Email:           rowMap["email"],
			Company:         rowMap["company"],
			Position:        rowMap["position"],
			WorkExperience:  rowMap["workExperience"],
			EntryDate:       rowMap["entryDate"],
			ProjectStartDate: rowMap["projectStartDate"],
			OnProjectStatus: rowMap["onProjectStatus"],
			Salary:          rowMap["salary"],
			Location:        rowMap["location"],
			Remark:          rowMap["remark"],
			Status:          rowMap["status"],
		}
		if rowMap["sort"] != "" {
			if sort, err := strconv.Atoi(rowMap["sort"]); err == nil {
				req.Sort = sort
			}
		}
		if req.Status == "" {
			req.Status = "active"
		}

		_, err := api.personnelService.Create(req)
		if err != nil {
			failCount++
			failRows = append(failRows, dto.ImportFailRow{
				Row:    headerRowIndex + 2 + idx,
				Data:   rowMap["name"],
				Reason: err.Error(),
			})
		} else {
			success++
		}
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: dto.ImportResultResp{
		Total:   len(dataRows),
		Success: success,
		Failed:  failCount,
		FailRows: failRows,
	}})
}

// parseExcelDate 解析Excel日期单元格
func parseExcelDate(val string) string {
	if val == "" {
		return ""
	}
	// 尝试作为数字序列号解析
	if num, err := strconv.ParseFloat(val, 64); err == nil {
		if t, err := excelize.ExcelDateToTime(num, false); err == nil {
			return t.Format("2006-01-02")
		}
	}
	// 已经是日期字符串，直接返回
	return val
}
