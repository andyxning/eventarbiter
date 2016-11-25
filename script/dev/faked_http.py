#!/usr/bin/env python
# -*- coding: utf-8 -*-

from flask import Flask, request
import json
import pprint

app = Flask(__name__)

@app.route('/', methods=['POST'])
def get():
    pprint.pprint(json.loads(request.data))
    return "OK"

if __name__ == '__main__':
    app.run(debug=True, port=3086)
