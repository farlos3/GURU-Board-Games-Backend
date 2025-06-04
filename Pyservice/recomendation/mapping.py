boardgame_mapping = {
    "settings": {
        "analysis": {
            "normalizer": {
                "lowercase_normalizer": {
                    "type": "custom",
                    "filter": ["lowercase", "asciifolding"]
                }
            },
            "analyzer": {
                "search_analyzer": {
                    "tokenizer": "standard",
                    "filter": ["lowercase", "asciifolding", "stop"]
                },
                "ngram_analyzer": {
                    "tokenizer": "ngram_tokenizer",
                    "filter": ["lowercase", "asciifolding"]
                }
            },
            "tokenizer": {
                "ngram_tokenizer": {
                    "type": "ngram",
                    "min_gram": 1,
                    "max_gram": 3,
                    "token_chars": ["letter", "digit"]
                }
            }
        }
    },
    "mappings": {
        "properties": {
            "title": {
                "type": "text",
                "analyzer": "search_analyzer",
                "fields": {
                    "keyword": {
                        "type": "keyword",
                        "normalizer": "lowercase_normalizer"
                    },
                    "ngram": {
                        "type": "text",
                        "analyzer": "ngram_analyzer"
                    },
                    "suggest": {
                        "type": "completion"
                    }
                }
            },
            "description": {
                "type": "text",
                "analyzer": "search_analyzer"
            },
            "min_players": {
                "type": "integer"
            },
            "max_players": {
                "type": "integer"
            },
            "play_time_min": {
                "type": "integer"
            },
            "play_time_max": {
                "type": "integer"
            },
            "categories": {
                "type": "text",
                "analyzer": "search_analyzer",
                "fields": {
                    "keyword": {
                        "type": "keyword",
                        "normalizer": "lowercase_normalizer"
                    }
                }
            },
            "rating_avg": {
                "type": "float"
            },
            "rating_count": {
                "type": "integer"
            },
            "popularity_score": {
                "type": "float"
            },
            "image_url": {
                "type": "keyword",
                "index": False
            },
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

# ปรับปรุง user_action_mappings
user_action_mappings = {
    "settings": {
        "analysis": {
            "normalizer": {
                "lowercase_normalizer": {
                    "type": "custom",
                    "filter": ["lowercase", "asciifolding"]
                }
            }
        }
    },
    "mappings": {
        "properties": {
            "user_id": {
                "type": "keyword"
            },
            "boardgame_id": {
                "type": "keyword"
            },
            "action_type": {
                "type": "keyword",
                "normalizer": "lowercase_normalizer"
            },
            "action_value": {
                "type": "float",
                "null_value": 0
            },
            "action_detail": {
                "type": "text",
                "analyzer": "standard",
                "index": True  # เปลี่ยนเป็น True เผื่อต้องการค้นหา review
            },
            "action_time": {
                "type": "date",
                "format": "strict_date_optional_time||epoch_millis"
            },
            # เพิ่ม fields สำหรับ analytics
            "session_id": {
                "type": "keyword"
            },
            "ip_address": {
                "type": "ip"
            },
            "user_agent": {
                "type": "text",
                "index": False
            }
        }
    }
}