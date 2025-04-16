import etcd3
import json
import time
from typing import Dict
from Log.log import logger

class EtcdService:
    def __init__(self, host: str = "localhost", port: int = 2379):
        """初始化 etcd 客户端"""
        self.client = etcd3.client(host=host, port=port)
        self.lease_ttl = 10  # 租约时间（秒）

    def register_service(self, service_name: str, service_info: Dict, lease_ttl: int = None):
        """注册服务到 etcd"""
        try:
            lease_ttl = lease_ttl or self.lease_ttl
            lease = self.client.lease(lease_ttl)
            service_key = f"/services/{service_name}/{service_info['host']}:{service_info['port']}"
            service_value = json.dumps(service_info)
            self.client.put(service_key, service_value, lease)
            logger.info(f"Registered service: {service_name} at {service_key}")

            # 保持租约续期
            def keep_alive():
                while True:
                    try:
                        lease.refresh()
                        time.sleep(lease_ttl // 2)
                    except Exception as e:
                        logger.error(f"Failed to keep lease alive: {e}")
                        break

            import threading
            keep_alive_thread = threading.Thread(target=keep_alive, daemon=True)
            keep_alive_thread.start()
            return lease
        except Exception as e:
            logger.error(f"Failed to register service {service_name}: {e}")
            raise

    def deregister_service(self, service_name: str, host: str, port: int):
        """注销服务"""
        try:
            service_key = f"/services/{service_name}/{host}:{port}"
            self.client.delete(service_key)
            logger.info(f"Deregistered service: {service_name} at {service_key}")
        except Exception as e:
            logger.error(f"Failed to deregister service {service_name}: {e}")

    def discover_services(self, service_name: str) -> list:
        """发现指定服务的所有实例"""
        try:
            service_prefix = f"/services/{service_name}/"
            services = []
            for value, metadata in self.client.get_prefix(service_prefix):
                service_info = json.loads(value.decode('utf-8'))
                services.append(service_info)
            logger.info(f"Discovered {len(services)} instances for {service_name}")
            return services
        except Exception as e:
            logger.error(f"Failed to discover services {service_name}: {e}")
            return []