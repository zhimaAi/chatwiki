// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package i18n

import (
	"chatwiki/internal/app/chatwiki/define"
	"regexp"
	"strings"
	"sync"
)

// BilingualPattern bilingual regex pattern with lazy loading from i18n
type BilingualPattern struct {
	zhPatternKey string         // i18n key for zh-CN pattern
	enPatternKey string         // i18n key for en-US pattern
	zhPattern    *regexp.Regexp // compiled zh-CN regex
	enPattern    *regexp.Regexp // compiled en-US regex
	once         sync.Once      // ensure patterns are compiled only once
}

// BilingualMatchResult bilingual match result
type BilingualMatchResult struct {
	Matched bool     // whether matched
	Lang    string   // matched language (zh-CN or en-US)
	Matches []string // matched groups
}

// NewBilingualPattern create bilingual pattern with i18n keys
func NewBilingualPattern(zhPatternKey, enPatternKey string) *BilingualPattern {
	return &BilingualPattern{
		zhPatternKey: zhPatternKey,
		enPatternKey: enPatternKey,
	}
}

// compile compiles regex patterns from i18n (lazy loading)
func (bp *BilingualPattern) compile() {
	bp.once.Do(func() {
		zhPattern := Show(define.LangZhCn, bp.zhPatternKey)
		enPattern := Show(define.LangEnUs, bp.enPatternKey)
		bp.zhPattern = regexp.MustCompile(zhPattern)
		bp.enPattern = regexp.MustCompile(enPattern)
	})
}

// Match performs bilingual matching, returns match result and language
func (bp *BilingualPattern) Match(text string) BilingualMatchResult {
	bp.compile() // ensure patterns are compiled

	// try zh-CN match first
	if matches := bp.zhPattern.FindStringSubmatch(text); len(matches) > 0 {
		return BilingualMatchResult{
			Matched: true,
			Lang:    define.LangZhCn,
			Matches: matches,
		}
	}

	// try en-US match
	if matches := bp.enPattern.FindStringSubmatch(text); len(matches) > 0 {
		return BilingualMatchResult{
			Matched: true,
			Lang:    define.LangEnUs,
			Matches: matches,
		}
	}

	return BilingualMatchResult{
		Matched: false,
		Lang:    "",
		Matches: nil,
	}
}

// MatchWithLang performs bilingual matching with preferred language
func (bp *BilingualPattern) MatchWithLang(text string, preferredLang string) BilingualMatchResult {
	bp.compile() // ensure patterns are compiled

	if preferredLang == define.LangEnUs {
		// try en-US match first
		if matches := bp.enPattern.FindStringSubmatch(text); len(matches) > 0 {
			return BilingualMatchResult{
				Matched: true,
				Lang:    define.LangEnUs,
				Matches: matches,
			}
		}
		// fallback to zh-CN match
		if matches := bp.zhPattern.FindStringSubmatch(text); len(matches) > 0 {
			return BilingualMatchResult{
				Matched: true,
				Lang:    define.LangZhCn,
				Matches: matches,
			}
		}
	} else {
		// default to zh-CN match first
		return bp.Match(text)
	}

	return BilingualMatchResult{
		Matched: false,
		Lang:    "",
		Matches: nil,
	}
}

// BilingualKeyword bilingual keyword with i18n key
type BilingualKeyword struct {
	i18nKey string // i18n key for both zh-CN and en-US
}

// NewBilingualKeyword create bilingual keyword with i18n key
func NewBilingualKeyword(i18nKey string) *BilingualKeyword {
	return &BilingualKeyword{
		i18nKey: i18nKey,
	}
}

// getKeywords gets keywords from i18n
func (bk *BilingualKeyword) getKeywords() (zhKeyword, enKeyword string) {
	zhKeyword = Show(define.LangZhCn, bk.i18nKey)
	enKeyword = Show(define.LangEnUs, bk.i18nKey)
	return
}

// Match matches keyword, returns whether matched and language
// en-US keyword matching is case-insensitive
func (bk *BilingualKeyword) Match(text string) (matched bool, lang string) {
	zhKeyword, enKeyword := bk.getKeywords()
	if text == zhKeyword {
		return true, define.LangZhCn
	}
	if strings.EqualFold(text, enKeyword) {
		return true, define.LangEnUs
	}
	return false, ""
}

// MatchWithLang matches keyword with preferred language
// en-US keyword matching is case-insensitive
func (bk *BilingualKeyword) MatchWithLang(text string, preferredLang string) (matched bool, lang string) {
	zhKeyword, enKeyword := bk.getKeywords()
	if preferredLang == define.LangEnUs {
		if strings.EqualFold(text, enKeyword) {
			return true, define.LangEnUs
		}
		if text == zhKeyword {
			return true, define.LangZhCn
		}
	} else {
		return bk.Match(text)
	}
	return false, ""
}

// Predefined bilingual regex patterns - payment related
// en-US patterns use (?i) flag for case-insensitive matching (configured in ini file)
var (
	// PaymentCountPackagePattern count package pattern: xxx[10 times]--99 yuan
	PaymentCountPackagePattern = NewBilingualPattern(
		"payment_count_package_pattern",
		"payment_count_package_pattern",
	)

	// PaymentDurationPackagePattern duration package pattern: xxx[30 days]--199 yuan
	PaymentDurationPackagePattern = NewBilingualPattern(
		"payment_duration_package_pattern",
		"payment_duration_package_pattern",
	)

	// PaymentAuthCodeContentPattern auth code content pattern: ###xxxC/D random###
	PaymentAuthCodeContentPattern = regexp.MustCompile(`###.+?###`)
)

// Predefined bilingual keywords - payment related
var (
	// KeywordAuthCode auth code keyword
	KeywordAuthCode = NewBilingualKeyword("keyword_auth_code")

	// KeywordMyBenefits my benefits keyword
	KeywordMyBenefits = NewBilingualKeyword("keyword_my_benefits")
)

// GetPaymentCountPackageFormat gets count package format string
// format: name [count times]--price yuan
func GetPaymentCountPackageFormat(lang string) string {
	if lang == define.LangEnUs {
		return "%s [%d times]--%s yuan"
	}
	return "%s\u3010%d\u6b21\u3011--%s\u5143"
}

// GetPaymentDurationPackageFormat gets duration package format string
// format: name [duration days]--price yuan
func GetPaymentDurationPackageFormat(lang string) string {
	if lang == define.LangEnUs {
		return "%s [%d days]--%s yuan"
	}
	return "%s\u3010%d\u5929\u3011--%s\u5143"
}
