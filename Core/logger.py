import logging
from Core.config import Settings

logging.basicConfig(level=getattr(logging, Settings.LOG_LEVEL))
getLogger = logging.getLogger
