package model

const TableDictParagraph = "dict_paragraph_new"

// DictParagraphModel 单词表
type DictParagraphModel struct {
	ID            int    `gorm:"type:int; primaryKey; autoIncrement; unsigned; not null" json:"id"`
	ElasticId     string `gorm:"type:varchar(64); primaryKey; not null" json:"elastic_id"`
	ArticleId     int    `gorm:"type:int; unsigned; not null" json:"article_id"`
	ByteCount     int    `gorm:"type:int; unsigned; not null" json:"byte_count"`
	Fre           string `gorm:"type:varchar(64);not null" json:"fre"`
	Fkgl          string `gorm:"type:varchar(64);not null" json:"fkgl"`
	SchoolLvClass int    `gorm:"type:int; unsigned; not null" json:"school_lv_class"`
	SchoolLvlName string `gorm:"type:varchar(128);not null" json:"school_lvl_name"`
	TechWordLv    string `gorm:"type:varchar(128);not null" json:"tech_word_lv"`
	CefrWordLv    string `gorm:"type:varchar(128);not null" json:"cefr_word_lv"`
	Status        int    `gorm:"type:int;not null" json:"status"`
}
