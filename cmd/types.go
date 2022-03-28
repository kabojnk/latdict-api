package main

import "github.com/gin-gonic/gin"

type Pagination struct {
	PageNum       int `json:"pageNum"`
	PageSize      int `json:"pageSize"`
	TotalNumPages int `json:"TotalNumPages"`
}

type QueryFilter struct {
	QueryText     string   `json:"query"`
	PartsOfSpeech []string `json:"partsOfSpeech"`
	Ages          []string `json:"ages"`
	Commonalities []string `json:"commonality"`
	Geographies   []string `json:"geographies"`
	Conjugations  []string `json:"conjugations"`
	Voices        []string `json:"voices"`
	Declensions   []string `json:"declensions"`
	Genders       []string `json:"genders"`
}

type Entry struct {
	Id               string `json:"id"`
	Lemma            string `json:"lemma"`
	CommonalityScore int    `json:"commonalityScore" db:"commonality_score"`
	Orthography      string `json:"orthography"`
	Speech           string `json:"speech"`
}

type Route struct {
	Path    string
	Handler gin.HandlerFunc
}

//
//type MysqlEntry struct {
//	Guid             string
//	CommonalityScore sql.NullString `db:"commonality_score"`
//	Key              sql.NullString
//	Orthography      sql.NullString
//	Speech           sql.NullString
//	Id               int
//}
//
//type PgEntry struct {
//	Uuid             string
//	CommonalityScore sql.NullInt32 `db:"commonality_score"`
//	Lemma            string
//	Orthography      sql.NullString
//	Speech           sql.NullString
//	Id               int
//}
//
//type MysqlGrammarValue struct {
//	Guid      string
//	EntryGuid string `db:"entry_guid""`
//	Key       sql.NullString
//	Value     sql.NullString
//	Id        int
//	EntryId   int
//}
//
//type MysqlSense struct {
//	Guid      string
//	EntryGuid string `db:"entry_guid"`
//	Order     int
//	Sense     string
//	Id        int
//	EntryId   int
//}
//
//type MysqlSearchKeyEnglish struct {
//	Guid             string
//	EntryGuid        string `db:"entry_guid"`
//	Key              string
//	CommonalityScore sql.NullString `db:"commonality_score"`
//	Id               int
//	EntryId          int
//}
//
//type MysqlSearchKeyLatin struct {
//	Guid             string
//	EntryGuid        string `db:"entry_guid"`
//	Key              string
//	CommonalityScore sql.NullString `db:"commonality_score"`
//	Id               int
//	EntryId          int
//}
//
//type MysqlAdditionalInfo struct {
//	Age              sql.NullString
//	CommonalityScore sql.NullString `db:"commonality_score"`
//	Context          sql.NullString
//	EntryGuid        string `db:"entry_guid"`
//	Frequency        sql.NullString
//	Geography        sql.NullString
//	Guid             string
//	Source           sql.NullString
//	Id               int
//}
//
//type EntryAuxiliaryData struct {
//	GrammarValues    []MysqlGrammarValue
//	AdditionalInfos  []MysqlAdditionalInfo
//	Senses           []MysqlSense
//	SearchKeysLatin  []MysqlSearchKeyLatin
//	SearchKeyEnglish []MysqlSearchKeyEnglish
//}
