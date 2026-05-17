# duptwo 数据登记平台 - 客户端 API 示例

> 本文档提供多语言客户端调用示例，涵盖认证、上传、查询等核心场景。

## 目录

- [认证流程](#认证流程)
- [上传记录](#上传记录)
- [项目管理](#项目管理)
- [人员管理](#人员管理)
- [字段配置](#字段配置)
- [数据同步](#数据同步)
- [完整脚本示例](#完整脚本示例)

---

## 基础信息

| 项目 | 值 |
|------|-----|
| 基础 URL | `http://localhost:8080`（生产环境替换为实际地址） |
| 认证方式 | Bearer Token (JWT) |
| 内容类型 | `application/json` |
| MFA | 支持 TOTP 双因素认证 |

**默认账号**: `admin` / `admin123`

---

## 认证流程

### 1. 用户登录

**请求**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

**响应（有 MFA）**
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "",
    "expireAt": 0,
    "user": {
      "id": 1,
      "username": "admin",
      "nickname": "管理员",
      "mfaEnabled": true
    },
    "mfaRequired": true
  }
}
```

**响应（无 MFA）**
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expireAt": 1745760000,
    "user": { "id": 1, "username": "admin", "mfaEnabled": false },
    "mfaRequired": false
  }
}
```

### 2. MFA 验证（如果已启用）

```bash
curl -X POST http://localhost:8080/api/auth/mfa/verify \
  -H "Content-Type: application/json" \
  -d '{
    "userId": 1,
    "code": "123456"
  }'
```

### 3. 获取当前用户信息

```bash
curl http://localhost:8080/api/auth/current \
  -H "Authorization: Bearer <token>"
```

### 4. 修改密码

```bash
curl -X POST http://localhost:8080/api/auth/change-password \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "oldPassword": "admin123",
    "newPassword": "NewPass@123"
  }'
```

### 5. 更新个人资料

```bash
curl -X PUT http://localhost:8080/api/auth/profile \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "nickname": "新昵称",
    "email": "newemail@example.com",
    "phone": "13800138000"
  }'
```

---

## 上传记录

### 公共上传（无需认证）

```bash
curl -X POST http://localhost:8080/public/upload-records \
  -F "file=@/path/to/data.xlsx" \
  -F "dataType=excel" \
  -F "projectName=项目A" \
  -F "uploader=张三"
```

### 认证上传

```bash
curl -X POST http://localhost:8080/api/upload-records \
  -H "Authorization: Bearer <token>" \
  -F "file=@/path/to/data.xlsx" \
  -F "dataType=excel" \
  -F "projectName=项目A" \
  -F "uploader=张三" \
  -F "remark=批次1数据"
```

### 查询上传记录列表

```bash
curl "http://localhost:8080/api/upload-records?page=1&pageSize=20" \
  -H "Authorization: Bearer <token>"
```

### 查询单条记录

```bash
curl http://localhost:8080/api/upload-records/1 \
  -H "Authorization: Bearer <token>"
```

### 更新上传记录

```bash
curl -X PUT http://localhost:8080/api/upload-records \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "status": "processed",
    "remark": "已审核"
  }'
```

### 删除上传记录

```bash
curl -X DELETE http://localhost:8080/api/upload-records/1 \
  -H "Authorization: Bearer <token>"
```

### 导出上传记录

```bash
curl -o export.xlsx \
  "http://localhost:8080/api/upload-records/export?format=xlsx&page=1&pageSize=100" \
  -H "Authorization: Bearer <token>"
```

### 下载导入模板

```bash
curl -o template.xlsx \
  http://localhost:8080/api/upload-records/template \
  -H "Authorization: Bearer <token>"
```

### 批量导入

```bash
curl -X POST http://localhost:8080/api/upload-records/import \
  -H "Authorization: Bearer <token>" \
  -F "file=@/path/to/template_filled.xlsx"
```

### 获取统计数据

```bash
curl http://localhost:8080/api/upload-records/statistics \
  -H "Authorization: Bearer <token>"
```

---

## 项目管理

### 创建项目

```bash
curl -X POST http://localhost:8080/api/projects \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "智慧城市数据平台",
    "code": "smart-city",
    "description": "智慧城市多源数据采集与分析",
    "status": "active",
    "stage": "running",
    "projectPerson": "李经理",
    "opsPerson": "王运维",
    "developerPerson": "张开发",
    "testerPerson": "赵测试",
    "businessPerson": "刘商务",
    "compliancePerson": "陈合规",
    "solution": "智慧城市综合解决方案",
    "companyAddr": "北京市海淀区",
    "projectPeriod": "2025-2027",
    "onsiteStations": [
      { "location": "数据中心A", "person": "驻场人员A", "phone": "13800001111" },
      { "location": "数据中心B", "person": "驻场人员B", "phone": "13800002222" }
    ]
  }'
```

### 查询项目列表

```bash
curl "http://localhost:8080/api/projects?page=1&pageSize=20&keyword=智慧" \
  -H "Authorization: Bearer <token>"
```

### 查询所有项目（下拉选择用）

```bash
curl http://localhost:8080/api/projects/simple \
  -H "Authorization: Bearer <token>"
```

### 查询看板视图

```bash
curl http://localhost:8080/api/projects/kanban \
  -H "Authorization: Bearer <token>"
```

### 更新项目

```bash
curl -X PUT http://localhost:8080/api/projects \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "name": "智慧城市数据平台（更新）",
    "status": "paused",
    "stage": "paused"
  }'
```

### 删除项目

```bash
curl -X DELETE http://localhost:8080/api/projects/1 \
  -H "Authorization: Bearer <token>"
```

---

## 人员管理

### 创建人员

```bash
curl -X POST http://localhost:8080/api/personnels \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "张三",
    "gender": "男",
    "age": 30,
    "phone": "13800138000",
    "email": "zhangsan@example.com",
    "company": "XX科技有限公司",
    "department": "技术部",
    "position": "高级工程师",
    "onboardDate": "2024-01-15",
    "projectId": 1,
    "location": "北京市",
    "skills": ["Python", "Go", "K8s"],
    "status": "active",
    "remark": "核心开发成员"
  }'
```

### 查询人员列表

```bash
curl "http://localhost:8080/api/personnels?page=1&pageSize=20&keyword=张" \
  -H "Authorization: Bearer <token>"
```

### 导出人员列表

```bash
curl -o personnel.xlsx \
  "http://localhost:8080/api/personnels/export?format=xlsx" \
  -H "Authorization: Bearer <token>"
```

### 更新人员

```bash
curl -X PUT http://localhost:8080/api/personnels \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "position": "技术总监",
    "status": "active"
  }'
```

---

## 字段配置

### 创建字段配置

```bash
curl -X POST http://localhost:8080/api/field-configs \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "数据来源",
    "code": "dataSource",
    "type": "select",
    "options": ["API", "手动录入", "批量导入", "爬虫"],
    "enabled": true,
    "sort": 1
  }'
```

### 查询所有启用的字段配置

```bash
curl http://localhost:8080/api/field-configs/all \
  -H "Authorization: Bearer <token>"
```

---

## 数据同步

> 数据同步系统采用 Center + Agent 架构，支持多站点数据汇聚。

### 基础信息

| 项目 | 值 |
|------|-----|
| 认证方式（管理） | Bearer Token (JWT) |
| 认证方式（Agent） | X-API-Key 请求头 |

### 站点管理（需 JWT 认证）

#### 获取站点列表

```bash
curl "http://localhost:8080/api/v1/sync/stations?page=1&pageSize=20" \
  -H "Authorization: Bearer <token>"
```

#### 创建站点

```bash
curl -X POST http://localhost:8080/api/v1/sync/stations \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "北京数据中心",
    "code": "bj-dc-01",
    "url": "http://bj-agent:8080",
    "status": "active",
    "description": "北京主数据中心"
  }'
```

#### 更新站点

```bash
curl -X PUT http://localhost:8080/api/v1/sync/stations/1 \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "北京数据中心（新）",
    "status": "paused"
  }'
```

#### 删除站点

```bash
curl -X DELETE http://localhost:8080/api/v1/sync/stations/1 \
  -H "Authorization: Bearer <token>"
```

### Agent 注册（无需认证）

#### 注册站点

```bash
curl -X POST http://localhost:8080/api/v1/sync/register \
  -H "Content-Type: application/json" \
  -d '{
    "stationId": "",
    "stationName": "上海站点",
    "url": "http://sh-agent:8080",
    "password": "agent_secret_password"
  }'
```

**响应示例**：

```json
{
  "code": 200,
  "message": "注册成功",
  "data": {
    "stationId": "1",
    "apiKey": "sk_sync_abc123xyz789...",
    "stationName": "上海站点"
  }
}
```

### Agent 上传记录（API Key 认证）

#### 上传记录

```bash
curl -X POST http://localhost:8080/api/v1/sync/upload \
  -H "X-API-Key: sk_sync_abc123xyz789..." \
  -H "Content-Type: application/json" \
  -d '{
    "records": [
      {
        "dataType": "sensor",
        "projectName": "环境监测",
        "uploader": "agent",
        "data": {
          "temperature": 25.5,
          "humidity": 60,
          "location": "上海浦东"
        }
      },
      {
        "dataType": "sensor",
        "projectName": "环境监测",
        "uploader": "agent",
        "data": {
          "temperature": 26.0,
          "humidity": 58,
          "location": "上海浦东"
        }
      }
    ]
  }'
```

### 同步历史（需 JWT 认证）

#### 查询同步历史

```bash
curl "http://localhost:8080/api/v1/sync/history?page=1&pageSize=20&stationId=1&startDate=2025-01-01&endDate=2025-12-31" \
  -H "Authorization: Bearer <token>"
```

#### 获取同步详情

```bash
curl http://localhost:8080/api/v1/sync/history/1/details \
  -H "Authorization: Bearer <token>"
```

#### 获取同步状态

```bash
curl http://localhost:8080/api/v1/sync/status \
  -H "Authorization: Bearer <token>"
```

---

## 完整脚本示例

### Bash 完整调用脚本

```bash
#!/bin/bash
# duptwo_api.sh - duptwo API 完整调用示例
# 用法: ./duptwo_api.sh <command> [args...]

BASE_URL="${API_URL:-http://localhost:8080}"
TOKEN=""

# 颜色输出
RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; NC='\033[0m'

log_ok()  { echo -e "${GREEN}[OK]${NC}  $*"; }
log_err() { echo -e "${RED}[ERR]${NC} $*"; }
log_info(){ echo -e "${YELLOW}[INFO]${NC} $*"; }

# 发送请求
request() {
  local method=$1
  local path=$2
  local data=$3
  local auth="${TOKEN:+ -H "Authorization: Bearer $TOKEN"}"

  if [ -n "$data" ]; then
    curl -s -X "$method" "${BASE_URL}${path}" \
      -H "Content-Type: application/json" \
      $auth \
      -d "$data"
  else
    curl -s -X "$method" "${BASE_URL}${path}" $auth
  fi
}

# 登录
cmd_login() {
  log_info "登录中..."
  local resp=$(request POST "/api/auth/login" \
    "{\"username\":\"$1\",\"password\":\"$2\"}")

  local code=$(echo "$resp" | python3 -c "import sys,json; print(json.load(sys.stdin).get('code',''))" 2>/dev/null)
  local mfa=$(echo "$resp" | python3 -c "import sys,json; print(json.load(sys.stdin).get('data',{}).get('mfaRequired',''))" 2>/dev/null)
  local token=$(echo "$resp" | python3 -c "import sys,json; print(json.load(sys.stdin).get('data',{}).get('token',''))" 2>/dev/null)

  if [ "$code" = "200" ]; then
    if [ "$mfa" = "True" ]; then
      log_ok "需要 MFA 验证，请输入验证码:"
      read -p "验证码: " mfa_code
      local user_id=$(echo "$resp" | python3 -c "import sys,json; print(json.load(sys.stdin).get('data',{}).get('user',{}).get('id',''))" 2>/dev/null)
      local mfa_resp=$(request POST "/api/auth/mfa/verify" \
        "{\"userId\":$user_id,\"code\":\"$mfa_code\"}")
      TOKEN=$(echo "$mfa_resp" | python3 -c "import sys,json; print(json.load(sys.stdin).get('data',{}).get('token',''))" 2>/dev/null)
      if [ -n "$TOKEN" ]; then
        log_ok "MFA 验证成功"
      else
        log_err "MFA 验证失败"
      fi
    else
      TOKEN="$token"
      log_ok "登录成功"
    fi
  else
    local msg=$(echo "$resp" | python3 -c "import sys,json; print(json.load(sys.stdin).get('message',''))" 2>/dev/null)
    log_err "登录失败: $msg"
  fi
}

# 查询上传记录
cmd_list_records() {
  log_info "查询上传记录..."
  request GET "/api/upload-records?page=1&pageSize=20" | python3 -m json.tool 2>/dev/null || echo "$resp"
}

# 创建项目
cmd_create_project() {
  log_info "创建项目..."
  local data="{\"name\":\"$1\",\"code\":\"$2\",\"status\":\"active\",\"stage\":\"planning\"}"
  request POST "/api/projects" "$data" | python3 -m json.tool 2>/dev/null
}

# 查询项目列表
cmd_list_projects() {
  log_info "查询项目列表..."
  request GET "/api/projects?page=1&pageSize=20" | python3 -m json.tool 2>/dev/null
}

# 获取统计数据
cmd_stats() {
  log_info "获取上传统计数据..."
  request GET "/api/upload-records/statistics" | python3 -m json.tool 2>/dev/null
}

# 帮助
cmd_help() {
  echo "用法: $0 <command> [args...]"
  echo ""
  echo "命令:"
  echo "  login <username> <password>     登录"
  echo "  list-records                    查询上传记录"
  echo "  list-projects                   查询项目列表"
  echo "  create-project <name> <code>     创建项目"
  echo "  stats                           获取统计数据"
  echo "  help                            显示帮助"
  echo ""
  echo "环境变量:"
  echo "  API_URL   API 基础地址（默认 http://localhost:8080）"
}

case "$1" in
  login)          cmd_login "$2" "$3" ;;
  list-records)   cmd_list_records ;;
  list-projects)  cmd_list_projects ;;
  create-project) cmd_create_project "$2" "$3" ;;
  stats)          cmd_stats ;;
  *)              cmd_help ;;
esac
```

```bash
# 使用示例
chmod +x duptwo_api.sh

# 登录
./duptwo_api.sh login admin admin123

# 查询记录（登录后 TOKEN 已保存）
./duptwo_api.sh list-records

# 创建项目
./duptwo_api.sh create-project "新数据中心项目" "new-dc-001"

# 查看统计
./duptwo_api.sh stats
```

### Python 调用脚本

```python
#!/usr/bin/env python3
"""
duptwo_client.py - duptwo 数据登记平台 Python 客户端
pip install requests pyotp
"""

import requests
import json
import sys
from typing import Optional, Dict, Any

BASE_URL = "http://localhost:8080"


class DupTwoClient:
    def __init__(self, base_url: str = BASE_URL):
        self.base_url = base_url.rstrip("/")
        self.token: Optional[str] = None
        self.session = requests.Session()
        self.session.headers.update({"Content-Type": "application/json"})

    def _request(self, method: str, path: str, data: Dict = None) -> Dict[str, Any]:
        url = f"{self.base_url}{path}"
        headers = {}
        if self.token:
            headers["Authorization"] = f"Bearer {self.token}"

        if method == "GET":
            resp = self.session.get(url, headers=headers, params=data)
        elif method == "POST":
            resp = self.session.post(url, headers=headers, json=data)
        elif method == "PUT":
            resp = self.session.put(url, headers=headers, json=data)
        elif method == "DELETE":
            resp = self.session.delete(url, headers=headers)
        else:
            raise ValueError(f"Unsupported method: {method}")

        resp.raise_for_status()
        return resp.json()

    # ===== 认证 =====
    def login(self, username: str, password: str) -> Dict[str, Any]:
        """登录，返回登录结果"""
        result = self._request("POST", "/api/auth/login", {"username": username, "password": password})
        if result.get("code") != 200:
            raise Exception(result.get("message", "登录失败"))

        data = result.get("data", {})
        if data.get("mfaRequired"):
            print(f"[INFO] 账号已启用 MFA，请输入验证码: ", end="", flush=True)
            code = input().strip()
            mfa_result = self._request("POST", "/api/auth/mfa/verify", {
                "userId": data["user"]["id"],
                "code": code
            })
            if mfa_result.get("code") != 200:
                raise Exception(mfa_result.get("message", "MFA 验证失败"))
            self.token = mfa_result["data"]["token"]
            print("[INFO] MFA 验证成功")
        else:
            self.token = data.get("token", "")

        print(f"[OK] 登录成功，Token: {self.token[:20]}...")
        return result

    def get_current_user(self) -> Dict[str, Any]:
        """获取当前用户信息"""
        return self._request("GET", "/api/auth/current")

    def update_profile(self, nickname: str = None, email: str = None, phone: str = None) -> Dict[str, Any]:
        """更新个人资料"""
        data = {}
        if nickname is not None: data["nickname"] = nickname
        if email is not None: data["email"] = email
        if phone is not None: data["phone"] = phone
        return self._request("PUT", "/api/auth/profile", data)

    def change_password(self, old: str, new: str) -> Dict[str, Any]:
        """修改密码"""
        return self._request("POST", "/api/auth/change-password", {
            "oldPassword": old,
            "newPassword": new
        })

    # ===== 上传记录 =====
    def list_upload_records(self, page: int = 1, page_size: int = 20, **kwargs) -> Dict[str, Any]:
        """查询上传记录"""
        params = {"page": page, "pageSize": page_size, **kwargs}
        return self._request("GET", "/api/upload-records", params)

    def upload_file(self, file_path: str, data_type: str = "excel",
                    project_name: str = "", uploader: str = "", remark: str = "") -> Dict[str, Any]:
        """上传文件"""
        with open(file_path, "rb") as f:
            files = {"file": f}
            data = {"dataType": data_type}
            if project_name: data["projectName"] = project_name
            if uploader: data["uploader"] = uploader
            if remark: data["remark"] = remark

            url = f"{self.base_url}/api/upload-records"
            headers = {}
            if self.token:
                headers["Authorization"] = f"Bearer {self.token}"

            resp = self.session.post(url, headers=headers, files=files, data=data)
            resp.raise_for_status()
            return resp.json()

    def get_statistics(self) -> Dict[str, Any]:
        """获取统计数据"""
        return self._request("GET", "/api/upload-records/statistics")

    def export_records(self, format: str = "xlsx", **kwargs) -> bytes:
        """导出记录"""
        params = {"format": format, **kwargs}
        url = f"{self.base_url}/api/upload-records/export"
        headers = {}
        if self.token:
            headers["Authorization"] = f"Bearer {self.token}"
        resp = self.session.get(url, headers=headers, params=params)
        resp.raise_for_status()
        return resp.content

    # ===== 项目 =====
    def list_projects(self, page: int = 1, page_size: int = 20, **kwargs) -> Dict[str, Any]:
        """查询项目列表"""
        params = {"page": page, "pageSize": page_size, **kwargs}
        return self._request("GET", "/api/projects", params)

    def create_project(self, name: str, code: str, **kwargs) -> Dict[str, Any]:
        """创建项目"""
        data = {"name": name, "code": code, **kwargs}
        return self._request("POST", "/api/projects", data)

    def get_project(self, project_id: int) -> Dict[str, Any]:
        """获取项目详情"""
        return self._request("GET", f"/api/projects/{project_id}")

    def update_project(self, project_id: int, **kwargs) -> Dict[str, Any]:
        """更新项目"""
        data = {"id": project_id, **kwargs}
        return self._request("PUT", "/api/projects", data)

    def delete_project(self, project_id: int) -> Dict[str, Any]:
        """删除项目"""
        return self._request("DELETE", f"/api/projects/{project_id}")

    # ===== 人员 =====
    def list_personnels(self, page: int = 1, page_size: int = 20, **kwargs) -> Dict[str, Any]:
        """查询人员列表"""
        params = {"page": page, "pageSize": page_size, **kwargs}
        return self._request("GET", "/api/personnels", params)

    def create_personnel(self, name: str, **kwargs) -> Dict[str, Any]:
        """创建人员"""
        data = {"name": name, **kwargs}
        return self._request("POST", "/api/personnels", data)


# ===== 使用示例 =====
if __name__ == "__main__":
    client = DupTwoClient("http://localhost:8080")

    # 1. 登录
    try:
        client.login("admin", "admin123")
    except Exception as e:
        print(f"[ERROR] {e}")
        sys.exit(1)

    # 2. 查看统计数据
    stats = client.get_statistics()
    print(f"[INFO] 总记录数: {stats.get('data', {}).get('total', 0)}")

    # 3. 查询项目列表
    projects = client.list_projects(page=1, page_size=10)
    print(f"[INFO] 项目数: {projects.get('data', {}).get('total', 0)}")

    # 4. 创建项目
    new_project = client.create_project(
        name="Python 客户端测试项目",
        code="py-test-001",
        status="active",
        stage="planning",
        description="通过 Python 客户端创建的项目"
    )
    print(f"[INFO] 新项目ID: {new_project.get('data', {}).get('id')}")

    # 5. 上传文件
    # result = client.upload_file(
    #     "/path/to/data.xlsx",
    #     data_type="excel",
    #     project_name="测试项目",
    #     uploader="API用户"
    # )

    # 6. 导出数据
    # excel_data = client.export_records(format="xlsx", page=1, pageSize=500)
    # with open("export.xlsx", "wb") as f:
    #     f.write(excel_data)
    # print("[INFO] 导出完成: export.xlsx")
```

### JavaScript / Node.js 调用

```javascript
/**
 * duptwo_client.js - duptwo API Node.js 客户端
 * 用法: node duptwo_client.js <command> [args]
 * 依赖: npm install axios form-data node-fetch
 */

const axios = require("axios");
const FormData = require("form-data");
const fs = require("fs");
const readline = require("readline");

const BASE_URL = process.env.API_URL || "http://localhost:8080";

class DupTwoClient {
  constructor(baseUrl = BASE_URL) {
    this.baseUrl = baseUrl;
    this.token = null;
    this.client = axios.create({ baseURL: baseUrl });
  }

  setToken(token) { this.token = token; }

  get headers() {
    const h = { "Content-Type": "application/json" };
    if (this.token) h["Authorization"] = `Bearer ${this.token}`;
    return h;
  }

  async request(method, path, data = null, isForm = false) {
    const config = { headers: this.headers };
    if (isForm) {
      delete config.headers["Content-Type"];
      config.headers["Content-Type"] = "multipart/form-data";
    }
    let resp;
    if (method === "GET") resp = await this.client.get(path, { ...config, params: data });
    else if (method === "POST") resp = await this.client.post(path, data, config);
    else if (method === "PUT") resp = await this.client.put(path, data, config);
    else if (method === "DELETE") resp = await this.client.delete(path, config);
    return resp.data;
  }

  async login(username, password) {
    const result = await this.request("POST", "/api/auth/login", { username, password });
    if (result.code !== 200) throw new Error(result.message);

    const data = result.data;
    if (data.mfaRequired) {
      const rl = readline.createInterface({ input: process.stdin, output: process.stdout });
      const code = await new Promise(resolve => rl.question("[INFO] 需要 MFA 验证码: ", resolve));
      rl.close();
      const mfaResult = await this.request("POST", "/api/auth/mfa/verify", {
        userId: data.user.id, code
      });
      if (mfaResult.code !== 200) throw new Error(mfaResult.message);
      this.setToken(mfaResult.data.token);
      console.log("[OK] MFA 验证成功");
    } else {
      this.setToken(data.token);
      console.log("[OK] 登录成功");
    }
    return result;
  }

  async listProjects(page = 1, pageSize = 20) {
    return this.request("GET", "/api/projects", { page, pageSize });
  }

  async createProject(name, code, extra = {}) {
    return this.request("POST", "/api/projects", { name, code, ...extra });
  }

  async listUploadRecords(page = 1, pageSize = 20) {
    return this.request("GET", "/api/upload-records", { page, pageSize });
  }

  async getStatistics() {
    return this.request("GET", "/api/upload-records/statistics");
  }

  async uploadFile(filePath, meta = {}) {
    const form = new FormData();
    form.append("file", fs.createReadStream(filePath));
    Object.entries(meta).forEach(([k, v]) => form.append(k, v));
    const resp = await this.client.post("/api/upload-records", form, {
      headers: { ...form.getHeaders(), Authorization: `Bearer ${this.token}` }
    });
    return resp.data;
  }
}

async function main() {
  const client = new DupTwoClient();
  const cmd = process.argv[2];

  try {
    switch (cmd) {
      case "login": {
        const username = process.argv[3] || "admin";
        const password = process.argv[4] || "admin123";
        await client.login(username, password);
        break;
      }
      case "projects": {
        await client.login("admin", "admin123");
        const result = await client.listProjects();
        console.log(JSON.stringify(result, null, 2));
        break;
      }
      case "create-project": {
        await client.login("admin", "admin123");
        const name = process.argv[3] || "Node测试项目";
        const code = process.argv[4] || "node-test-001";
        const result = await client.createProject(name, code, { status: "active", stage: "planning" });
        console.log("[OK] 项目创建成功:", JSON.stringify(result.data));
        break;
      }
      case "stats": {
        await client.login("admin", "admin123");
        const result = await client.getStatistics();
        console.log(JSON.stringify(result, null, 2));
        break;
      }
      default:
        console.log("用法:");
        console.log("  node duptwo_client.js login [user] [pass]");
        console.log("  node duptwo_client.js projects");
        console.log("  node duptwo_client.js create-project [name] [code]");
        console.log("  node duptwo_client.js stats");
    }
  } catch (err) {
    console.error("[ERROR]", err.message);
    process.exit(1);
  }
}

main();
```

### cURL 常用命令速查

```bash
# === 认证 ===
# 登录
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# 获取当前用户
curl http://localhost:8080/api/auth/current \
  -H "Authorization: Bearer YOUR_TOKEN"

# === 上传记录 ===
# 查询列表
curl "http://localhost:8080/api/upload-records?page=1&pageSize=20" \
  -H "Authorization: Bearer YOUR_TOKEN"

# 上传文件
curl -X POST http://localhost:8080/api/upload-records \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@data.xlsx" \
  -F "dataType=excel" \
  -F "projectName=项目A" \
  -F "uploader=API"

# 导出 Excel
curl -o export.xlsx \
  "http://localhost:8080/api/upload-records/export?format=xlsx" \
  -H "Authorization: Bearer YOUR_TOKEN"

# 下载模板
curl -o template.xlsx \
  http://localhost:8080/api/upload-records/template \
  -H "Authorization: Bearer YOUR_TOKEN"

# === 项目 ===
# 创建
curl -X POST http://localhost:8080/api/projects \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"新项目","code":"new-proj","status":"active","stage":"planning"}'

# 查询列表
curl "http://localhost:8080/api/projects?page=1&pageSize=20" \
  -H "Authorization: Bearer YOUR_TOKEN"

# 更新
curl -X PUT http://localhost:8080/api/projects \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"id":1,"status":"paused"}'

# 删除
curl -X DELETE http://localhost:8080/api/projects/1 \
  -H "Authorization: Bearer YOUR_TOKEN"

# === 人员 ===
# 创建
curl -X POST http://localhost:8080/api/personnels \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"李四","phone":"13900001111","status":"active"}'

# 列表
curl "http://localhost:8080/api/personnels?page=1&pageSize=20" \
  -H "Authorization: Bearer YOUR_TOKEN"

# === 字段配置 ===
# 创建
curl -X POST http://localhost:8080/api/field-configs \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"紧急程度","code":"urgency","type":"select","options":["高","中","低"],"enabled":true}'

# 所有启用的字段
curl http://localhost:8080/api/field-configs/all \
  -H "Authorization: Bearer YOUR_TOKEN"

# === 数据同步 - 站点管理 ===
# 站点列表
curl "http://localhost:8080/api/v1/sync/stations?page=1&pageSize=20" \
  -H "Authorization: Bearer YOUR_TOKEN"

# 创建站点
curl -X POST http://localhost:8080/api/v1/sync/stations \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"站点A","code":"station-a","url":"http://agent-a:8080","status":"active"}'

# Agent 注册
curl -X POST http://localhost:8080/api/v1/sync/register \
  -H "Content-Type: application/json" \
  -d '{"stationName":"站点B","url":"http://agent-b:8080","password":"secret"}'

# Agent 上传记录
curl -X POST http://localhost:8080/api/v1/sync/upload \
  -H "X-API-Key: sk_sync_xxx" \
  -H "Content-Type: application/json" \
  -d '{"records":[{"dataType":"sensor","data":{"temp":25}}]}'

# 同步历史
curl "http://localhost:8080/api/v1/sync/history?page=1&pageSize=20" \
  -H "Authorization: Bearer YOUR_TOKEN"

# 同步状态
curl http://localhost:8080/api/v1/sync/status \
  -H "Authorization: Bearer YOUR_TOKEN"

# === 日志 ===
# 操作日志
curl "http://localhost:8080/api/audit/operation-logs?page=1&pageSize=20" \
  -H "Authorization: Bearer YOUR_TOKEN"

# 登录日志
curl "http://localhost:8080/api/audit/login-logs?page=1&pageSize=20" \
  -H "Authorization: Bearer YOUR_TOKEN"
```
