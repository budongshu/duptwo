package user

import (
	"datauptwo/app/dto"
	"datauptwo/app/service"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type UserApi struct {
	userService  *service.UserService
	auditService *service.AuditService
}

func NewUserApi() *UserApi {
	return &UserApi{
		userService:  service.NewUserService(),
		auditService: service.NewAuditService(),
	}
}

// Create 创建用户
// @Summary 创建用户
// @Tags User
// @Security Bearer
// @Accept json
// @Param request body dto.UserCreateReq true "用户信息"
// @Success 200 {object} dto.Response
// @Router /api/users [post]
func (api *UserApi) Create(c *gin.Context) {
	var req dto.UserCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	user, err := api.userService.Create(req)
	if err != nil {
		// 根据错误类型返回合适的HTTP状态码
		errMsg := err.Error()
		if errMsg == "用户名已存在" || errMsg == "检查用户名失败" {
			c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: errMsg})
		} else {
			c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: errMsg})
		}
		return
	}

	// 记录操作日志
	api.auditService.LogOperation(
		api.getUserID(c), api.getUsername(c),
		"用户管理", "create", "User", user.ID,
		user.Username, c.ClientIP(), c.GetHeader("User-Agent"), nil,
	)

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "创建成功", Data: user})
}

// List 用户列表
// @Summary 用户列表
// @Tags User
// @Security Bearer
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param keyword query string false "搜索关键词"
// @Param status query string false "状态"
// @Param roleId query int false "角色ID"
// @Param groupId query int false "用户组ID"
// @Success 200 {object} dto.Response
// @Router /api/users [get]
func (api *UserApi) List(c *gin.Context) {
	var req dto.UserListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		req.Page = 1
		req.PageSize = 20
	}

	result, err := api.userService.List(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	// 记录操作日志
	api.auditService.LogOperation(
		api.getUserID(c), api.getUsername(c),
		"用户管理", "view", "User", 0,
		"查看用户列表", c.ClientIP(), c.GetHeader("User-Agent"), nil,
	)

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: result})
}

// GetAll 获取所有用户
// @Summary 获取所有用户
// @Tags User
// @Security Bearer
// @Success 200 {object} dto.Response
// @Router /api/users/all [get]
func (api *UserApi) GetAll(c *gin.Context) {
	users, err := api.userService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: users})
}

// GetByID 获取用户详情
// @Summary 获取用户详情
// @Tags User
// @Security Bearer
// @Param id path int true "用户ID"
// @Success 200 {object} dto.Response
// @Router /api/users/{id} [get]
func (api *UserApi) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	user, err := api.userService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{Code: 404, Message: "用户不存在"})
		return
	}

	// 记录操作日志
	api.auditService.LogOperation(
		api.getUserID(c), api.getUsername(c),
		"用户管理", "view", "User", uint(id),
		user.Username, c.ClientIP(), c.GetHeader("User-Agent"), nil,
	)

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: user})
}

// Update 更新用户
// @Summary 更新用户
// @Tags User
// @Security Bearer
// @Accept json
// @Param request body dto.UserUpdateReq true "用户信息"
// @Success 200 {object} dto.Response
// @Router /api/users [put]
func (api *UserApi) Update(c *gin.Context) {
	var req dto.UserUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	// 先获取变更前的数据
	oldUser, _ := api.userService.GetByID(req.ID)

	user, err := api.userService.Update(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	// 记录操作日志（带变更明细）
	detail := &service.OperationDetail{
		After: user,
	}
	if oldUser != nil {
		detail.Before = oldUser
		// 定义字段中文标签
		fieldLabels := map[string]string{
			"username":  "用户名",
			"nickname":  "昵称",
			"email":     "邮箱",
			"phone":     "手机号",
			"status":    "状态",
			"roleId":    "角色ID",
			"groupId":   "用户组ID",
			"roleName":  "角色名称",
			"groupName": "用户组名称",
		}
		detail.Changes = service.GetChanges(*oldUser, *user, fieldLabels)
	}
	api.auditService.LogOperationWithDetail(
		api.getUserID(c), api.getUsername(c),
		"用户管理", "update", "User", req.ID,
		user.Username, c.ClientIP(), c.GetHeader("User-Agent"), detail,
	)

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "更新成功", Data: user})
}

// Delete 删除用户
// @Summary 删除用户
// @Tags User
// @Security Bearer
// @Param id path int true "用户ID"
// @Success 200 {object} dto.Response
// @Router /api/users/{id} [delete]
func (api *UserApi) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	// 先获取用户信息用于日志
	user, _ := api.userService.GetByID(uint(id))
	username := ""
	if user != nil {
		username = user.Username
	}

	if err := api.userService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	// 记录操作日志
	api.auditService.LogOperation(
		api.getUserID(c), api.getUsername(c),
		"用户管理", "delete", "User", uint(id),
		username, c.ClientIP(), c.GetHeader("User-Agent"), nil,
	)

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "删除成功"})
}

// BatchDelete 批量删除用户
// @Summary 批量删除用户
// @Tags User
// @Security Bearer
// @Accept json
// @Param request body dto.BatchDeleteReq true "用户ID列表"
// @Success 200 {object} dto.Response
// @Router /api/users/batch-delete [post]
func (api *UserApi) BatchDelete(c *gin.Context) {
	var req dto.BatchDeleteReq
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "请选择要删除的用户"})
		return
	}

	if err := api.userService.BatchDelete(req.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: fmt.Sprintf("成功删除 %d 个用户", len(req.IDs))})
}

// BatchUpdateRole 批量更新用户角色
// @Summary 批量更新用户角色
// @Tags User
// @Security Bearer
// @Accept json
// @Param request body dto.BatchUpdateRoleReq true "批量更新角色信息"
// @Success 200 {object} dto.Response
// @Router /api/users/batch-update-role [post]
func (api *UserApi) BatchUpdateRole(c *gin.Context) {
	var req dto.BatchUpdateRoleReq
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "请选择要更新的用户"})
		return
	}
	if req.RoleID == 0 {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "请选择要分配的角色"})
		return
	}

	if err := api.userService.BatchUpdateRole(req.IDs, req.RoleID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	// 记录操作日志
	api.auditService.LogOperation(
		api.getUserID(c), api.getUsername(c),
		"用户管理", "batch_update_role", "User", 0,
		fmt.Sprintf("批量更新 %d 个用户的角色", len(req.IDs)), c.ClientIP(), c.GetHeader("User-Agent"), nil,
	)

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: fmt.Sprintf("成功更新 %d 个用户的角色", len(req.IDs))})
}

func (api *UserApi) getUserID(c *gin.Context) uint {
	if id, exists := c.Get("userId"); exists {
		return id.(uint)
	}
	return 0
}

// getUsername 获取当前用户名
func (api *UserApi) getUsername(c *gin.Context) string {
	if username, exists := c.Get("username"); exists {
		return username.(string)
	}
	return ""
}

// ResetPassword 重置用户密码
// @Summary 重置用户密码
// @Tags User
// @Security Bearer
// @Accept json
// @Param request body dto.ResetPasswordReq true "重置密码信息"
// @Success 200 {object} dto.Response
// @Router /api/users/reset-password [post]
func (api *UserApi) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	err := api.userService.ResetPassword(req.UserID, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	// 记录操作日志
	api.auditService.LogOperation(
		api.getUserID(c), api.getUsername(c),
		"用户管理", "reset_password", "User", req.UserID,
		"", c.ClientIP(), c.GetHeader("User-Agent"), nil,
	)

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "密码重置成功"})
}

// ResetMFA 重置用户MFA
// @Summary 重置用户MFA
// @Tags User
// @Security Bearer
// @Accept json
// @Param request body dto.ResetMFAReq true "重置MFA信息"
// @Success 200 {object} dto.Response
// @Router /api/users/reset-mfa [post]
func (api *UserApi) ResetMFA(c *gin.Context) {
	var req dto.ResetMFAReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	err := api.userService.ResetMFA(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	// 记录操作日志
	api.auditService.LogOperation(
		api.getUserID(c), api.getUsername(c),
		"用户管理", "reset_mfa", "User", req.UserID,
		"", c.ClientIP(), c.GetHeader("User-Agent"), nil,
	)

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "MFA重置成功"})
}

// GenerateMFASecret 生成MFA密钥
// @Summary 生成MFA密钥
// @Tags User
// @Security Bearer
// @Param id query int true "用户ID"
// @Success 200 {object} dto.Response
// @Router /api/users/generate-mfa-secret [get]
func (api *UserApi) GenerateMFASecret(c *gin.Context) {
	userIDStr := c.Query("id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "用户ID不能为空"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "无效的用户ID"})
		return
	}

	secret, qrCode, err := api.userService.GenerateMFASecret(uint(userID))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: dto.GenerateMFASecretResp{
		Secret: secret,
		QRCode: qrCode,
	}})
}

// AdminEnableMFA 管理员为用户启用MFA
// @Summary 管理员为用户启用MFA
// @Tags User
// @Security Bearer
// @Accept json
// @Param request body dto.AdminEnableMFAReq true "启用MFA信息"
// @Success 200 {object} dto.Response
// @Router /api/users/admin-enable-mfa [post]
func (api *UserApi) AdminEnableMFA(c *gin.Context) {
	var req dto.AdminEnableMFAReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	err := api.userService.AdminEnableMFA(req.UserID, req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	// 记录操作日志
	api.auditService.LogOperation(
		api.getUserID(c), api.getUsername(c),
		"用户管理", "enable_mfa", "User", req.UserID,
		"", c.ClientIP(), c.GetHeader("User-Agent"), nil,
	)

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "MFA启用成功"})
}

// Export 导出用户Excel
// @Summary 导出用户Excel
// @Tags User
// @Security Bearer
// @Param keyword query string false "关键词"
// @Param status query string false "状态"
// @Success 200 {file} file
// @Router /api/users/export [get]
func (api *UserApi) Export(c *gin.Context) {
	var req dto.UserListReq
	c.ShouldBindQuery(&req)

	list, err := api.userService.ListForExport(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: err.Error()})
		return
	}

	f := excelize.NewFile()
	defer f.Close()

	sheet := "用户列表"
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

	headers := []string{"序号", "用户名", "昵称", "邮箱", "手机号", "角色", "用户组", "状态", "MFA", "创建时间"}
	headerWidths := map[string]float64{"A": 6, "B": 14, "C": 12, "D": 22, "E": 14, "F": 12, "G": 12, "H": 8, "I": 8, "J": 18}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}
	for col, w := range headerWidths {
		f.SetColWidth(sheet, col, col, w)
	}

	// 获取角色和用户组映射
	roleMap, groupMap := api.userService.GetRoleAndGroupMaps()

	statusMap := map[string]string{"active": "启用", "inactive": "禁用", "": "启用"}
	mfaMap := map[bool]string{true: "已启用", false: "未启用"}

	for idx, u := range list {
		row := idx + 2
		rowData := []interface{}{
			row,
			u.Username,
			u.Nickname,
			u.Email,
			u.Phone,
			roleMap[u.RoleID],
			groupMap[u.GroupID],
			statusMap[u.Status],
			mfaMap[u.MFAEnabled],
			u.CreatedAt.Format("2006-01-02 15:04:05"),
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
	fileName := "用户列表_" + time.Now().Format("20060102150405") + ".xlsx"
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
// @Tags User
// @Security Bearer
// @Success 200 {object} dto.Response
// @Router /api/users/template [get]
func (api *UserApi) GetTemplate(c *gin.Context) {
	fields := []dto.ImportTemplateField{
		{Field: "用户名", Code: "username", Required: true, Type: "string", MaxLength: 64, Example: "zhangsan"},
		{Field: "密码", Code: "password", Required: true, Type: "password", MaxLength: 64, Example: "Abc123@!"},
		{Field: "昵称", Code: "nickname", Required: false, Type: "string", MaxLength: 64, Example: "张三"},
		{Field: "邮箱", Code: "email", Required: false, Type: "email", MaxLength: 128, Example: "zhangsan@example.com"},
		{Field: "手机号", Code: "phone", Required: false, Type: "string", MaxLength: 32, Example: "13800138000"},
		{Field: "状态", Code: "status", Required: false, Type: "select", Options: "active:启用,inactive:禁用", Example: "active"},
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: dto.ImportTemplateResp{Fields: fields, SheetName: "用户列表", Title: "用户导入模板"}})
}

// Import 批量导入用户
// @Summary 批量导入用户
// @Tags User
// @Security Bearer
// @Accept multipart/form-data
// @Param file formData file true "Excel文件"
// @Success 200 {object} dto.Response
// @Router /api/users/import [post]
func (api *UserApi) Import(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "请上传文件"})
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
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "无法读取Excel文件"})
		return
	}
	defer f.Close()

	rows, err := f.GetRows("用户列表")
	if err != nil {
		rows, err = f.GetRows("Sheet1")
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "Excel中没有数据"})
			return
		}
	}

	if len(rows) < 2 {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "Excel中没有数据行"})
		return
	}

	// 动态查找表头行（跳过标题行和说明行）
	headerRowIndex := -1
	var headerRow []string
	for i, row := range rows {
		for _, cell := range row {
			cellClean := strings.TrimSpace(cell)
			if cellClean == "用户名" || strings.ReplaceAll(cellClean, " ", "") == "用户名*" {
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
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "未找到有效的表头行（需包含「用户名」列）"})
		return
	}

	// 建立列索引映射（去掉末尾的 * 标记）
	headerMap := make(map[string]int)
	for idx, colName := range headerRow {
		headerMap[strings.TrimSuffix(strings.TrimSpace(colName), " *")] = idx
	}

	// 检查必需列
	requiredCols := []string{"用户名", "密码"}
	for _, col := range requiredCols {
		if _, ok := headerMap[col]; !ok {
			c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: fmt.Sprintf("缺少必需列: %s", col)})
			return
		}
	}

	// 获取角色和用户组映射
	roleMap, groupMap := api.userService.GetRoleAndGroupMaps()
	roleNameToID := make(map[string]uint)
	for id, name := range roleMap {
		roleNameToID[name] = id
	}
	groupNameToID := make(map[string]uint)
	for id, name := range groupMap {
		groupNameToID[name] = id
	}

	var success, failed int
	var failRows []dto.ImportFailRow

	for i := headerRowIndex + 1; i < len(rows); i++ {
		row := rows[i]
		rowNum := i + 1
		getVal := func(key string) string {
			if col, ok := headerMap[key]; ok && col < len(row) {
				return strings.TrimSpace(row[col])
			}
			return ""
		}

		username := getVal("用户名")
		password := getVal("密码")
		nickname := getVal("昵称")
		email := getVal("邮箱")
		phone := getVal("手机号")
		status := getVal("状态")
		roleName := getVal("角色")
		groupName := getVal("用户组")

		if username == "" || password == "" {
			failed++
			failRows = append(failRows, dto.ImportFailRow{Row: rowNum, Data: username, Reason: "用户名和密码不能为空"})
			continue
		}

		if status == "" || status == "启用" {
			status = "active"
		} else {
			status = "inactive"
		}

		var roleID uint
		if roleName != "" {
			roleID = roleNameToID[roleName]
		}
		var groupID uint
		if groupName != "" {
			groupID = groupNameToID[groupName]
		}

		req := dto.UserCreateReq{
			Username: username,
			Password: password,
			Nickname: nickname,
			Email:    email,
			Phone:    phone,
			Status:   status,
			RoleID:   roleID,
			GroupID:  groupID,
		}

		_, err := api.userService.Create(req)
		if err != nil {
			failed++
			failRows = append(failRows, dto.ImportFailRow{Row: rowNum, Data: username, Reason: err.Error()})
		} else {
			success++
		}
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: dto.ImportResultResp{
		Total:   success + failed,
		Success: success,
		Failed:  failed,
		FailRows: failRows,
	}})
}

// Preview 预览导入数据（返回识别的记录数）
func (api *UserApi) Preview(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "请上传文件"})
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
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "无法读取Excel文件"})
		return
	}
	defer f.Close()

	rows, err := f.GetRows("用户列表")
	if err != nil {
		rows, err = f.GetRows("Sheet1")
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "Excel中没有数据"})
			return
		}
	}

	headerRowIndex := -1
	var headerRow []string
	for i, row := range rows {
		for _, cell := range row {
			cellClean := strings.TrimSpace(cell)
			if cellClean == "用户名" || strings.ReplaceAll(cellClean, " ", "") == "用户名*" {
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
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "未找到有效的表头行（需包含「用户名」列）"})
		return
	}

	headerMap := make(map[string]int)
	for idx, colName := range headerRow {
		headerMap[strings.TrimSuffix(strings.TrimSpace(colName), " *")] = idx
	}

	usernameCol, hasUsername := headerMap["用户名"]
	dataRowCount := 0
	for i := headerRowIndex + 1; i < len(rows); i++ {
		if hasUsername && usernameCol < len(rows[i]) && strings.TrimSpace(rows[i][usernameCol]) != "" {
			dataRowCount++
		}
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Data: map[string]interface{}{
		"total":     dataRowCount,
		"sheetName": "用户列表",
		"fields":    headerRow,
	}})
}
