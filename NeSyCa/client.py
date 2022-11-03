import random
from pathlib import Path
from time import sleep

import requests

from Core.config import Settings
from Core.logger import getLogger

logger = getLogger(__name__)


class NeSyCaClient:
    def __init__(
        self,
        bank_path: Path = None,
        band_min: float = None,
        band_max: float = None,
        length_min: float = None,
        length_max: float = None,
        batch_min: int = None,
        batch_max: int = None,
    ):
        self.bank_path = bank_path or Settings.URL_BANK_PATH
        self.band_min = band_min or Settings.BAND_MIN
        self.band_max = band_max or Settings.BAND_MAX
        self.length_min = length_min or Settings.LENGTH_MIN
        self.length_max = length_max or Settings.LENGTH_MAX
        self.batch_min = batch_min or Settings.BATCH_SIZE_MIN
        self.batch_max = batch_max or Settings.BATCH_SIZE_MAX

        with self.bank_path.open() as f_:
            self.bank = [x.rstrip("\n") for x in f_.read().split("\n")]

    def rand__target(self):
        v_ = random.choice(self.bank)
        logger.debug(f"rand__target: {v_}")
        return v_

    def rand__band(self):
        v_ = random.random() * (self.band_max - self.band_min) + self.band_min
        logger.debug(f"rand__band: {v_}")
        return v_

    def rand__length(self, available):
        v_ = (
            random.random() * (self.length_max - self.length_min)
            + self.length_min
        )
        if v_ > available:
            v_ = available
        logger.debug(f"rand__length: {v_}")
        return int(v_)

    def rand__batch(self):
        v_ = (
            random.random() * (self.batch_max - self.batch_min)
            + self.batch_min
        )
        logger.debug(f"rand__batch: {v_}")
        return int(v_)

    def send(self):
        t_ = self.rand__target()
        b_ = self.rand__band()
        bs_ = self.rand__batch()
        with requests.Session() as sess:
            with sess.get(t_, stream=True) as r_:
                cl_ = r_.headers["Content-length"]
                cl_select = self.rand__length(int(cl_))
                logger.info(
                    f"reading {cl_select} out of {cl_} ({bs_} byte batches - {b_} delay)"
                )
                for c_, _ in enumerate(r_.iter_content(bs_)):
                    sleep(b_)
                    if (c_ + 1) * bs_ > cl_select:
                        break
                    logger.debug(f"{(c_ + 1) * bs_}<{cl_select}. continue ...")

        logger.info("done")
