from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import List
from datetime import datetime
import uvicorn
import logging
import os
from dotenv import load_dotenv
from recomendation.service import recommendation_service, UserAction, Boardgame

# Load environment variables
load_dotenv()

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

@app.post("/api/recommendations")
async def get_recommendations(user_id: str, limit: int = 10):
    try:
        recommendations = recommendation_service.get_recommendations(user_id, limit)
        return {"boardgames": recommendations}
    except Exception as e:
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