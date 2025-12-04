package service

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/anugrahsputra/go-quran-api/domain/model"
	"github.com/anugrahsputra/go-quran-api/repository"
)

type SearchAyahService interface {
	IndexQuran() error
	Search(query string, page, limit int) ([]model.Ayah, int, error)
}

type searchAyahService struct {
	quranRepo  repository.IQuranRepository
	searchRepo repository.SearchAyahRepository
}

func NewSearchAyahService(quranRepo repository.IQuranRepository, searchRepo repository.SearchAyahRepository) SearchAyahService {
	return &searchAyahService{quranRepo: quranRepo, searchRepo: searchRepo}
}

func (s *searchAyahService) IndexQuran() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	const (
		totalSurahs    = 114
		batchAyahLimit = 1000 // Index in batches of ~1000 ayahs to avoid memory issues
		maxRetries     = 3
		retryDelay     = 2 * time.Second
	)

	var (
		allAyahs         []model.Ayah
		successCount     int
		failureCount     int
		totalAyahs       int
		emptyTranslation int
		startTime        = time.Now()
	)

	log.Printf("Starting Quran indexing process for %d surahs...", totalSurahs)

	for i := 1; i <= totalSurahs; i++ {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return fmt.Errorf("indexing cancelled or timed out: %w", ctx.Err())
		default:
		}

		// Fetch surah with retry logic
		var detailSurah model.DetailSurahApi
		var err error

		for attempt := 1; attempt <= maxRetries; attempt++ {
			detailSurah, err = s.quranRepo.GetSurahDetail(ctx, i, 0, 300)
			if err == nil {
				break
			}

			if attempt < maxRetries {
				log.Printf("Attempt %d/%d failed for surah %d: %v. Retrying in %v...",
					attempt, maxRetries, i, err, retryDelay)
				time.Sleep(retryDelay)
			}
		}

		if err != nil {
			failureCount++
			log.Printf("ERROR: Failed to fetch surah %d after %d attempts: %v", i, maxRetries, err)

			// Continue with next surah instead of failing completely
			// This allows partial indexing if some surahs fail
			continue
		}

		if len(detailSurah.Data) == 0 {
			log.Printf("WARNING: Surah %d returned no data", i)
			failureCount++
			continue
		}

		// Process verses
		surahAyahCount := 0
		for _, verse := range detailSurah.Data {
			// Validate verse data
			if verse.SurahID <= 0 || verse.Ayah <= 0 {
				log.Printf("WARNING: Invalid verse data in surah %d: SurahID=%d, Ayah=%d",
					i, verse.SurahID, verse.Ayah)
				continue
			}

			ayah := model.Ayah{
				SurahNumber: verse.SurahID,
				AyahNumber:  verse.Ayah,
				Text:        verse.Arabic,
				Latin:       verse.Latin,
				Translation: verse.Translation,
			}

			// Track empty translations
			if ayah.Translation == "" {
				emptyTranslation++
			}

			allAyahs = append(allAyahs, ayah)
			surahAyahCount++
			totalAyahs++

			// Log sample data for first surah
			if i == 1 && len(allAyahs) <= 3 {
				translationPreview := verse.Translation
				if len(translationPreview) > 50 {
					translationPreview = translationPreview[:50] + "..."
				}
				log.Printf("Sample: Surah %d, Ayah %d, Translation: %s",
					verse.SurahID, verse.Ayah, translationPreview)
			}
		}

		successCount++

		// Progress logging every 10 surahs or at milestones
		if i%10 == 0 || i == totalSurahs {
			progress := float64(i) / float64(totalSurahs) * 100
			elapsed := time.Since(startTime)
			estimatedTotal := time.Duration(float64(elapsed) / progress * 100)
			remaining := estimatedTotal - elapsed

			log.Printf("Progress: %d/%d surahs (%.1f%%) | %d ayahs indexed | "+
				"Elapsed: %v | Estimated remaining: %v | "+
				"Success: %d | Failed: %d",
				i, totalSurahs, progress, totalAyahs,
				elapsed.Round(time.Second), remaining.Round(time.Second),
				successCount, failureCount)
		}

		// Batch indexing to avoid memory issues
		// Index in batches when we reach the limit or at the end
		if len(allAyahs) >= batchAyahLimit || i == totalSurahs {
			batchNum := (i / 10) + 1
			if err := s.searchRepo.Index(allAyahs); err != nil {
				return fmt.Errorf("failed to index batch %d (up to surah %d): %w", batchNum, i, err)
			}
			log.Printf("Indexed batch %d: %d ayahs (up to surah %d)", batchNum, len(allAyahs), i)
			allAyahs = allAyahs[:0] // Clear slice but keep capacity
		}
	}

	// Final statistics
	duration := time.Since(startTime)
	log.Printf("Indexing completed!")
	log.Printf("Statistics:")
	log.Printf("  - Total surahs processed: %d/%d", successCount, totalSurahs)
	log.Printf("  - Failed surahs: %d", failureCount)
	log.Printf("  - Total ayahs indexed: %d", totalAyahs)
	log.Printf("  - Ayahs with empty translation: %d (%.1f%%)",
		emptyTranslation, float64(emptyTranslation)/float64(totalAyahs)*100)
	log.Printf("  - Total duration: %v", duration.Round(time.Second))
	log.Printf("  - Average: %.2f ayahs/second", float64(totalAyahs)/duration.Seconds())

	if failureCount > 0 {
		return fmt.Errorf("indexing completed with %d failed surahs out of %d total",
			failureCount, totalSurahs)
	}

	return nil
}

func (s *searchAyahService) Search(query string, page, limit int) ([]model.Ayah, int, error) {
	searchResult, err := s.searchRepo.Search(query, page, limit)
	if err != nil {
		return nil, 0, err
	}

	totalResults := int(searchResult.Total)
	log.Printf("Search query: %s, Total results: %d, Hits found: %d (page: %d, limit: %d)",
		query, totalResults, len(searchResult.Hits), page, limit)

	var ayahs []model.Ayah
	for i, hit := range searchResult.Hits {
		// Safely extract fields with type checking
		var surahNumber, ayahNumber int
		var text, latin, translation string

		// Debug: Log what we have in the hit
		if i == 0 {
			log.Printf("DEBUG: First hit - ID: %s, Fields count: %d, Fields: %+v", hit.ID, len(hit.Fields), hit.Fields)
		}

		// Use hit.Fields which should contain the stored fields
		fields := hit.Fields
		if fields == nil || len(fields) == 0 {
			log.Printf("Warning: No fields found in hit - ID: %s, Score: %f", hit.ID, hit.Score)
			// Parse the ID to extract surah and ayah numbers as fallback
			// ID format is "surah:ayah" (e.g., "1:1")
			parts := splitDocumentID(hit.ID)
			if len(parts) == 2 {
				if sn, err := strconv.Atoi(parts[0]); err == nil {
					surahNumber = sn
				}
				if an, err := strconv.Atoi(parts[1]); err == nil {
					ayahNumber = an
				}
				// If we got surah and ayah from ID, create a minimal ayah
				if surahNumber > 0 && ayahNumber > 0 {
					ayahs = append(ayahs, model.Ayah{
						SurahNumber: surahNumber,
						AyahNumber:  ayahNumber,
						Text:        "", // Will be empty if fields aren't available
						Latin:       "",
						Translation: "",
					})
					log.Printf("Warning: Created minimal ayah from ID for %s", hit.ID)
				}
			}
			continue
		}

		// Extract SurahNumber
		if sn, ok := fields["SurahNumber"]; ok {
			switch v := sn.(type) {
			case float64:
				surahNumber = int(v)
			case int:
				surahNumber = v
			case int64:
				surahNumber = int(v)
			case int32:
				surahNumber = int(v)
			}
		}

		// Extract AyahNumber
		if an, ok := fields["AyahNumber"]; ok {
			switch v := an.(type) {
			case float64:
				ayahNumber = int(v)
			case int:
				ayahNumber = v
			case int64:
				ayahNumber = int(v)
			case int32:
				ayahNumber = int(v)
			}
		}

		// Extract Text
		if t, ok := fields["Text"]; ok {
			if textStr, ok := t.(string); ok {
				text = textStr
			}
		}

		// Extract Latin field
		if l, ok := fields["Latin"]; ok {
			if latinStr, ok := l.(string); ok {
				latin = latinStr
			} else {
				log.Printf("Warning: Latin field exists but is not a string, type: %T, value: %v", l, l)
			}
		}

		// Extract Translation field
		if t, ok := fields["Translation"]; ok {
			if translationStr, ok := t.(string); ok {
				translation = translationStr
			}
		}

		// Only add ayah if we have valid data (at least surah and ayah numbers)
		if surahNumber > 0 && ayahNumber > 0 {
			ayahs = append(ayahs, model.Ayah{
				SurahNumber: surahNumber,
				AyahNumber:  ayahNumber,
				Text:        text,
				Latin:       latin,
				Translation: translation,
			})
		} else {
			log.Printf("Warning: Skipping hit with incomplete data - ID: %s, SurahNumber: %v, AyahNumber: %v, Fields keys: %v",
				hit.ID, fields["SurahNumber"], fields["AyahNumber"], getFieldKeys(fields))
		}
	}

	return ayahs, totalResults, nil
}

// Helper function to get keys from a map for logging
func getFieldKeys(fields map[string]interface{}) []string {
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	return keys
}

// Helper function to split document ID (format: "surah:ayah")
func splitDocumentID(id string) []string {
	return strings.Split(id, ":")
}
