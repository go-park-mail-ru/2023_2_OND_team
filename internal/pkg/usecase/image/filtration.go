package image

import (
	"context"
	"errors"
	"strings"

	vision "cloud.google.com/go/vision/v2/apiv1"
	pb "cloud.google.com/go/vision/v2/apiv1/visionpb"
	validate "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/validation"
)

var (
	maxAnnotationsNumber int32 = 15
	explicitLabels             = []string{"goose", "duck"}
	ErrExplicitImage           = errors.New("image content doesn't comply with service policy")
)

type ImageFilter interface {
	Filter(ctx context.Context, imgBytes []byte, explicitLabels []string) error
}

type googleVision struct {
	visionClient *vision.ImageAnnotatorClient
	censor       validate.ProfanityCensor
}

func NewFilter(client *vision.ImageAnnotatorClient, censor validate.ProfanityCensor) *googleVision {
	return &googleVision{client, censor}
}

func CheckAnnotations(annotation *pb.SafeSearchAnnotation) bool {
	return annotation.GetAdult() >= pb.Likelihood_LIKELY ||
		annotation.GetMedical() >= pb.Likelihood_LIKELY ||
		annotation.GetRacy() >= pb.Likelihood_LIKELY ||
		annotation.GetViolence() >= pb.Likelihood_LIKELY ||
		annotation.GetSpoof() >= pb.Likelihood_LIKELY
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

func getTextDescription(resp *pb.AnnotateImageResponse) string {
	annotations := resp.GetTextAnnotations()
	if len(annotations) == 0 {
		return ""
	}
	return annotations[0].GetDescription()
}

func CheckExplicit(resp *pb.AnnotateImageResponse, explicitLabels []string, censor validate.ProfanityCensor) error {
	if CheckCertainLabels(explicitLabels, GetImageLabels(resp.GetLabelAnnotations())) ||
		CheckAnnotations(resp.GetSafeSearchAnnotation()) ||
		censor.IsProfane(getTextDescription(resp)) {
		return ErrExplicitImage
	}
	return nil
}

func (filter *googleVision) Filter(ctx context.Context, imgBytes []byte, explicitLabels []string) error {
	req := &pb.BatchAnnotateImagesRequest{
		Requests: []*pb.AnnotateImageRequest{
			{
				Image: &pb.Image{Content: imgBytes},
				Features: []*pb.Feature{
					{Type: pb.Feature_LABEL_DETECTION, MaxResults: maxAnnotationsNumber},
					{Type: pb.Feature_SAFE_SEARCH_DETECTION, MaxResults: maxAnnotationsNumber},
					{Type: pb.Feature_TEXT_DETECTION},
				},
			},
		},
	}
	resp, err := filter.visionClient.BatchAnnotateImages(ctx, req)
	if err != nil {
		return err
	}
	return CheckExplicit(resp.GetResponses()[0], explicitLabels, filter.censor)
}
