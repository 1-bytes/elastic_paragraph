package main

import (
	"context"
	"elastic_paragraph/bootstrap"
	"elastic_paragraph/model"
	"elastic_paragraph/pkg/elastic"
	"encoding/json"
	"log"
	"strconv"
)

func main() {
	bootstrap.Setup()
	client := elastic.GetInstance()
	step := 4000
	scroll := client.Scroll("dict_article").Size(step)
	for {
		do, err := scroll.Do(context.Background())
		if err != nil {
			panic(err)
		}
		if len(do.Hits.Hits) <= 0 {
			break
		}
		source := Source{}
		var paragraphModel []model.DictParagraphModel
		for _, hits := range do.Hits.Hits {
			_ = json.Unmarshal(hits.Source, &source)
			ArticleID, err := strconv.Atoi(source.Id)
			if err != nil {
				log.Printf("error: conversion source.id to int failed, source.id: %s", source.Id)
				continue
			}
			paragraphModel = append(paragraphModel, model.DictParagraphModel{
				ElasticId: hits.Id,
				ArticleId: ArticleID,
				ByteCount: len(source.Paragraph.EN),
			})
		}
		db := bootstrap.DB
		tx := db.Table(model.TableDictParagraph).Create(&paragraphModel)
		if tx.Error != nil {
			log.Printf("error: create paragraph model failed, err: %s", tx.Error)
			continue
		}
	}
}

type Source struct {
	Id           string `json:"id"`
	SourceDomain string `json:"source_domain"`
	SourceUrl    string `json:"source_url"`
	Paragraph    struct {
		EN string `json:"EN"`
	} `json:"paragraph"`
}
