import grpc
from hayovanloon.bartender.v1 import bartender_pb2

from server import BartenderServicer


class FakeContext(grpc.ServicerContext):

    def __init__(self) -> None:
        super().__init__()
        self.code = grpc.StatusCode.OK

    def set_code(self, code: grpc.StatusCode):
        self.code = code


class TestBartenderServicer(object):
    # TODO(hvl): use proper testing harness
    # TODO(hvl): extend with non-happy tests

    def test_CreateBeer(self):
        serv = BartenderServicer()

        input = bartender_pb2.CreateBeerRequest(brand="HertogJan")
        input.beer.name = "Pils"
        expected = bartender_pb2.Beer(name="Pils")

        _ = serv.CreateBeer(input, context=FakeContext())
        actual_types = serv._cache.get('HertogJan')

        as_dict = {t.name: t for t in actual_types}
        assert as_dict == {expected.name: expected}

    def test_GetBeer(self):
        serv = BartenderServicer()

        setup = bartender_pb2.CreateBeerRequest(brand="HertogJan")
        setup.beer.name = "Pils"
        _ = serv.CreateBeer(setup, context=FakeContext())

        input = bartender_pb2.GetBeerRequest(brand="HertogJan", name="Pils")
        expected = bartender_pb2.Beer(name="Pils")

        actual = serv.GetBeer(input, context=FakeContext())

        assert actual == expected

    def test_ListBeers(self):
        serv = BartenderServicer()

        setup1 = bartender_pb2.CreateBeerRequest(brand="Heineken")
        setup1.beer.name = "Pils"
        _ = serv.CreateBeer(setup1, context=FakeContext())
        setup2 = bartender_pb2.CreateBeerRequest(brand="HertogJan")
        setup2.beer.name = "Pils"
        _ = serv.CreateBeer(setup2, context=FakeContext())

        input = bartender_pb2.ListBeersRequest()
        expected = bartender_pb2.ListBeersResponse()
        b1 = expected.brands.add()
        b1.name = "Heineken"
        b1.types.extend([bartender_pb2.Beer(name="Pils")])
        b2 = expected.brands.add()
        b2.name = "HertogJan"
        b2.types.extend([bartender_pb2.Beer(name="Pils")])

        actual = serv.ListBeers(input, context=FakeContext())

        assert actual == expected
