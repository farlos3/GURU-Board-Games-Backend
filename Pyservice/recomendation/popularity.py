from collections import defaultdict

def calculate_popularity_scores(user_actions_data):
    """
    Calculates popularity scores for boardgames based on user actions.

    Args:
        user_actions_data: An iterable containing user action dictionaries.
                           Each dictionary should have at least 'boardgame_id', 'action_type',
                           and optionally 'action_value' (for 'rate' and 'page_view').

    Returns:
        A dictionary where keys are boardgame_ids and values are their popularity scores.
    """
    # เก็บข้อมูลแบบ boardgame_id -> data
    stats = defaultdict(lambda: {
        "rating_sum": 0.0,
        "rating_count": 0,
        "like_count": 0,
        "favorite_count": 0,
        "total_page_view_seconds": 0
    })

    for action in user_actions_data:
        bg_id = action["boardgame_id"]
        action_type = action["action_type"]
        value = action.get("action_value", 0)

        if action_type == "rate":
            stats[bg_id]["rating_sum"] += value
            stats[bg_id]["rating_count"] += 1
        elif action_type == "like":
            stats[bg_id]["like_count"] += 1
        elif action_type == "favorite":
            stats[bg_id]["favorite_count"] += 1
        elif action_type == "page_view":
            stats[bg_id]["total_page_view_seconds"] += value

    normalized_data = {}
    max_rating_avg = 0
    max_like = 0
    max_fav = 0
    max_page_view = 0

    for bg_id, data in stats.items():
        rating_count = data["rating_count"]
        rating_avg = data["rating_sum"] / rating_count if rating_count > 0 else 0

        # เก็บค่าสูงสุด
        max_rating_avg = max(max_rating_avg, rating_avg)
        max_like = max(max_like, data["like_count"])
        max_fav = max(max_fav, data["favorite_count"])
        max_page_view = max(max_page_view, data["total_page_view_seconds"])

        normalized_data[bg_id] = {
            "rating_avg": rating_avg,
            "like_count": data["like_count"],
            "favorite_count": data["favorite_count"],
            "page_view": data["total_page_view_seconds"]
        }

    popularity_scores = {}
    for bg_id, data in normalized_data.items():
        norm_rating = data["rating_avg"] / max_rating_avg if max_rating_avg > 0 else 0
        norm_like = data["like_count"] / max_like if max_like > 0 else 0
        norm_fav = data["favorite_count"] / max_fav if max_fav > 0 else 0
        norm_page = data["page_view"] / max_page_view if max_page_view > 0 else 0

        # กำหนดน้ำหนัก (weight)
        popularity_score = (
            norm_rating * 0.4 +
            norm_like * 0.2 +
            norm_fav * 0.2 +
            norm_page * 0.2
        )

        popularity_scores[bg_id] = popularity_score

        print(f"🧩 Boardgame ID: {bg_id}")
        print(f"   📊 Normalized Rating Avg: {norm_rating:.2f}")
        print(f"   ❤️ Normalized Like: {norm_like:.2f}")
        print(f"   💾 Normalized Favorite: {norm_fav:.2f}")
        print(f"   ⏱️ Normalized Page View: {norm_page:.2f}")
        print(f"   🔥 Final Popularity Score: {popularity_score:.3f}")

    return popularity_scores