package model

// 段落内容
type DictArticleParagraph struct {
	ID            uint   `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	ElasticID     string `gorm:"column:elastic_id;default:0;NOT NULL"`      // ElasticSearch ID
	ArticleID     uint   `gorm:"column:article_id;default:0;NOT NULL"`      // 文章ID
	ByteCount     uint   `gorm:"column:byte_count;default:0;NOT NULL"`      // 段落长度
	Fre           string `gorm:"column:fre;NOT NULL"`                       // Flesch Reading Ease.FRE数值越高，文章就越简单，可读性也越高。
	Fkgl          string `gorm:"column:fkgl;NOT NULL"`                      // Flesch–Kincaid Grade Level.FKGL数值越高，文章就越复杂，文章的可读性也就越低。
	SchoolLvClass uint   `gorm:"column:school_lv_class;default:0;NOT NULL"` // 依据fre给出的评级
	SchoolLvlName string `gorm:"column:school_lvl_name;NOT NULL"`           // slv对应的学校等级
	TechWordLv    string `gorm:"column:tech_word_lv;NOT NULL"`              // 人教评级
	CefrWordLv    string `gorm:"column:cefr_word_lv;NOT NULL"`              // cefr评级
	Status        uint   `gorm:"column:status;default:1"`                   // 状态：1-正常；0-删除
}

func (m *DictArticleParagraph) TableName() string {
	return "dict_article_paragraph"
}
