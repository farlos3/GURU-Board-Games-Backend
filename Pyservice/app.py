import logging
from connection import client

logger = logging.getLogger(__name__)

def main():
    """Main application function"""
    logger.info("🚀 Starting Pyservice application...")
    
    if client is None:
        logger.error("❌ Elasticsearch client is not available")
        print("=" * 90)
        return False
    
    try:
        logger.info("📡 Using Elasticsearch client...")
        
        target_indices = ["boardgame", "user_action"]
        
        try:
            for index_name in target_indices:
                if client.indices.exists(index=index_name):
                    search_result = client.search(
                        index=index_name,
                        body={"size": 1}
                    )
                    logger.info(f"✅ Successfully accessed index '{index_name}'")
                else:
                    logger.warning(f"⚠️  Index '{index_name}' does not exist")
        
        except Exception as e:
            logger.error(f"💥 Error accessing indices: {e}")
            print("=" * 90)
            return False

        print("=" * 90)
        return True
        
    except Exception as e:
        logger.error(f"💥 Application error: {e}")
        print("=" * 90)
        return False

if __name__ == "__main__":
    success = main()
    if not success:
        exit(1)