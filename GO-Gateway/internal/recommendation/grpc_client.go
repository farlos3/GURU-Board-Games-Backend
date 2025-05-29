package recommendation

import (
    "context"
    "log"
    "time"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    recommendationpb "guru-game/gRPC"
)

type GRPCRecommendationClient struct {
    conn   *grpc.ClientConn
    client recommendationpb.RecommendationServiceClient
}

// สร้าง gRPC client connection
func NewGRPCRecommendationClient(address string) (*GRPCRecommendationClient, error) {
    log.Printf("🔌 [gRPC Client] Connecting to %s...", address)
    
    // สร้าง connection โดยไม่ใช้ TLS (insecure)
    conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Printf("❌ [gRPC Client] Failed to connect: %v", err)
        return nil, err
    }
    
    // สร้าง client จาก connection
    client := recommendationpb.NewRecommendationServiceClient(conn)
    
    log.Printf("✅ [gRPC Client] Connected to %s successfully", address)
    
    return &GRPCRecommendationClient{
        conn:   conn,
        client: client,
    }, nil
}

// ส่งข้อมูล boardgames ทั้งหมดไปยัง Python server
func (c *GRPCRecommendationClient) SendAllBoardgames(boardgames []*recommendationpb.Boardgame) error {
    log.Printf("📤 [gRPC Client] Sending %d boardgames to Python server...", len(boardgames))
    
    // สร้าง context กับ timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // สร้าง request
    req := &recommendationpb.BoardgamesRequest{
        Boardgames: boardgames,
    }
    
    // เรียก gRPC method
    resp, err := c.client.SendAllBoardgames(ctx, req)
    if err != nil {
        log.Printf("❌ [gRPC Client] Failed to send boardgames: %v", err)
        return err
    }
    
    log.Printf("✅ [gRPC Client] Response from Python server:")
    log.Printf("   - Success: %t", resp.Success)
    log.Printf("   - Message: %s", resp.Message)
    
    return nil
}

// ขอ recommendations จาก Python server
func (c *GRPCRecommendationClient) GetRecommendations(userID string, limit int32) ([]*recommendationpb.Boardgame, error) {
    log.Printf("📤 [gRPC Client] Requesting recommendations for user %s (limit: %d)...", userID, limit)
    
    // สร้าง context กับ timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    // สร้าง request
    req := &recommendationpb.RecommendationRequest{
        UserId: userID,
        Limit:  limit,
    }
    
    // เรียก gRPC method
    resp, err := c.client.GetRecommendations(ctx, req)
    if err != nil {
        log.Printf("❌ [gRPC Client] Failed to get recommendations: %v", err)
        return nil, err
    }
    
    log.Printf("✅ [gRPC Client] Received %d recommendations from Python server", len(resp.Boardgames))
    
    return resp.Boardgames, nil
}

// ปิด connection
func (c *GRPCRecommendationClient) Close() error {
    log.Printf("🔌 [gRPC Client] Closing connection...")
    return c.conn.Close()
}