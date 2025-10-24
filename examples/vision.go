package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/PixiGPT/pixigpt-go/client"
	"github.com/joho/godotenv"
)

// Helper to create string pointer
func ptrString(s string) *string {
	return &s
}

func main() {
	// Load .env for testing convenience
	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("Warning: .env not found, using environment variables")
	}

	apiKey := os.Getenv("PIXIGPT_API_KEY")
	baseURL := os.Getenv("PIXIGPT_BASE_URL")

	if apiKey == "" || baseURL == "" {
		log.Fatal("Missing required environment variables: PIXIGPT_API_KEY, PIXIGPT_BASE_URL")
	}

	// Create client
	c := client.New(apiKey, baseURL)
	ctx := context.Background()

	// Example 1: Image Analysis
	fmt.Println("=== Image Analysis ===")
	imgResp, err := c.AnalyzeImage(ctx, client.VisionAnalyzeRequest{
		ImageURL:   "https://rub.soulkyn.com/b6e4b172-b8d8-4552-9361-abd56a1e84ee.webp",
		UserPrompt: ptrString("Describe this image in detail."),
	})
	if err != nil {
		log.Printf("Image analysis failed: %v", err)
	} else {
		fmt.Printf("Analysis: %s\n", imgResp.Result)
		fmt.Printf("Tokens: %d total\n\n", imgResp.Usage.TotalTokens)
	}

	// Example 2: Tag Generation
	fmt.Println("=== Tag Generation ===")
	tagsResp, err := c.AnalyzeImageForTags(ctx, client.VisionTagsRequest{
		ImageURL: "https://rub.soulkyn.com/b6e4b172-b8d8-4552-9361-abd56a1e84ee.webp",
	})
	if err != nil {
		log.Printf("Tag generation failed: %v", err)
	} else {
		fmt.Printf("Tags: %s\n", tagsResp.Result)
		fmt.Printf("Tokens: %d total\n\n", tagsResp.Usage.TotalTokens)
	}

	// Example 3: OCR Text Extraction
	fmt.Println("=== OCR Text Extraction ===")
	ocrResp, err := c.ExtractText(ctx, client.VisionOCRRequest{
		ImageURL: "https://qianwen-res.oss-accelerate.aliyuncs.com/Qwen3-VL/qwen3vl_4b_8b_text_instruct.jpg",
	})
	if err != nil {
		log.Printf("OCR failed: %v", err)
	} else {
		fmt.Printf("Extracted text: %s\n", ocrResp.Result)
		fmt.Printf("Tokens: %d total\n\n", ocrResp.Usage.TotalTokens)
	}

	// Example 4: Video Analysis
	fmt.Println("=== Video Analysis ===")
	videoResp, err := c.AnalyzeVideo(ctx, client.VisionVideoRequest{
		VideoURL:   "https://rub.soulkyn.com/d8d917b0-bb37-43d6-9d3f-19eee1547065.mp4",
		UserPrompt: ptrString("Describe what happens in this video."),
	})
	if err != nil {
		log.Printf("Video analysis failed: %v", err)
	} else {
		fmt.Printf("Video analysis: %s\n", videoResp.Result)
		fmt.Printf("Tokens: %d total\n\n", videoResp.Usage.TotalTokens)
	}

	// Example 5: Text Moderation
	fmt.Println("=== Text Moderation ===")
	textModResp, err := c.ModerateText(ctx, client.ModerationTextRequest{
		Prompt: "Generate a beautiful landscape with mountains and sunset",
	})
	if err != nil {
		log.Printf("Text moderation failed: %v", err)
	} else {
		fmt.Printf("Category: %s (score: %.2f)\n", textModResp.Category, textModResp.Score)
		fmt.Printf("Tokens: %d total\n\n", textModResp.Usage.TotalTokens)
	}

	// Example 6: Image Moderation
	fmt.Println("=== Image Moderation ===")
	imgModResp, err := c.ModerateMedia(ctx, client.ModerationMediaRequest{
		MediaURL: "https://rub.soulkyn.com/b6e4b172-b8d8-4552-9361-abd56a1e84ee.webp",
		IsVideo:  false,
	})
	if err != nil {
		log.Printf("Image moderation failed: %v", err)
	} else {
		fmt.Printf("Category: %s (score: %.2f)\n", imgModResp.Category, imgModResp.Score)
		fmt.Printf("Tokens: %d total\n\n", imgModResp.Usage.TotalTokens)
	}

	// Example 7: Video Moderation
	fmt.Println("=== Video Moderation ===")
	videoModResp, err := c.ModerateMedia(ctx, client.ModerationMediaRequest{
		MediaURL: "https://rub.soulkyn.com/d8d917b0-bb37-43d6-9d3f-19eee1547065.mp4",
		IsVideo:  true,
	})
	if err != nil {
		log.Printf("Video moderation failed: %v", err)
	} else {
		fmt.Printf("Category: %s (score: %.2f)\n", videoModResp.Category, videoModResp.Score)
		fmt.Printf("Tokens: %d total\n", videoModResp.Usage.TotalTokens)
	}
}
