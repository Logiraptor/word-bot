from ctypes import *
from collections import namedtuple

_wordbot = cdll.LoadLibrary('./wordbot.so')

_makeContext = _wordbot.MakeContext
_makeContext.restype = c_longlong
_freeContext = _wordbot.FreeContext
_freeContext.restype = c_longlong
_printContext = _wordbot.PrintContext
_printContext.restype = c_longlong
_generateMoves = _wordbot.GenerateMoves
_generateMoves.argtypes = [
    c_longlong,
    POINTER(POINTER(c_longlong)),
    POINTER(c_longlong),
]
_convertToTensor = _wordbot.ConvertToTensor
_convertToTensor.argtypes = [
    c_longlong,
    POINTER(POINTER(c_double)),
    POINTER(c_longlong),
]
_freeContextBuffer = _wordbot.FreeContextBuffer
_freeContextBuffer.argtypes = [
    POINTER(POINTER(c_longlong)),
]
_freeTensorBuffer = _wordbot.FreeTensorBuffer
_freeTensorBuffer.argtypes = [
    POINTER(POINTER(c_double)),
]


Result = namedtuple('Result', ['winner'])


class GameContext(object):
    def __init__(self, key):
        self.ctx = key

    @classmethod
    def _from_key(cls, key):
        return GameContext(key)

    @classmethod
    def make(cls):
        return GameContext(_makeContext())

    def dump(self):
        _printContext(self.ctx)

    def free(self):
        _freeContext(self.ctx)

    def getMoves(self):
        element_ptr = POINTER(c_longlong)()
        element_len = c_longlong()
        _generateMoves(self.ctx, pointer(element_ptr), pointer(element_len))
        output = [GameContext._from_key(element_ptr[i])
                  for i in range(element_len.value)]
        return output

    def getTensor(self):
        element_ptr = POINTER(c_double)()
        element_len = c_longlong()
        _convertToTensor(self.ctx, pointer(element_ptr), pointer(element_len))
        output = [element_ptr[i] for i in range(element_len.value)]
        return output

    def result(self):
        # TODO: load result from _wordbot
        return Result(winner=True)
