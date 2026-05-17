package service

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/big"
	"strings"
	"sync"
	"time"
)

// CaptchaStore 验证码存储
type CaptchaStore struct {
	mu      sync.RWMutex
	captchas map[string]*captchaData
}

type captchaData struct {
	Answer  string
	Expires time.Time
}

var captchaStore = &CaptchaStore{
	captchas: make(map[string]*captchaData),
}

// CaptchaResp 验证码响应
type CaptchaResp struct {
	CaptchaID    string `json:"captchaId"`
	CaptchaImage string `json:"captchaImage"` // base64 encoded PNG image
	Enabled      bool   `json:"enabled"`
}

// GenerateCaptcha 生成验证码
func (s *AuthService) GenerateCaptcha() (*CaptchaResp, error) {
	// 生成两个数字和运算符
	a, _ := rand.Int(rand.Reader, big.NewInt(9))
	b, _ := rand.Int(rand.Reader, big.NewInt(9))
	op := "+"
	result := new(big.Int)

	if a.Int64() < b.Int64() {
		a, b = b, a
		if n, _ := rand.Int(rand.Reader, big.NewInt(2)); n.Int64() == 1 {
			op = "-"
			result.Sub(a, b)
		} else {
			result.Add(a, b)
		}
	} else {
		if n, _ := rand.Int(rand.Reader, big.NewInt(10)); n.Int64() < 7 {
			op = "+"
			result.Add(a, b)
		} else {
			op = "-"
			result.Sub(a, b)
		}
	}

	expression := fmt.Sprintf("%s %s %s = ?", a.String(), op, b.String())
	answer := result.String()

	// 生成唯一ID
	idBytes := make([]byte, 16)
	rand.Read(idBytes)
	captchaID := fmt.Sprintf("%x", idBytes)

	// 存储验证码
	captchaStore.mu.Lock()
	captchaStore.captchas[captchaID] = &captchaData{
		Answer:  answer,
		Expires: time.Now().Add(5 * time.Minute),
	}
	captchaStore.mu.Unlock()

	// 生成验证码图片
	imgBase64, err := generateCaptchaImage(expression)
	if err != nil {
		return nil, err
	}

	return &CaptchaResp{
		CaptchaID:    captchaID,
		CaptchaImage: imgBase64,
		Enabled:      true,
	}, nil
}

// generateCaptchaImage 生成简洁清晰的验证码图片（PNG），返回 base64 编码
func generateCaptchaImage(expression string) (string, error) {
	width, height := 280, 80

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 纯白背景
	bg := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	drawFill(img, 0, 0, width, height, bg)

	// 每个像素块大小：3x3
	px := 3
	// 字符间距
	spacing := 8
	// 计算起始X，让表达式居中
	charCount := len(expression)
	charWidth := 5*px + spacing
	totalWidth := charCount*charWidth - spacing
	startX := (width - totalWidth) / 2
	startY := (height - 7*px) / 2

	// 文字颜色：深黑色
	charColor := color.RGBA{R: 10, G: 10, B: 10, A: 255}

	for i, ch := range expression {
		cx := startX + i*charWidth
		drawChar(img, cx, startY, string(ch), px, charColor)
	}

	// 编码为 PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// drawFill 用颜色填充矩形区域
func drawFill(img *image.RGBA, x1, y1, x2, y2 int, c color.RGBA) {
	for y := y1; y < y2; y++ {
		for x := x1; x < x2; x++ {
			img.Set(x, y, c)
		}
	}
}

// drawLine 绘制线段（Bresenham）
func drawLine(img *image.RGBA, x1, y1, x2, y2 int, c color.RGBA) {
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	sx, sy := -1, -1
	if x1 < x2 {
		sx = 1
	}
	if y1 < y2 {
		sy = 1
	}
	err := dx - dy
	for {
		if x1 >= 0 && x1 < img.Bounds().Dx() && y1 >= 0 && y1 < img.Bounds().Dy() {
			img.Set(x1, y1, c)
		}
		if x1 == x2 && y1 == y2 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x1 += sx
		}
		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// drawChar 用像素块绘制单个字符，支持 0-9、+、-、=、?
// 网格: 5列(x=0~4) x 7行(y=0~6)，每个像素块 px 大小
func drawChar(img *image.RGBA, x, y int, ch string, px int, c color.RGBA) {
	type pixel struct{ gx, gy int }
	var pixels []pixel

	switch ch {
	case "0":
		pixels = []pixel{
			{1, 0}, {2, 0}, {3, 0},
			{0, 1}, {4, 1},
			{0, 2}, {4, 2},
			{0, 3}, {4, 3},
			{0, 4}, {4, 4},
			{0, 5}, {4, 5},
			{1, 6}, {2, 6}, {3, 6},
		}
	case "1":
		pixels = []pixel{
			{2, 0},
			{1, 1}, {2, 1},
			{2, 2}, {3, 2},
			{2, 3},
			{2, 4},
			{1, 5}, {2, 5}, {3, 5},
			{2, 6},
		}
	case "2":
		pixels = []pixel{
			{1, 0}, {2, 0}, {3, 0},
			{0, 1}, {4, 1},
			{3, 2},
			{2, 3}, {3, 3},
			{1, 4},
			{0, 5},
			{0, 6}, {1, 6}, {2, 6}, {3, 6}, {4, 6},
		}
	case "3":
		pixels = []pixel{
			{1, 0}, {2, 0}, {3, 0},
			{0, 1}, {4, 1},
			{3, 2},
			{1, 3}, {2, 3}, {3, 3},
			{4, 4},
			{0, 5}, {4, 5},
			{1, 6}, {2, 6}, {3, 6},
		}
	case "4":
		pixels = []pixel{
			{0, 0}, {4, 0},
			{0, 1}, {4, 1},
			{0, 2}, {4, 2},
			{0, 3}, {1, 3}, {2, 3}, {3, 3}, {4, 3},
			{4, 4},
			{4, 5},
			{4, 6},
		}
	case "5":
		pixels = []pixel{
			{0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 0},
			{0, 1},
			{0, 2}, {1, 2}, {2, 2}, {3, 2},
			{4, 3},
			{0, 4}, {4, 4},
			{0, 5}, {4, 5},
			{1, 6}, {2, 6}, {3, 6},
		}
	case "6":
		pixels = []pixel{
			{1, 0}, {2, 0}, {3, 0},
			{0, 1}, {4, 1},
			{0, 2},
			{0, 3}, {1, 3}, {2, 3}, {3, 3},
			{0, 4}, {4, 4},
			{0, 5}, {4, 5},
			{1, 6}, {2, 6}, {3, 6},
		}
	case "7":
		pixels = []pixel{
			{0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 0},
			{4, 1},
			{3, 2},
			{3, 3},
			{2, 4},
			{2, 5},
			{1, 6},
		}
	case "8":
		pixels = []pixel{
			{1, 0}, {2, 0}, {3, 0},
			{0, 1}, {4, 1},
			{0, 2}, {4, 2},
			{1, 3}, {2, 3}, {3, 3},
			{0, 4}, {4, 4},
			{0, 5}, {4, 5},
			{1, 6}, {2, 6}, {3, 6},
		}
	case "9":
		pixels = []pixel{
			{1, 0}, {2, 0}, {3, 0},
			{0, 1}, {4, 1},
			{0, 2}, {4, 2},
			{1, 3}, {2, 3}, {3, 3}, {4, 3},
			{4, 4},
			{4, 5},
			{1, 6}, {2, 6}, {3, 6},
		}
	case "+":
		pixels = []pixel{
			{2, 1}, {2, 2}, {2, 3}, {2, 4}, {2, 5},
			{1, 3}, {2, 3}, {3, 3},
		}
	case "-":
		pixels = []pixel{
			{1, 3}, {2, 3}, {3, 3},
		}
	case "=":
		pixels = []pixel{
			{1, 2}, {2, 2}, {3, 2},
			{1, 4}, {2, 4}, {3, 4},
		}
	case "?":
		pixels = []pixel{
			{1, 0}, {2, 0}, {3, 0},
			{0, 1}, {4, 1},
			{3, 2},
			{2, 3}, {3, 3},
			{2, 4},
			{2, 5},
		}
	case " ":
		// 空格不画
	default:
		return
	}

	// 绘制每个像素块
	for _, p := range pixels {
		for dy := 0; dy < px; dy++ {
			for dx := 0; dx < px; dx++ {
				img.Set(x+p.gx*px+dx, y+p.gy*px+dy, c)
			}
		}
	}
}

// ValidateCaptcha 验证验证码
func (s *AuthService) ValidateCaptcha(captchaID, answer string) error {
	captchaStore.mu.Lock()
	defer captchaStore.mu.Unlock()

	data, exists := captchaStore.captchas[captchaID]
	if !exists {
		return errors.New("验证码已过期，请刷新")
	}

	if time.Now().After(data.Expires) {
		delete(captchaStore.captchas, captchaID)
		return errors.New("验证码已过期，请刷新")
	}

	answer = strings.TrimSpace(answer)
	if data.Answer != answer {
		return errors.New("验证码错误")
	}

	delete(captchaStore.captchas, captchaID)
	return nil
}

// IsCaptchaEnabled 检查是否启用验证码
func (s *AuthService) IsCaptchaEnabled() bool {
	return true
}
