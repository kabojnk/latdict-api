package frequency

type Frequency string

const (
	UNKNOWN                       Frequency = "X"
	VERY_FREQUENT_TOP_1K                    = "A"
	FREQUENT_TOP_2K                         = "B"
	TOP_10K                                 = "C"
	TOP_20K                                 = "D"
	MAX_3_CITATIONS                         = "E"
	ONLY_CITED_IN_OLD_OR_LS                 = "F"
	ONLY_CITATION_IS_INSCRIPTION            = "I"
	PRESENTLY_NOT_MUCH_USED                 = "M"
	ONLY_IN_PLINY_NATURAL_HISTORY           = "N"
)
