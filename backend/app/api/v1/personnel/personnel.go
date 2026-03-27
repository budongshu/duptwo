package api

import (
	"net/http"
	"strconv"
	"time"

	"datauptwo/app/dto"
	"datauptwo/app/service"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type PersonnelApi struct {
	personnelService *service.PersonnelService
}

func NewPersonnelApi() *PersonnelApi {
	return &PersonnelApi{
		personnelService: service.NewPersonnelService(),
	}
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

	if err := api.personnelService.Update(req); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

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

	if err := api.personnelService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

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
