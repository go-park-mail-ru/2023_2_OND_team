package v1

import (
	"context"

	pb "github.com/go-park-mail-ru/2023_2_OND_team/internal/roll/api"
	roll "github.com/go-park-mail-ru/2023_2_OND_team/internal/roll/internal/pkg/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func FromAnswerProtoToService(ans *pb.RollAnswer, userID int) roll.RollAnswer {
	return roll.RollAnswer{
		UserID:     userID,
		RollID:     int(ans.RollID),
		QuestionID: int(ans.QuestionID),
		Answer:     ans.Answer,
	}
}

func (h *ServerGRPC) FillRoll(ctx context.Context, answers *pb.RollAnswers) (*pb.Nothing, error) {

	currUserID, _ := ctx.Value("userID").(int)
	if currUserID == 0 {
		return &pb.Nothing{}, status.Error(codes.Unauthenticated, "no_auth")
	}

	rollAnswers := make([]roll.RollAnswer, 0)
	for _, answer := range answers.Answers {
		rollAnswers = append(rollAnswers, FromAnswerProtoToService(answer, currUserID))
	}

	return &pb.Nothing{}, h.rollService.FillRoll(ctx, rollAnswers)
}

func HasAnsweredFromServiceToProto(hasAnswered bool) *pb.HasAnswered {
	return &pb.HasAnswered{HasAnswered: hasAnswered}
}

func (h *ServerGRPC) HasUserFilledRoll(context.Context, *pb.RollUserData) (*pb.HasAnswered, error) {

	return nil, status.Errorf(codes.Unimplemented, "method HasUserFilledRoll not implemented")
}
func (h *ServerGRPC) GetHistStat(context.Context, *pb.HisStatRequest) (*pb.RollStatHist, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Form1Stat not implemented")
}
