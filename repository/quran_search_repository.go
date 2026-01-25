package repository

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/anugrahsputra/go-quran-api/domain/model"
	"github.com/blevesearch/bleve/v2"
	_ "github.com/blevesearch/bleve/v2/analysis/analyzer/custom"
	_ "github.com/blevesearch/bleve/v2/analysis/analyzer/standard"
	_ "github.com/blevesearch/bleve/v2/analysis/lang/ar"
	_ "github.com/blevesearch/bleve/v2/analysis/token/lowercase"
	_ "github.com/blevesearch/bleve/v2/analysis/token/ngram"
	_ "github.com/blevesearch/bleve/v2/analysis/tokenizer/unicode"
)

type QuranSearchRepository interface {
	Index(ayahs []model.Ayah) error
	Search(query string, page, limit int) (*bleve.SearchResult, error)
	GetDocument(id string) (map[string]any, error)
	GetDocCount() (uint64, error)
	IsHealthy() bool
}

type quranSearchRepository struct {
	index bleve.Index
	path  string
}

func NewQuranSearchRepository(indexPath string) (QuranSearchRepository, error) {
	index, err := bleve.Open(indexPath)
	if err == bleve.ErrorIndexPathDoesNotExist {
		index, err = createNewIndex(indexPath)
		if err != nil {
			return nil, err
		}
		log.Printf("Created new search index at %s", indexPath)
	} else if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "resource temporarily unavailable") || strings.Contains(errStr, "locked") || strings.Contains(errStr, "already open") {
			return nil, fmt.Errorf("failed to open search index at %s: index is locked by another process: %w", indexPath, err)
		}

		log.Printf("Warning: Failed to open existing index at %s: %v. Attempting to recreate...", indexPath, err)
		if removeErr := os.RemoveAll(indexPath); removeErr != nil {
			return nil, fmt.Errorf("failed to remove corrupted index at %s: %w", indexPath, removeErr)
		}
		index, err = createNewIndex(indexPath)
		if err != nil {
			return nil, fmt.Errorf("failed to recreate index at %s: %w", indexPath, err)
		}
		log.Printf("Recreated search index at %s", indexPath)
	} else {
		docCount, err := index.DocCount()
		if err != nil {
			log.Printf("Warning: Could not get document count: %v", err)
		} else {
			log.Printf("Opened existing search index with %d documents", docCount)
		}
	}

	return &quranSearchRepository{index: index, path: indexPath}, nil
}

func createNewIndex(indexPath string) (bleve.Index, error) {
	mapping := bleve.NewIndexMapping()

	err := mapping.AddCustomTokenFilter("ngram_filter", map[string]interface{}{
		"type": "ngram",
		"min":  3.0,
		"max":  4.0,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add ngram filter: %w", err)
	}

	err = mapping.AddCustomAnalyzer("ngram_analyzer", map[string]interface{}{
		"type":          "custom",
		"tokenizer":     "unicode",
		"token_filters": []string{"to_lower", "ngram_filter"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add ngram analyzer: %w", err)
	}

	ayahMapping := bleve.NewDocumentMapping()

	surahNumberFieldMapping := bleve.NewNumericFieldMapping()
	surahNumberFieldMapping.Store = true
	ayahMapping.AddFieldMappingsAt("SurahNumber", surahNumberFieldMapping)

	ayahNumberFieldMapping := bleve.NewNumericFieldMapping()
	ayahNumberFieldMapping.Store = true
	ayahMapping.AddFieldMappingsAt("AyahNumber", ayahNumberFieldMapping)

	textFieldMapping := bleve.NewTextFieldMapping()
	textFieldMapping.Store = true
	textFieldMapping.Analyzer = "standard"
	ayahMapping.AddFieldMappingsAt("Text", textFieldMapping)

	latinStd := bleve.NewTextFieldMapping()
	latinStd.Store = true
	latinStd.Analyzer = "standard"
	ayahMapping.AddFieldMappingsAt("Latin", latinStd)

	latinNgram := bleve.NewTextFieldMapping()
	latinNgram.Store = false
	latinNgram.Analyzer = "ngram_analyzer"
	ayahMapping.AddFieldMappingsAt("Latin_ngram", latinNgram)

	transStd := bleve.NewTextFieldMapping()
	transStd.Store = true
	transStd.Analyzer = "standard"
	ayahMapping.AddFieldMappingsAt("Translation", transStd)

	transNgram := bleve.NewTextFieldMapping()
	transNgram.Store = false
	transNgram.Analyzer = "ngram_analyzer"
	ayahMapping.AddFieldMappingsAt("Translation_ngram", transNgram)

	tafsirStd := bleve.NewTextFieldMapping()
	tafsirStd.Store = true
	tafsirStd.Analyzer = "standard"
	ayahMapping.AddFieldMappingsAt("Tafsir", tafsirStd)

	topicStd := bleve.NewTextFieldMapping()
	topicStd.Store = true
	topicStd.Analyzer = "standard"
	ayahMapping.AddFieldMappingsAt("Topic", topicStd)

	mapping.DefaultAnalyzer = "standard"
	mapping.DefaultMapping = ayahMapping

	return bleve.New(indexPath, mapping)
}

func (r *quranSearchRepository) Index(ayahs []model.Ayah) error {
	batch := r.index.NewBatch()
	indexedCount := 0

	for _, ayah := range ayahs {
		id := strconv.Itoa(ayah.SurahNumber) + ":" + strconv.Itoa(ayah.AyahNumber)

		doc := map[string]any{
			"SurahNumber":       ayah.SurahNumber,
			"AyahNumber":        ayah.AyahNumber,
			"Text":              ayah.Text,
			"Latin":             ayah.Latin,
			"Latin_ngram":       ayah.Latin,
			"Translation":       ayah.Translation,
			"Translation_ngram": ayah.Translation,
			"Tafsir":            ayah.Tafsir,
			"Topic":             ayah.Topic,
		}

		batch.Index(id, doc)
		indexedCount++
	}

	log.Printf("Indexed %d ayahs to %s", indexedCount, r.path)
	return r.index.Batch(batch)
}

func (r *quranSearchRepository) Search(query string, page, limit int) (*bleve.SearchResult, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	queryLower := query

	translationMatch := bleve.NewMatchQuery(queryLower)
	translationMatch.SetField("Translation")
	translationMatch.SetBoost(5.0)

	tafsirMatch := bleve.NewMatchQuery(queryLower)
	tafsirMatch.SetField("Tafsir")
	tafsirMatch.SetBoost(3.0)

	topicMatch := bleve.NewMatchQuery(queryLower)
	topicMatch.SetField("Topic")
	topicMatch.SetBoost(4.0)

	latinMatch := bleve.NewMatchQuery(queryLower)
	latinMatch.SetField("Latin")
	latinMatch.SetBoost(5.0)

	translationNgram := bleve.NewMatchQuery(queryLower)
	translationNgram.SetField("Translation_ngram")
	translationNgram.SetBoost(1.0)

	latinNgram := bleve.NewMatchQuery(queryLower)
	latinNgram.SetField("Latin_ngram")
	latinNgram.SetBoost(1.0)

	disjunctionQuery := bleve.NewDisjunctionQuery(
		translationMatch,
		tafsirMatch,
		topicMatch,
		latinMatch,
		translationNgram,
		latinNgram,
	)
	disjunctionQuery.SetMin(1)

	offset := (page - 1) * limit

	searchRequest := bleve.NewSearchRequest(disjunctionQuery)
	searchRequest.Fields = []string{"SurahNumber", "AyahNumber", "Text", "Latin", "Translation", "Tafsir", "Topic"}
	searchRequest.Size = limit
	searchRequest.From = offset
	searchRequest.IncludeLocations = false

	result, err := r.index.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	if result.Total == 0 {
		log.Printf("Search '%s' returned 0 results", query)
	}

	return result, nil
}

func (r *quranSearchRepository) GetDocument(id string) (map[string]any, error) {
	return nil, nil
}

func (r *quranSearchRepository) GetDocCount() (uint64, error) {
	return r.index.DocCount()
}

func (r *quranSearchRepository) IsHealthy() bool {
	_, err := r.index.DocCount()
	return err == nil
}