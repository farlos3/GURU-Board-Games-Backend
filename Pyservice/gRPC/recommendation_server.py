# Pyservice/gRPC/recommendation_server.py
import grpc
from concurrent import futures
import logging
from datetime import datetime

# Import generated protobuf files
import recommendation_pb2
import recommendation_pb2_grpc

# Setup logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

class RecommendationServiceImpl(recommendation_pb2_grpc.RecommendationServiceServicer):
    def __init__(self):
        self.boardgames_data = []
        self.user_actions = []
        logger.info("🚀 RecommendationService initialized")
    
    def SendAllBoardgames(self, request, context):
        logger.info(f"📥 [gRPC Server] Received `SendAllBoardgames` request")
        logger.info(f"   Number of boardgames: {len(request.boardgames)}")
        
        try:
            self.boardgames_data = []
            for boardgame in request.boardgames:
                bg_data = {
                    'id': boardgame.id,
                    'title': boardgame.title,
                    'description': boardgame.description,
                    'min_players': boardgame.min_players,
                    'max_players': boardgame.max_players,
                    'play_time_min': boardgame.play_time_min,
                    'play_time_max': boardgame.play_time_max,
                    'categories': boardgame.categories,
                    'rating_avg': boardgame.rating_avg,
                    'rating_count': boardgame.rating_count,
                    'popularity_score': boardgame.popularity_score,
                    'image_url': boardgame.image_url
                }
                self.boardgames_data.append(bg_data)
       
            logger.info(f"✅ [gRPC Server] Successfully received {len(request.boardgames)} boardgames")
            print("=" * 115)
            
            return recommendation_pb2.Response(
                success=True,
                message=f"Successfully received {len(request.boardgames)} boardgames"
            )
        except Exception as e:
            logger.error(f"❌ [gRPC Server] Error processing boardgames: {str(e)}")
            return recommendation_pb2.Response(
                success=False,
                message=f"Error processing boardgames: {str(e)}"
            )
    
    def SendUserAction(self, request, context):
        """
        📥 รับข้อมูล user action สำหรับการเก็บพฤติกรรมผู้ใช้
        ใช้สำหรับการฝึก ML หรือส่งไปยัง Elasticsearch ในอนาคต

        🔹 ข้อมูลที่ได้รับ:
            - user_id: ไอดีผู้ใช้
            - boardgame_id: ไอดีของบอร์ดเกม
            - action_type: ประเภทการกระทำ (เช่น like, play, rate)
            - action_value: ค่าที่เกี่ยวข้องกับ action เช่น คะแนน 5 ดาว
            - action_time: เวลาที่เกิด action นั้น

        🔸 ตอนนี้ข้อมูลจะถูกเก็บไว้ในหน่วยความจำ (in-memory)
        คือ self.user_actions (เป็น list ธรรมดา ยังไม่เชื่อมต่อฐานข้อมูลจริง)
        """
        logger.info(f"📥 [gRPC Server] Received SendUserAction request")
        logger.info(f"   - User ID: {request.user_id}")
        logger.info(f"   - Boardgame ID: {request.boardgame_id}")
        logger.info(f"   - Action Type: {request.action_type}")
        logger.info(f"   - Action Value: {request.action_value}")
        
        try:
            # แปลง timestamp
            action_time = datetime.fromtimestamp(
                request.action_time.seconds + request.action_time.nanos / 1e9
            )
            
            # เก็บข้อมูล user action ไว้ใน memory (mock)
            action_data = {
                'user_id': request.user_id,
                'boardgame_id': request.boardgame_id,
                'action_type': request.action_type,
                'action_value': request.action_value,
                'action_time': action_time
            }
            self.user_actions.append(action_data)
            
            logger.info(f"   - Action Time: {action_time}")
            logger.info(f"✅ [gRPC Server] User action stored successfully")
            
            return recommendation_pb2.Response(
                success=True,
                message="User action received and processed successfully"
            )
            
        except Exception as e:
            logger.error(f"❌ [gRPC Server] Error processing user action: {str(e)}")
            return recommendation_pb2.Response(
                success=False,
                message=f"Error processing user action: {str(e)}"
            )

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    recommendation_pb2_grpc.add_RecommendationServiceServicer_to_server(
        RecommendationServiceImpl(), server
    )
    port = '8001'
    server.add_insecure_port(f'[::]:{port}')
    server.start()
    logger.info(f"🚀 gRPC Server started on port {port}")
    logger.info("🔌 Ready to receive requests from Go service...")
    
    try:
        server.wait_for_termination()
    except KeyboardInterrupt:
        logger.info("🛑 Server stopped by user")
        server.stop(0)

if __name__ == '__main__':
    serve()