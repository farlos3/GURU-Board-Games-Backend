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

# Define weights for different user actions and scoring components
ACTION_WEIGHTS = {
    "like": 1.0,
    "favorite": 2,
    "rating_multiplier": 0.5, 
    "category_match": 3.0,
    "player_count_match": 0.2,
    "play_time_match": 0.2,
    "rating_avg_consideration": 1,
    "popularity_consideration": 0,
    "similarity_impact": 2.0 # Weight for the similarity-based score component
}

# Helper function to calculate similarity between two boardgames (based on categories)
def calculate_similarity(bg1: Boardgame, bg2: Boardgame) -> float:
    """Calculates similarity between two boardgames based on categories (Jaccard Index)."""
    if not bg1.categories or not bg2.categories:
        return 0.0 # No categories, no similarity based on this

    categories1 = set([cat.strip() for cat in bg1.categories.split(",") if cat.strip()])
    categories2 = set([cat.strip() for cat in bg2.categories.split(",") if cat.strip()])

    if not categories1 or not categories2:
        return 0.0 # One or both have no valid categories

    intersection = categories1.intersection(categories2)
    union = categories1.union(categories2)

    if not union:
        return 0.0 # Should not happen if both sets are not empty, but just in case

    return len(intersection) / len(union)

class RecommendationService:
    def __init__(self):
        self.boardgames: List[Boardgame] = []
        self.user_actions: List[UserAction] = []
        # สร้าง indices ใน Elasticsearch
        try:
            create_indices()
            logger.info("✅ Elasticsearch indices created successfully")
            print("==========================================================================================")
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
            print(f"\n🔍 ===== Generating Recommendations for User {user_id} =====")

            # Get user's actions from Elasticsearch if not provided
            if user_actions is None:
                user_actions = self.get_user_actions(user_id)

            # Get all boardgames
            all_boardgames = self.get_all_boardgames()
            logger.info(f"🎲 Total boardgames in system: {len(all_boardgames)}")

            if len(all_boardgames) == 0:
                logger.error("❌ No boardgames found in the system!")
                return []

            # Create a dictionary for quick lookup of Boardgame objects by ID
            boardgame_by_id = {str(bg.id): bg for bg in all_boardgames}

            # Calculate preference score for each boardgame the user has interacted with
            user_boardgame_preference_scores = {}
            user_preferences = {
                'categories': set(user_categories) if user_categories else set(),
                'player_counts': set(),
                'play_times': set(),
            } # Ratings are now directly used for preference score

            for action in user_actions:
                bg_id_str = action.boardgame_id
                boardgame = boardgame_by_id.get(bg_id_str)

                if not boardgame:
                    logger.warning(f"⚠️ Boardgame {bg_id_str} from user action not found in all_boardgames")
                    continue

                logger.info(f"📝 Processing action: {action.action_type} for boardgame {boardgame.title}")

                # Update user preferences (for characteristic matching)
                if action.action_type == "like" or action.action_type == "favorite":
                    if not user_categories and boardgame.categories:
                        categories = [cat.strip() for cat in boardgame.categories.split(",") if cat.strip()]
                        user_preferences['categories'].update(categories)
                        logger.info(f"  🏷️ Added categories to preferences: {categories}")

                    user_preferences['player_counts'].add(boardgame.min_players)
                    user_preferences['player_counts'].add(boardgame.max_players)
                    logger.info(f"  👥 Added player count range: {boardgame.min_players}-{boardgame.max_players}")

                    user_preferences['play_times'].add(boardgame.play_time_min)
                    user_preferences['play_times'].add(boardgame.play_time_max)
                    logger.info(f"  ⏱️ Added play time range: {boardgame.play_time_min}-{boardgame.play_time_max}")

                # Calculate preference score for this specific boardgame based on actions
                current_preference_score = user_boardgame_preference_scores.get(bg_id_str, 0.0)
                if action.action_type == "like":
                    current_preference_score += ACTION_WEIGHTS["like"]
                    logger.info(f"  👍 Added {ACTION_WEIGHTS['like']} for like action")
                elif action.action_type == "favorite":
                    current_preference_score += ACTION_WEIGHTS["favorite"]
                    logger.info(f"  ⭐ Added {ACTION_WEIGHTS['favorite']} for favorite action")
                elif action.action_type == "rating" and action.action_value > 0:
                    rating_score = ACTION_WEIGHTS["rating_multiplier"] * (action.action_value / 5.0)
                    current_preference_score += rating_score
                    logger.info(f"  ⭐ Added {rating_score:.2f} for rating {action.action_value}")

                user_boardgame_preference_scores[bg_id_str] = current_preference_score
                logger.info(f"  Current preference score for {boardgame.title}: {current_preference_score:.2f}")

            # Score each boardgame based on user preferences and similarity
            boardgame_scores = {}
            interacted_boardgame_ids = set(user_boardgame_preference_scores.keys())

            for boardgame in all_boardgames:
                bg_id_str = str(boardgame.id)

                # Skip boardgames the user has already interacted with (don't recommend them again)
                if bg_id_str in interacted_boardgame_ids:
                    logger.info(f"Skipping boardgame {boardgame.id} (already interacted)")
                    continue

                score = 0.0

                # --- 1. Characteristic-Based Scoring (Existing Logic) ---

                # Category matching
                if boardgame.categories and user_preferences['categories']:
                    categories = set([cat.strip() for cat in boardgame.categories.split(",") if cat.strip()])
                    if categories:
                        matching_categories = user_preferences['categories'].intersection(categories)
                        if matching_categories:
                            # Use min of boardgame categories or user preferred categories count for denominator
                            # to avoid division by zero or very small numbers if user_preferences['categories'] is small
                            denominator = min(len(categories), len(user_preferences['categories'])) if min(len(categories), len(user_preferences['categories'])) > 0 else 1
                            category_score = ACTION_WEIGHTS["category_match"] * (len(matching_categories) / denominator)
                            score += category_score
                            logger.info(f"  🎯 Category match: {matching_categories} (score: {category_score:.2f})")


                # Player count matching
                if user_preferences['player_counts']:
                    min_players_match = any(boardgame.min_players <= count <= boardgame.max_players
                                         for count in user_preferences['player_counts'])
                    if min_players_match:
                        score += ACTION_WEIGHTS["player_count_match"]
                        logger.info(f"  👥 Player count match (score: {ACTION_WEIGHTS['player_count_match']:.2f})")

                # Play time matching
                if user_preferences['play_times']:
                    play_time_match = any(boardgame.play_time_min <= time <= boardgame.play_time_max
                                       for time in user_preferences['play_times'])
                    if play_time_match:
                        score += ACTION_WEIGHTS["play_time_match"]
                        logger.info(f"  ⏱️ Play time match (score: {ACTION_WEIGHTS['play_time_match']:.2f})")

                # Rating consideration
                if boardgame.rating_avg > 0:
                    rating_score = ACTION_WEIGHTS["rating_avg_consideration"] * (boardgame.rating_avg / 5.0)
                    score += rating_score

                # Popularity consideration
                if boardgame.popularity_score > 0:
                    popularity_score = ACTION_WEIGHTS["popularity_consideration"] * (boardgame.popularity_score / 100.0)
                    score += popularity_score

                # --- 2. Similarity-Based Scoring (New Logic) ---
                similarity_score_component = 0.0
                if user_boardgame_preference_scores:
                    logger.info("  Calculating similarity-based score component:")
                    for interacted_bg_id_str, preference_score in user_boardgame_preference_scores.items():
                        interacted_boardgame = boardgame_by_id.get(interacted_bg_id_str)
                        if interacted_boardgame:
                            similarity = calculate_similarity(boardgame, interacted_boardgame)
                            similarity_score_component += similarity * preference_score * ACTION_WEIGHTS["similarity_impact"]
                            # logger.info(f"    Similarity with {interacted_boardgame.title} ({interacted_bg_id_str}): {similarity:.2f}, Preference Score: {preference_score:.2f}, Component: {similarity * preference_score * ACTION_WEIGHTS['similarity_impact']:.2f}") # Too verbose

                score += similarity_score_component
                if similarity_score_component > 0:
                     logger.info(f"  ✨ Similarity-based score component added: {similarity_score_component:.2f}")


                boardgame_scores[bg_id_str] = score

            # Sort boardgames by score
            sorted_boardgames = sorted(boardgame_scores.items(), key=lambda x: x[1], reverse=True)

            recommendations = []
            logger.info(f"✅ Generated {len(sorted_boardgames[:limit])} recommendations")
            for boardgame_id_str, score in sorted_boardgames[:limit]:
                # Retrieve the full Boardgame object using the ID
                boardgame = boardgame_by_id.get(boardgame_id_str)
                if boardgame:
                    recommendations.append(boardgame)
                    logger.info(f"Boardgame {boardgame.id} - {boardgame.title} - Score: {score:.2f}")

            logger.info("===========================================")
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
            # Convert categories to list if it's a string
            if isinstance(categories, str):
                try:
                    # Try to parse as JSON if it's a string representation of a list
                    import json
                    categories = json.loads(categories)
                except json.JSONDecodeError:
                    # If not JSON, treat as comma-separated string
                    categories = [cat.strip() for cat in categories.split(',')]
            
            # Clean up categories - remove any JSON-like formatting
            cleaned_categories = []
            for cat in categories:
                if isinstance(cat, str):
                    # Remove any JSON-like formatting
                    cat = cat.strip('[]"\'')
                    if cat:
                        cleaned_categories.append(cat.lower())
            
            if cleaned_categories:
                query["bool"]["must"].append({
                    "bool": {
                        "should": [
                            {"terms": {"categories": cleaned_categories}},
                            {"match": {"categories": " ".join(cleaned_categories)}}
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
