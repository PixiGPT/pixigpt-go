package client

import (
	"context"
	"encoding/json"
)

// AnalyzeImage analyzes an image and returns a detailed description.
//
// The server downloads and preprocesses the image (resize, convert to JPEG).
// For soulkyn.com URLs, Cloudflare bypass is automatically applied.
//
// Example:
//
//	resp, err := client.AnalyzeImage(ctx, VisionAnalyzeRequest{
//	    ImageURL: "https://example.com/image.jpg",
//	    UserPrompt: ptrString("Describe this in detail."),
//	})
func (c *Client) AnalyzeImage(ctx context.Context, req VisionAnalyzeRequest) (*VisionAnalyzeResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	var resp VisionAnalyzeResponse
	if err := c.doRequest(ctx, "POST", "/vision/analyze", body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// AnalyzeImageForTags generates comma-separated tags for an image.
//
// Returns short tags suitable for categorization and search.
//
// Example:
//
//	resp, err := client.AnalyzeImageForTags(ctx, VisionTagsRequest{
//	    ImageURL: "https://example.com/image.jpg",
//	})
func (c *Client) AnalyzeImageForTags(ctx context.Context, req VisionTagsRequest) (*VisionTagsResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	var resp VisionTagsResponse
	if err := c.doRequest(ctx, "POST", "/vision/tags", body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// ExtractText performs OCR on an image and returns extracted text.
//
// Preserves structure (tables, lists, hierarchy) and uses high detail mode.
//
// Example:
//
//	resp, err := client.ExtractText(ctx, VisionOCRRequest{
//	    ImageURL: "https://example.com/document.jpg",
//	})
func (c *Client) ExtractText(ctx context.Context, req VisionOCRRequest) (*VisionOCRResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	var resp VisionOCRResponse
	if err := c.doRequest(ctx, "POST", "/vision/ocr", body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// AnalyzeVideo analyzes a video and returns a description of the content.
//
// Videos must be under 10MB. The server performs size check via HEAD request.
// For soulkyn.com URLs, Cloudflare bypass is automatically applied.
//
// Example:
//
//	resp, err := client.AnalyzeVideo(ctx, VisionVideoRequest{
//	    VideoURL: "https://example.com/video.mp4",
//	    UserPrompt: ptrString("Describe what happens."),
//	})
func (c *Client) AnalyzeVideo(ctx context.Context, req VisionVideoRequest) (*VisionVideoResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	var resp VisionVideoResponse
	if err := c.doRequest(ctx, "POST", "/vision/video", body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// ModerateText classifies text content into 11 categories with confidence scores.
//
// Categories: UNDERAGE_SEXUAL (priority), JAILBREAK, SUICIDE_SELF_HARM, PII,
// COPYRIGHT_VIOLATION, VIOLENT, ILLEGAL_ACTS, UNETHICAL, HATE_SPEECH,
// SEXUAL_ADULT, SAFE.
//
// Score ranges: 1.00 = perfect match, 0.90-0.99 = very strong, 0.70-0.89 = strong,
// 0.50-0.69 = moderate, 0.00-0.49 = weak.
//
// Example:
//
//	resp, err := client.ModerateText(ctx, ModerationTextRequest{
//	    Prompt: "text to moderate",
//	})
func (c *Client) ModerateText(ctx context.Context, req ModerationTextRequest) (*ModerationResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	var resp ModerationResponse
	if err := c.doRequest(ctx, "POST", "/moderations", body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// ModerateMedia classifies image or video content into 11 categories with confidence scores.
//
// Same categories as ModerateText but with visual assessment.
// SEXUAL_ADULT = visible genitals OR active sex acts only.
// SAFE = cleavage, lingerie, bikinis, clothed, suggestive.
//
// Example:
//
//	resp, err := client.ModerateMedia(ctx, ModerationMediaRequest{
//	    MediaURL: "https://example.com/image.jpg",
//	    IsVideo: false,
//	})
func (c *Client) ModerateMedia(ctx context.Context, req ModerationMediaRequest) (*ModerationResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	var resp ModerationResponse
	if err := c.doRequest(ctx, "POST", "/moderations/media", body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
