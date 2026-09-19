"""Minimal mitmproxy event adapter for Broker v1.
Only emits metadata; body redaction belongs to Broker worker in this first version.
"""
import json
from mitmproxy import http, ctx

def request(flow: http.HTTPFlow) -> None:
    ctx.log.info(json.dumps({"type":"request","flow_id":flow.id,"method":flow.request.method,"url":flow.request.pretty_url}))

def response(flow: http.HTTPFlow) -> None:
    ctx.log.info(json.dumps({"type":"response","flow_id":flow.id,"status":flow.response.status_code,"url":flow.request.pretty_url}))
