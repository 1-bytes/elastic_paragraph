package main

import (
	"context"
	"elastic_paragraph/bootstrap"
	"elastic_paragraph/config"
	"elastic_paragraph/model"
	"elastic_paragraph/pkg/elastic"
	"elastic_paragraph/pkg/fetcher"
	"encoding/json"
	"fmt"
	elasticSearch "github.com/olivere/elastic/v7"
	"github.com/panjf2000/ants/v2"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
)

func main() {
	bootstrap.Setup()
	// FullEStoMySQL()
	//PartMySQL()
	defer ants.Release()
	runTimes := 100000
	var wg sync.WaitGroup
	p, _ := ants.NewPoolWithFunc(1, func(_ interface{}) {
		PartMySQL()
		wg.Done()
	})
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		p.Invoke(0)
	}
	wg.Wait()
}

// PartMySQL 根据MySQL里面的数据差异化进行评级
func PartMySQL() {
	client := elastic.GetInstance()
	// 难度评级接口
	api := "http://127.0.0.1:9002/difficultyAnalysis"
	step := 1000
	db := bootstrap.DB
	// 从 MySQL 中分批取出 Fre 为空的数据
	var paragraphModel []model.DictArticleParagraph
	tx := db.Table("dict_article_paragraph").
		Where("fre = ?", "").
		Limit(step).
		Find(&paragraphModel)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		panic(fmt.Errorf("error: get paragraph model failed, err: %s", tx.Error))
	}
	// 拿着mysql当中的 elastic_id 去 elasticsearch 中查询
	var paragraphListJson config.ParagraphListStruct
	for _, paragraph := range paragraphModel {
		//fmt.Println(key, paragraph.ElasticID)
		idsQuery := elasticSearch.NewIdsQuery()
		idsQuery.Ids(paragraph.ElasticID)
		do, err := client.Search("dict_article").
			Size(1).
			Query(idsQuery).
			Do(context.Background())
		if err != nil || do.Hits.TotalHits.Value == 0 {
			panic(fmt.Errorf("error: get paragraph data is null, elasticsearch id: %s", paragraph.ElasticID))
		}
		source := config.SourceStruct{}
		_ = json.Unmarshal(do.Hits.Hits[0].Source, &source)
		paragraphListJson.ParagraphList = append(paragraphListJson.ParagraphList, config.ParagraphList{
			ParagraphListId: 0,
			Paragraph:       source.Paragraph.EN,
		})
	}
	// 将数据提交给难度评级的接口
	bodyJson, err := json.Marshal(paragraphListJson)
	respJson, err := fetcher.Fetch(http.MethodPost, api, bodyJson)
	if err != nil {
		panic(fmt.Errorf("error: Get word difficultyAnalysis failed: %v", err))
	}
	var resp []config.RespStruct
	_ = json.Unmarshal(respJson, &resp)
	// 将难度评级的结果和 ES 里面的数据结合
	for key, value := range resp {
		paragraphModel[key].Fre = fmt.Sprintf("%f", value.DifficultyAnalysis.Fre)
		paragraphModel[key].Fkgl = fmt.Sprintf("%f", value.DifficultyAnalysis.Fkgl)
		paragraphModel[key].SchoolLvClass = uint(value.DifficultyAnalysis.SchoolLvClass)
		paragraphModel[key].SchoolLvlName = value.DifficultyAnalysis.SchoolLvlName
		paragraphModel[key].TechWordLv = value.DifficultyAnalysis.TeachWordLv
		paragraphModel[key].CefrWordLv = value.DifficultyAnalysis.CefrWordLv
		paragraphModel[key].Status = 1
	}
	// 插入到 MySQL 当中
	tx = db.Table(paragraphModel[0].TableName()).Save(&paragraphModel)
	if tx.Error != nil {
		log.Printf("error: update paragraph model failed, err: %s", tx.Error)
		return
	}
	log.Println("success: create paragraph model success / step: ", step)
	return
}

// FullEStoMySQL 根据ES里面的数据 全量进行难度评级
func FullEStoMySQL() {
	client := elastic.GetInstance()
	// 难度评级接口
	api := "http://127.0.0.1:9002/difficultyAnalysis"
	step := 4000
	scroll := client.Scroll("dict_article").Size(step)
	for {
		// 从 es 当中滚动查询数据
		do, err := scroll.Do(context.Background())
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		// 命中数据为 0 则退出
		if len(do.Hits.Hits) <= 0 {
			break
		}

		paragraphListJson := config.ParagraphListStruct{}
		var paragraphModel []model.DictArticleParagraph
		for _, hits := range do.Hits.Hits {
			source := config.SourceStruct{}
			_ = json.Unmarshal(hits.Source, &source)
			// ----------------- 段落难度评级数据插入到 MySQL 数据库 -----------------
			ArticleID, err := strconv.Atoi(source.Id)
			if err != nil {
				log.Printf("error: conversion source.id to int failed, source.id: %s", source.Id)
				continue
			}
			// 需要评级难度的数据组合到一起 准备发送给接口
			paragraphListJson.ParagraphList = append(paragraphListJson.ParagraphList, config.ParagraphList{
				ParagraphListId: 0,
				Paragraph:       source.Paragraph.EN,
			})
			// ES 里面存储的一些基本信息先保存一份后面会用到
			paragraphModel = append(paragraphModel, model.DictArticleParagraph{
				ElasticID: hits.Id,
				ArticleID: uint(ArticleID),
				ByteCount: uint(len(source.Paragraph.EN)),
			})
		}
		// 将数据提交给难度评级的接口
		bodyJson, err := json.Marshal(paragraphListJson)
		respJson, err := fetcher.Fetch(http.MethodPost, api, bodyJson)
		if err != nil {
			fmt.Errorf("error: Get word difficultyAnalysis failed: %v", err)
			continue
		}
		var resp []config.RespStruct
		err = json.Unmarshal(respJson, &resp)
		// 将难度评级的结果和 ES 里面的数据结合
		for key, value := range resp {
			paragraphModel[key].Fre = fmt.Sprintf("%f", value.DifficultyAnalysis.Fre)
			paragraphModel[key].Fkgl = fmt.Sprintf("%f", value.DifficultyAnalysis.Fkgl)
			paragraphModel[key].SchoolLvClass = uint(value.DifficultyAnalysis.SchoolLvClass)
			paragraphModel[key].SchoolLvlName = value.DifficultyAnalysis.SchoolLvlName
			paragraphModel[key].TechWordLv = value.DifficultyAnalysis.TeachWordLv
			paragraphModel[key].CefrWordLv = value.DifficultyAnalysis.CefrWordLv
			paragraphModel[key].Status = 1
		}

		// 插入到 MySQL 当中
		db := bootstrap.DB
		tx := db.Table(paragraphModel[0].TableName()).Create(&paragraphModel)
		if tx.Error != nil {
			log.Printf("error: create paragraph model failed, err: %s", tx.Error)
			continue
		}
		log.Println("success: create paragraph model success / step: ", step)
	}
	log.Println("done.")
}
