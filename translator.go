package main

import (
	"context"
	"fmt"

	"github.com/EYERCORD/deepl-sdk-go"
	"github.com/EYERCORD/deepl-sdk-go/params"
	"github.com/EYERCORD/deepl-sdk-go/types"
)

type DeepLTranslator struct {
	Client *deepl.Client
}

func NewDeepLTranslator(apiKey string) (*DeepLTranslator, error) {
	client, err := deepl.NewClient(apiKey, "free")
	if err != nil {
		return nil, err
	}
	return &DeepLTranslator{
		Client: client,
	}, nil
}

type DeepLTranslationResult struct {
	DetectedLanguage string
	Result           string
}

func (translator *DeepLTranslator) Translate(text string, targetLanguage types.TargetLangCode, sourceLanguage types.SourceLangCode) (*DeepLTranslationResult, error) {
	var res *types.TranslateTextResponse
	var errRes *types.ErrorResponse
	var err error
	if sourceLanguage != "" {
		res, errRes, err = translator.Client.TranslateText(context.TODO(), &params.TranslateTextParams{
			TargetLang: targetLanguage,
			SourceLang: sourceLanguage,
			Text: []string{
				text,
			},
		})
	} else {
		res, errRes, err = translator.Client.TranslateText(context.TODO(), &params.TranslateTextParams{
			TargetLang: targetLanguage,
			Text: []string{
				text,
			},
		})
	}
	if err != nil {
		if errRes != nil {
			return nil, fmt.Errorf("%d %s", errRes.StatusCode, errRes.Message)
		}
		return nil, err
	}
	return &DeepLTranslationResult{
		DetectedLanguage: string(res.Translations[0].DetectedSourceLanguage),
		Result:           res.Translations[0].Text,
	}, nil
}
