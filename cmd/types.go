package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
)

/**
 * Contains genera models for use with the API.

 * A note on the DB* and API* model dichotomy: these separate models because they might need to change between what is
 * pulled from the DB and what is ultimately sent to the API. This could mean omitting some fields, or in cases of
 * fields of `sql.Null*` types, we don't want to flatten the structure of those values with defaults.
 */

type Pagination struct {
	PageNum       int `json:"pageNum"`
	PageSize      int `json:"pageSize"`
	TotalNumPages int `json:"TotalNumPages"`
}

type DBEntry struct {
	ID               string         `json:"id"`
	UUID             string         `json:"uuid"`
	Lemma            string         `json:"lemma"`
	CommonalityScore int            `json:"commonalityScore" db:"commonality_score"`
	Orthography      sql.NullString `json:"orthography"`
	Speech           string         `json:"speech"`
}

type APIEntry struct {
	UUID             string `json:"uuid"`
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
