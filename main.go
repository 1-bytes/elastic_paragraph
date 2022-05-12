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
	"io"
	"log"
	"net/http"
	"strconv"
)

func main() {
	bootstrap.Setup()
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

		var paragraphModel []model.DictParagraphModel
		for _, hits := range do.Hits.Hits {
			source := config.SourceStruct{}
			_ = json.Unmarshal(hits.Source, &source)
			// ----------------- 段落难度评级数据插入到 MySQL 数据库 -----------------
			paragraphListJson := config.ParagraphListStruct{ParagraphList: []config.ParagraphList{
				{
					ParagraphListId: 0,
					Paragraph:       source.Paragraph.EN,
				},
			}}
			bodyJson, err := json.Marshal(paragraphListJson)
			respJson, err := fetcher.Fetch(http.MethodPost, api, bodyJson)
			if err != nil {
				fmt.Errorf("error: Get word difficultyAnalysis failed: %v", err)
				continue
			}
			var resp []config.RespStruct
			err = json.Unmarshal(respJson, &resp)

			// ----------------- 查出的数据插入到 MySQL 数据库 -----------------
			ArticleID, err := strconv.Atoi(source.Id)
			if err != nil {
				log.Printf("error: conversion source.id to int failed, source.id: %s", source.Id)
				continue
			}
			paragraphModel = append(paragraphModel, model.DictParagraphModel{
				ElasticId:     hits.Id,
				ArticleId:     ArticleID,
				ByteCount:     len(source.Paragraph.EN),
				Fre:           fmt.Sprintf("%f", resp[0].DifficultyAnalysis.Fre),
				Fkgl:          fmt.Sprintf("%f", resp[0].DifficultyAnalysis.Fkgl),
				SchoolLvClass: resp[0].DifficultyAnalysis.SchoolLvClass,
				SchoolLvlName: resp[0].DifficultyAnalysis.SchoolLvlName,
				TechWordLv:    resp[0].DifficultyAnalysis.TeachWordLv,
				CefrWordLv:    resp[0].DifficultyAnalysis.CefrWordLv,
				Status:        1,
			})
		}
		db := bootstrap.DB
		tx := db.Table(model.TableDictParagraph).Create(&paragraphModel)
		if tx.Error != nil {
			log.Printf("error: create paragraph model failed, err: %s", tx.Error)
			continue
		}
		log.Println("success: create paragraph model success / step: ", step)
	}
	log.Println("done.")
}
