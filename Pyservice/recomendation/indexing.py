from connection.connection import client
from .setting import boardgame_index_name, user_action_index_name, boardgame_settings, user_action_settings

def create_indices():
    """
    Checks if Elasticsearch indices exist and creates them if they don't.
    """
    if not client.indices.exists(index=boardgame_index_name):
        client.indices.create(index=boardgame_index_name, body=boardgame_settings)
        client.indices.put_alias(index=boardgame_index_name, name="boardgame_alias")

    if not client.indices.exists(index=user_action_index_name):
        client.indices.create(index=user_action_index_name, body=user_action_settings)
        client.indices.put_alias(index=user_action_index_name, name="user_action_alias")

# ต้องเอาข้อมูล API จาก GO มาประมวลผล