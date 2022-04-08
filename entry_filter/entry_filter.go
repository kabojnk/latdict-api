package entry_filter

import (
	"github.com/kabojnk/latdict-api/query_param"
	"net/url"
	"strconv"
)

type EntryFilter struct {
	IncludeSenses         bool `json:"includeSenses"`
	IncludeGrammarInfo    bool `json:"IncludeGrammarInfo"`
	IncludeAdditionalInfo bool `json:"IncludeAdditionalInfo"`
	IncludeInflections    bool `json:"IncludeInflections"`
}

func (entryFilter *EntryFilter) InitWithQueryString(values url.Values) {
	entryFilter.IncludeSenses, _ = strconv.ParseBool(values.Get(query_param.INCLUDE_SENSES))
	entryFilter.IncludeAdditionalInfo, _ = strconv.ParseBool(values.Get(query_param.INCLUDE_ADDITIONA_INFO))
	entryFilter.IncludeGrammarInfo, _ = strconv.ParseBool(values.Get(query_param.INCLUDE_GRAMMAR_INFO))
	entryFilter.IncludeInflections, _ = strconv.ParseBool(values.Get(query_param.INCLUDE_INFLECTIONS))
}
