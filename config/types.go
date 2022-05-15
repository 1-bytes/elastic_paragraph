package config

type SourceStruct struct {
	Id           string `json:"id"`
	SourceDomain string `json:"source_domain"`
	SourceUrl    string `json:"source_url"`
	Paragraph    struct {
		EN string `json:"en"`
	} `json:"paragraph"`
}

type RespStruct struct {
	Id                 int `json:"id"`
	DifficultyAnalysis struct {
		Fre                float64     `json:"fre"`
		Fkgl               float64     `json:"fkgl"`
		SchoolLvClass      int         `json:"schoolLvClass"`
		SchoolLvlName      string      `json:"schoolLvlName"`
		TeachWordLv        string      `json:"teachWordLv"`
		CefrWordLv         string      `json:"cefrWordLv"`
		SentenceCount      int         `json:"sentenceCount"`
		TeachNewWord       int         `json:"teachNewWord"`
		CefrNewWord        int         `json:"cefrNewWord"`
		WordCount          int         `json:"wordCount"`
		WordCountUnique    int         `json:"wordCountUnique"`
		TeachLvList        []string    `json:"teachLvList"`
		CefrLvList         []string    `json:"cefrLvList"`
		TeachWordLvInfo    interface{} `json:"teachWordLvInfo"`
		CefrWordLvInfo     interface{} `json:"cefrWordLvInfo"`
		LemmatizedWordList []string    `json:"lemmatizedWordList"`
	} `json:"difficultyAnalysis"`
}

type ParagraphListStruct struct {
	ParagraphList []ParagraphList `json:"paragraphList"`
}
type ParagraphList struct {
	ParagraphListId int    `json:"paragraphListId"`
	Paragraph       string `json:"paragraph"`
}
