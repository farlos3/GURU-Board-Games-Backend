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

def search_boardgames(
    search_query: Optional[str] = None,
    player_count: Optional[int] = None,
    play_time: Optional[int] = None,
    categories: Optional[List[str]] = None,
    limit: int = 10,
    page: int = 1
) -> List[Boardgame]:
    """Search boardgames with improved query logic (works with existing mapping)."""
    try:
        query = {
            "bool": {
                "must": [],
                "should": [],
                "minimum_should_match": 0
            }
        }

        if search_query:
            search_query = search_query.strip()
            
            if len(search_query) <= 2:
                # สำหรับคำสั้น ใช้หลายกลยุทธ์
                query["bool"]["should"].extend([
                    # Wildcard search - หาทุกที่ในชื่อ
                    {
                        "wildcard": {
                            "title": {
                                "value": f"*{search_query.lower()}*",
                                "boost": 3,
                                "case_insensitive": True
                            }
                        }
                    },
                    # Prefix search - หาคำที่ขึ้นต้น
                    {
                        "prefix": {
                            "title": {
                                "value": search_query.lower(),
                                "boost": 2,
                                "case_insensitive": True
                            }
                        }
                    },
                    # Fuzzy search - รองรับ typo
                    {
                        "fuzzy": {
                            "title": {
                                "value": search_query,
                                "fuzziness": "AUTO",
                                "boost": 1
                            }
                        }
                    },
                    # Match search - หาใน description และ categories
                    {
                        "multi_match": {
                            "query": search_query,
                            "fields": ["description^0.5", "categories^1"],
                            "fuzziness": "AUTO",
                            "boost": 0.5
                        }
                    }
                ])
                query["bool"]["minimum_should_match"] = 1
            else:
                # สำหรับคำยาว ใช้ multi_match หลัก
                query["bool"]["must"].append({
                    "multi_match": {
                        "query": search_query,
                        "fields": [
                            "title^4",        # น้ำหนักสูงสุดให้ title
                            "description^1",   # น้ำหนักปกติให้ description
                            "categories^2"     # น้ำหนักสูงให้ categories
                        ],
                        "fuzziness": "AUTO",
                        "operator": "or",
                        "minimum_should_match": "75%"
                    }
                })

        # Player count filter
        if player_count is not None and player_count > 0:
            query["bool"]["must"].extend([
                {"range": {"min_players": {"lte": player_count}}},
                {"range": {"max_players": {"gte": player_count}}}
            ])

        # Play time filter  
        if play_time is not None and play_time > 0:
            query["bool"]["must"].extend([
                {"range": {"play_time_min": {"lte": play_time}}},
                {"range": {"play_time_max": {"gte": play_time}}}
            ])

        # Categories filter
        if categories:
            # ลองทั้ง exact match และ partial match เผื่อ categories ไม่ได้เป็น keyword
            query["bool"]["must"].append({
                "bool": {
                    "should": [
                        {"terms": {"categories": [cat.lower() for cat in categories]}},
                        {"terms": {"categories": categories}},  # กรณี case-sensitive
                        {"match": {"categories": " ".join(categories)}}  # partial match
                    ],
                    "minimum_should_match": 1
                }
            })

        # Handle empty query - แสดงเกมยอดนิยม
        if not query["bool"]["must"] and not query["bool"]["should"]:
            query = {"match_all": {}}

        # Calculate pagination
        from_ = max(0, (page - 1) * limit)

        # Execute search
        response = client.search(
            index=boardgame_index_name,
            body={
                "query": query,
                "size": limit,
                "from": from_,
                "sort": [
                    {"_score": {"order": "desc"}},
                    {"popularity_score": {"order": "desc", "missing": "_last"}},
                    {"rating_avg": {"order": "desc", "missing": "_last"}}
                ],
                # เพิ่ม highlighting เพื่อดูว่าตรงกับอะไร
                "highlight": {
                    "fields": {
                        "title": {
                            "pre_tags": ["<mark>"],
                            "post_tags": ["</mark>"],
                            "fragment_size": 150
                        },
                        "description": {
                            "pre_tags": ["<mark>"],
                            "post_tags": ["</mark>"],
                            "fragment_size": 150,
                            "number_of_fragments": 1
                        }
                    }
                }
            }
        )

        boardgames = []
        for hit in response['hits']['hits']:
            boardgame_data = hit['_source']
            
            # เพิ่มข้อมูล debug
            boardgame_data['_search_score'] = hit['_score']
            if 'highlight' in hit:
                boardgame_data['_highlights'] = hit['highlight']
                
            boardgames.append(boardgame_data)

        logger.info(f"✅ Search completed: query='{search_query}', player_count={player_count}, play_time={play_time}, categories={categories}, results={len(boardgames)}")
        
        # Debug logging สำหรับกรณีไม่เจอผลลัพธ์
        if len(boardgames) == 0:
            logger.warning(f"⚠️ No results found. Query used: {query}")
            
            # ลองค้นหาแบบง่ายๆ เพื่อดูว่ามีข้อมูลหรือไม่
            if search_query:
                debug_response = client.search(
                    index=boardgame_index_name,
                    body={
                        "query": {"match_all": {}},
                        "size": 3
                    }
                )
                logger.info(f"Sample data in index: {[hit['_source'].get('title', 'No title') for hit in debug_response['hits']['hits']]}")

        return boardgames

    except Exception as e:
        logger.error(f"❌ Search error: {str(e)}")
        logger.error(f"Query that caused error: {query if 'query' in locals() else 'Query not constructed'}")
        return []
