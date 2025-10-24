package client

import (
	"context"
	"os"
	"testing"
)

func TestCreateEmbedding(t *testing.T) {
	baseURL := os.Getenv("PIXIGPT_BASE_URL")
	apiKey := os.Getenv("PIXIGPT_API_KEY")
	if baseURL == "" || apiKey == "" {
		t.Skip("PIXIGPT_BASE_URL or PIXIGPT_API_KEY not set")
	}

	client := New(apiKey, baseURL)
	ctx := context.Background()

	// Test 1: Single embedding
	t.Run("Single embedding", func(t *testing.T) {
		resp, err := client.CreateEmbedding(ctx, EmbeddingRequest{
			Input: "The quick brown fox jumps over the lazy dog",
		})
		if err != nil {
			t.Fatalf("CreateEmbedding failed: %v", err)
		}

		if len(resp.Data) != 1 {
			t.Fatalf("Expected 1 embedding, got %d", len(resp.Data))
		}

		if len(resp.Data[0].Embedding) == 0 {
			t.Fatal("Embedding vector is empty")
		}

		t.Logf("Embedding dimensions: %d", len(resp.Data[0].Embedding))
		t.Logf("First 5 values: %v", resp.Data[0].Embedding[:5])
		t.Logf("Usage: %d prompt tokens, %d total", resp.Usage.PromptTokens, resp.Usage.TotalTokens)
	})

	// Test 2: Batch embeddings
	t.Run("Batch embeddings", func(t *testing.T) {
		texts := []string{
			"Artificial intelligence is transforming technology",
			"Machine learning models process vast amounts of data",
			"Neural networks are inspired by biological neurons",
			"Deep learning requires significant computational resources",
		}

		resp, err := client.CreateEmbedding(ctx, EmbeddingRequest{
			Input: texts,
		})
		if err != nil {
			t.Fatalf("Batch CreateEmbedding failed: %v", err)
		}

		if len(resp.Data) != len(texts) {
			t.Fatalf("Expected %d embeddings, got %d", len(texts), len(resp.Data))
		}

		t.Logf("Generated %d embeddings", len(resp.Data))
		t.Logf("Dimensions: %d", len(resp.Data[0].Embedding))
		t.Logf("Usage: %d total tokens", resp.Usage.TotalTokens)
	})
}

func TestRerank(t *testing.T) {
	baseURL := os.Getenv("PIXIGPT_BASE_URL")
	apiKey := os.Getenv("PIXIGPT_API_KEY")
	if baseURL == "" || apiKey == "" {
		t.Skip("PIXIGPT_BASE_URL or PIXIGPT_API_KEY not set")
	}

	client := New(apiKey, baseURL)
	ctx := context.Background()

	query := "machine learning algorithms"
	documents := []string{
		"Machine learning is a subset of artificial intelligence that focuses on data-driven predictions",
		"Cats are popular pets known for their independence and playful nature",
		"Supervised learning algorithms learn from labeled training data to make predictions",
		"The weather forecast predicts rain tomorrow afternoon",
		"Neural networks use layers of interconnected nodes to process information",
		"Pizza is a traditional Italian dish with cheese and tomato sauce",
	}

	t.Run("Rerank documents", func(t *testing.T) {
		resp, err := client.Rerank(ctx, RerankRequest{
			Query:     query,
			Documents: documents,
			TopK:      3,
		})
		if err != nil {
			t.Fatalf("Rerank failed: %v", err)
		}

		if len(resp.Results) == 0 {
			t.Fatal("No rerank results returned")
		}

		t.Logf("Top %d results:", len(resp.Results))
		for i, result := range resp.Results {
			t.Logf("  %d. [%.3f] %s", i+1, result.RelevanceScore, result.Document)
		}
		t.Logf("Usage: %d total tokens", resp.Usage.TotalTokens)

		// Verify results are sorted by score (descending)
		for i := 1; i < len(resp.Results); i++ {
			if resp.Results[i].RelevanceScore > resp.Results[i-1].RelevanceScore {
				t.Errorf("Results not sorted: score[%d]=%.3f > score[%d]=%.3f",
					i, resp.Results[i].RelevanceScore, i-1, resp.Results[i-1].RelevanceScore)
			}
		}
	})
}
