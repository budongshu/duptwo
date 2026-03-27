# ============================================================
# duptwo 客户端上传示例
# 支持 Bash/cURL、Python、Go 等多种语言
# ============================================================

# ========== 配置区 ==========
API_BASE="${DATA_REGISTRY_API:-http://localhost:8080}"
SERIAL_NO=""

# ============================================================
# 1. Bash / cURL 上传
# ============================================================

# --- 1.1 公开接口上传（无需认证）---
upload_bash() {
    local data_type="$1"
    local project_name="${2:-}"
    local dest_path="$3"
    local file_size="$4"
    local uploader="$5"
    local remark="${6:-}"

    curl -s -X POST "${API_BASE}/public/upload-records" \
        -H "Content-Type: application/json" \
        -d "{
            \"dataType\": \"${data_type}\",
            \"projectName\": \"${project_name}\",
            \"destPath\": \"${dest_path}\",
            \"fileSize\": ${file_size},
            \"uploader\": \"${uploader}\",
            \"status\": \"completed\",
            \"remark\": \"${remark}\"
        }" | python3 -m json.tool
}

# --- 1.2 认证后上传（Bearer Token）---
upload_auth_bash() {
    local token="$1"
    local data_type="$2"
    local project_name="${3:-}"
    local file_path="$4"
    local file_size="$5"
    local uploader="$6"

    curl -s -X POST "${API_BASE}/api/upload-records" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer ${token}" \
        -d "{
            \"dataType\": \"${data_type}\",
            \"projectName\": \"${project_name}\",
            \"destPath\": \"${file_path}\",
            \"fileSize\": ${file_size},
            \"uploader\": \"${uploader}\",
            \"status\": \"pending\"
        }" | python3 -m json.tool
}

# --- 1.3 带动态字段上传 ---
upload_with_fields_bash() {
    local data_type="$1"
    local file_path="$2"
    local uploader="$3"

    # 动态字段通过 data 传入
    curl -s -X POST "${API_BASE}/public/upload-records" \
        -H "Content-Type: application/json" \
        -d "{
            \"dataType\": \"${data_type}\",
            \"destPath\": \"${file_path}\",
            \"fileSize\": 104857600,
            \"uploader\": \"${uploader}\",
            \"status\": \"completed\",
            \"data\": {
                \"patientId\": \"P20240001\",
                \"modality\": \"CT\",
                \"studyDate\": \"2024-03-15\",
                \"department\": \"放射科\"
            }
        }" | python3 -m json.tool
}

# --- 1.4 登录获取 Token ---
login_bash() {
    local username="${1:-admin}"
    local password="${2:-admin123}"

    curl -s -X POST "${API_BASE}/api/auth/login" \
        -H "Content-Type: application/json" \
        -d "{
            \"username\": \"${username}\",
            \"password\": \"${password}\"
        }" | python3 -m json.tool
}

# --- 1.5 根据流水号查询 ---
get_by_serial_bash() {
    local serial_no="$1"

    curl -s "${API_BASE}/public/upload-records/${serial_no}" | python3 -m json.tool
}

# --- 1.6 批量上传（循环）---
batch_upload_bash() {
    local data_type="$1"
    local uploader="$2"
    local count="${3:-10}"

    for i in $(seq 1 "$count"); do
        echo "上传第 ${i} 条..."
        result=$(curl -s -X POST "${API_BASE}/public/upload-records" \
            -H "Content-Type: application/json" \
            -d "{
                \"dataType\": \"${data_type}\",
                \"destPath\": \"/data/batch/file_${i}.dat\",
                \"fileSize\": $((RANDOM * 1024)),
                \"uploader\": \"${uploader}\",
                \"status\": \"completed\",
                \"remark\": \"批量上传第 ${i} 条\"
            }")
        echo "$result" | python3 -c "import sys,json; d=json.load(sys.stdin); print(f'  流水号: {d[\"data\"][\"serialNo\"]}')"
    done
    echo "批量上传完成！"
}


# ============================================================
# 2. Python 上传示例
# ============================================================

cat << 'PYTHON_EOF'
# --- 2.1 基础上传 ---
import requests
import json

API_BASE = "http://localhost:8080"

def upload_record(data_type, file_path, file_size, uploader, project_name="", remark=""):
    """上传记录"""
    url = f"{API_BASE}/public/upload-records"
    payload = {
        "dataType": data_type,
        "projectName": project_name,
        "destPath": file_path,
        "fileSize": file_size,
        "uploader": uploader,
        "status": "completed",
        "remark": remark,
    }
    resp = requests.post(url, json=payload, timeout=10)
    resp.raise_for_status()
    return resp.json()

def upload_with_dynamic_fields(data_type, file_path, file_size, uploader, fields: dict):
    """上传记录（带动态字段）"""
    url = f"{API_BASE}/public/upload-records"
    payload = {
        "dataType": data_type,
        "destPath": file_path,
        "fileSize": file_size,
        "uploader": uploader,
        "status": "completed",
        "data": fields,  # 动态字段
    }
    resp = requests.post(url, json=payload, timeout=10)
    resp.raise_for_status()
    return resp.json()

def query_by_serial(serial_no):
    """根据流水号查询"""
    url = f"{API_BASE}/public/upload-records/{serial_no}"
    resp = requests.get(url, timeout=10)
    resp.raise_for_status()
    return resp.json()

# --- 示例调用 ---
if __name__ == "__main__":
    # 上传基础记录
    result = upload_record(
        data_type="SSD-001",
        file_path="/data/medical/ct_001.dcm",
        file_size=104857600,
        uploader="张三",
        project_name="医疗影像项目",
        remark="CT扫描数据"
    )
    serial_no = result["data"]["serialNo"]
    print(f"上传成功！流水号: {serial_no}")

    # 上传带动态字段
    result = upload_with_dynamic_fields(
        data_type="SSD-001",
        file_path="/data/medical/mri_001.dcm",
        file_size=209715200,
        uploader="李四",
        fields={
            "patientId": "P20240002",
            "modality": "MRI",
            "studyDate": "2024-03-16",
            "radiologist": "王医生",
        }
    )
    print(f"带字段上传成功！流水号: {result['data']['serialNo']}")

    # 查询记录
    record = query_by_serial(serial_no)
    print(f"查询结果: {record['data']['destPath']}")


# --- 2.2 认证上传 ---
def login(username, password):
    url = f"{API_BASE}/api/auth/login"
    resp = requests.post(url, json={"username": username, "password": password}, timeout=10)
    resp.raise_for_status()
    data = resp.json()
    return data["data"]["token"]

def upload_auth(token, **kwargs):
    url = f"{API_BASE}/api/upload-records"
    headers = {"Authorization": f"Bearer {token}"}
    resp = requests.post(url, json=kwargs, headers=headers, timeout=10)
    resp.raise_for_status()
    return resp.json()

# 登录获取 token
token = login("admin", "admin123")

# 认证上传
result = upload_auth(
    token,
    dataType="SSD-002",
    destPath="/data/financial/report.xlsx",
    fileSize=5242880,
    uploader="财务部",
    status="pending",
    data={"department": "财务", "year": 2024, "quarter": "Q1"}
)
print(f"认证上传成功！流水号: {result['data']['serialNo']}")


# --- 2.3 批量上传（带进度） ---
import time

def batch_upload(data_type, uploader, count=100):
    """批量上传并显示进度"""
    success = 0
    failed = 0
    results = []

    for i in range(1, count + 1):
        try:
            result = upload_record(
                data_type=data_type,
                file_path=f"/data/batch/file_{i:04d}.dat",
                file_size=1024 * (i % 100 + 1),
                uploader=uploader,
                remark=f"批量第{i}条"
            )
            serial_no = result["data"]["serialNo"]
            results.append(serial_no)
            success += 1
            print(f"\r进度: {i}/{count}  成功: {success}  失败: {failed}  最新: {serial_no}", end="")
        except Exception as e:
            failed += 1
            print(f"\r进度: {i}/{count}  成功: {success}  失败: {failed}  错误: {e}", end="")

        if i % 10 == 0:
            time.sleep(0.1)  # 避免请求过快

    print()  # 换行
    return results

# 执行批量上传（100条）
print("开始批量上传...")
serials = batch_upload("SSD-001", "系统管理员", count=100)
print(f"\n批量上传完成！成功: {len(serials)} 条")
PYTHON_EOF


# ============================================================
# 3. Go 上传示例
# ============================================================

cat << 'GO_EOF'
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const APIBase = "http://localhost:8080"

type UploadRecordReq struct {
	DataType    string                 `json:"dataType"`
	ProjectName string                `json:"projectName,omitempty"`
	FilePath    string                 `json:"destPath"`
	FileSize    int64                  `json:"fileSize"`
	Uploader    string                 `json:"uploader"`
	Status      string                 `json:"status"`
	Remark      string                 `json:"remark,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty"`
}

type APIResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func uploadRecord(req UploadRecordReq) (string, error) {
	body, _ := json.Marshal(req)
	resp, err := http.Post(APIBase+"/public/upload-records",
		"application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	var result APIResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return "", err
	}

	// 解析流水号
	var record map[string]interface{}
	json.Unmarshal(result.Data, &record)
	return record["serialNo"].(string), nil
}

func main() {
	// 基础上传
	serial, err := uploadRecord(UploadRecordReq{
		DataType:    "SSD-001",
		ProjectName: "测试项目",
		FilePath:    "/data/test.csv",
		FileSize:    1024000,
		Uploader:    "系统",
		Status:      "completed",
		Remark:      "Go SDK 测试",
	})
	if err != nil {
		fmt.Printf("上传失败: %v\n", err)
		return
	}
	fmt.Printf("上传成功！流水号: %s\n", serial)

	// 带动态字段
	serial, err = uploadRecord(UploadRecordReq{
		DataType: "SSD-001",
		FilePath: "/data/test2.csv",
		FileSize: 2048000,
		Uploader: "系统",
		Status:   "completed",
		Data: map[string]interface{}{
			"patientId":  "P001",
			"modality":    "CT",
			"studyDate":   "2024-03-15",
			"radiologist": "医生A",
		},
	})
	fmt.Printf("动态字段上传: %s (err=%v)\n", serial, err)
}
GO_EOF


# ============================================================
# 4. 使用示例（直接在终端运行）
# ============================================================

# --- 4.1 一行命令上传 ---
# 公开接口上传（无需认证）
curl -X POST "${API_BASE}/public/upload-records" \
  -H "Content-Type: application/json" \
  -d '{"dataType":"SSD-001","destPath":"/data/test.csv","fileSize":1024000,"uploader":"测试用户","status":"completed"}'

# --- 4.2 完整示例：登录 + 上传 + 查询 ---
echo "=== duptwo API 使用演示 ==="

echo "1. 登录..."
LOGIN_RESP=$(curl -s -X POST "${API_BASE}/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}')
TOKEN=$(echo "$LOGIN_RESP" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['token'])")
echo "   Token: ${TOKEN:0:20}..."

echo "2. 上传记录..."
UPLOAD_RESP=$(curl -s -X POST "${API_BASE}/api/upload-records" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${TOKEN}" \
  -d '{
    "dataType": "SSD-001",
    "projectName": "测试项目",
    "destPath": "/data/demo/file.csv",
    "fileSize": 1048576,
    "uploader": "API测试",
    "status": "completed",
    "remark": "API演示上传"
  }')
SERIAL=$(echo "$UPLOAD_RESP" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['serialNo'])")
echo "   流水号: ${SERIAL}"

echo "3. 查询记录..."
curl -s "${API_BASE}/public/upload-records/${SERIAL}" | python3 -m json.tool

echo ""
echo "=== 演示完成 ==="


# ============================================================
# 5. 数据格式说明
# ============================================================
#
# 字段说明:
#   dataType     磁盘标签（必填，字符串）
#   projectName  项目名称（选填）
#   destPath     文件路径（必填）
#   fileSize     文件大小字节数（必填，整数）
#   uploader     上传人（必填）
#   status       状态：pending | processing | completed | failed
#   remark       备注（选填）
#   data         动态字段（JSON 对象，字段由 FieldConfig 配置决定）
#
# 动态字段示例（取决于系统配置的字段）:
#   {
#     "patientId": "P001",       -- 患者ID
#     "modality": "CT",           -- 检查类型
#     "studyDate": "2024-03-15", -- 检查日期
#     "department": "放射科"       -- 科室
#   }
#
