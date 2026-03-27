// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"errors"
	"strings"

	"chatwiki/internal/app/chatwiki/define"
)

func NormalizeLibrarySearchType(searchType int, librarySearchType string) (string, error) {
	if searchType != define.SearchTypeMixed && searchType != define.SearchTypeFullText {
		return ``, nil
	}
	librarySearchType = strings.TrimSpace(librarySearchType)
	if librarySearchType == `` {
		return define.LibrarySearchTypeFullText, nil
	}
	switch librarySearchType {
	case define.LibrarySearchTypeFullText, define.LibrarySearchTypeKeyword:
		return librarySearchType, nil
	default:
		return ``, errors.New(`invalid library_search_type`)
	}
}
