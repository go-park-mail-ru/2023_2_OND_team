package image

import (
	"context"
	"errors"
	"strings"

	pb "cloud.google.com/go/vision/v2/apiv1/visionpb"
)

var (
	maxAnnotationsNumber int32 = 15
	explicitLabels             = []string{"goose", "duck"}
	ErrExplicitImage           = errors.New("Image content doesn't comply with service policy")
)

func CheckAnnotations(annotation *pb.SafeSearchAnnotation) bool {
	if annotation.GetAdult() >= pb.Likelihood_LIKELY ||
		annotation.GetMedical() >= pb.Likelihood_LIKELY ||
		annotation.GetRacy() >= pb.Likelihood_LIKELY ||
		annotation.GetViolence() >= pb.Likelihood_LIKELY ||
		annotation.GetSpoof() >= pb.Likelihood_LIKELY {
		return true
	}
	return false
}

func GetImageLabels(annotations []*pb.EntityAnnotation) []string {
	imgLabels := make([]string, 0, len(annotations))
	for _, label := range annotations {
		imgLabels = append(imgLabels, label.GetDescription())
	}
	return imgLabels
}

func CheckCertainLabels(explicitLabels, imgLabels []string) bool {
	for _, label := range explicitLabels {
		if HasExplicitLabel(label, imgLabels) {
			return true
		}
	}
	return false
}

func HasExplicitLabel(explicitLabel string, imgLabels []string) bool {
	for _, label := range imgLabels {
		if strings.Contains(strings.ToLower(label), strings.ToLower(explicitLabel)) {
			return true
		}
	}
	return false
}

func CheckExplicit(resp *pb.AnnotateImageResponse, explicitLabels []string) error {
	if CheckCertainLabels(explicitLabels, GetImageLabels(resp.GetLabelAnnotations())) ||
		CheckAnnotations(resp.GetSafeSearchAnnotation()) {
		return ErrExplicitImage
	}
	return nil
}

func (img *imageCase) FilterImage(ctx context.Context, imgBytes []byte, explicitLabels []string) error {
	req := &pb.BatchAnnotateImagesRequest{
		Requests: []*pb.AnnotateImageRequest{
			{
				Image: &pb.Image{Content: imgBytes},
				Features: []*pb.Feature{
					{Type: pb.Feature_LABEL_DETECTION, MaxResults: maxAnnotationsNumber},
					{Type: pb.Feature_SAFE_SEARCH_DETECTION, MaxResults: maxAnnotationsNumber},
				},
			},
		},
	}
	resp, err := img.visionClient.BatchAnnotateImages(ctx, req)
	if err != nil {
		return err
	}

	return CheckExplicit(resp.GetResponses()[0], explicitLabels)
}
