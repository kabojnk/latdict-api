package query_filter

import (
	"github.com/kabojnk/latdict-api/query_param"
	"net/url"
	"strconv"
	"strings"
)

type QueryFilter struct {
	QueryText       string   `json:"query"`
	Language        string   `json:"language"`
	NeedsExactMatch bool     `json:"needsExactMatch"`
	PartsOfSpeech   []string `json:"partsOfSpeech"`
	Ages            []string `json:"ages"`
	Frequencies     []string `json:"frequencies"`
	Geographies     []string `json:"geographies"`
	Sources         []string `json:"sources"`
	Subjects        []string `json:"subjects"`
	Conjugations    []string `json:"conjugations"`
	Voices          []string `json:"voices"`
	Declensions     []string `json:"declensions"`
	Genders         []string `json:"genders"`
	IncludeSenses   bool     `json:"includeSenses"`
}

func (queryFilter *QueryFilter) InitWithQueryString(values url.Values) {
	queryFilter.QueryText = strings.TrimSpace(values.Get(query_param.QUERY))
	queryFilter.Language = strings.TrimSpace(values.Get(query_param.LANGUAGE))
	if len(queryFilter.Language) == 0 {
		queryFilter.Language = "latin"
	}
	queryFilter.NeedsExactMatch, _ = strconv.ParseBool(values.Get(query_param.NEEDS_EXACT_MATCH))
	queryFilter.PartsOfSpeech = strings.Split(strings.TrimSpace(values.Get(query_param.PARTS_OF_SPEECH)), ",")
	queryFilter.Ages = strings.Split(strings.TrimSpace(values.Get(query_param.AGES)), ",")
	queryFilter.Frequencies = strings.Split(strings.TrimSpace(values.Get(query_param.FREQUENCIES)), ",")
	queryFilter.Geographies = strings.Split(strings.TrimSpace(values.Get(query_param.GEOGRAPHIES)), ",")
	queryFilter.Sources = strings.Split(strings.TrimSpace(values.Get(query_param.SOURCES)), ",")
	queryFilter.Subjects = strings.Split(strings.TrimSpace(values.Get(query_param.SUBJECTS)), ",")
	queryFilter.Conjugations = strings.Split(strings.TrimSpace(values.Get(query_param.CONJUGATIONS)), ",")
	queryFilter.Voices = strings.Split(strings.TrimSpace(values.Get(query_param.VOICES)), ",")
	queryFilter.Declensions = strings.Split(strings.TrimSpace(values.Get(query_param.DECLENSIONS)), ",")
	queryFilter.Genders = strings.Split(strings.TrimSpace(values.Get(query_param.GENDERS)), ",")
	queryFilter.IncludeSenses, _ = strconv.ParseBool(values.Get(query_param.INCLUDE_SENSES))
}
