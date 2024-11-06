#!/usr/bin python
# encoding: utf-8
# 使用python向gdownloader中发送请求的示例

import socket
import logging
import platform

from base64 import b64encode
from os import environ, path
from typing import Any, Union
from logging.config import dictConfig

from dotenv.main import DotEnv
from pydantic import BaseModel, Field
from scrapy.http.headers import Headers
from aredis import StrictRedis, StrictRedisCluster

ENV_PATH = path.join(path.dirname(path.dirname(path.abspath(__file__))), "env_demo")
DEFAULT_LOG_FORMAT = "%(asctime)s - pid:%(process)d - %(filename)s [line: %(lineno)d] [%(levelname)s] ----- %(message)s"
DEFAULT_LOG_DATEFMT = "%Y-%m-%d %H:%M:%S"


class Request(BaseModel):
    """
    请求模型
    """
    URL: str
    Method: str = "GET"
    Depth: int = 0
    Body: Union[str, bytes] = b''
    ID: int = 0
    Ctx: dict[str, Any] = dict()
    Headers: dict[str, Any] = dict()
    ProxyURL: str = ""
    

class Config(BaseModel):
    Spider_Debug: bool = False
    # redis
    # 项目使用的redis key的前缀
    Redis_Prefix: str = "example"
    Redis_Host: str = "127.0.0.1"
    Redis_Port: int = 6379
    Redis_Mode: str = ""
    Redis_DB: int = Field(default=0, ge=0, le=15)
    Redis_Password: str = ""
    Downloader_MaxTopicSize: int = 10000
    # 日志
    Log_Level: str = "DEBUG"
    Log_Format: str = DEFAULT_LOG_FORMAT

    @classmethod
    def load_envs(cls, env_path=ENV_PATH, encoding="utf-8"):
        """
        加载env然后返回一个config实例
        """
        if not path.isfile(env_path):
            raise ValueError("%s no exist or not a file" % env_path)
        env_manager = DotEnv(
            dotenv_path=env_path, verbose=True, encoding=encoding
        )
        sys_envs = dict(environ.copy())
        envs = env_manager.dict()
        sys_envs.update(envs)
        return cls(**sys_envs)

    @classmethod
    def get_instance(cls):
        """
        获取config实例，单例
        """
        if not hasattr(cls, "instance"):
            setattr(cls, "instance", cls.load_envs())
        return getattr(cls, "instance")

    def get(self, configuration, default: Any = None):
        if hasattr(self, configuration):
            return getattr(self, configuration)
        else:
            self.get_logger().warning(
                "Try to get an unexpected option: %s" % configuration
            )
            return default

    @classmethod
    def get_logging_config(cls) -> dict[str, Any]:
        config :Config = cls.get_instance()
        log_level = config.Log_Level.upper()
        return {
                "version": 1,
                "disable_existing_loggers": False,
                "formatters": {
                    "default": {
                        'format': config.Log_Format,
                        'datefmt': DEFAULT_LOG_DATEFMT
                    }
                },

                "handlers": {
                    "console": {
                        "class": "logging.StreamHandler",
                        "level": "DEBUG",
                        "formatter": "default",
                        "stream": "ext://sys.stdout"
                    },
                },

                "loggers": {
                    "simple": {
                        'handlers': ['console'],
                        'level': log_level,
                        'propagate': False
                    }
                },

                "root": {
                    'handlers': ['console'],
                    'level': config.Spider_Debug and "DEBUG" or "WARNING"
                }
            }

    @classmethod
    def get_logger(cls) -> logging.Logger:
        if not hasattr(cls, "logger"):
            dictConfig(cls.get_logging_config())
            setattr(cls, "logger", logging.getLogger("simple"))
        return getattr(cls, "logger")


class Sender(object):
    """
    请求发送器
    """

    def __init__(
            self,
            redis_cli: Union[StrictRedis, StrictRedisCluster],
            urls_queue: str,
            reqs_queue: str,
            default_topic: str = "default",
            default_clear_head: bool = True,
            default_gzip_compress: bool = True

        ) -> None:
        self.redis_cli = redis_cli
        self.urls_queue = urls_queue
        self.reqs_queue = reqs_queue
        self.default_topic = default_topic
        self.default_clear_head = default_clear_head
        self.default_gzip_compress = default_gzip_compress

    def _tostring(self, x):
        """
        用于将header中的值转为string类型
        """
        if isinstance(x, bytes):
            return x.decode("utf-8")
        if isinstance(x, str):
            return x
        if isinstance(x, int):
            return str(x)
        raise TypeError(f"Unsupported value type: {type(x)}")


    def normvalue(self, value: Any):
        """
        将header的值转为list类型
        """
        """Normalize values to bytes"""
        if value is None:
            value = []
        elif isinstance(value, (str, bytes)):
            value = [value]
        elif not hasattr(value, "__iter__"):
            value = [value]

        return [self._tostring(x) for x in value]


    def serialize_headers(self, headers: dict[str, Any]) -> dict[str, list[Any]]:
        """
        格式化headers
        将headers中所有header的值转为list
        """
        headers_list = dict()
        for k, v in headers.items():
            headers_list[k] = v and self.normvalue(v) or list()
        return  headers_list   

    async def add_request(self, req: Request):
        """
        添加一个求情到队列
        """
        if req.Ctx.get("topic") is None:
            req.Ctx["topic"] = self.default_topic
        if req.Ctx.get("clearHead") is None:
            req.Ctx["clearHead"] = self.default_clear_head
        if req.Ctx.get("gzipCompress") is None:
            req.Ctx["gzipCompress"]= self.default_gzip_compress
        if req.Headers is not None:
            req.Headers = self.serialize_headers(req.Headers)
        if req.Body:
            req.Body = b64encode(req.Body).decode(encoding="utf-8")
        await self.redis_cli.rpush(self.reqs_queue, req.model_dump_json())


class ARedis(object):
    instance = None

    # socket保活选项
    DEFAULT_SOCKET_KEEPALIVE_OPTIONS = {
        socket.TCP_KEEPINTVL: 10,
        socket.TCP_KEEPCNT: 3
    }
    if platform.system() != "Darwin":
        # mac不支持此参数
        DEFAULT_SOCKET_KEEPALIVE_OPTIONS[socket.TCP_KEEPIDLE] = 30

    DEFAULT_REDIS_PARAMS = {
    "encoding": "utf-8",
    'socket_timeout': 30,
    'socket_connect_timeout': 30,
    'socket_keepalive': True,
    'socket_keepalive_options': DEFAULT_SOCKET_KEEPALIVE_OPTIONS,
    'retry_on_timeout': True,
    'max_connections': 0,
    'decode_responses': True
}

    @classmethod
    def get_instance(cls) -> Union[StrictRedis, StrictRedisCluster]:
        if cls.instance is None:
            params = cls.DEFAULT_REDIS_PARAMS.copy()
            params["host"] = config.get("Redis_Host")
            params["port"] = config.get("Redis_Port")
            params["password"] = config.get("Redis_Password")
            cluster = config.get("Redis_Mode") == "cluster"
            if cluster:
                redis_cls = StrictRedisCluster
                params["skip_full_coverage_check"] = True
            else:
                redis_cls = StrictRedis
                params["db"] = config.get("Redis_DB")
            cls.instance = redis_cls(**params)
        return cls.instance


config: Config = Config.get_instance()
logger = config.get_logger()



async def main():
    sender = Sender(
        redis_cli=ARedis.get_instance(),
        urls_queue=config.Redis_Prefix + ":start_urls",
        reqs_queue=config.Redis_Prefix + ":queue",
        default_topic=config.Redis_Prefix + ":items"
    )
    body = '{"invoke_info":{"pos_1":[{}],"pos_2":[{}],"pos_3":[{}]}}'
    req = Request(
        URL="https://ug.baidu.com/mcp/pc/pcsearch",
        Method="POST",
        Body=body.encode("utf-8"),
        Headers={
            "Content-Type": "application/json",
            "Origin": "https://www.baidu.com",
		    "Referer": "https://www.baidu.com"
        }
    )
    await sender.add_request(req)

        

if __name__ == "__main__":
    import asyncio
    loop = asyncio.get_event_loop()
    loop.run_until_complete(main())
    loop.close()
