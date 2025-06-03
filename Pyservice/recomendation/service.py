from typing import List, Optional
from datetime import datetime
import logging
from pydantic import BaseModel
from .indexing import create_indices
from .setting import boardgame_index_name, user_action_index_name
from connection.connection import client

# ตั้งค่า logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class UserAction(BaseModel):
    user_id: str
    boardgame_id: str
    action_type: str  # like, view, play, rate
    action_value: float  # rating score (1-5) หรือ 1 สำหรับ like/view/play
    action_detail: Optional[str] = None  # review message หรือข้อมูลเพิ่มเติม
    action_time: datetime = datetime.now()

class Boardgame(BaseModel):
    id: int
    title: str
    description: str
    min_players: int
    max_players: int
    play_time_min: int
    play_time_max: int
    categories: str
    rating_avg: float
    rating_count: int
    popularity_score: float
    image_url: str

class RecommendationService:
    def __init__(self):
        self.boardgames: List[Boardgame] = []
        self.user_actions: List[UserAction] = []
        # สร้าง indices ใน Elasticsearch
        try:
            create_indices()
            logger.info("✅ Elasticsearch indices created successfully")
        except Exception as e:
            logger.error(f"❌ Failed to create Elasticsearch indices: {e}")

    def add_user_action(self, action: UserAction) -> bool:
        """Add a new user action to the system and Elasticsearch"""
        try:
            # ตรวจสอบ action_type
            if action.action_type not in ["like", "view", "play", "rate"]:
                logger.error(f"❌ Invalid action_type: {action.action_type}")
                return False

            # ตรวจสอบ action_value
            if action.action_type == "rate" and not (1 <= action.action_value <= 5):
                logger.error(f"❌ Invalid rating value: {action.action_value}")
                return False
            elif action.action_type in ["like", "view", "play"] and action.action_value != 1:
                logger.error(f"❌ Invalid action_value for {action.action_type}: {action.action_value}")
                return False

            # บันทึก action
            self.user_actions.append(action)
            
            # อัพเดท popularity_score ของบอร์ดเกม
            self._update_boardgame_popularity(action)

            # บันทึกลง Elasticsearch
            response = client.index(
                index=user_action_index_name,
                body=action.dict()
            )
            logger.info(f"✅ User action added to Elasticsearch: {response['_id']}")
            return True
        except Exception as e:
            logger.error(f"❌ Error adding user action: {e}")
            return False

    def _update_boardgame_popularity(self, action: UserAction) -> None:
        """Update boardgame popularity score based on user action"""
        try:
            # ค้นหาบอร์ดเกม
            response = client.get(
                index=boardgame_index_name,
                id=action.boardgame_id
            )
            boardgame = Boardgame(**response['_source'])

            # คำนวณ popularity score
            weight = {
                "like": 2.0,
                "view": 0.5,
                "play": 1.5,
                "rate": 1.0
            }

            # อัพเดท popularity score
            boardgame.popularity_score += weight[action.action_type] * action.action_value

            # บันทึกลง Elasticsearch
            client.index(
                index=boardgame_index_name,
                id=action.boardgame_id,
                body=boardgame.dict()
            )
            logger.info(f"✅ Updated popularity score for boardgame {action.boardgame_id}")
        except Exception as e:
            logger.error(f"❌ Error updating boardgame popularity: {e}")

    def get_user_actions(self, user_id: str) -> List[UserAction]:
        """Get all actions for a specific user"""
        try:
            response = client.search(
                index=user_action_index_name,
                body={
                    "query": {
                        "term": {
                            "user_id": user_id
                        }
                    },
                    "sort": [
                        {"action_time": {"order": "desc"}}
                    ]
                }
            )
            
            actions = []
            for hit in response['hits']['hits']:
                actions.append(UserAction(**hit['_source']))
            
            logger.info(f"✅ Retrieved {len(actions)} actions for user {user_id}")
            return actions
        except Exception as e:
            logger.error(f"❌ Error getting user actions: {e}")
            return []

    def get_boardgame_actions(self, boardgame_id: str) -> List[UserAction]:
        """Get all actions for a specific boardgame"""
        try:
            response = client.search(
                index=user_action_index_name,
                body={
                    "query": {
                        "term": {
                            "boardgame_id": boardgame_id
                        }
                    },
                    "sort": [
                        {"action_time": {"order": "desc"}}
                    ]
                }
            )
            
            actions = []
            for hit in response['hits']['hits']:
                actions.append(UserAction(**hit['_source']))
            
            logger.info(f"✅ Retrieved {len(actions)} actions for boardgame {boardgame_id}")
            return actions
        except Exception as e:
            logger.error(f"❌ Error getting boardgame actions: {e}")
            return []

    def get_recommendations(self, user_id: str, limit: int = 10) -> List[Boardgame]:
        """Get boardgame recommendations for a user"""
        try:
            # TODO: Implement recommendation algorithm using Elasticsearch
            # For now, return top rated boardgames
            sorted_games = sorted(self.boardgames, key=lambda x: x.rating_avg, reverse=True)
            logger.info(f"✅ Got {len(sorted_games[:limit])} recommendations for user {user_id}")
            return sorted_games[:limit]
        except Exception as e:
            logger.error(f"❌ Error getting recommendations: {e}")
            return []

    def update_boardgames(self, boardgames: List[Boardgame]) -> bool:
        """Update the boardgames list and Elasticsearch"""
        try:
            self.boardgames = boardgames
            # อัพเดทข้อมูลใน Elasticsearch
            for bg in boardgames:
                response = client.index(
                    index=boardgame_index_name,
                    id=str(bg.id),
                    body=bg.dict()
                )
                logger.info(f"✅ Boardgame {bg.id} updated in Elasticsearch: {response['_id']}")
            return True
        except Exception as e:
            logger.error(f"❌ Error updating boardgames: {e}")
            return False

    def get_all_boardgames(self) -> List[Boardgame]:
        """Get all boardgames from Elasticsearch"""
        try:
            # ดึงข้อมูลจาก Elasticsearch
            response = client.search(
                index=boardgame_index_name,
                body={
                    "query": {
                        "match_all": {}
                    },
                    "size": 10000  # เพิ่มขนาดเพื่อดึงข้อมูลมากกว่า 10 รายการ
                }
            )
            
            boardgames = []
            for hit in response['hits']['hits']:
                boardgames.append(Boardgame(**hit['_source']))
            logger.info(f"✅ Retrieved {len(boardgames)} boardgames from Elasticsearch")
            return boardgames
        except Exception as e:
            logger.error(f"❌ Error getting boardgames from Elasticsearch: {e}")
            logger.info("⚠️ Falling back to in-memory boardgames")
            return self.boardgames  # ถ้าเกิด error ให้ใช้ข้อมูลจาก memory แทน

    def get_popular_boardgames(self, limit: int) -> List[Boardgame]:
        """Get top N popular boardgames based on popularity_score"""
        try:
            # ดึงข้อมูลจาก Elasticsearch
            response = client.search(
                index=boardgame_index_name,
                body={
                    "size": limit,
                    "sort": [
                        {"popularity_score": {"order": "desc"}}
                    ]
                }
            )
            
            boardgames = []
            for hit in response['hits']['hits']:
                boardgames.append(Boardgame(**hit['_source']))
            
            logger.info(f"✅ Retrieved {len(boardgames)} popular boardgames from Elasticsearch")
            return boardgames
        except Exception as e:
            logger.error(f"❌ Error getting popular boardgames from Elasticsearch: {e}")
            logger.info("⚠️ Falling back to in-memory boardgames")
            # ถ้าเกิด error ให้ใช้ข้อมูลจาก memory แทน
            sorted_games = sorted(self.boardgames, key=lambda x: x.popularity_score, reverse=True)
            return sorted_games[:limit]

# Create a singleton instance
recommendation_service = RecommendationService() 