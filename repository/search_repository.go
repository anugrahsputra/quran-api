package repository

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/anugrahsputra/go-quran-api/domain/model"
	"github.com/blevesearch/bleve/v2"
	_ "github.com/blevesearch/bleve/v2/analysis/lang/ar"
)

type SearchAyahRepository interface {
	Index(ayahs []model.Ayah) error
	Search(query string, page, limit int) (*bleve.SearchResult, error)
	GetDocument(id string) (map[string]any, error)
	GetDocCount() (uint64, error)
	IsHealthy() bool
}

type searchRepository struct {
	index bleve.Index
	path  string
}

func NewSearchAyahRepository(indexPath string) (SearchAyahRepository, error) {
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

	return &searchRepository{index: index, path: indexPath}, nil
}

// createNewIndex creates a new Bleve index with the proper mapping
func createNewIndex(indexPath string) (bleve.Index, error) {
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

	latinFieldMapping := bleve.NewTextFieldMapping()
	latinFieldMapping.Store = true
	latinFieldMapping.Analyzer = "standard"
	ayahMapping.AddFieldMappingsAt("Latin", latinFieldMapping)

	translationFieldMapping := bleve.NewTextFieldMapping()
	translationFieldMapping.Store = true
	translationFieldMapping.Analyzer = "standard"
	ayahMapping.AddFieldMappingsAt("Translation", translationFieldMapping)

	mapping := bleve.NewIndexMapping()
	mapping.DefaultAnalyzer = "standard"
	mapping.DefaultMapping = ayahMapping

	return bleve.New(indexPath, mapping)
}

func (r *searchRepository) Index(ayahs []model.Ayah) error {
	batch := r.index.NewBatch()
	indexedCount := 0
	emptyLatinCount := 0

	for _, ayah := range ayahs {
		id := strconv.Itoa(ayah.SurahNumber) + ":" + strconv.Itoa(ayah.AyahNumber)

		// log sample data for first few ayahs
		if indexedCount < 3 {
			latinPreview := ayah.Latin
			if len(latinPreview) > 50 {
				latinPreview = latinPreview[:50] + "..."
			}
			log.Printf("Indexing sample: ID=%s, Surah=%d, Ayah=%d, Latin length=%d, Latin preview=%s",
				id, ayah.SurahNumber, ayah.AyahNumber, len(ayah.Latin), latinPreview)
		}

		if ayah.Latin == "" {
			emptyLatinCount++
		}

		// index as map to ensure field names match the mapping
		batch.Index(id, map[string]any{
			"SurahNumber": ayah.SurahNumber,
			"AyahNumber":  ayah.AyahNumber,
			"Text":        ayah.Text,
			"Latin":       ayah.Latin,
			"Translation": ayah.Translation,
		})
		indexedCount++
	}

	log.Printf("Indexed %d ayahs (empty Latin fields: %d)", indexedCount, emptyLatinCount)
	return r.index.Batch(batch)
}

func (r *searchRepository) Search(query string, page, limit int) (*bleve.SearchResult, error) {
	// Validate and set defaults for pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10 // Default limit
	}
	if limit > 100 {
		limit = 100 // Max limit to prevent performance issues
	}

	docCount, err := r.index.DocCount()
	if err != nil {
		log.Printf("Warning: Could not get document count: %v", err)
	} else {
		log.Printf("Index contains %d documents", docCount)
	}

	queryLower := query

	translationQuery := bleve.NewMatchQuery(queryLower)
	translationQuery.SetField("Translation")

	translationWildcardQuery := bleve.NewWildcardQuery("*" + queryLower + "*")
	translationWildcardQuery.SetField("Translation")

	translationPrefixQuery := bleve.NewPrefixQuery(queryLower)
	translationPrefixQuery.SetField("Translation")

	disjunctionQuery := bleve.NewDisjunctionQuery(
		translationQuery,
		translationWildcardQuery,
		translationPrefixQuery,
	)
	disjunctionQuery.SetMin(1) // At least one should match

	// Calculate offset for pagination
	offset := (page - 1) * limit

	searchRequest := bleve.NewSearchRequest(disjunctionQuery)
	searchRequest.Fields = []string{"SurahNumber", "AyahNumber", "Text", "Latin", "Translation"}
	searchRequest.Size = limit
	searchRequest.From = offset
	searchRequest.IncludeLocations = false // We don't need location data

	log.Printf("Executing search request in repository with query: %s, page: %d, limit: %d, offset: %d",
		query, page, limit, offset)
	result, err := r.index.Search(searchRequest)
	if err != nil {
		log.Printf("Search failed: %v", err)
		return nil, err
	}

	log.Printf("Search found %d total results, %d hits (page %d, limit %d)",
		result.Total, len(result.Hits), page, limit)

	if result.Total == 0 {
		log.Printf("No results with standard queries. Trying fuzzy query on Translation...")
		translationFuzzyQuery := bleve.NewFuzzyQuery(queryLower)
		translationFuzzyQuery.SetField("Translation")
		translationFuzzyQuery.SetFuzziness(1) // Allow 1 character difference

		fuzzyRequest := bleve.NewSearchRequest(translationFuzzyQuery)
		fuzzyRequest.Fields = []string{"SurahNumber", "AyahNumber", "Text", "Latin", "Translation"}
		fuzzyRequest.Size = limit
		fuzzyRequest.From = offset
		fuzzyRequest.IncludeLocations = false
		fuzzyResult, fuzzyErr := r.index.Search(fuzzyRequest)
		if fuzzyErr == nil && fuzzyResult.Total > 0 {
			log.Printf("Fuzzy query found %d results", fuzzyResult.Total)
			return fuzzyResult, nil
		}
	}

	return result, nil
}

// GetDocument retrieves a document by ID from the index
// This is a fallback when hit.Fields is empty
func (r *searchRepository) GetDocument(id string) (map[string]interface{}, error) {
	// For now, return nil - we'll rely on hit.Fields being populated
	// If Fields are stored and specified in search request, they should be available
	// This method can be enhanced later if needed
	return nil, nil
}

func (r *searchRepository) GetDocCount() (uint64, error) {
	return r.index.DocCount()
}

func (r *searchRepository) IsHealthy() bool {
	_, err := r.index.DocCount()
	return err == nil
}
