#!/usr/bin/python
# -*- coding: utf-8 -*-

import uuid
import random
import hashlib

random.seed(uuid.getnode())
flags = [hashlib.md5(random.randbytes(32)).hexdigest() for _ in range(8 * 6)]
template = open('docker-compose.template.yaml', 'r').read()
text = template.format(*flags)
open('docker-compose.yaml', 'w').write(text)
