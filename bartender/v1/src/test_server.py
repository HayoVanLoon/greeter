from hayovanloon.bartender.v1 import bartender_pb2

from server import BartenderServicer


class TestBartenderServicer:
    # TODO(hvl): use proper testing harness
    # TODO(hvl): extend with non-happy tests

    def test_CreateBeer(self):
        serv = BartenderServicer()

        input = bartender_pb2.CreateBeerRequest(brand="HertogJan")
        input.beer.name = "Pils"
        expected = bartender_pb2.Beer(name="Pils")

        _ = serv.CreateBeer(input, context=None)
        actual_types = serv._cache.get('HertogJan')

        as_dict = {t.name: t for t in actual_types}
        assert as_dict == {expected.name: expected}

    def test_GetBeer(self):
        serv = BartenderServicer()

        setup = bartender_pb2.CreateBeerRequest(brand="HertogJan")
        setup.beer.name = "Pils"
        _ = serv.CreateBeer(setup, context=None)

        input = bartender_pb2.GetBeerRequest(brand="HertogJan", name="Pils")
        expected = bartender_pb2.Beer(name="Pils")

        actual = serv.GetBeer(input, context=None)

        assert actual == expected
