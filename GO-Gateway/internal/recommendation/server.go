package recommendation

import (
	"context"
	"log"
	"fmt"
	recommendationpb "guru-game/gRPC"
)

// สร้าง struct RecommendationServer ที่ implement interface ของ gRPC
type RecommendationServer struct {
	recommendationpb.UnimplementedRecommendationServiceServer
}

// ตัวอย่าง implement ฟังก์ชันใน gRPC service
func (s *RecommendationServer) GetRecommendations(ctx context.Context, req *recommendationpb.RecommendationRequest) (*recommendationpb.RecommendationResponse, error) {
	log.Printf("📤 [gRPC Server] GetRecommendations called")
	log.Printf("   - User ID: %s", req.UserId)
	log.Printf("   - Limit: %d", req.Limit)
	
	// TODO: ใส่ logic สำหรับสร้าง recommendation ตรงนี้
	
	log.Printf("✅ [gRPC Server] Returning empty recommendations for now")
	
	// ตัวอย่าง return ค่าเปล่า
	return &recommendationpb.RecommendationResponse{
		Boardgames: []*recommendationpb.Boardgame{}, // empty array สำหรับตอนนี้
	}, nil
}

func (s *RecommendationServer) SendAllBoardgames(ctx context.Context, req *recommendationpb.BoardgamesRequest) (*recommendationpb.Response, error) {
	log.Printf("📥 [gRPC Server] SendAllBoardgames called")
	log.Printf("   - Received %d boardgames from client", len(req.Boardgames))
	
	// Log รายละเอียดของ boardgames ที่ได้รับ (แสดงแค่ 5 อันแรก)
	for i, bg := range req.Boardgames {
		if i < 5 {
			log.Printf("   - Boardgame %d: %s (ID: %d)", i+1, bg.Title, bg.Id)
			log.Printf("     Players: %d-%d, Time: %d-%d min, Rating: %.2f", 
				bg.MinPlayers, bg.MaxPlayers, bg.PlayTimeMin, bg.PlayTimeMax, bg.RatingAvg)
		}
	}
	
	if len(req.Boardgames) > 5 {
		log.Printf("   ... and %d more boardgames", len(req.Boardgames)-5)
	}
	
	// TODO: ทำอะไรกับข้อมูลที่ได้รับ
	// เช่น:
	// - บันทึกลง database สำหรับ recommendation system
	// - ส่งไป ML model เพื่อ train
	// - ประมวลผล categories, ratings สำหรับ algorithm
	
	log.Printf("✅ [gRPC Server] Successfully processed %d boardgames", len(req.Boardgames))
	
	return &recommendationpb.Response{
		Success: true,
		Message: fmt.Sprintf("Successfully received and processed %d boardgames", len(req.Boardgames)),
	}, nil
}