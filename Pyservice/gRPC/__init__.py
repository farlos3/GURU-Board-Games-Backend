"""
gRPC package for recommendation service
"""

# Import generated protobuf modules
try:
    from . import recommendation_pb2
    from . import recommendation_pb2_grpc
    from . import recommendation_server
    
    __all__ = [
        'recommendation_pb2',
        'recommendation_pb2_grpc', 
        'recommendation_server'
    ]
except ImportError:
    # If protobuf files haven't been generated yet
    print("⚠️  Protobuf files not found. Please run:")
    print("   python -m grpc_tools.protoc --proto_path=. --python_out=gRPC --grpc_python_out=gRPC recommendation.proto")
    __all__ = []