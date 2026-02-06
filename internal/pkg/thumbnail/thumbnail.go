// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.
package thumbnail

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/gen2brain/go-fitz"
	"github.com/go-pdf/fpdf"
	"github.com/tealeg/xlsx"
)

// GenerateThumbnail generates a thumbnail (memory processing version)
// input: content (original file binary), filename (used to identify file type, e.g., "test.docx")
// return: (thumbnail PNG binary, error message)
func GenerateThumbnail(content []byte, filename string) ([]byte, string, error) {
	// 1. Get file extension
	ext := filepath.Ext(filename)

	// 2. Core: Convert any format to PDF binary data (in-memory processing)
	var pdfBytes []byte
	var err error

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
		return nil, "", fmt.Errorf("unsupported file type: %s", ext)
	}

	if err != nil {
		return nil, "", fmt.Errorf("convert to pdf failed: %w", err)
	}

	// 3. Render PDF first page as image binary
	thumbBytes, err := renderPdfBytesToImageBytes(pdfBytes)
	if err != nil {
		return nil, "", fmt.Errorf("render image failed: %w", err)
	}

	nameWithoutExt := strings.TrimSuffix(filename, ext)
	fileName := nameWithoutExt + "_thumb.png"

	return thumbBytes, fileName, nil
}

// =========================================================
// Rendering Layer: PDF -> Image Bytes (PNG)
// =========================================================

func renderPdfBytesToImageBytes(pdfContent []byte) ([]byte, error) {
	doc, err := fitz.NewFromMemory(pdfContent)
	if err != nil {
		return nil, fmt.Errorf("fitz load pdf failed: %w", err)
	}
	defer doc.Close()

	if doc.NumPage() == 0 {
		return nil, fmt.Errorf("pdf has no pages")
	}

	// 1. Render the first page as a raw image
	img, err := doc.Image(0)
	if err != nil {
		return nil, fmt.Errorf("render page 0 failed: %w", err)
	}

	// 2. Get original image dimensions
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 3. Crop: Keep only the top half (1/2)
	// Use CropAnchor with the Top anchor point, width unchanged, height set to half
	topHalfImg := imaging.CropAnchor(img, width, height/2, imaging.Top)

	// 4. Smart Scaling: Fit the cropped top half into an 800x800 container
	// Note: Fit maintains aspect ratio. Since only the top half is used, the resulting thumbnail will show header content more clearly than a full-page thumbnail.
	dstImage := imaging.Fit(topHalfImg, 800, 800, imaging.Lanczos)

	// 5. Encode the image to PNG byte stream
	var buf bytes.Buffer
	err = imaging.Encode(&buf, dstImage, imaging.PNG)
	if err != nil {
		return nil, fmt.Errorf("encode image to png failed: %w", err)
	}

	return buf.Bytes(), nil
}

// =========================================================
// Conversion Layer: Any -> PDF (using go-pdf/fpdf)
// =========================================================

func newPdfGenerator() (*fpdf.Fpdf, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 20, 20)
	pdf.AddPage()

	// Must have simhei.ttf otherwise Chinese characters cannot be displayed
	// Note: Since this is now pure in-memory processing, the fontPath here still relies on the local filesystem
	// If complete file-independence is needed, change to AddUTF8FontFromBytes and bundle the font file into the program
	fontPath := "internal/pkg/thumbnail/fonts/SimHei.ttf"

	// Add simple error tolerance here to prevent panic. In practice, ensure the font exists.
	pdf.AddUTF8Font("SimHei", "", fontPath)
	pdf.SetFont("SimHei", "", 12)

	return pdf, nil
}

// DOCX to PDF
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

// XLSX to PDF
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

// CSV to PDF
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

// HTML to PDF
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

// Text to PDF
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

// OFD to PDF
func convertOfdToPdf(content []byte) ([]byte, error) {
	r, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		// Fallback
		return convertTextToPdf(string(content))
	}

	var sb strings.Builder
	fileFound := false

	for _, f := range r.File {
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

	if !fileFound || sb.Len() == 0 {
		return convertTextToPdf(string(content))
	}

	return convertTextToPdf(sb.String())
}

// =========================================================
// Helper Functions
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

func sanitizeText(s string) string {
	return strings.Map(func(r rune) rune {
		if r < 32 {
			if r == 10 || r == 13 || r == 9 {
				return r
			}
			return -1
		}
		if r > 0xFFFF {
			return -1
		}
		return r
	}, s)
}
