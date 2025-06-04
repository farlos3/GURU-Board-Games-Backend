from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import List, Optional
from datetime import datetime
import uvicorn
import logging
import os
from dotenv import load_dotenv
from recomendation.service import recommendation_service, UserAction, Boardgame, search_boardgames

# Load environment variables
load_dotenv()

logger = logging.getLogger(__name__)
logger.setLevel(logging.INFO)

# Get port from environment variable, default to 50051
PORT = int(os.getenv("PYTHON_SERVICE_PORT", "50051"))

app = FastAPI(title="Board Game Recommendation API")

class BoardgamesRequest(BaseModel):
    boardgames: List[Boardgame]

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

class RecommendationRequest(BaseModel):
    user_id: str
    limit: int = 10
    user_actions: Optional[List[UserAction]] = None
    user_categories: Optional[List[str]] = None

# Routes
@app.post("/api/actions")
async def send_user_action(action: UserAction):
    try:
        success = recommendation_service.add_user_action(action)
        if not success:
            raise HTTPException(status_code=500, detail="Failed to record user action")
        return {"success": True, "message": "Action recorded successfully"}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/recommendations")
async def get_recommendations(request: RecommendationRequest):
    try:
        logger.info("\n📥 ===== Received Recommendation Request =====")
        logger.info(f"👤 User ID: {request.user_id}")
        logger.info(f"🎯 Limit: {request.limit}")
        
        if request.user_actions:
            logger.info(f"📊 User Actions: {len(request.user_actions)}")
            for action in request.user_actions:
                logger.info(f"  - {action.action_type} for boardgame {action.boardgame_id}")
        
        if request.user_categories:
            logger.info(f"🏷️ User Categories: {request.user_categories}")
        
        # Get recommendations using the service
        recommendations = recommendation_service.get_recommendations(
            user_id=request.user_id,
            limit=request.limit,
            user_actions=request.user_actions,
            user_categories=request.user_categories
        )
        
        logger.info(f"✅ Generated {len(recommendations)} recommendations")
        logger.info("===========================================\n")
        
        # Return only the recommended boardgames
        return {"boardgames": recommendations}
    except Exception as e:
        logger.error(f"❌ Error in recommendations endpoint: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/api/boardgames")
async def send_all_boardgames(request: BoardgamesRequest):
    try:
        success = recommendation_service.update_boardgames(request.boardgames)
        if not success:
            raise HTTPException(status_code=500, detail="Failed to update boardgames")
        return {"success": True, "message": "Boardgames updated successfully"}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/api/boardgames")
async def get_all_boardgames():
    try:
        boardgames = recommendation_service.get_all_boardgames()
        return {"boardgames": boardgames}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/api/boardgames/popular")
async def get_popular_boardgames(limit: int = 5):
    try:
        boardgames = recommendation_service.get_popular_boardgames(limit)
        return {"boardgames": boardgames}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/api/actions/user/{user_id}")
async def get_user_actions(user_id: str):
    try:
        actions = recommendation_service.get_user_actions(user_id)
        return {"actions": actions}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/api/actions/boardgame/{boardgame_id}")
async def get_boardgame_actions(boardgame_id: str):
    try:
        actions = recommendation_service.get_boardgame_actions(boardgame_id)
        return {"actions": actions}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

# New endpoint for searching boardgames
@app.get("/api/search")
async def search_boardgames_endpoint(
    searchQuery: Optional[str] = None,
    playerCount: Optional[int] = None,
    playTime: Optional[int] = None,
    categories: Optional[str] = None,
    limit: int = 10,
    page: int = 1
):
    try:
        categories_list = categories.split(",") if categories else None

        results = search_boardgames(
            search_query=searchQuery,
            player_count=playerCount,
            play_time=playTime,
            categories=categories_list,
            limit=limit,
            page=page
        )
        return results
    except Exception as e:
        logger.error(f"❌ Error in /api/search endpoint: {e}")
        # Return an empty list on error to match expected frontend format
        return []

def main():
    """Main entry point for Python ML Service"""
    try:
        logging.info(f"🚀 Starting Python service on port {PORT}")
        uvicorn.run(app, host="0.0.0.0", port=PORT)
    except Exception as e:
        logging.error(f"❌ Failed to start server: {e}")
        raise

if __name__ == '__main__':
    main()