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

    def get_recommendations(
        self, 
        user_id: str, 
        limit: int = 10,
        user_actions: Optional[List[UserAction]] = None,
        user_categories: Optional[List[str]] = None
    ) -> List[Boardgame]:
        """Get boardgame recommendations for a user based on their behavior"""
        try:
            logger.info(f"\n🔍 ===== Generating Recommendations for User {user_id} =====")
            
            # Get user's actions from Elasticsearch if not provided
            if user_actions is None:
                user_actions = self.get_user_actions(user_id)
            logger.info(f"📊 Found {len(user_actions)} user actions")
            
            # Get all boardgames
            all_boardgames = self.get_all_boardgames()
            logger.info(f"🎲 Total boardgames in system: {len(all_boardgames)}")
            
            if len(all_boardgames) == 0:
                logger.error("❌ No boardgames found in the system!")
                return []
            
            # Log some sample boardgames
            logger.info("\n📋 Sample Boardgames:")
            for bg in all_boardgames[:3]:
                logger.info(f"  - ID: {bg.id}, Title: {bg.title}, Categories: {bg.categories}")
            
            # Create a scoring dictionary for boardgames
            boardgame_scores = {}
            
            # Process user actions to build preferences
            user_preferences = {
                'categories': set(user_categories) if user_categories else set(),
                'player_counts': set(),
                'play_times': set(),
                'ratings': {}
            }
            
            logger.info("\n👤 User Preferences:")
            logger.info(f"Initial Categories: {user_preferences['categories']}")
            
            for action in user_actions:
                logger.info(f"\n📝 Processing action: {action.action_type} for boardgame {action.boardgame_id}")
                
                # Find the corresponding boardgame
                boardgame = next((bg for bg in all_boardgames if str(bg.id) == action.boardgame_id), None)
                if not boardgame:
                    logger.warning(f"⚠️ Boardgame {action.boardgame_id} not found")
                    continue
                
                logger.info(f"Found boardgame: {boardgame.title}")
                
                # Update user preferences based on action
                if action.action_type == "like" or action.action_type == "favorite":
                    # Add categories to preferences if not already provided
                    if not user_categories and boardgame.categories:
                        categories = [cat.strip() for cat in boardgame.categories.split(",")]
                        user_preferences['categories'].update(categories)
                        logger.info(f"🏷️ Added categories to preferences: {categories}")
                    
                    # Add player count range
                    user_preferences['player_counts'].add(boardgame.min_players)
                    user_preferences['player_counts'].add(boardgame.max_players)
                    logger.info(f"👥 Added player count range: {boardgame.min_players}-{boardgame.max_players}")
                    
                    # Add play time range
                    user_preferences['play_times'].add(boardgame.play_time_min)
                    user_preferences['play_times'].add(boardgame.play_time_max)
                    logger.info(f"⏱️ Added play time range: {boardgame.play_time_min}-{boardgame.play_time_max}")
                
                elif action.action_type == "rating":
                    # Store rating for this boardgame
                    user_preferences['ratings'][action.boardgame_id] = action.action_value
                    logger.info(f"⭐ Stored rating: {action.action_value} for boardgame {action.boardgame_id}")
            
            logger.info("\n📊 User Preferences Summary:")
            logger.info(f"Categories: {user_preferences['categories']}")
            logger.info(f"Player Counts: {user_preferences['player_counts']}")
            logger.info(f"Play Times: {user_preferences['play_times']}")
            logger.info(f"Ratings: {user_preferences['ratings']}")
            
            # Score each boardgame based on user preferences
            logger.info("\n🎯 Scoring Boardgames:")
            for boardgame in all_boardgames:
                score = 0.0
                
                # Skip boardgames the user has already interacted with
                if str(boardgame.id) in [action.boardgame_id for action in user_actions]:
                    logger.info(f"Skipping boardgame {boardgame.id} (already interacted)")
                    continue
                
                logger.info(f"\nScoring boardgame {boardgame.id} - {boardgame.title}")
                
                # Category matching (weight: 2.0)
                if boardgame.categories and user_preferences['categories']:
                    categories = [cat.strip() for cat in boardgame.categories.split(",")]
                    matching_categories = user_preferences['categories'].intersection(set(categories))
                    if matching_categories:
                        category_score = 2.0 * (len(matching_categories) / len(categories))
                        score += category_score
                        logger.info(f"🎯 Category match: {matching_categories} (score: {category_score:.2f})")
                
                # Player count matching (weight: 1.5)
                if user_preferences['player_counts']:
                    min_players_match = any(boardgame.min_players <= count <= boardgame.max_players 
                                         for count in user_preferences['player_counts'])
                    if min_players_match:
                        score += 0.5
                        logger.info(f"👥 Player count match (score: 1.50)")
                
                # Play time matching (weight: 1.5)
                if user_preferences['play_times']:
                    play_time_match = any(boardgame.play_time_min <= time <= boardgame.play_time_max 
                                       for time in user_preferences['play_times'])
                    if play_time_match:
                        score += 0.5
                        logger.info(f"⏱️ Play time match (score: 1.50)")
                
                # Rating consideration (weight: 1.0)
                if boardgame.rating_avg > 0:
                    rating_score = 0.5 * (boardgame.rating_avg / 5.0)
                    score += rating_score
                    logger.info(f"⭐ Rating consideration: {boardgame.rating_avg} (score: {rating_score:.2f})")
                
                # Popularity consideration (weight: 1.0)
                if boardgame.popularity_score > 0:
                    popularity_score = 0.5 * (boardgame.popularity_score / 100.0)
                    score += popularity_score
                    logger.info(f"🔥 Popularity consideration: {boardgame.popularity_score} (score: {popularity_score:.2f})")
                
                boardgame_scores[boardgame.id] = score
                logger.info(f"Total score: {score:.2f}")
            
            # Sort boardgames by score
            sorted_boardgames = sorted(boardgame_scores.items(), key=lambda x: x[1], reverse=True)

            # สร้าง dictionary สำหรับค้นหา Boardgame object จาก ID ได้ง่ายๆ
            boardgame_by_id = {str(bg.id): bg for bg in all_boardgames} # ใช้ str(bg.id) ให้ตรงกับ key ใน boardgame_scores

            recommendations = []
            for boardgame_id, score in sorted_boardgames[:limit]:
                # ดึง Boardgame object จาก dictionary
                boardgame = boardgame_by_id.get(str(boardgame_id))
                if boardgame:
                    recommendations.append(boardgame)
                    # แก้ไขบรรทัด log ให้ใช้ boardgame object ที่ถูกต้อง
                    logger.info(f"{len(recommendations)}. Boardgame {boardgame.id} - {boardgame.title} - Score: {score:.2f}")

            logger.info("===========================================\n")
            return recommendations
            
        except Exception as e:
            logger.error(f"❌ Error getting recommendations: {e}")
            logger.exception("Detailed error:")
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
