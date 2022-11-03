from functools import lru_cache
from pathlib import Path

from pydantic import BaseSettings, Field

_BASE_PATH = Path(__file__).parent.parent


class ConfigCls(BaseSettings):
    BASE_PATH: Path = Field(default=_BASE_PATH, const=True)
    URL_BANK_PATH: Path = _BASE_PATH / "storage/bank"
    BAND_MIN: float = 0.001
    BAND_MAX: float = 0.01
    LENGTH_MIN: float = 100 * 1024
    LENGTH_MAX: float = 10 * 1024**2
    BATCH_SIZE_MIN: int = 1024
    BATCH_SIZE_MAX: int = 100 * 1024
    LOG_LEVEL: str = "WARNING"

    class Config:
        env_prefix = "NESYCA__"
        env_file = str(_BASE_PATH / ".env")


@lru_cache
def _settings():
    return ConfigCls()


Settings = _settings()
