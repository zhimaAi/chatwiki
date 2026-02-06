// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"regexp"
	"strings"

	"github.com/beevik/etree"
)

// OMMLToLatexConverter is an OMML to LaTeX converter
type OMMLToLatexConverter struct {
	// Define mapping tables needed for conversion
	charMap      map[string]string
	funcMap      map[string]string
	fracTypes    map[string]string
	delimiterMap map[string]string
	accentMap    map[string]string
	barPosMap    map[string]string
}

// NewOMMLToLatexConverter creates a new converter instance
func NewOMMLToLatexConverter() *OMMLToLatexConverter {
	converter := &OMMLToLatexConverter{
		charMap:      make(map[string]string),
		funcMap:      make(map[string]string),
		fracTypes:    make(map[string]string),
		delimiterMap: make(map[string]string),
		accentMap:    make(map[string]string),
		barPosMap:    make(map[string]string),
	}

	// Initialize mapping tables
	converter.initMaps()
	return converter
}

// initMaps initializes various symbol mapping tables
func (c *OMMLToLatexConverter) initMaps() {
	c.charMap = getCharMap()
	// Function name mapping
	c.funcMap = getFuncMap()

	// Fraction type mapping
	c.fracTypes["bar"] = "\\frac{%s}{%s}"
	c.fracTypes["skw"] = "^{%s}/_{%s}"
	c.fracTypes["noBar"] = "\\genfrac{}{}{0pt}{}{%s}{%s}"
	c.fracTypes["lin"] = "{%s}/{%s}"

	// Bracket mapping
	c.delimiterMap["("] = "("
	c.delimiterMap[")"] = ")"
	c.delimiterMap["["] = "["
	c.delimiterMap["]"] = "]"
	c.delimiterMap["{"] = "\\{"
	c.delimiterMap["}"] = "\\}"
	c.delimiterMap["|"] = "|"
	c.delimiterMap["∥"] = "\\|"
	c.delimiterMap["⟨"] = "\\langle"
	c.delimiterMap["⟩"] = "\\rangle"
	c.delimiterMap["⌊"] = "\\lfloor"
	c.delimiterMap["⌋"] = "\\rfloor"
	c.delimiterMap["⌈"] = "\\lceil"
	c.delimiterMap["⌉"] = "\\rceil"
	c.delimiterMap["⟦"] = "\\llbracket"
	c.delimiterMap["⟧"] = "\\rrbracket"

	// Accent symbol mapping
	c.accentMap["\u0302"] = "\\hat{%s}"                // ˆ
	c.accentMap["\u0303"] = "\\tilde{%s}"              // ˜
	c.accentMap["\u0304"] = "\\bar{%s}"                // ¯
	c.accentMap["\u0305"] = "\\overbar{%s}"            // ‾
	c.accentMap["\u0306"] = "\\breve{%s}"              // ˘
	c.accentMap["\u0307"] = "\\dot{%s}"                // ˙
	c.accentMap["\u0308"] = "\\ddot{%s}"               // ¨
	c.accentMap["\u030c"] = "\\check{%s}"              // ˇ
	c.accentMap["\u20d7"] = "\\vec{%s}"                // vector arrow
	c.accentMap["\u20e1"] = "\\overleftrightarrow{%s}" // bidirectional arrow
	c.accentMap["\u23dc"] = "\\overparen{%s}"          // arc superscript
	c.accentMap["\u23de"] = "\\overbrace{%s}"          // overbrace
	c.accentMap["\u23dd"] = "\\underparen{%s}"         // arc subscript
	c.accentMap["\u23df"] = "\\underbrace{%s}"         // underbrace
	c.accentMap["\u030a"] = "\\mathring{%s}"           // °
	c.accentMap["\u0300"] = "\\grave{%s}"              // `
	c.accentMap["\u0301"] = "\\acute{%s}"              // ´
	c.accentMap["\u030b"] = "\\H{%s}"                  // ˝
	c.accentMap["\u0323"] = "\\underdot{%s}"           // ̣
	c.accentMap["\u0327"] = "\\c{%s}"                  // ¸
	c.accentMap["\u0328"] = "\\k{%s}"                  // ˛

	// Overline/underline position mapping
	c.barPosMap["top"] = "\\overline{%s}"
	c.barPosMap["bot"] = "\\underline{%s}"
}

// ConvertOMMLToLatex converts OMML XML string to LaTeX
func (c *OMMLToLatexConverter) ConvertOMMLToLatex(ommlXML string) (string, error) {
	result := ConvertOMMLToLatexString(ommlXML)
	return result, nil
}

// ConvertOMMLToLatexString converts OMML XML string to LaTeX
func ConvertOMMLToLatexString(oMathXML string) string {
	// Preprocessing: remove line breaks and extra spaces between XML tags, but preserve content within tags
	doc := etree.NewDocument()
	if err := doc.ReadFromString(oMathXML); err != nil {
		return ""
	}
	root := doc.Root()
	result := processOMathElement(root)
	// Now need to clean up extra spaces produced by XML parsing, this is necessary because XML parsing introduces extra spaces
	result = cleanupSpaces(result)
	return result
}

// cleanupSpaces cleans up extra spaces while preserving meaningful spaces
func cleanupSpaces(latex string) string {
	// Remove redundant consecutive spaces, but preserve single space
	latex = regexp.MustCompile(`\s+`).ReplaceAllString(latex, " ")
	// Remove leading and trailing spaces
	return strings.TrimSpace(latex)
}

// processOMathElement processes math elements
func processOMathElement(element *etree.Element) string {
	if element.Tag == "oMath" {
		var result strings.Builder
		for _, child := range element.Child {
			if elem, ok := child.(*etree.Element); ok {
				result.WriteString(processElement(elem))
			} else if charData, ok := child.(*etree.CharData); ok {
				result.WriteString(string(charData.Data))
			}
		}
		return result.String()
	}
	return processElement(element)
}

// processElement processes an element
func processElement(element *etree.Element) string {
	switch element.Tag {
	case "r":
		return processRun(element)
	case "f":
		return processFraction(element)
	case "func":
		return processFunction(element)
	case "sSub":
		return processSubscript(element)
	case "sSup":
		return processSuperscript(element)
	case "sSubSup":
		return processSubSup(element)
	case "acc":
		return processAccent(element)
	case "bar":
		return processBar(element)
	case "d":
		return processDelimiter(element)
	case "rad":
		return processRadical(element)
	case "eqArr":
		return processEqArray(element)
	case "limLow":
		return processLimLow(element)
	case "limUpp":
		return processLimUpp(element)
	case "m":
		return processMatrix(element)
	case "mr":
		return processMatrixRow(element)
	case "nary":
		return processNary(element)
	case "groupChr":
		return processGroupChar(element)
	case "phant":
		return processPhantom(element)
	case "sPre":
		return processSPre(element)
	case "t":
		return processText(element)
	default:
		var result strings.Builder
		for _, child := range element.Child {
			if elem, ok := child.(*etree.Element); ok {
				result.WriteString(processElement(elem))
			} else if charData, ok := child.(*etree.CharData); ok {
				result.WriteString(string(charData.Data))
			}
		}
		return result.String()
	}
}

// processRun processes text runs
func processRun(element *etree.Element) string {
	var result strings.Builder
	for _, child := range element.Child {
		if elem, ok := child.(*etree.Element); ok {
			result.WriteString(processElement(elem))
		} else if charData, ok := child.(*etree.CharData); ok {
			result.WriteString(string(charData.Data))
		}
	}
	return result.String()
}

// processText processes text elements
func processText(element *etree.Element) string {
	text := strings.TrimSpace(element.Text())

	// Check if there is a style attribute (e.g., p indicates function name)
	var style string
	for _, attr := range element.Attr {
		if attr.Key == "val" && attr.Space == "m:sty" {
			style = attr.Value
			break
		}
	}

	// If it's function name style, handle specially
	if style == "p" {
		// Check if it's a known function name
		if latexFunc, exists := getFuncMap()[text]; exists {
			return latexFunc
		}
	}

	// Try to map characters to LaTeX commands
	converted := ""
	for _, char := range text {
		charStr := string(char)
		if charStr == "⁡" { // Function application character, ignore
			continue
		}
		if latexChar, exists := getCharMap()[charStr]; exists {
			converted += latexChar + " "
		} else {
			converted += charStr
		}
	}
	// Note: Don't TrimSpace here, otherwise segment concatenation like "m×n" will become "m\\timesn".
	// Unified space normalization is handled by cleanupSpaces at the end of ConvertOMMLToLatexString.
	return converted
}

// processFraction processes fractions
func processFraction(element *etree.Element) string {
	var numerator, denominator string

	for _, child := range element.Child {
		if elem, ok := child.(*etree.Element); ok {
			switch elem.Tag {
			case "num":
				for _, numChild := range elem.Child {
					if numElem, ok := numChild.(*etree.Element); ok {
						numerator += processElement(numElem)
					} else if charData, ok := numChild.(*etree.CharData); ok {
						numerator += string(charData.Data)
					}
				}
			case "den":
				for _, denChild := range elem.Child {
					if denElem, ok := denChild.(*etree.Element); ok {
						denominator += processElement(denElem)
					} else if charData, ok := denChild.(*etree.CharData); ok {
						denominator += string(charData.Data)
					}
				}
			}
		}
	}
	return "\\frac{" + numerator + "}{" + denominator + "}"
}

// processFunction processes functions
func processFunction(element *etree.Element) string {
	var functionName, functionArg string

	for _, child := range element.Child {
		if elem, ok := child.(*etree.Element); ok {
			switch elem.Tag {
			case "fName":
				for _, nameChild := range elem.Child {
					if nameElem, ok := nameChild.(*etree.Element); ok {
						functionName += processElement(nameElem)
					} else if charData, ok := nameChild.(*etree.CharData); ok {
						functionName += string(charData.Data)
					}
				}
			case "e":
				for _, argChild := range elem.Child {
					if argElem, ok := argChild.(*etree.Element); ok {
						functionArg += processElement(argElem)
					} else if charData, ok := argChild.(*etree.CharData); ok {
						functionArg += string(charData.Data)
					}
				}
			}
		}
	}
	return strings.TrimSpace(functionName) + "{" + strings.TrimSpace(functionArg) + "}"
}

// processSubscript processes subscripts
func processSubscript(element *etree.Element) string {
	var base, subscript string

	for _, child := range element.Child {
		if elem, ok := child.(*etree.Element); ok {
			switch elem.Tag {
			case "e":
				for _, baseChild := range elem.Child {
					if baseElem, ok := baseChild.(*etree.Element); ok {
						base += processElement(baseElem)
					} else if charData, ok := baseChild.(*etree.CharData); ok {
						base += string(charData.Data)
					}
				}
			case "sub":
				for _, subChild := range elem.Child {
					if subElem, ok := subChild.(*etree.Element); ok {
						subscript += processElement(subElem)
					} else if charData, ok := subChild.(*etree.CharData); ok {
						subscript += string(charData.Data)
					}
				}
			}
		}
	}
	return strings.TrimSpace(base) + "_{" + strings.TrimSpace(subscript) + "}"
}

// processSuperscript processes superscripts
func processSuperscript(element *etree.Element) string {
	var base, superscript string

	for _, child := range element.Child {
		if elem, ok := child.(*etree.Element); ok {
			switch elem.Tag {
			case "e":
				for _, baseChild := range elem.Child {
					if baseElem, ok := baseChild.(*etree.Element); ok {
						base += processElement(baseElem)
					} else if charData, ok := baseChild.(*etree.CharData); ok {
						base += string(charData.Data)
					}
				}
			case "sup":
				for _, supChild := range elem.Child {
					if supElem, ok := supChild.(*etree.Element); ok {
						superscript += processElement(supElem)
					} else if charData, ok := supChild.(*etree.CharData); ok {
						superscript += string(charData.Data)
					}
				}
			}
		}
	}
	return base + "^{" + strings.TrimSpace(superscript) + "}"
}

// processSubSup processes subscripts and superscripts
func processSubSup(element *etree.Element) string {
	var base, subscript, superscript string

	for _, child := range element.Child {
		if elem, ok := child.(*etree.Element); ok {
			switch elem.Tag {
			case "e":
				for _, baseChild := range elem.Child {
					if baseElem, ok := baseChild.(*etree.Element); ok {
						base += processElement(baseElem)
					} else if charData, ok := baseChild.(*etree.CharData); ok {
						base += string(charData.Data)
					}
				}
			case "sub":
				for _, subChild := range elem.Child {
					if subElem, ok := subChild.(*etree.Element); ok {
						subscript += processElement(subElem)
					} else if charData, ok := subChild.(*etree.CharData); ok {
						subscript += string(charData.Data)
					}
				}
			case "sup":
				for _, supChild := range elem.Child {
					if supElem, ok := supChild.(*etree.Element); ok {
						superscript += processElement(supElem)
					} else if charData, ok := supChild.(*etree.CharData); ok {
						superscript += string(charData.Data)
					}
				}
			}
		}
	}
	return base + "_{" + strings.TrimSpace(subscript) + "}^{" + strings.TrimSpace(superscript) + "}"
}

// processAccent processes accent marks
func processAccent(element *etree.Element) string {
	var base string
	var accentChar string

	for _, child := range element.Child {
		if elem, ok := child.(*etree.Element); ok {
			switch elem.Tag {
			case "e":
				for _, baseChild := range elem.Child {
					if baseElem, ok := baseChild.(*etree.Element); ok {
						base += processElement(baseElem)
					} else if charData, ok := baseChild.(*etree.CharData); ok {
						base += string(charData.Data)
					}
				}
			case "accPr": // Process accent mark attributes
				for _, accPrChild := range elem.Child {
					if accPrElem, ok := accPrChild.(*etree.Element); ok {
						if accPrElem.Tag == "chr" {
							// Find the val attribute of chr element
							for _, attr := range accPrElem.Attr {
								if attr.Key == "val" {
									accentChar = attr.Value
									break
								}
							}
						}
					}
				}
			}
		}
	}

	// If accent character found, generate different LaTeX formats based on character type
	if accentChar != "" {
		switch accentChar {
		case "⇀": // Right harpoon arrow
			return "\\mathord{ \\buildrel{ \\lower3pt \\hbox{$ \\scriptscriptstyle \\rightharpoonup$}} \\over " + base + "}"
		case "↼": // Left harpoon arrow
			return "\\mathord{ \\buildrel{ \\lower3pt \\hbox{$ \\scriptscriptstyle \\leftharpoonup$}} \\over " + base + "}"
		case "⇀↽": // Bidirectional arrow
			return "\\mathord{ \\buildrel{ \\lower3pt \\hbox{$ \\scriptscriptstyle \\rightleftharpoons$}} \\over " + base + "}"
		case "→", "⟶": // Right arrow variant
			return "\\mathord{ \\buildrel{ \\lower3pt \\hbox{$ \\scriptscriptstyle \\rightarrow$}} \\over " + base + "}"
		case "←", "⟵": // Left arrow variant
			return "\\mathord{ \\buildrel{ \\lower3pt \\hbox{$ \\scriptscriptstyle \\leftarrow$}} \\over " + base + "}"
		case "↑":
			return "\\mathord{ \\buildrel{ \\lower3pt \\hbox{$ \\scriptscriptstyle \\uparrow$}} \\over " + base + "}"
		case "↓":
			return "\\mathord{ \\buildrel{ \\lower3pt \\hbox{$ \\scriptscriptstyle \\downarrow$}} \\over " + base + "}"
		case "^", "ˆ": // Hat symbol
			return "\\hat{" + base + "}"
		case "˜", "~": // Tilde
			return "\\tilde{" + base + "}"
		case "˙": // Dot
			return "\\dot{" + base + "}"
		case "¨": // Double dot
			return "\\ddot{" + base + "}"
		case "‾", "¯": // Overline
			return "\\overline{" + base + "}"
		case "⃗": // Vector arrow
			return "\\vec{" + base + "}"
		default:
			// If it's another symbol, try to find mapping
			if latexCmd, exists := getCharMap()[accentChar]; exists {
				return "\\" + latexCmd + "{" + base + "}"
			}
			// By default, return format with accent symbol
			return "\\accent{" + accentChar + "}{" + base + "}"
		}
	}

	// If no accent character found, return base content directly
	return base
}

// processBar processes overlines/underlines
func processBar(element *etree.Element) string {
	var base string

	for _, child := range element.Child {
		if elem, ok := child.(*etree.Element); ok {
			switch elem.Tag {
			case "e":
				for _, baseChild := range elem.Child {
					if baseElem, ok := baseChild.(*etree.Element); ok {
						base += processElement(baseElem)
					} else if charData, ok := baseChild.(*etree.CharData); ok {
						base += string(charData.Data)
					}
				}
			}
		}
	}
	return base
}

// processDelimiter processes delimiters
func processDelimiter(element *etree.Element) string {
	var content string
	var leftBracket, rightBracket string = "", ""

	// Check if there's a dPr child element to define left/right brackets
	for _, child := range element.Child {
		if elem, ok := child.(*etree.Element); ok {
			if elem.Tag == "dPr" {
				// Iterate through dPr child elements, looking for begChr and endChr
				for _, dPrChild := range elem.Child {
					if dPrElem, ok := dPrChild.(*etree.Element); ok {
						if dPrElem.Tag == "begChr" {
							// Get the val attribute value of begChr
							for _, attr := range dPrElem.Attr {
								if attr.Key == "val" {
									val := attr.Value
									leftBracket = convertBracket(val, true) // true indicates left bracket
									break
								}
							}
						} else if dPrElem.Tag == "endChr" {
							// Get the val attribute value of endChr
							for _, attr := range dPrElem.Attr {
								if attr.Key == "val" {
									val := attr.Value
									rightBracket = convertBracket(val, false) // false indicates right bracket
									break
								}
							}
						}
					}
				}
			} else if elem.Tag == "e" {
				for _, contentChild := range elem.Child {
					if contentElem, ok := contentChild.(*etree.Element); ok {
						content += processElement(contentElem)
					} else if charData, ok := contentChild.(*etree.CharData); ok {
						content += string(charData.Data)
					}
				}
			}
		}
	}

	// If no special left/right bracket definition found, default to parentheses
	if leftBracket == "" && rightBracket == "" {
		leftBracket = "\\left("
		rightBracket = "\\right)"
	}

	// If left/right brackets found, use them to wrap content
	if leftBracket != "" || rightBracket != "" {
		return leftBracket + content + rightBracket
	}

	// If no special left/right bracket definition found, return content directly
	return content
}

// convertBracket converts bracket values to LaTeX format
func convertBracket(val string, isLeft bool) string {
	// Empty value maps to dot (.)
	if val == "" {
		if isLeft {
			return "\\left."
		} else {
			return "\\right."
		}
	}

	// Determine prefix based on whether it's a left bracket
	var prefix string
	if isLeft {
		prefix = "\\left"
	} else {
		prefix = "\\right"
	}

	// Special handling for curly braces: { and }
	switch val {
	case "{":
		return prefix + "\\{"
	case "}":
		return prefix + "\\}"
	default:
		// Other symbols add prefix as-is
		return prefix + val
	}
}

// processRadical processes radicals
func processRadical(element *etree.Element) string {
	var base, degree string

	for _, child := range element.Child {
		if elem, ok := child.(*etree.Element); ok {
			switch elem.Tag {
			case "e":
				for _, baseChild := range elem.Child {
					if baseElem, ok := baseChild.(*etree.Element); ok {
						base += processElement(baseElem)
					} else if charData, ok := baseChild.(*etree.CharData); ok {
						base += string(charData.Data)
					}
				}
			case "deg":
				for _, degChild := range elem.Child {
					if degElem, ok := degChild.(*etree.Element); ok {
						degree += processElement(degElem)
					} else if charData, ok := degChild.(*etree.CharData); ok {
						degree += string(charData.Data)
					}
				}
			}
		}
	}

	if degree != "" {
		return "\\sqrt[" + strings.TrimSpace(degree) + "]{" + strings.TrimSpace(base) + "}"
	}
	return "\\sqrt{" + strings.TrimSpace(base) + "}"
}

// processEqArray processes equation arrays
func processEqArray(element *etree.Element) string {
	var items []string
	var hasRightBrace bool // Mark whether an independent right brace is encountered

	for _, child := range element.Child {
		if elem, ok := child.(*etree.Element); ok {
			switch elem.Tag {
			case "e":
				var itemContent strings.Builder
				for _, itemChild := range elem.Child {
					if itemElem, ok := itemChild.(*etree.Element); ok {
						itemContent.WriteString(processElement(itemElem))
					} else if charData, ok := itemChild.(*etree.CharData); ok {
						itemContent.WriteString(string(charData.Data))
					}
				}
				if itemContent.Len() > 0 {
					items = append(items, itemContent.String())
				}
			case "r":
				// Check if it's a standalone right brace, if so mark it instead of processing as normal text
				for _, rChild := range elem.Child {
					if charData, ok := rChild.(*etree.CharData); ok {
						text := strings.TrimSpace(string(charData.Data))
						if text == "}" {
							hasRightBrace = true
						}
					}
				}
			}
		}
	}

	if len(items) > 0 {
		// Build matrix form equation system
		matrixContent := strings.Join(items, "\\\\\n")
		result := "\\left.\\begin{matrix} " + matrixContent + " \\end{matrix}\\right\\}"
		// If an independent right brace is encountered, don't add it
		if !hasRightBrace {
			return result
		} else {
			// If there's an independent right brace, only return the matrix part, let subsequent processing handle the right brace
			return "\\left.\\begin{matrix} " + matrixContent + " \\end{matrix}\\right\\}"
		}
	}
	return ""
}

// processLimLow processes lower limits
func processLimLow(element *etree.Element) string {
	var operator, limit string

	for _, child := range element.Child {
		if elem, ok := child.(*etree.Element); ok {
			switch elem.Tag {
			case "e":
				for _, opChild := range elem.Child {
					if opElem, ok := opChild.(*etree.Element); ok {
						operator += processElement(opElem)
					} else if charData, ok := opChild.(*etree.CharData); ok {
						operator += string(charData.Data)
					}
				}
			case "lim":
				for _, limChild := range elem.Child {
					if limElem, ok := limChild.(*etree.Element); ok {
						limit += processElement(limElem)
					} else if charData, ok := limChild.(*etree.CharData); ok {
						limit += string(charData.Data)
					}
				}
			}
		}
	}

	// If operator is "lim", use \lim \limits format
	if strings.TrimSpace(operator) == "lim" || strings.TrimSpace(operator) == "\\lim" {
		return "\\lim \\limits_{" + limit + "}"
	}

	return operator + "_{" + limit + "}"
}

// processLimUpp processes upper limits
func processLimUpp(element *etree.Element) string {
	var operator, limit string

	for _, child := range element.Child {
		if elem, ok := child.(*etree.Element); ok {
			switch elem.Tag {
			case "e":
				for _, opChild := range elem.Child {
					if opElem, ok := opChild.(*etree.Element); ok {
						operator += processElement(opElem)
					} else if charData, ok := opChild.(*etree.CharData); ok {
						operator += string(charData.Data)
					}
				}
			case "lim":
				for _, limChild := range elem.Child {
					if limElem, ok := limChild.(*etree.Element); ok {
						limit += processElement(limElem)
					} else if charData, ok := limChild.(*etree.CharData); ok {
						limit += string(charData.Data)
					}
				}
			}
		}
	}

	// Special handling for vectors and unit vectors
	if limit == "^" {
		return "\\hat{" + operator + "}"
	} else if strings.Contains(limit, "⇀") || strings.Contains(limit, "→") {
		return "\\vec{" + operator + "}"
	}

	return "\\overset{" + limit + "}{" + operator + "}"
}

// processMatrix processes matrices
func processMatrix(element *etree.Element) string {
	var rows []string

	for _, child := range element.Child {
		if elem, ok := child.(*etree.Element); ok {
			if elem.Tag == "mr" { // matrix row
				var rowContent strings.Builder
				for _, rowChild := range elem.Child {
					if rowElem, ok := rowChild.(*etree.Element); ok {
						if rowContent.Len() > 0 {
							rowContent.WriteString(" & ")
						}
						rowContent.WriteString(processElement(rowElem))
					} else if charData, ok := rowChild.(*etree.CharData); ok {
						if rowContent.Len() > 0 {
							rowContent.WriteString(" & ")
						}
						rowContent.WriteString(string(charData.Data))
					}
				}
				if rowContent.Len() > 0 {
					rows = append(rows, rowContent.String())
				}
			}
		}
	}

	if len(rows) > 0 {
		matrixContent := strings.Join(rows, " \\\\\n ")
		return "\\begin{matrix} " + matrixContent + " \\end{matrix}"
	}
	return ""
}

// processMatrixRow processes matrix rows
func processMatrixRow(element *etree.Element) string {
	var elements []string

	for _, child := range element.Child {
		if elem, ok := child.(*etree.Element); ok {
			if elem.Tag == "e" {
				var elementContent strings.Builder
				for _, elemChild := range elem.Child {
					if elemElem, ok := elemChild.(*etree.Element); ok {
						elementContent.WriteString(processElement(elemElem))
					} else if charData, ok := elemChild.(*etree.CharData); ok {
						elementContent.WriteString(string(charData.Data))
					}
				}
				if elementContent.Len() > 0 {
					elements = append(elements, elementContent.String())
				}
			}
		}
	}

	return strings.Join(elements, "&")
}

// processNary processes n-ary operators
func processNary(element *etree.Element) string {
	var sub, sup, content string
	naryChar := "\\sum" // Default to summation symbol

	for _, child := range element.Child {
		if elem, ok := child.(*etree.Element); ok {
			switch elem.Tag {
			case "naryPr": // n-ary operator properties
				naryChar = processNaryPr(elem)
			case "sub":
				for _, subChild := range elem.Child {
					if subElem, ok := subChild.(*etree.Element); ok {
						sub += processElement(subElem)
					} else if charData, ok := subChild.(*etree.CharData); ok {
						sub += string(charData.Data)
					}
				}
			case "sup":
				for _, supChild := range elem.Child {
					if supElem, ok := supChild.(*etree.Element); ok {
						sup += processElement(supElem)
					} else if charData, ok := supChild.(*etree.CharData); ok {
						sup += string(charData.Data)
					}
				}
			case "e":
				for _, contentChild := range elem.Child {
					if contentElem, ok := contentChild.(*etree.Element); ok {
						content += processElement(contentElem)
					} else if charData, ok := contentChild.(*etree.CharData); ok {
						content += string(charData.Data)
					}
				}
			}
		}
	}

	// According to test expectations, don't use \limits form, use standard subscript/superscript form
	if sub != "" && sup != "" {
		return naryChar + "_{" + strings.TrimSpace(sub) + "}^{" + strings.TrimSpace(sup) + "} " + strings.TrimSpace(content)
	} else if sub != "" {
		return naryChar + "_{" + strings.TrimSpace(sub) + "} " + strings.TrimSpace(content)
	} else if sup != "" {
		return naryChar + "^{" + strings.TrimSpace(sup) + "} " + strings.TrimSpace(content)
	}
	return naryChar + " " + strings.TrimSpace(content)
}

// processNaryPr processes n-ary operator properties
func processNaryPr(element *etree.Element) string {
	// First check if there's a chr element defining the operator character
	for _, child := range element.Child {
		if elem, ok := child.(*etree.Element); ok {
			if elem.Tag == "chr" {
				for _, attr := range elem.Attr {
					if attr.Key == "val" {
						// If it's not a special symbol, try to find the corresponding LaTeX command
						if latexCmd, exists := getCharMap()[attr.Value]; exists {
							return latexCmd + " "
						}
						return attr.Value
					}
				}
			}
		}
	}

	// If no chr element found, infer from other properties
	// Check limLoc attribute to help determine operator type
	for _, child := range element.Child {
		if elem, ok := child.(*etree.Element); ok {
			if elem.Tag == "limLoc" {
				for _, attr := range elem.Attr {
					if attr.Key == "val" {
						// For limLoc=subSup, usually indicates integral
						if attr.Value == "subSup" {
							return "\\int" + " "
						}
					}
				}
			}
		}
	}

	// Check if it's an integral-related context
	// If naryPr has subHide=1, it might indicate integral symbol
	for _, child := range element.Child {
		if elem, ok := child.(*etree.Element); ok {
			if elem.Tag == "subHide" {
				for _, attr := range elem.Attr {
					if attr.Key == "val" && attr.Value == "1" {
						return "\\int" + " "
					}
				}
			}
			if elem.Tag == "supHide" {
				for _, attr := range elem.Attr {
					if attr.Key == "val" && attr.Value == "1" {
						return "\\int" + " "
					}
				}
			}
		}
	}

	return "\\sum " // Default value
}

// processGroupChar processes group characters
func processGroupChar(element *etree.Element) string {
	var content string

	for _, child := range element.Child {
		if elem, ok := child.(*etree.Element); ok {
			switch elem.Tag {
			case "e":
				for _, contentChild := range elem.Child {
					if contentElem, ok := contentChild.(*etree.Element); ok {
						content += processElement(contentElem)
					} else if charData, ok := contentChild.(*etree.CharData); ok {
						content += string(charData.Data)
					}
				}
			}
		}
	}
	return content
}

// processPhantom processes phantom elements
func processPhantom(element *etree.Element) string {
	var content string

	for _, child := range element.Child {
		if elem, ok := child.(*etree.Element); ok {
			for _, contentChild := range elem.Child {
				if contentElem, ok := contentChild.(*etree.Element); ok {
					content += processElement(contentElem)
				} else if charData, ok := contentChild.(*etree.CharData); ok {
					content += string(charData.Data)
				}
			}
		}
	}
	return "\\phantom{" + content + "}"
}

// processSPre processes prescripts (pre-subscripts and pre-superscripts)
func processSPre(element *etree.Element) string {
	var preSub, preSup, base string

	for _, child := range element.Child {
		if elem, ok := child.(*etree.Element); ok {
			switch elem.Tag {
			case "sub":
				for _, subChild := range elem.Child {
					if subElem, ok := subChild.(*etree.Element); ok {
						preSub += processElement(subElem)
					} else if charData, ok := subChild.(*etree.CharData); ok {
						preSub += string(charData.Data)
					}
				}
			case "sup":
				for _, supChild := range elem.Child {
					if supElem, ok := supChild.(*etree.Element); ok {
						preSup += processElement(supElem)
					} else if charData, ok := supChild.(*etree.CharData); ok {
						preSup += string(charData.Data)
					}
				}
			case "e":
				for _, baseChild := range elem.Child {
					if baseElem, ok := baseChild.(*etree.Element); ok {
						base += processElement(baseElem)
					} else if charData, ok := baseChild.(*etree.CharData); ok {
						base += string(charData.Data)
					}
				}
			}
		}
	}

	result := ""
	if preSub != "" {
		result += "_{" + preSub + "}"
	}
	if preSup != "" {
		result += "^{" + preSup + "}"
	}
	result += base
	return result
}

// getCharMap is a helper function that returns the character mapping table
func getCharMap() map[string]string {
	charMap := make(map[string]string)
	// Greek letters and special symbols mapping
	greekLetters := map[string]string{
		"α": "\\alpha", "β": "\\beta", "γ": "\\gamma", "ε": "\\varepsilon",
		"ζ": "\\zeta", "η": "\\eta", "ι": "\\iota", "κ": "\\kappa",
		"λ": "\\lambda", "μ": "\\mu", "ν": "\\nu", "ξ": "\\xi", "ο": "\\omicron",
		"π": "\\pi", "ρ": "\\rho", "σ": "\\sigma", "τ": "\\tau", "υ": "\\upsilon",
		"φ": "\\phi", "χ": "\\chi", "ψ": "\\psi", "ω": "\\omega", "∂": "\\partial",
		"϶": "\\epsilon", "ϰ": "\\varkappa", "ϕ": "\\varphi",
		"ϱ": "\\varrho", "ϖ": "\\varpi", "ϑ": "\\vartheta", "θ": "\\theta",
		"ϝ": "\\digamma", "Ϝ": "\\Digamma",
	}

	// Merge into charMap
	for k, v := range greekLetters {
		charMap[k] = v
	}

	// Arrows
	arrows := map[string]string{
		"→": "\\rightarrow", "←": "\\leftarrow", "↑": "\\uparrow", "↓": "\\downarrow",
		"↔": "\\leftrightarrow", "↕": "\\updownarrow", "↖": "\\nwarrow", "↗": "\\nearrow",
		"↘": "\\searrow", "↙": "\\swarrow", "↚": "\\nleftarrow", "↛": "\\nrightarrow",
		"↜": "\\twoheadleftarrow", "↝": "\\twoheadrightarrow", "↞": "\\leftarrowtail",
		"↟": "\\upharpoonleft", "↠": "\\rightarrowtail", "⇀": "\\rightharpoonup", "⇁": "\\rightharpoondown",
		"⇄": "\\leftrightarrows", "⇅": "\\updownarrows", "⇆": "\\leftleftarrows", "⇇": "\\rightrightarrows",
		"⇉": "\\rightrightarrows", "⇋": "\\leftrightharpoons", "↼": "\\leftharpoonup", "↽": "\\leftharpoondown",
		"⇂": "\\downharpoonright", "⇃": "\\downharpoonleft", "↾": "\\upharpoonright",
		"↿": "\\\\upharpoonleft", "↷": "\\curvearrowright", "↻": "\\circlearrowright",
	}
	for k, v := range arrows {
		charMap[k] = v
	}

	// Relation symbols - add appropriate spaces when mapping
	relSymbols := map[string]string{
		"≤": "\\leq", "≥": "\\geq", "≠": "\\neq", "∼": "\\sim",
		"≡": "\\equiv", "∝": "\\propto", "⊂": "\\subset", "⊃": "\\supset",
		"⊆": "\\subseteq", "⊇": "\\supseteq", "∈": "\\in", "∉": "\\notin",
		"∋": "\\ni", "∀": "\\forall", "∃": "\\exists", "¬": "\\neg",
		"⊥": "\\perp", "∠": "\\angle", "∨": "\\vee", "∧": "\\wedge",
		"⇒": "\\Rightarrow", "⇐": "\\Leftarrow", "⇔": "\\Leftrightarrow",
		"→": "\\to", "←": "\\gets", "↔": "\\leftrightarrow",
		"≫": "\\gg", "≪": "\\ll",
		"≻": "\\succ", "≺": "\\prec", "≼": "\\preceq", "≽": "\\succeq",
		"⊤": "\\top", "⊢": "\\vdash", "⊨": "\\models", "⌣": "\\smile",
		"⌢": "\\frown", "◃": "\\triangleleft", "▹": "\\triangleright",
		"⋪": "\\ntriangleleft", "⋫": "\\ntriangleright", "⊴": "\\unlhd", "⊵": "\\unrhd",
		"⊏": "\\sqsubset", "⊐": "\\sqsupset", "⊑": "\\sqsubseteq", "⊒": "\\sqsupseteq",
		"⊊": "\\subsetneq", "⊋": "\\supsetneq", "⊈": "\\nsubseteq", "⊉": "\\nsupseteq",
		"∤": "\\nmid", "∥": "\\parallel", "∦": "\\nparallel",
		"≂": "\\eqsim",
		"≑": "\\Doteq", "≒": "\\fallingdotseq", "≜": "\\triangleq",
		"≖": "\\eqcirc", "≗": "\\circeq",
		// The following are duplicates, placed at the end so they will override previous identical keys
		">": "\\ge", "<": "\\le", "≎": "\\Bumpeq",
		"≈": "\\approx", "≐": "\\doteq",
		"≓": "\\risingdotseq", "≏": "\\bumpeq",
	}
	for k, v := range relSymbols {
		charMap[k] = v
	}

	// Operators
	operators := map[string]string{
		"±": "\\pm", "∓": "\\mp", "×": "\\times", "÷": "\\div", "∗": "\\ast",
		"⋆": "\\star", "∘": "\\circ", "∙": "\\bullet", "∝": "\\propto",
		"∞": "\\infty", "∅": "\\emptyset", "∇": "\\nabla", "√": "\\surd",
		"∑": "\\sum", "∏": "\\prod", "∐": "\\coprod", "∫": "\\int",
		"∮": "\\oint", "⋂": "\\bigcap", "⋃": "\\bigcup", "⨿": "\\amalg",
		"⊎": "\\uplus", "⊔": "\\sqcup", "⊓": "\\sqcap", "∧": "\\wedge",
		"∨": "\\vee", "⋁": "\\bigvee", "⋀": "\\bigwedge",
		"∖": "\\setminus",
		"≀": "\\wr", "⋄": "\\diamond", "⋈": "\\bowtie", "⋉": "\\ltimes",
		"⋊": "\\rtimes", "⋋": "\\leftthreetimes", "⋌": "\\rightthreetimes",
	}
	for k, v := range operators {
		charMap[k] = v
	}

	return charMap
}

// getFuncMap is a helper function that returns the function mapping table
func getFuncMap() map[string]string {
	// Function name mapping
	funcMap := map[string]string{
		"sin": "\\sin ", "cos": "\\cos ", "tan": "\\tan ", "log": "\\log ",
		"ln": "\\ln ", "exp": "\\exp ", "lim": "\\lim ", "max": "\\max ",
		"min": "\\min ", "gcd": "\\gcd ", "det": "\\det ", "ker": "\\ker ",
		"dim": "\\dim ", "hom": "\\hom ", "cot": "\\cot ", "sec": "\\sec ",
		"csc": "\\csc ", "arcsin": "\\arcsin ", "arccos": "\\arccos ",
		"arctan": "\\arctan ", "Pr": "\\Pr ", "sup": "\\sup ", "inf": "\\inf ",
		"liminf": "\\liminf ", "limsup": "\\limsup ",
	}
	return funcMap
}
