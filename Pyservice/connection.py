import os
import logging
from elasticsearch import Elasticsearch
from dotenv import load_dotenv

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

logging.getLogger('elasticsearch').setLevel(logging.WARNING)
logging.getLogger('urllib3').setLevel(logging.WARNING)

load_dotenv()

def create_elasticsearch_client():
    """สร้างและคืนค่า Elasticsearch client"""
    endpoint = os.getenv('ELASTICSEARCH_ENDPOINT')
    api_key = os.getenv('ELASTICSEARCH_API_KEY')

    if not endpoint:
        logger.error("❌ ELASTICSEARCH_ENDPOINT not found in environment variables")
        print("=" * 90)
        return None

    if not api_key:
        logger.error("❌ ELASTICSEARCH_API_KEY not found in environment variables")
        print("=" * 90)
        return None

    logger.info(f"🔧 Connecting to Elasticsearch endpoint: {endpoint}")
    logger.info("🔑 API key loaded successfully")

    try:
        client = Elasticsearch(
            endpoint,
            api_key=api_key,
        )
        
        logger.info("🔄 Testing connection...")
        
        if client.ping():
            logger.info("✅ Successfully connected to Elasticsearch!")
            
            try:
                cluster_info = client.info()
                logger.info(f"📊 Cluster name: {cluster_info['cluster_name']}")
                logger.info(f"🏷️  Version: {cluster_info['version']['number']}")
            except Exception as e:
                logger.warning(f"⚠️  Could not retrieve cluster info: {e}")
            
            print("=" * 90)
            return client
        else:
            logger.error("❌ Failed to connect to Elasticsearch")
            print("=" * 90)
            return None
            
    except Exception as e:
        logger.error(f"💥 Connection error: {e}")
        logger.error("🔍 Please check your endpoint URL and API key")
        print("=" * 90)
        return None

client = create_elasticsearch_client()