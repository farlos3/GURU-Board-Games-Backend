package recommendation

import (
    "strconv"
    "github.com/gofiber/fiber/v2"
    recommendationpb "guru-game/gRPC"
    "guru-game/internal/boardgame/service_board" 
    "guru-game/models"
)

type RecommendationHandler struct {
    grpcClient       *GRPCRecommendationClient
    boardgameService *service_board.BoardgameService
}

func NewRecommendationHandler(client *GRPCRecommendationClient, bgService *service_board.BoardgameService) *RecommendationHandler {
    return &RecommendationHandler{
        grpcClient:       client,
        boardgameService: bgService,
    }
}

// ส่งข้อมูล boardgames ทั้งหมดไปยัง Python ML service
func (h *RecommendationHandler) HandleSendAllBoardgames(c *fiber.Ctx) error {
    // ดึงข้อมูล boardgames จาก DB ผ่าน service
    boardgames, err := h.boardgameService.GetAllBoardgames()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "failed to fetch boardgames from database",
        })
    }

    // แปลงข้อมูลเป็น protobuf objects
    pbBoardgames := convertToPBBoardgames(boardgames)

    // ส่งข้อมูลบอร์ดเกมทั้งหมดไปยัง Python ML service ผ่าน gRPC
    err = h.grpcClient.SendAllBoardgames(pbBoardgames)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "failed to send boardgames to Python ML service",
            "details": err.Error(),
        })
    }

    return c.JSON(fiber.Map{
        "success": true,
        "message": "All boardgames sent to Python ML service successfully",
        "count":   len(pbBoardgames),
    })
}

// ขอ recommendations จาก Python ML service
func (h *RecommendationHandler) HandleGetRecommendations(c *fiber.Ctx) error {
    // รับพารามิเตอร์จาก URL
    userID := c.Query("user_id")
    limitStr := c.Query("limit", "10") // default 10
    
    if userID == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "user_id is required",
        })
    }
    
    limit, err := strconv.Atoi(limitStr)
    if err != nil || limit <= 0 {
        limit = 10
    }
    
    // ขอ recommendations จาก Python ML service
    recommendations, err := h.grpcClient.GetRecommendations(userID, int32(limit))
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "failed to get recommendations from Python ML service",
            "details": err.Error(),
        })
    }
    
    // แปลง protobuf กลับเป็น JSON response
    result := make([]map[string]interface{}, len(recommendations))
    for i, rec := range recommendations {
        result[i] = map[string]interface{}{
            "id":              rec.Id,
            "title":           rec.Title,
            "description":     rec.Description,
            "min_players":     rec.MinPlayers,
            "max_players":     rec.MaxPlayers,
            "play_time_min":   rec.PlayTimeMin,
            "play_time_max":   rec.PlayTimeMax,
            "categories":      rec.Categories,
            "rating_avg":      rec.RatingAvg,
            "rating_count":    rec.RatingCount,
            "popularity_score": rec.PopularityScore,
            "image_url":       rec.ImageUrl,
        }
    }
    
    return c.JSON(fiber.Map{
        "success": true,
        "user_id": userID,
        "recommendations": result,
        "count": len(result),
    })
}

// แปลง slice ของ models.BoardGame เป็น slice ของ protobuf Boardgame
func convertToPBBoardgames(boardgames []models.BoardGame) []*recommendationpb.Boardgame {
    var pbBoardgames []*recommendationpb.Boardgame
    for _, bg := range boardgames {
        pbBoardgames = append(pbBoardgames, &recommendationpb.Boardgame{
            Id:              int32(bg.ID),
            Title:           bg.Title,
            Description:     bg.Description,
            MinPlayers:      int32(bg.MinPlayers),
            MaxPlayers:      int32(bg.MaxPlayers),
            PlayTimeMin:     int32(bg.PlayTimeMin),
            PlayTimeMax:     int32(bg.PlayTimeMax),
            Categories:      bg.Categories,
            RatingAvg:       bg.RatingAvg,
            RatingCount:     int32(bg.RatingCount),
            PopularityScore: bg.PopularityScore,
            ImageUrl:        bg.ImageURL,
        })
    }
    return pbBoardgames
}