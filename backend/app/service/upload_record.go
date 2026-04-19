package service

import (
	"crypto/rand"
	"datauptwo/app/dto"
	"datauptwo/app/model"
	"datauptwo/app/repo"
	"datauptwo/global"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type UploadRecordService struct {
	uploadRecordRepo *repo.UploadRecordRepo
	projectRepo     *repo.ProjectRepo
}

func NewUploadRecordService() *UploadRecordService {
	return &UploadRecordService{
		uploadRecordRepo: repo.NewUploadRecordRepo(),
		projectRepo:     repo.NewProjectRepo(),
	}
}

// ensureProjectExists 校验项目是否存在（项目必须先创建，再录入数据）
// 返回项目ID（仅当 projectName 非空时有效）和 error
func (s *UploadRecordService) ensureProjectExists(projectName string) (*uint, error) {
	if projectName == "" {
		return nil, nil // 空项目名允许（表示未归类）
	}
	project, err := s.projectRepo.GetByName(projectName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("项目「" + projectName + "」不存在，请先在项目管理中创建该项目")
		}
		return nil, err
	}
	return &project.ID, nil
}

// Create 创建上传记录
func (s *UploadRecordService) Create(req dto.UploadRecordCreateReq) (*model.UploadRecord, error) {
	// 生成流水号
	serialNo := s.generateSerialNo(req.DiskLabel)

	// 处理动态字段数据
	dataJSON := ""
	if req.Data != nil && len(req.Data) > 0 {
		if b, err := json.Marshal(req.Data); err == nil {
			dataJSON = string(b)
		}
	}

	// 解析并设置项目ID：优先使用传入的 ProjectID，否则按 ProjectName 查找
	var projectID *uint
	if req.ProjectID != nil && *req.ProjectID > 0 {
		projectID = req.ProjectID
	} else if req.ProjectName != "" {
		// 按名称查找项目，校验项目是否存在
		pid, err := s.ensureProjectExists(req.ProjectName)
		if err != nil {
			return nil, errors.New("项目「" + req.ProjectName + "」不存在，请先在项目管理中创建该项目")
		}
		projectID = pid
	}

	record := &model.UploadRecord{
		SerialNo:    serialNo,
		DiskLabel:   req.DiskLabel,
		ProjectID:   projectID,
		ProjectName: req.ProjectName,
		DestPath:    req.DestPath,
		FileSize:    int64(math.Round(req.FileSize)), // float64 四舍五入后转 int64
		Uploader:    req.Uploader,
		Status:      req.Status,
		Remark:      req.Remark,
		Data:        dataJSON,
	}

	// 如果传入了创建时间（批量导入时可选），使用该时间
	if req.CreatedAt != "" {
		// 尝试多种日期格式
		for _, format := range []string{"2006-01-02", "2006-01-02 15:04:05", "2006-01-02T15:04:05Z07:00"} {
			if t, err := time.Parse(format, req.CreatedAt); err == nil {
				record.CreatedAt = t
				record.UpdatedAt = t
				break
			}
		}
	}

	// 默认状态
	if record.Status == "" {
		record.Status = "pending"
	}

	if err := s.uploadRecordRepo.Create(record); err != nil {
		return nil, err
	}

	return record, nil
}

// List 分页列表
func (s *UploadRecordService) List(req dto.UploadRecordListReq) (*dto.PageResult, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 20
	}

	records, total, err := s.uploadRecordRepo.List(req)
	if err != nil {
		return nil, err
	}

	items := make([]dto.UploadRecordResp, len(records))
	for i, record := range records {
		items[i] = s.toRecordResp(&record)
	}

	return &dto.PageResult{Total: total, Items: items}, nil
}

// GetByID 获取详情
func (s *UploadRecordService) GetByID(id uint) (*dto.UploadRecordResp, error) {
	record, err := s.uploadRecordRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("记录不存在")
	}

	resp := s.toRecordResp(record)
	return &resp, nil
}

// Update 更新记录
func (s *UploadRecordService) Update(req dto.UploadRecordUpdateReq) (*model.UploadRecord, error) {
	record, err := s.uploadRecordRepo.GetByID(req.ID)
	if err != nil {
		return nil, errors.New("记录不存在")
	}

	record.Status = req.Status
	record.Remark = req.Remark

	// 上传完成后可补充文件大小
	if req.FileSize != nil {
		record.FileSize = *req.FileSize
	}

	// 处理动态字段数据
	if req.Data != nil {
		if b, err := json.Marshal(req.Data); err == nil {
			record.Data = string(b)
		}
	}

	if err := s.uploadRecordRepo.Update(record); err != nil {
		return nil, err
	}

	return record, nil
}

// UpdateBySerialNo 根据流水号更新上传记录（公开接口，无需认证）
func (s *UploadRecordService) UpdateBySerialNo(serialNo string, req dto.UploadRecordPublicUpdateReq) (*model.UploadRecord, error) {
	record, err := s.uploadRecordRepo.GetBySerialNo(serialNo)
	if err != nil {
		return nil, errors.New("记录不存在")
	}

	record.Status = req.Status
	record.Remark = req.Remark

	// 可选：更新文件大小
	if req.FileSize != nil {
		record.FileSize = *req.FileSize
	}

	if err := s.uploadRecordRepo.Update(record); err != nil {
		return nil, err
	}

	return record, nil
}

// Delete 删除记录
func (s *UploadRecordService) Delete(id uint) error {
	return s.uploadRecordRepo.Delete(id)
}

// BatchDelete 批量删除记录
func (s *UploadRecordService) BatchDelete(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return s.uploadRecordRepo.BatchDelete(ids)
}

// BatchUpdateStatus 批量更新记录状态
func (s *UploadRecordService) BatchUpdateStatus(ids []uint, status string) error {
	if len(ids) == 0 {
		return nil
	}
	return s.uploadRecordRepo.BatchUpdateStatus(ids, status)
}

// GetStatistics 获取统计数据
func (s *UploadRecordService) GetStatistics(startDate, endDate, projectName, diskLabel, status, uploader string) (*dto.UploadRecordStatisticsResp, error) {
	now := time.Now()
	today := now.Format("2006-01-02")

	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	weekStart := now.AddDate(0, 0, -weekday+1).Format("2006-01-02")
	monthStart := fmt.Sprintf("%d-%02d-01", now.Year(), now.Month())


	if startDate == "" {
		startDate = now.AddDate(0, 0, -29).Format("2006-01-02")
	}
	if endDate == "" {
		endDate = today
	}

	// 构建筛选条件
	filter := repo.StatisticsFilter{
		ProjectName: projectName,
		DiskLabel:   diskLabel,
		Status:      status,
		Uploader:    uploader,
	}
	filterNoDate := repo.StatisticsFilter{
		ProjectName: projectName,
		DiskLabel:   diskLabel,
		Status:      status,
		Uploader:    uploader,
	}

	resp := &dto.UploadRecordStatisticsResp{}

	resp.TodayCount, _ = s.uploadRecordRepo.CountByDate(today, filter)
	resp.TodaySize, _ = s.uploadRecordRepo.SumFileSizeByDate(today, filter)
	resp.TodaySizeStr = formatFileSize(resp.TodaySize)

	resp.WeekCount, _ = s.uploadRecordRepo.CountByDateRange(weekStart, today, filter)
	resp.WeekSize, _ = s.uploadRecordRepo.SumFileSizeByDateRange(weekStart, today, filter)
	resp.WeekSizeStr = formatFileSize(resp.WeekSize)

	resp.MonthCount, _ = s.uploadRecordRepo.CountByDateRange(monthStart, today, filter)
	resp.MonthSize, _ = s.uploadRecordRepo.SumFileSizeByDateRange(monthStart, today, filter)
	resp.MonthSizeStr = formatFileSize(resp.MonthSize)

	totalCount, _ := s.uploadRecordRepo.CountTotal(filterNoDate)
	totalSize, _ := s.uploadRecordRepo.SumFileSizeByDateRange("1970-01-01", today, filterNoDate)
	resp.TotalCount = totalCount
	resp.TotalSize = totalSize
	resp.TotalSizeStr = formatFileSize(totalSize)

	// trend 查询应用日期范围 + 其他筛选条件
	trend, _ := s.uploadRecordRepo.GetTrend(startDate, endDate, filter)
	resp.Trend = trend

	// byStatus / byDiskLabel / byProject 应用全部筛选条件（含日期范围）
	filterWithDate := repo.StatisticsFilter{
		ProjectName: projectName,
		DiskLabel:   diskLabel,
		StartDate:   startDate,
		EndDate:    endDate,
		Status:      status,
		Uploader:    uploader,
	}
	byStatus, _ := s.uploadRecordRepo.CountByStatus(filterWithDate)
	resp.ByStatus = byStatus

	byDiskLabel, _ := s.uploadRecordRepo.CountByDiskLabel(filterWithDate)
	resp.ByDiskLabel = byDiskLabel

	byProject, _ := s.uploadRecordRepo.CountByProject(filterWithDate)
	resp.ByProject = byProject

	return resp, nil
}

// GetDiskLabelStatuses 获取所有磁盘标签及其综合状态（支持日期范围筛选）
func (s *UploadRecordService) GetDiskLabelStatuses(startDate, endDate string) ([]dto.DiskLabelStatus, error) {
	// 获取每个标签的总体记录数（带日期筛选）
	labels, err := s.uploadRecordRepo.GetDiskLabelStatusAll(startDate, endDate)
	if err != nil {
		return nil, err
	}
	// 获取每个标签下各状态的详细数量（带日期筛选）
	detail, err := s.uploadRecordRepo.GetDiskLabelStatusDetail(startDate, endDate)
	if err != nil {
		return nil, err
	}
	// 计算综合状态
	for i := range labels {
		d := detail[labels[i].DiskLabel]
		if d == nil {
			labels[i].Status = "pending"
			continue
		}
		total := d["pending"] + d["processing"] + d["completed"] + d["failed"]
		if total == 0 {
			labels[i].Status = "pending"
			continue
		}
		// 全部完成 → completed
		if d["failed"] > 0 {
			// 有失败
			if d["completed"] > 0 {
				labels[i].Status = "mixed" // 混合状态
			} else {
				labels[i].Status = "failed"
			}
		} else if d["completed"] == total {
			labels[i].Status = "completed"
		} else {
			labels[i].Status = "pending" // 处理中或待处理
		}
	}
	return labels, nil
}

// GetRecent 获取最近上传记录
func (s *UploadRecordService) GetRecent(limit int, projectName, diskLabel, status, uploader string) ([]dto.UploadRecordResp, error) {
	if limit < 1 {
		limit = 10
	}

	filter := repo.StatisticsFilter{
		ProjectName: projectName,
		DiskLabel:   diskLabel,
		Status:      status,
		Uploader:    uploader,
	}
	records, err := s.uploadRecordRepo.GetRecentRecords(limit, filter)
	if err != nil {
		return nil, err
	}

	items := make([]dto.UploadRecordResp, len(records))
	for i, record := range records {
		items[i] = s.toRecordResp(&record)
	}

	return items, nil
}

// GetUploaderList 获取所有上传者
func (s *UploadRecordService) GetUploaderList() ([]string, error) {
	return s.uploadRecordRepo.GetUploaderList()
}

// GetImportTemplate 获取导入模板字段定义
func (s *UploadRecordService) GetImportTemplate() *dto.ImportTemplateResp {
	fields := []dto.ImportTemplateField{
		{Field: "磁盘标签", Code: "diskLabel", Required: true, Type: "text", MaxLength: 64, Example: "原始数据/成果数据"},
		{Field: "项目名称", Code: "projectName", Required: false, Type: "text", MaxLength: 128, Example: "XX项目"},
		{Field: "目标路径", Code: "destPath", Required: true, Type: "text", MaxLength: 512, Example: "/data/output/2025"},
		{Field: "文件大小(字节)", Code: "fileSize", Required: true, Type: "number", MaxLength: 0, Example: "1048576"},
		{Field: "上传人", Code: "uploader", Required: true, Type: "text", MaxLength: 64, Example: "张三"},
		{Field: "上传状态", Code: "status", Required: true, Type: "select", Options: "pending,processing,completed,failed", MaxLength: 0, Example: "pending（待处理）/processing（处理中）/completed（已完成）/failed（失败）"},
		{Field: "创建时间", Code: "createdAt", Required: false, Type: "date", MaxLength: 0, Example: "2025-01-15 10:30:00（文本格式，不可用日期格式）"},
		{Field: "备注", Code: "remark", Required: false, Type: "text", MaxLength: 512, Example: "这是一批测试数据"},
	}
	return &dto.ImportTemplateResp{
		Fields:    fields,
		SheetName: "上传记录导入",
		Title:     "上传记录批量导入模板",
	}
}

// ImportRecords 批量导入上传记录
func (s *UploadRecordService) ImportRecords(rows []map[string]string) *dto.ImportResultResp {
	result := &dto.ImportResultResp{
		Total:    len(rows),
		FailRows: []dto.ImportFailRow{},
	}

	// 收集所有不重复的项目名称，预校验并记录项目ID
	projectNames := make(map[string]bool)
	for _, row := range rows {
		if pn := row["projectName"]; pn != "" {
			projectNames[pn] = true
		}
	}
	// 预校验：收集不存在的项目名称及对应的错误信息
	invalidProjects := make(map[string]string)
	// 预校验：收集所有存在的项目名称对应的 projectID
	validProjectIDs := make(map[string]*uint)
	for projectName := range projectNames {
		projectID, err := s.ensureProjectExists(projectName)
		if err != nil {
			invalidProjects[projectName] = "项目「" + projectName + "」不存在，请先在项目管理中创建该项目"
		} else {
			validProjectIDs[projectName] = projectID
		}
	}
	// 预先将无效项目的行标记为失败（不在此处 return，保证主循环一定执行，Success 计数正确）
	for i, row := range rows {
		rowNum := i + 2
		serialized, _ := json.Marshal(row)
		if reason, ok := invalidProjects[row["projectName"]]; ok {
			result.FailRows = append(result.FailRows, dto.ImportFailRow{Row: rowNum, Data: string(serialized), Reason: reason})
		}
	}

	for i, row := range rows {
		rowNum := i + 2 // Excel row number (1 = header)
		serialized, _ := json.Marshal(row)

		// 跳过已预标记为无效项目的行（避免重复加入 FailRows）
		if _, isInvalid := invalidProjects[row["projectName"]]; isInvalid {
			continue
		}

		// 必填字段校验
		if v := row["diskLabel"]; v == "" {
			result.FailRows = append(result.FailRows, dto.ImportFailRow{Row: rowNum, Data: string(serialized), Reason: "磁盘标签（diskLabel）不能为空"})
			continue
		}
		if v := row["destPath"]; v == "" {
			result.FailRows = append(result.FailRows, dto.ImportFailRow{Row: rowNum, Data: string(serialized), Reason: "目标路径（destPath）不能为空"})
			continue
		}
		if v := row["uploader"]; v == "" {
			result.FailRows = append(result.FailRows, dto.ImportFailRow{Row: rowNum, Data: string(serialized), Reason: "上传人（uploader）不能为空"})
			continue
		}
		if v := row["status"]; v == "" {
			result.FailRows = append(result.FailRows, dto.ImportFailRow{Row: rowNum, Data: string(serialized), Reason: "上传状态（status）不能为空"})
			continue
		}

		// 解析文件大小（支持整数和浮点数，统一四舍五入）
		var fileSize int64
		if v := row["fileSize"]; v != "" {
			// 先尝试解析为 float64，再四舍五入转为 int64
			var parsedFloat float64
			if _, err := fmt.Sscanf(v, "%f", &parsedFloat); err != nil {
				// 尝试解析为整数
				var parsedInt int64
				if _, err := fmt.Sscanf(v, "%d", &parsedInt); err != nil {
					result.FailRows = append(result.FailRows, dto.ImportFailRow{Row: rowNum, Data: string(serialized), Reason: fmt.Sprintf("文件大小（fileSize）格式错误，无法解析为数字: %s", v)})
					continue
				}
				fileSize = parsedInt
			} else {
				fileSize = int64(math.Round(parsedFloat)) // 四舍五入
			}
		}
		if fileSize <= 0 {
			result.FailRows = append(result.FailRows, dto.ImportFailRow{Row: rowNum, Data: string(serialized), Reason: "文件大小（fileSize）必须大于0"})
			continue
		}

		// 状态值校验
		status := row["status"]
		if status != "pending" && status != "processing" && status != "completed" && status != "failed" {
			result.FailRows = append(result.FailRows, dto.ImportFailRow{Row: rowNum, Data: string(serialized), Reason: fmt.Sprintf("上传状态（status）值非法: %s，支持值: pending/processing/completed/failed", status)})
			continue
		}

		// 构建创建请求（传入预查好的 ProjectID）
		req := dto.UploadRecordCreateReq{
			DiskLabel:   row["diskLabel"],
			ProjectID:   validProjectIDs[row["projectName"]], // 预查好的项目ID（可能为 nil）
			ProjectName: row["projectName"],
			DestPath:    row["destPath"],
			FileSize:    float64(fileSize), // int64 转 float64 兼容 DTO
			Uploader:    row["uploader"],
			Status:      status,
			Remark:      row["remark"],
		}

		// 解析可选的创建时间（支持多种日期格式）
		if createdAtStr := row["createdAt"]; createdAtStr != "" {
			formats := []string{
				"2006-01-02",
				"2006-01-02 15:04:05",
				"2006/01/02",
				"2006/01/02 15:04:05",
				"2006/1/2 15:04",
				"2006-1-2 15:04",
				"01/02/2006",
				"1/2/2006",
				"2006/1/2",
				"2006-1-2",
				"1/2/2006 15:04",
				"01/02/2006 15:04:05",
				"2006-01-02T15:04:05Z07:00",
				"2006-01-02T15:04:05",
				"2006-01-02T15:04",
			}
			for _, format := range formats {
				if t, err := time.Parse(format, createdAtStr); err == nil {
					req.CreatedAt = t.Format("2006-01-02 15:04:05")
					break
				}
			}
		}

		if _, err := s.Create(req); err != nil {
			result.FailRows = append(result.FailRows, dto.ImportFailRow{Row: rowNum, Data: string(serialized), Reason: "创建记录失败: " + err.Error()})
			continue
		}

		result.Success++
	}

	result.Failed = len(result.FailRows)
	return result
}

// GetBySerialNo 根据流水号获取记录
func (s *UploadRecordService) GetBySerialNo(serialNo string) (*dto.UploadRecordResp, error) {
	record, err := s.uploadRecordRepo.GetBySerialNo(serialNo)
	if err != nil {
		return nil, errors.New("记录不存在")
	}

	resp := s.toRecordResp(record)
	return &resp, nil
}

// Export 导出上传记录为 Excel
func (s *UploadRecordService) Export(req dto.UploadRecordListReq) ([]byte, error) {
	records, err := s.uploadRecordRepo.ListAllForExport(req)
	if err != nil {
		return nil, err
	}

	// 使用 excelize 生成 Excel
	f := excelize.NewFile()
	defer f.Close()

	// 创建工作表
	sheet := "上传记录"
	index, err := f.NewSheet(sheet)
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(index)
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
	headers := []string{"序号", "流水号", "磁盘标签", "项目名称", "目标路径", "文件大小", "上传人", "状态", "备注", "上传时间", "更新时间"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	// 数据行
	row := 2
	for i, record := range records {
		resp := s.toRecordResp(&record)

		// 序号
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), resp.SerialNo)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), resp.DiskLabel)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), resp.ProjectName)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), resp.DestPath)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), resp.FileSizeStr)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), resp.Uploader)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), resp.StatusText)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", row), resp.Remark)
		f.SetCellValue(sheet, fmt.Sprintf("J%d", row), resp.CreatedAt)
		f.SetCellValue(sheet, fmt.Sprintf("K%d", row), resp.UpdatedAt)

		row++
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
	if len(records) > 0 {
		startCell, _ := excelize.CoordinatesToCellName(1, 2)
		endCell, _ := excelize.CoordinatesToCellName(11, row-1)
		f.SetCellStyle(sheet, startCell, endCell, dataStyle)
	}

	// 设置列宽
	colWidths := []float64{6, 24, 14, 20, 30, 12, 12, 10, 20, 20, 20}
	colLetters := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K"}
	for i, w := range colWidths {
		f.SetColWidth(sheet, colLetters[i], colLetters[i], w)
	}

	// 设置行高
	f.SetRowHeight(sheet, 1, 32)

	// 生成文件
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// generateSerialNo 生成8位UUID流水号
func (s *UploadRecordService) generateSerialNo(diskLabel string) string {
	// 格式：8位小写字母+数字随机字符
	const charset = "0123456789abcdefghijklmnopqrstuvwxyz"
	var rnd [8]byte
	cryptoRandRead(rnd[:])

	var serialNo [8]byte
	for i := 0; i < 8; i++ {
		serialNo[i] = charset[rnd[i]%36]
	}

	return string(serialNo[:])
}

// cryptoRandRead 用crypto/rand填充切片， panic on failure
func cryptoRandRead(b []byte) {
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand read failed: " + err.Error())
	}
}

// toRecordResp 转换为响应格式（时间格式化为 YYYY-MM-DD HH:mm:ss）
func (s *UploadRecordService) toRecordResp(record *model.UploadRecord) dto.UploadRecordResp {
	statusText := ""
	switch record.Status {
	case "pending":
		statusText = "待处理"
	case "processing":
		statusText = "处理中"
	case "completed":
		statusText = "已完成"
	case "failed":
		statusText = "失败"
	}

	// 解析动态字段数据
	var data map[string]interface{}
	if record.Data != "" {
		if err := json.Unmarshal([]byte(record.Data), &data); err != nil {
			global.AppLogger.Error("解析记录动态字段失败, id=%d, serial_no=%s, err=%v", record.ID, record.SerialNo, err)
		}
	}

	return dto.UploadRecordResp{
		ID:          record.ID,
		SerialNo:    record.SerialNo,
		DiskLabel:   record.DiskLabel,
		ProjectID:   record.ProjectID,
		ProjectName: record.ProjectName,
		DestPath:    record.DestPath,
		FileSize:    record.FileSize,
		FileSizeStr: formatFileSize(record.FileSize),
		Uploader:    record.Uploader,
		Status:      record.Status,
		StatusText:  statusText,
		Remark:      record.Remark,
		Data:        data,
		CreatedAt:   record.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   record.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// formatFileSize 格式化文件大小
func formatFileSize(size int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)

	if size >= TB {
		return fmt.Sprintf("%.2f TB", float64(size)/TB)
	}
	if size >= GB {
		return fmt.Sprintf("%.2f GB", float64(size)/GB)
	}
	if size >= MB {
		return fmt.Sprintf("%.2f MB", float64(size)/MB)
	}
	if size >= KB {
		return fmt.Sprintf("%.2f KB", float64(size)/KB)
	}
	return fmt.Sprintf("%d B", size)
}
