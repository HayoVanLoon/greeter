#!/usr/bin/env python3
# Copyright 2020 Hayo van Loon
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import os
import time
from concurrent import futures

import grpc
from grpc_reflection.v1alpha import reflection

from hayovanloon.bartender.v1 import bartender_pb2, bartender_pb2_grpc


##############################################################
# Start of custom code
##############################################################
class BartenderServicer(bartender_pb2_grpc.BartenderServicer):

    def __init__(self) -> None:
        super().__init__()
        # For short-lived demo purposes only, use real storage for real applications
        self._cache = {}

    def CreateBeer(self,
                   create_beer_request: bartender_pb2.CreateBeerRequest,
                   context: grpc.ServicerContext) -> bartender_pb2.Beer:
        beer = create_beer_request.beer
        if not beer or not beer.name or not create_beer_request.brand:
            context.set_details("Well ain't that cute? But it's WRONG!")
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            return bartender_pb2.Beer()

        types = self._cache.get(create_beer_request.brand, [])
        if beer.name in types:
            context.set_details("Plagiarism!")
            context.set_code(grpc.StatusCode.ALREADY_EXISTS)
            return bartender_pb2.Beer()

        types.append(beer)
        self._cache[create_beer_request.brand] = types

        return beer

    def GetBeer(self,
                get_beer_request: bartender_pb2.GetBeerRequest,
                context: grpc.ServicerContext) -> bartender_pb2.Beer:
        types = self._cache.get(get_beer_request.brand, [])
        if not types:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            return bartender_pb2.Beer()

        beer = [b for b in types if b.name == get_beer_request.name]
        if not beer:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            return bartender_pb2.Beer()
        return beer[0]

    def ListBeers(self,
                  list_beers_request: bartender_pb2.ListBeersRequest,
                  context: grpc.ServicerContext) -> bartender_pb2.Beer:
        resp = bartender_pb2.ListBeersResponse()
        for b, types in self._cache.items():
            brand = resp.brands.add()
            brand.name = b
            brand.types.extend(types)
        return resp

##############################################################
# End of custom code
##############################################################


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