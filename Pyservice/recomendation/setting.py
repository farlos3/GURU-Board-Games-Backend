from .mapping import boardgame_mapping, user_action_mappings

boardgame_index_name = "boardgame"
user_action_index_name = "user_action"

boardgame_settings = {
    "settings": {
        "number_of_shards": 1,
        "number_of_replicas": 0
    },
    "mappings": boardgame_mapping["mappings"]
}

user_action_settings = {
    "settings": {
        "number_of_shards": 1,
        "number_of_replicas": 0
    },
    "mappings": user_action_mappings["mappings"]
}