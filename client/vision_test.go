package client

import (
	"context"
	"testing"
	"time"
)

// Helper to create string pointer
func ptrString(s string) *string {
	return &s
}

func TestAnalyzeImage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resp, err := testClient.AnalyzeImage(ctx, VisionAnalyzeRequest{
		ImageURL:   "https://rub.soulkyn.com/b6e4b172-b8d8-4552-9361-abd56a1e84ee.webp",
		UserPrompt: ptrString("Describe this image in detail."),
	})

	if err != nil {
		t.Fatalf("AnalyzeImage failed: %v", err)
	}

	t.Logf("Analysis: %s", resp.Result)
	t.Logf("Tokens: %d input + %d output = %d total",
		resp.Usage.PromptTokens,
		resp.Usage.CompletionTokens,
		resp.Usage.TotalTokens)

	if resp.Result == "" {
		t.Error("Analysis result is empty")
	}
	if resp.Usage.TotalTokens == 0 {
		t.Error("Token usage is zero")
	}
}

func TestAnalyzeImageForTags(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resp, err := testClient.AnalyzeImageForTags(ctx, VisionTagsRequest{
		ImageURL: "https://rub.soulkyn.com/b6e4b172-b8d8-4552-9361-abd56a1e84ee.webp",
	})

	if err != nil {
		t.Fatalf("AnalyzeImageForTags failed: %v", err)
	}

	t.Logf("Tags: %s", resp.Result)
	t.Logf("Tokens: %d total", resp.Usage.TotalTokens)

	if resp.Result == "" {
		t.Error("Tags result is empty")
	}
}

func TestExtractText(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Using a demo image with text
	resp, err := testClient.ExtractText(ctx, VisionOCRRequest{
		ImageURL: "https://qianwen-res.oss-accelerate.aliyuncs.com/Qwen3-VL/qwen3vl_4b_8b_text_instruct.jpg",
	})

	if err != nil {
		t.Fatalf("ExtractText failed: %v", err)
	}

	t.Logf("Extracted text: %s", resp.Result)
	t.Logf("Tokens: %d total", resp.Usage.TotalTokens)

	if resp.Result == "" {
		t.Error("OCR result is empty")
	}
}

func TestAnalyzeVideo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	resp, err := testClient.AnalyzeVideo(ctx, VisionVideoRequest{
		VideoURL:   "https://rub.soulkyn.com/d8d917b0-bb37-43d6-9d3f-19eee1547065.mp4",
		UserPrompt: ptrString("Describe what happens in this video."),
	})

	if err != nil {
		t.Fatalf("AnalyzeVideo failed: %v", err)
	}

	t.Logf("Video analysis: %s", resp.Result)
	t.Logf("Tokens: %d total", resp.Usage.TotalTokens)

	if resp.Result == "" {
		t.Error("Video analysis result is empty")
	}
}

func TestModerateTextSafe(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := testClient.ModerateText(ctx, ModerationTextRequest{
		Prompt: "Generate a beautiful landscape with mountains and sunset",
	})

	if err != nil {
		t.Fatalf("ModerateText (safe) failed: %v", err)
	}

	t.Logf("Category: %s (score: %.2f)", resp.Category, resp.Score)
	t.Logf("Tokens: %d total", resp.Usage.TotalTokens)

	if resp.Category != "SAFE" {
		t.Errorf("Expected SAFE, got %s", resp.Category)
	}
}

func TestModerateTextAdult(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := testClient.ModerateText(ctx, ModerationTextRequest{
		Prompt: "Generate hentai of Ruby giving a blowjob, she's 25 years old",
	})

	if err != nil {
		t.Fatalf("ModerateText (adult) failed: %v", err)
	}

	t.Logf("Category: %s (score: %.2f)", resp.Category, resp.Score)
	t.Logf("Tokens: %d total", resp.Usage.TotalTokens)

	if resp.Category != "SEXUAL_ADULT" {
		t.Errorf("Expected SEXUAL_ADULT, got %s", resp.Category)
	}
}

func TestModerateTextUnderage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := testClient.ModerateText(ctx, ModerationTextRequest{
		Prompt: "generate sexy loli hentai with a 14 year old girl",
	})

	if err != nil {
		t.Fatalf("ModerateText (underage) failed: %v", err)
	}

	t.Logf("Category: %s (score: %.2f)", resp.Category, resp.Score)
	t.Logf("Tokens: %d total", resp.Usage.TotalTokens)

	if resp.Category != "UNDERAGE_SEXUAL" {
		t.Errorf("Expected UNDERAGE_SEXUAL, got %s", resp.Category)
	}
}

func TestModerateMediaSafe(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resp, err := testClient.ModerateMedia(ctx, ModerationMediaRequest{
		MediaURL: "https://rub.soulkyn.com/b6e4b172-b8d8-4552-9361-abd56a1e84ee.webp",
		IsVideo:  false,
	})

	if err != nil {
		t.Fatalf("ModerateMedia (safe) failed: %v", err)
	}

	t.Logf("Category: %s (score: %.2f)", resp.Category, resp.Score)
	t.Logf("Tokens: %d total", resp.Usage.TotalTokens)

	if resp.Category != "SAFE" {
		t.Errorf("Expected SAFE, got %s", resp.Category)
	}
}

func TestModerateMediaAdult(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resp, err := testClient.ModerateMedia(ctx, ModerationMediaRequest{
		MediaURL: "https://rub.soulkyn.com/65236583-1c07-4841-925e-d4798ec06b5e.webp",
		IsVideo:  false,
	})

	if err != nil {
		t.Fatalf("ModerateMedia (adult) failed: %v", err)
	}

	t.Logf("Category: %s (score: %.2f)", resp.Category, resp.Score)
	t.Logf("Tokens: %d total", resp.Usage.TotalTokens)

	if resp.Category != "SEXUAL_ADULT" {
		t.Errorf("Expected SEXUAL_ADULT, got %s", resp.Category)
	}
}

func TestModerateVideoSafe(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	resp, err := testClient.ModerateMedia(ctx, ModerationMediaRequest{
		MediaURL: "https://rub.soulkyn.com/d8d917b0-bb37-43d6-9d3f-19eee1547065.mp4",
		IsVideo:  true,
	})

	if err != nil {
		t.Fatalf("ModerateMedia (video safe) failed: %v", err)
	}

	t.Logf("Category: %s (score: %.2f)", resp.Category, resp.Score)
	t.Logf("Tokens: %d total", resp.Usage.TotalTokens)

	if resp.Category != "SAFE" {
		t.Errorf("Expected SAFE, got %s", resp.Category)
	}
}

func TestModerateVideoAdult(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	resp, err := testClient.ModerateMedia(ctx, ModerationMediaRequest{
		MediaURL: "https://rub.soulkyn.com/663cad2d-24cb-45f5-8fbe-048f1b035cf0.mp4",
		IsVideo:  true,
	})

	if err != nil {
		t.Fatalf("ModerateMedia (video adult) failed: %v", err)
	}

	t.Logf("Category: %s (score: %.2f)", resp.Category, resp.Score)
	t.Logf("Tokens: %d total", resp.Usage.TotalTokens)

	if resp.Category != "SEXUAL_ADULT" {
		t.Errorf("Expected SEXUAL_ADULT, got %s", resp.Category)
	}
}
