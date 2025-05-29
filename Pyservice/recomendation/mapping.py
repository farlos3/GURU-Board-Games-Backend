boardgame_mapping = {
  "mappings": {
    "properties": {
      "title": {
        "type": "text",
        "analyzer": "standard"
      },
      "description": { "type": "text", "analyzer": "standard" },
      "min_players": { "type": "integer" },
      "max_players": { "type": "integer" },
      "play_time_min": { "type": "integer" },
      "play_time_max": { "type": "integer" },
      "categories": { "type": "keyword" },
      "rating_avg": { "type": "float" },
      "rating_count": { "type": "integer" },
      "popularity_score": { "type": "float" },
      "image_url": { "type": "keyword", "index": False },
      "created_at": {
        "type": "date",
        "format": "strict_date_optional_time||epoch_millis"
      },
      "updated_at": {
        "type": "date",
        "format": "strict_date_optional_time||epoch_millis"
      }
    }
  }
}

user_action_mappings = {
  "mappings": {
    "properties": {
      "user_id": { "type": "keyword" },            # รหัสผู้ใช้
      "boardgame_id": { "type": "keyword" },       # รหัสบอร์ดเกม
      "action_type": { "type": "keyword" },        # ประเภทการกระทำ เช่น rate, like, play

      "action_value": {
        "type": "float", 
        "null_value": 0                            # ค่าที่ใช้เวลาว่าง เช่น rate: 4.5, play: 1
      },

      "action_detail": { 
        "type": "text", 
        "index": False 
      },  # ใช้เก็บคอมเมนต์ หรือข้อมูลเสริม (เช่น review message)

      "action_time": {
        "type": "date",
        "format": "strict_date_optional_time||epoch_millis"
      },
      
    }
  }
}
