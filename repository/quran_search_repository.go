package repository

import (
	"fmt"
	"log"
	"os"
	"strconv"

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
		// Index doesn't exist, create a new one
		index, err = createNewIndex(indexPath)
		if err != nil {
			return nil, err
		}
		log.Printf("Created new search index at %s", indexPath)
	} else if err != nil {
		// Index exists but there's an error (possibly corrupted)
		log.Printf("Warning: Failed to open existing index: %v. Attempting to recreate...", err)
		// Try to remove the corrupted index directory and create a new one
		if removeErr := os.RemoveAll(indexPath); removeErr != nil {
			log.Printf("Warning: Failed to remove corrupted index: %v", removeErr)
		}
		index, err = createNewIndex(indexPath)
		if err != nil {
			return nil, fmt.Errorf("failed to recreate index: %w", err)
		}
		log.Printf("Recreated search index at %s", indexPath)
	} else {
		// Index exists and opened successfully, check document count
		docCount, err := index.DocCount()
		if err != nil {
			log.Printf("Warning: Could not get document count: %v", err)
		} else {
			log.Printf("Opened existing search index with %d documents", docCount)
		}
	}

	return &quranSearchRepository{index: index, path: indexPath}, nil
}

// createNewIndex creates a new Bleve index with optimized N-gram mapping
func createNewIndex(indexPath string) (bleve.Index, error) {
	mapping := bleve.NewIndexMapping()

	// 1. Define custom token filter for N-grams (min 3, max 4)
	// This allows efficient partial matching (e.g., "rahm" matches "rahman")
	err := mapping.AddCustomTokenFilter("ngram_filter", map[string]interface{}{
		"type": "ngram",
		"min":  3.0,
		"max":  4.0,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add ngram filter: %w", err)
	}

	// 2. Define custom analyzer using the ngram filter
	err = mapping.AddCustomAnalyzer("ngram_analyzer", map[string]interface{}{
		"type":          "custom",
		"tokenizer":     "unicode",
		"token_filters": []string{"to_lower", "ngram_filter"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add ngram analyzer: %w", err)
	}

	// 3. Define Document Mapping
	ayahMapping := bleve.NewDocumentMapping()

	// Numeric fields
	surahNumberFieldMapping := bleve.NewNumericFieldMapping()
	surahNumberFieldMapping.Store = true
	ayahMapping.AddFieldMappingsAt("SurahNumber", surahNumberFieldMapping)

	ayahNumberFieldMapping := bleve.NewNumericFieldMapping()
	ayahNumberFieldMapping.Store = true
	ayahMapping.AddFieldMappingsAt("AyahNumber", ayahNumberFieldMapping)

	// Arabic Text (Standard analysis)
	textFieldMapping := bleve.NewTextFieldMapping()
	textFieldMapping.Store = true
	textFieldMapping.Analyzer = "standard" // or specific arabic analyzer if available/configured
	ayahMapping.AddFieldMappingsAt("Text", textFieldMapping)

	// Let's manually define the mappings for our specific strategy (Standard + Ngram fields)

	// Latin
	latinStd := bleve.NewTextFieldMapping()
	latinStd.Store = true
	latinStd.Analyzer = "standard"
	ayahMapping.AddFieldMappingsAt("Latin", latinStd)

	latinNgram := bleve.NewTextFieldMapping()
	latinNgram.Store = false // Source is same as Latin
	latinNgram.Analyzer = "ngram_analyzer"
	ayahMapping.AddFieldMappingsAt("Latin_ngram", latinNgram)

	// Translation
	transStd := bleve.NewTextFieldMapping()
	transStd.Store = true
	transStd.Analyzer = "standard"
	ayahMapping.AddFieldMappingsAt("Translation", transStd)

	transNgram := bleve.NewTextFieldMapping()
	transNgram.Store = false
	transNgram.Analyzer = "ngram_analyzer"
	ayahMapping.AddFieldMappingsAt("Translation_ngram", transNgram)

	// Tafsir
	tafsirStd := bleve.NewTextFieldMapping()
	tafsirStd.Store = true
	tafsirStd.Analyzer = "standard"
	ayahMapping.AddFieldMappingsAt("Tafsir", tafsirStd)

	tafsirNgram := bleve.NewTextFieldMapping()
	tafsirNgram.Store = false
	tafsirNgram.Analyzer = "ngram_analyzer"
	ayahMapping.AddFieldMappingsAt("Tafsir_ngram", tafsirNgram)

	// Topic
	topicStd := bleve.NewTextFieldMapping()
	topicStd.Store = true
	topicStd.Analyzer = "standard"
	ayahMapping.AddFieldMappingsAt("Topic", topicStd)

	topicNgram := bleve.NewTextFieldMapping()
	topicNgram.Store = false
	topicNgram.Analyzer = "ngram_analyzer"
	ayahMapping.AddFieldMappingsAt("Topic_ngram", topicNgram)

	mapping.DefaultAnalyzer = "standard"
	mapping.DefaultMapping = ayahMapping

	// Use this mapping
	// Note: We need to update Index() to populate *_ngram fields

	// Since I cannot change the instruction mid-flight to update Index() method effectively without
	// rewriting the whole file (which I am doing), I will implement the logic in Index() below.

	return bleve.New(indexPath, mapping)
}

func (r *quranSearchRepository) Index(ayahs []model.Ayah) error {
	batch := r.index.NewBatch()
	indexedCount := 0

	for _, ayah := range ayahs {
		id := strconv.Itoa(ayah.SurahNumber) + ":" + strconv.Itoa(ayah.AyahNumber)

		// Create document map
		// We explicitly populate the _ngram fields with the same content
		doc := map[string]any{
			"SurahNumber":       ayah.SurahNumber,
			"AyahNumber":        ayah.AyahNumber,
			"Text":              ayah.Text,
			"Latin":             ayah.Latin,
			"Latin_ngram":       ayah.Latin, // Duplicate data for ngram indexing
			"Translation":       ayah.Translation,
			"Translation_ngram": ayah.Translation,
			"Tafsir":            ayah.Tafsir,
			"Tafsir_ngram":      ayah.Tafsir,
			"Topic":             ayah.Topic,
			"Topic_ngram":       ayah.Topic,
		}

		batch.Index(id, doc)
		indexedCount++
	}

	log.Printf("Indexed %d ayahs to %s", indexedCount, r.path)
	return r.index.Batch(batch)
}

func (r *quranSearchRepository) Search(query string, page, limit int) (*bleve.SearchResult, error) {
	// Validate defaults
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

	// 1. Exact/Standard Matches (Higher Boost)
	// Matches whole words or standard tokens
	translationMatch := bleve.NewMatchQuery(queryLower)
	translationMatch.SetField("Translation")
	translationMatch.SetBoost(5.0)

	tafsirMatch := bleve.NewMatchQuery(queryLower)
	tafsirMatch.SetField("Tafsir")
	tafsirMatch.SetBoost(3.0) // Tafsir is less important than translation

	topicMatch := bleve.NewMatchQuery(queryLower)
	topicMatch.SetField("Topic")
	topicMatch.SetBoost(4.0)

	latinMatch := bleve.NewMatchQuery(queryLower)
	latinMatch.SetField("Latin")
	latinMatch.SetBoost(5.0)

	// 2. N-gram Matches (Lower Boost, for partials)
	// Replaces expensive WildcardQuery (*query*)
	translationNgram := bleve.NewMatchQuery(queryLower)
	translationNgram.SetField("Translation_ngram")
	translationNgram.SetBoost(1.0)

	tafsirNgram := bleve.NewMatchQuery(queryLower)
	tafsirNgram.SetField("Tafsir_ngram")
	tafsirNgram.SetBoost(0.5)

	topicNgram := bleve.NewMatchQuery(queryLower)
	topicNgram.SetField("Topic_ngram")
	topicNgram.SetBoost(0.8)

	latinNgram := bleve.NewMatchQuery(queryLower)
	latinNgram.SetField("Latin_ngram")
	latinNgram.SetBoost(1.0)

	// Combine all queries
	disjunctionQuery := bleve.NewDisjunctionQuery(
		translationMatch,
		tafsirMatch,
		topicMatch,
		latinMatch,
		translationNgram,
		tafsirNgram,
		topicNgram,
		latinNgram,
	)
	disjunctionQuery.SetMin(1)

	// Calculate offset
	offset := (page - 1) * limit

	searchRequest := bleve.NewSearchRequest(disjunctionQuery)
	// We only fetch stored fields (original text), not the _ngram fields
	searchRequest.Fields = []string{"SurahNumber", "AyahNumber", "Text", "Latin", "Translation", "Tafsir", "Topic"}
	searchRequest.Size = limit
	searchRequest.From = offset
	searchRequest.IncludeLocations = false // Performance optimization

	result, err := r.index.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	// Minimal logging for performance/debugging
	if result.Total == 0 {
		log.Printf("Search '%s' returned 0 results", query)
	}

	return result, nil
}

func (r *quranSearchRepository) GetDocument(id string) (map[string]any, error) {
	// Not implemented/used for now as we use hit.Fields
	return nil, nil
}

func (r *quranSearchRepository) GetDocCount() (uint64, error) {
	return r.index.DocCount()
}

func (r *quranSearchRepository) IsHealthy() bool {
	_, err := r.index.DocCount()
	return err == nil
}
