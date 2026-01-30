// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.
package thumbnail

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/gen2brain/go-fitz"
	"github.com/go-pdf/fpdf"
	"github.com/tealeg/xlsx"
)

// GenerateThumbnail 生成缩略图
// 返回: (生成的缩略图路径, 错误信息)
func GenerateThumbnail(inputPath, outputPath, thumbFolderPath string) (string, string, string, error) {

	var content []byte
	var err error

	// 1. 获取文件内容
	if isURL(inputPath) {
		content, err = downloadFile(inputPath)
		if err != nil {
			return "", "", "", fmt.Errorf("download failed: %w", err)
		}
		// 如果是 URL，默认下载到当前运行目录，或者你可以指定一个临时目录
	} else {
		content, err = ioutil.ReadFile(inputPath)
		if err != nil {
			return "", "", "", fmt.Errorf("read file failed: %w", err)
		}
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		err := os.MkdirAll(outputPath, 0755)
		if err != nil {
			return "", "", "", fmt.Errorf("create directory %s failed: %w", outputPath, err)
		}
	}

	// 2. 构造输出文件名
	baseName := filepath.Base(inputPath)
	ext := filepath.Ext(baseName)
	nameWithoutExt := strings.TrimSuffix(baseName, ext)
	fileName := nameWithoutExt + "_thumb.png"

	// 最终生成的缩略图路径
	thumbnailPath := filepath.Join(outputPath, fileName)

	// 3. 核心：将任意格式转换为 PDF 的二进制数据 (内存处理)
	var pdfBytes []byte

	switch strings.ToLower(ext) {
	case ".pdf":
		pdfBytes = content
	case ".docx":
		pdfBytes, err = convertDocxToPdf(content)
	case ".xlsx":
		pdfBytes, err = convertXlsxToPdf(content)
	case ".csv":
		pdfBytes, err = convertCsvToPdf(content)
	case ".txt", ".md":
		pdfBytes, err = convertTextToPdf(string(content))
	case ".html", ".htm":
		pdfBytes, err = convertHtmlToPdf(content)
	case ".ofd":
		pdfBytes, err = convertOfdToPdf(content)
	default:
		return "", "", "", fmt.Errorf("unsupported file type: %s", ext)
	}

	if err != nil {
		return "", "", "", fmt.Errorf("convert to pdf failed: %w", err)
	}

	// 4. 渲染 PDF 第一页为图片并保存
	err = renderPdfBytesToImage(pdfBytes, thumbnailPath)
	if err != nil {
		return "", "", "", fmt.Errorf("render image failed: %w", err)
	}

	fmt.Printf("Thumbnail generated: %s\n", thumbnailPath)

	// 5. 返回生成的路径
	return filepath.Join(thumbFolderPath, fileName), fileName, thumbnailPath, nil
}

// =========================================================
// 渲染层：PDF -> Image (保持不变)
// =========================================================

func renderPdfBytesToImage(pdfContent []byte, thumbnailPath string) error {
	doc, err := fitz.NewFromMemory(pdfContent)
	if err != nil {
		return fmt.Errorf("fitz load pdf failed: %w", err)
	}
	defer doc.Close()

	if doc.NumPage() == 0 {
		return fmt.Errorf("pdf has no pages")
	}

	// 渲染第一页
	img, err := doc.Image(0)
	if err != nil {
		return fmt.Errorf("render page 0 failed: %w", err)
	}

	// 智能缩放：放入 800x800 容器，保持比例
	dstImage := imaging.Fit(img, 800, 800, imaging.Lanczos)

	return imaging.Save(dstImage, thumbnailPath)
}

// =========================================================
// 转换层：Any -> PDF (使用 go-pdf/fpdf)
// =========================================================

func newPdfGenerator() (*fpdf.Fpdf, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 20, 20)
	pdf.AddPage()

	// 必须要有 simhei.ttf 否则中文无法显示
	fontPath := "internal/pkg/thumbnail/fonts/SimHei.ttf"

	// 这里的路径检查可以根据需要加，fpdf 如果找不到文件会 panic 或报错
	// 建议确保文件存在
	pdf.AddUTF8Font("SimHei", "", fontPath)
	pdf.SetFont("SimHei", "", 12)

	return pdf, nil
}

// DOCX 转 PDF
func convertDocxToPdf(content []byte) ([]byte, error) {
	text := extractDocxText(content)
	pdf, err := newPdfGenerator()
	if err != nil {
		return nil, err
	}

	text = sanitizeText(text)

	pdf.SetFont("SimHei", "", 16)
	pdf.Cell(0, 10, "Document Preview")
	pdf.Ln(12)
	pdf.SetFont("SimHei", "", 12)
	pdf.MultiCell(0, 8, text, "", "", false)

	return generatePdfBytes(pdf)
}

// XLSX 转 PDF
func convertXlsxToPdf(content []byte) ([]byte, error) {
	xlFile, err := xlsx.OpenBinary(content)
	if err != nil {
		return nil, err
	}

	pdf, err := newPdfGenerator()
	if err != nil {
		return nil, err
	}

	pdf.SetFont("SimHei", "", 10)

	for _, sheet := range xlFile.Sheets {
		pdf.SetFont("SimHei", "", 14)
		pdf.Cell(0, 10, "Sheet: "+sheet.Name)
		pdf.Ln(10)
		pdf.SetFont("SimHei", "", 10)

		colWidth := 30.0
		rowHeight := 8.0

		for i, row := range sheet.Rows {
			if i > 50 {
				break
			}
			for j, cell := range row.Cells {
				if j > 7 {
					break
				}
				val := cell.String()

				val = sanitizeText(val)

				if len([]rune(val)) > 15 {
					val = string([]rune(val)[:14]) + ".."
				}
				pdf.CellFormat(colWidth, rowHeight, val, "1", 0, "L", false, 0, "")
			}
			pdf.Ln(-1)
		}
		break
	}
	return generatePdfBytes(pdf)
}

// CSV 转 PDF
func convertCsvToPdf(content []byte) ([]byte, error) {
	if len(content) > 3 && content[0] == 0xEF && content[1] == 0xBB && content[2] == 0xBF {
		content = content[3:]
	}

	pdf, err := newPdfGenerator()
	if err != nil {
		return nil, err
	}

	pdf.SetFont("SimHei", "", 10)
	reader := csv.NewReader(bytes.NewReader(content))

	colWidth := 30.0
	rowHeight := 8.0
	rowCount := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return convertTextToPdf(string(content))
		}

		if rowCount > 50 {
			break
		}

		for i, val := range record {
			if i > 7 {
				break
			}

			val = sanitizeText(val)

			if len([]rune(val)) > 15 {
				val = string([]rune(val)[:14]) + ".."
			}
			pdf.CellFormat(colWidth, rowHeight, val, "1", 0, "L", false, 0, "")
		}
		pdf.Ln(-1)
		rowCount++
	}
	return generatePdfBytes(pdf)
}

// HTML 转 PDF
func convertHtmlToPdf(content []byte) ([]byte, error) {
	htmlStr := string(content)
	reScript := regexp.MustCompile(`(?si)<script.*?>.*?</script>`)
	htmlStr = reScript.ReplaceAllString(htmlStr, "")
	reStyle := regexp.MustCompile(`(?si)<style.*?>.*?</style>`)
	htmlStr = reStyle.ReplaceAllString(htmlStr, "")

	htmlStr = strings.ReplaceAll(htmlStr, "</div>", "\n")
	htmlStr = strings.ReplaceAll(htmlStr, "</p>", "\n")
	htmlStr = strings.ReplaceAll(htmlStr, "<br>", "\n")
	htmlStr = strings.ReplaceAll(htmlStr, "</h1>", "\n")

	reTags := regexp.MustCompile(`<[^>]*>`)
	text := reTags.ReplaceAllString(htmlStr, "")

	text = strings.ReplaceAll(text, "&nbsp;", " ")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")

	return convertTextToPdf(text)
}

// 文本 转 PDF
func convertTextToPdf(text string) ([]byte, error) {
	pdf, err := newPdfGenerator()
	if err != nil {
		return nil, err
	}

	text = sanitizeText(text)

	text = strings.ReplaceAll(text, "\t", "    ")
	text = strings.ReplaceAll(text, "\r", "")
	pdf.MultiCell(0, 8, text, "", "", false)

	return generatePdfBytes(pdf)
}

// OFD 转 PDF (尝试解压解析，失败则降级为文本)
func convertOfdToPdf(content []byte) ([]byte, error) {
	// 尝试作为 ZIP 打开
	r, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		// 【关键修复】：如果不是有效的 ZIP 文件，不要直接报错
		// 可能是原始 XML 格式或文件头损坏，尝试直接按纯文本处理
		// 这样至少能生成一张带有文字内容的缩略图
		return convertTextToPdf(string(content))
	}

	var sb strings.Builder
	fileFound := false

	// 遍历压缩包寻找文本内容
	for _, f := range r.File {
		// 通常 OFD 的核心文字在 Document.xml 或 Content.xml 中
		// 这里放宽匹配条件，提取所有 xml 内容可能会更保险，但目前保持你原有的逻辑
		if strings.HasSuffix(f.Name, "Content.xml") || strings.HasSuffix(f.Name, "Document.xml") {
			rc, err := f.Open()
			if err == nil {
				text := extractDocxTextSimple(rc)
				sb.WriteString(text + "\n")
				rc.Close()
				fileFound = true
			}
		}
	}

	// 如果解压成功了，但在 zip 里没找到预期的 xml 文件
	// 也降级为直接转换原始内容，防止生成空白 PDF
	if !fileFound || sb.Len() == 0 {
		// 尝试读取 zip 中第一个文件作为备选，或者直接渲染原始内容
		return convertTextToPdf(string(content))
	}

	return convertTextToPdf(sb.String())
}

// =========================================================
// 辅助函数
// =========================================================

func generatePdfBytes(pdf *fpdf.Fpdf) ([]byte, error) {
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func extractDocxText(content []byte) string {
	r, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		return ""
	}
	var docXML *zip.File
	for _, f := range r.File {
		if f.Name == "word/document.xml" {
			docXML = f
			break
		}
	}
	if docXML == nil {
		return ""
	}
	rc, _ := docXML.Open()
	defer rc.Close()
	return extractDocxTextSimple(rc)
}

func extractDocxTextSimple(r io.Reader) string {
	decoder := xml.NewDecoder(r)
	var sb strings.Builder
	for {
		t, err := decoder.Token()
		if err != nil {
			break
		}
		switch tok := t.(type) {
		case xml.StartElement:
			if tok.Name.Local == "p" {
				sb.WriteString("\n")
			}
		case xml.CharData:
			sb.WriteString(string(tok))
		}
	}
	return sb.String()
}

func isURL(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}

func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}

func sanitizeText(s string) string {
	return strings.Map(func(r rune) rune {
		// 1. 保留常用的控制字符：换行(\n=10)、回车(\r=13)、制表符(\t=9)
		if r < 32 {
			if r == 10 || r == 13 || r == 9 {
				return r
			}
			return -1 // 删除其他控制字符
		}

		// 2. 删除超过 0xFFFF 的字符 (主要是 Emoji)
		// SimHei 等常见中文字体通常不包含这些 Emoji，且 fpdf 处理 4字节字符会报错
		if r > 0xFFFF {
			return -1
		}

		return r
	}, s)
}
