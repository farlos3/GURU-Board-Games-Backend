import os
import sys
import logging

# เพิ่ม gRPC directory ใน Python path
current_dir = os.path.dirname(os.path.abspath(__file__))
grpc_dir = os.path.join(current_dir, 'gRPC')
sys.path.insert(0, grpc_dir)

def main():
    """Main entry point สำหรับ Python ML Service"""
    try:
        # Import และรัน server
        from gRPC.recommendation_server import serve
        serve()
    except ImportError as e:
        print(f"❌ Import Error: {e}")
        print("💡 Make sure you have generated protobuf files:")
        print("   python -m grpc_tools.protoc --proto_path=. --python_out=gRPC --grpc_python_out=gRPC recommendation.proto")
        sys.exit(1)
    except Exception as e:
        logging.error(f"❌ Failed to start server: {e}")
        sys.exit(1)

if __name__ == '__main__':
    main()