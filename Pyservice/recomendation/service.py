from typing import List, Optional
from datetime import datetime
import logging
from pydantic import BaseModel
from .indexing import create_indices
from .setting import boardgame_index_name, user_action_index_name
from connection.connection import client
from typing import Optional, List, Union

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

def preprocess_categories(categories: str) -> str:
    """Preprocess categories string to ensure consistent format.
    
    Args:
        categories: Comma-separated string of categories
        
    Returns:
        Preprocessed categories string
    """
    if not categories:
        return ""
        
    # Split by comma and clean each category
    category_list = [cat.strip().lower() for cat in categories.split(',')]
    
    # Remove empty categories and duplicates
    category_list = list(set(filter(None, category_list)))
    
    # Sort categories for consistency
    category_list.sort()
    
    # Join back with comma
    return ','.join(category_list)

def preprocess_boardgame(boardgame: Boardgame) -> Boardgame:
    """Preprocess boardgame data before indexing.
    
    Args:
        boardgame: Boardgame object to preprocess
        
    Returns:
        Preprocessed Boardgame object
    """
    # Create a copy of the boardgame to avoid modifying the original
    processed = boardgame.copy()
    
    # Preprocess categories
    if processed.categories:
        processed.categories = preprocess_categories(processed.categories)
    
    return processed

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
                # Preprocess boardgame data before indexing
                processed_bg = preprocess_boardgame(bg)
                
                response = client.index(
                    index=boardgame_index_name,
                    id=str(processed_bg.id),
                    body=processed_bg.dict()
                )
                logger.info(f"✅ Boardgame {processed_bg.id} updated in Elasticsearch: {response['_id']}")
                logger.info(f"  Categories: {processed_bg.categories}")
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

# แก้ไขฟังก์ชัน search_boardgames ในส่วนการจัดการ categories filter

def search_boardgames(
    search_query: Optional[str] = None,
    player_count: Optional[int] = None,
    play_time: Optional[int] = None,
    categories: Optional[List[str]] = None,
    limit: int = 10,
    page: int = 1,
    search_logic: str = "OR",  # OR logic for text search
    category_logic: str = "AND"  # AND logic for categories
) -> List[Boardgame]:
    """
    Search boardgames with flexible AND/OR logic for different search components.
    
    Args:
        search_logic: "OR" or "AND" - logic for text search terms
        category_logic: "OR" or "AND" - logic between categories
    """
    try:
        query = {
            "bool": {
                "must": [],     # Required conditions (AND logic)
                "should": [],   # Optional conditions (OR logic)  
                "filter": [],   # Exact filters (AND logic)
                "minimum_should_match": 0
            }
        }

        # === TEXT SEARCH SECTION ===
        if search_query:
            search_query = search_query.strip()
            logger.info(f"🔍 Text search with {search_logic} logic: '{search_query}'")
            
            # Create different search strategies
            text_search_queries = [
                # 1. Exact phrase match (highest priority)
                {
                    "match_phrase": {
                        "title": {
                            "query": search_query,
                            "boost": 5
                        }
                    }
                },
                # 2. Fuzzy match
                {
                    "fuzzy": {
                        "title": {
                            "value": search_query,
                            "fuzziness": "AUTO",
                            "prefix_length": 0,
                            "max_expansions": 50,
                            "boost": 4
                        }
                    }
                },
                # 3. Wildcard match
                {
                    "wildcard": {
                        "title": {
                            "value": f"*{search_query.lower()}*",
                            "case_insensitive": True,
                            "boost": 3
                        }
                    }
                },
                # 4. Prefix match
                {
                    "prefix": {
                        "title": {
                            "value": search_query.lower(),
                            "case_insensitive": True,
                            "boost": 2
                        }
                    }
                }
            ]
            
            if search_logic.upper() == "OR":
                # OR logic: any of the search strategies can match
                query["bool"]["should"].extend(text_search_queries)
                query["bool"]["minimum_should_match"] = max(1, query["bool"]["minimum_should_match"])
            else:  # AND logic
                # AND logic: all search strategies must match (usually too strict)
                # More practical: require the query to match in at least one way
                text_search_bool = {
                    "bool": {
                        "should": text_search_queries,
                        "minimum_should_match": 1
                    }
                }
                query["bool"]["must"].append(text_search_bool)

        # === CATEGORY SEARCH SECTION ===
        if categories:
            # Convert and clean categories
            if isinstance(categories, str):
                try:
                    import json
                    categories = json.loads(categories)
                except json.JSONDecodeError:
                    categories = [cat.strip() for cat in categories.split(',')]
            
            cleaned_categories = []
            for cat in categories:
                if isinstance(cat, str):
                    cat = cat.strip('[]"\'').lower()
                    if cat:
                        cleaned_categories.append(cat)
            
            if cleaned_categories:
                logger.info(f"🏷️ Category search with {category_logic} logic: {cleaned_categories}")
                
                # Create fuzzy search for each category
                category_queries = []
                for cat in cleaned_categories:
                    cat_search_strategies = [
                        # Exact match
                        {
                            "term": {
                                "categories.keyword": {
                                    "value": cat,
                                    "case_insensitive": True,
                                    "boost": 5
                                }
                            }
                        },
                        # Phrase match
                        {
                            "match_phrase": {
                                "categories": {
                                    "query": cat,
                                    "boost": 4
                                }
                            }
                        },
                        # Fuzzy match
                        {
                            "match": {
                                "categories": {
                                    "query": cat,
                                    "fuzziness": "AUTO",
                                    "boost": 3
                                }
                            }
                        },
                        # Wildcard match
                        {
                            "wildcard": {
                                "categories": {
                                    "value": f"*{cat}*",
                                    "case_insensitive": True,
                                    "boost": 2
                                }
                            }
                        }
                    ]
                    
                    # Each category must match at least one strategy
                    category_query = {
                        "bool": {
                            "should": cat_search_strategies,
                            "minimum_should_match": 1
                        }
                    }
                    category_queries.append(category_query)
                
                # Apply category logic
                if category_logic.upper() == "OR":
                    # OR logic: match any category
                    query["bool"]["should"].extend(category_queries)
                    query["bool"]["minimum_should_match"] = max(1, query["bool"]["minimum_should_match"])
                else:  # AND logic
                    # AND logic: must match all categories
                    query["bool"]["must"].extend(category_queries)
                
                logger.info(f"🔍 Added {len(category_queries)} category queries with {category_logic} logic")

        # === EXACT FILTERS SECTION (Always AND logic) ===
        
        # Player count filter
        if player_count is not None and player_count > 0:
            logger.info(f"👥 Player count filter: {player_count}")
            query["bool"]["filter"].extend([
                {"range": {"min_players": {"lte": player_count}}},
                {"range": {"max_players": {"gte": player_count}}}
            ])

        # Play time filter  
        if play_time is not None and play_time > 0:
            logger.info(f"⏱️ Play time filter: {play_time} minutes")
            query["bool"]["filter"].extend([
                {"range": {"play_time_min": {"lte": play_time}}},
                {"range": {"play_time_max": {"gte": play_time}}}
            ])

        # === HANDLE EMPTY QUERY ===
        if (not query["bool"]["must"] and 
            not query["bool"]["should"] and 
            not query["bool"]["filter"]):
            logger.info("📋 No search criteria provided, showing popular games")
            query = {"match_all": {}}

        # === EXECUTE SEARCH ===
        from_ = max(0, (page - 1) * limit)
        
        # Log the final query structure for debugging
        logger.info(f"🔧 Final query structure: must={len(query.get('bool', {}).get('must', []))}, should={len(query.get('bool', {}).get('should', []))}, filter={len(query.get('bool', {}).get('filter', []))}")

        response = client.search(
            index=boardgame_index_name,
            body={
                "query": query,
                "size": limit,
                "from": from_,
                "sort": [
                    {"_score": {"order": "desc"}},
                    {"popularity_score": {"order": "desc", "missing": "_last"}},
                    {"rating_avg": {"order": "desc", "missing": "_last"}},
                    {"id": {"order": "asc"}}
                ],
                "highlight": {
                    "fields": {
                        "title": {
                            "pre_tags": ["<mark>"],
                            "post_tags": ["</mark>"],
                            "fragment_size": 150,
                            "number_of_fragments": 1
                        },
                        "categories": {
                            "pre_tags": ["<mark>"],
                            "post_tags": ["</mark>"],
                            "fragment_size": 150,
                            "number_of_fragments": 3
                        }
                    }
                }
            }
        )

        # === PROCESS RESULTS ===
        boardgames = []
        for hit in response['hits']['hits']:
            boardgame_data = hit['_source']
            boardgame_data['_search_score'] = hit['_score']
            if 'highlight' in hit:
                boardgame_data['_highlights'] = hit['highlight']
                
            try:
                boardgame = Boardgame(**boardgame_data)
                boardgames.append(boardgame)
            except Exception as e:
                logger.warning(f"⚠️ Error creating Boardgame object: {e}")
                boardgames.append(boardgame_data)

        # === LOGGING ===
        search_type = f"text({search_logic})" if search_query else ""
        category_type = f"categories({category_logic})" if categories else ""
        filter_type = "filters" if (player_count or play_time) else ""
        
        search_description = " + ".join(filter(None, [search_type, category_type, filter_type]))
        if not search_description:
            search_description = "popular"
            
        logger.info(f"✅ {search_description.upper()} search completed: results={len(boardgames)}")
        
        # Log top results
        if boardgames:
            logger.info("🎯 Top search results:")
            for i, bg in enumerate(boardgames[:3]):
                if isinstance(bg, dict):
                    title = bg.get('title', 'Unknown')
                    score = bg.get('_search_score', 0)
                    cats = bg.get('categories', '')
                else:
                    title = bg.title
                    score = getattr(bg, '_search_score', 0)
                    cats = bg.categories
                
                logger.info(f"  {i+1}. {title} (Score: {score:.2f})")
        
        # Suggestions for no results
        if len(boardgames) == 0 and (search_query or categories):
            logger.info("💡 No results found. Try:")
            logger.info("  - Using broader search terms")
            logger.info("  - Switching logic (OR vs AND)")
            logger.info("  - Checking spelling")

        return boardgames

    except Exception as e:
        logger.error(f"❌ Search error: {str(e)}")
        logger.error(f"Query that caused error: {query if 'query' in locals() else 'Query not constructed'}")
        return []


# === CONVENIENCE FUNCTIONS ===

def search_boardgames_flexible_or(
    search_query: Optional[str] = None,
    categories: Optional[List[str]] = None,
    **kwargs
) -> List[Boardgame]:
    """Search with OR logic for both text and categories (more permissive)."""
    return search_boardgames(
        search_query=search_query,
        categories=categories,
        search_logic="OR",
        category_logic="OR",
        **kwargs
    )

def search_boardgames_strict_and(
    search_query: Optional[str] = None,
    categories: Optional[List[str]] = None,
    **kwargs
) -> List[Boardgame]:
    """Search with AND logic for both text and categories (more restrictive)."""
    return search_boardgames(
        search_query=search_query,
        categories=categories,
        search_logic="AND",
        category_logic="AND",
        **kwargs
    )

def search_boardgames_mixed(
    search_query: Optional[str] = None,
    categories: Optional[List[str]] = None,
    **kwargs
) -> List[Boardgame]:
    """Search with OR logic for text but AND logic for categories (balanced)."""
    return search_boardgames(
        search_query=search_query,
        categories=categories,
        search_logic="OR",
        category_logic="AND",
        **kwargs
    )
    
def find_similar_categories(search_categories: List[str], limit: int = 5) -> List[str]:
    """Find categories similar to the search terms using fuzzy matching"""
    try:
        similar_cats = []
        
        for search_cat in search_categories:
            response = client.search(
                index=boardgame_index_name,
                body={
                    "size": 0,
                    "aggs": {
                        "similar_categories": {
                            "terms": {
                                "field": "categories.keyword",
                                "include": f".*{search_cat.lower()}.*",
                                "size": limit
                            }
                        }
                    }
                }
            )
            
            if 'aggregations' in response and 'similar_categories' in response['aggregations']:
                for bucket in response['aggregations']['similar_categories']['buckets']:
                    if bucket['key'] not in similar_cats:
                        similar_cats.append(bucket['key'])
        
        return similar_cats[:limit]
        
    except Exception as e:
        logger.error(f"❌ Error finding similar categories: {e}")
        return []
    
# Enhanced function to get category suggestions
def get_category_suggestions(partial_category: str, limit: int = 10) -> List[str]:
    """Get category suggestions based on partial input with fuzzy matching"""
    try:
        # Use multiple strategies to find matching categories
        strategies = [
            # Prefix match
            {"prefix": {"categories.keyword": {"value": partial_category.lower(), "case_insensitive": True}}},
            # Wildcard match
            {"wildcard": {"categories.keyword": {"value": f"*{partial_category.lower()}*", "case_insensitive": True}}},
            # Fuzzy match
            {"fuzzy": {"categories.keyword": {"value": partial_category.lower(), "fuzziness": "AUTO"}}}
        ]
        
        all_suggestions = set()
        
        for strategy in strategies:
            response = client.search(
                index=boardgame_index_name,
                body={
                    "query": strategy,
                    "size": 0,
                    "aggs": {
                        "categories": {
                            "terms": {
                                "field": "categories.keyword",
                                "size": limit * 2  # Get more to filter later
                            }
                        }
                    }
                }
            )
            
            if 'aggregations' in response and 'categories' in response['aggregations']:
                for bucket in response['aggregations']['categories']['buckets']:
                    all_suggestions.add(bucket['key'])
        
        # Sort suggestions by relevance (those that start with the partial term first)
        suggestions = list(all_suggestions)
        suggestions.sort(key=lambda x: (
            0 if x.lower().startswith(partial_category.lower()) else 1,  # Prefix matches first
            len(x),  # Shorter terms first within each group
            x.lower()  # Alphabetical
        ))
        
        return suggestions[:limit]
        
    except Exception as e:
        logger.error(f"❌ Error getting category suggestions: {e}")
        return []

# เพิ่มฟังก์ชันสำหรับดูข้อมูล categories ทั้งหมดในระบบ
def get_all_categories() -> List[str]:
    """Get all unique categories from the system for debugging"""
    try:
        # ใช้ aggregation เพื่อดู categories ทั้งหมด
        response = client.search(
            index=boardgame_index_name,
            body={
                "size": 0,
                "aggs": {
                    "categories": {
                        "terms": {
                            "field": "categories.keyword",
                            "size": 100
                        }
                    }
                }
            }
        )
        
        categories = []
        if 'aggregations' in response and 'categories' in response['aggregations']:
            for bucket in response['aggregations']['categories']['buckets']:
                categories.append(bucket['key'])
        
        logger.info(f"📋 All categories in system: {categories}")
        return categories
        
    except Exception as e:
        logger.error(f"❌ Error getting categories: {e}")
        return []

# เพิ่มฟังก์ชันสำหรับตรวจสอบ mapping ของ categories field
def check_categories_mapping():
    """Check the mapping of categories field"""
    try:
        mapping_response = client.indices.get_mapping(
            index=boardgame_index_name
        )
        
        logger.info("🗺️ Categories field mapping:")
        categories_mapping = mapping_response[boardgame_index_name]['mappings']['properties'].get('categories', {})
        logger.info(f"  {categories_mapping}")
        
        return categories_mapping
        
    except Exception as e:
        logger.error(f"❌ Error checking mapping: {e}")
        return {}