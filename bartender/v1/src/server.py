#!/usr/bin/env python3
import os
import time
from concurrent import futures

import grpc
from grpc_reflection.v1alpha import reflection

from hayovanloon.bartender.v1 import bartender_pb2, bartender_pb2_grpc


class BartenderServicer(bartender_pb2_grpc.BartenderServicer):

    def __init__(self) -> None:
        super().__init__()
        self._cache = {}

    def CreateBeer(self, create_beer_request: bartender_pb2.CreateBeerRequest, context):
        beer = create_beer_request.beer
        if not beer or not beer.name or not create_beer_request.brand:
            context.set_details("Well ain't that cute? But it's WRONG!")
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            return bartender_pb2.Beer()

        key = (beer.name, create_beer_request.brand)
        if key in self._cache:
            context.set_details("Plagiarism!")
            context.set_code(grpc.StatusCode.ALREADY_EXISTS)
            return bartender_pb2.Beer()

        self._cache[key] = beer

        return beer

    def GetBeer(self, get_beer_request: bartender_pb2.GetBeerRequest, context):
        key = (get_beer_request.name, get_beer_request.brand)
        beer = self._cache.get(key)
        if not beer:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            return bartender_pb2.Beer()
        return beer


def serve(port, shutdown_grace_duration):
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))

    bartender_pb2_grpc.add_BartenderServicer_to_server(BartenderServicer(), server)
    reflection.enable_server_reflection((
        bartender_pb2.DESCRIPTOR.services_by_name['Bartender'].full_name,
        reflection.SERVICE_NAME,
    ), server)
    server.add_insecure_port('[::]:{}'.format(port))
    server.start()

    print('Listening on port {}'.format(port))

    try:
        while True:
            time.sleep(1000)
    except KeyboardInterrupt:
        server.stop(shutdown_grace_duration)


if __name__ == '__main__':
    port = os.environ.get('PORT')
    if not port:
        port = 8080

    serve(port, 5)