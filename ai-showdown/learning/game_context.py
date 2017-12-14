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
_getFinalScore = _wordbot.GetFinalScore
_getFinalScore.argtypes = [
    c_longlong
]
_getFinalScore.restype = c_longlong

Result = namedtuple('Result', ['winner', 'diff'])

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

    def get_moves(self):
        element_ptr = POINTER(c_longlong)()
        element_len = c_longlong()
        _generateMoves(self.ctx, pointer(element_ptr), pointer(element_len))
        output = [GameContext._from_key(element_ptr[i])
                  for i in range(element_len.value)]
        return output

    def get_tensor(self):
        element_ptr = POINTER(c_double)()
        element_len = c_longlong()
        _convertToTensor(self.ctx, pointer(element_ptr), pointer(element_len))
        output = [element_ptr[i] for i in range(element_len.value)]
        return output

    def result(self):
        final_score = int(_getFinalScore(self.ctx))
        return Result(winner=final_score > 0, diff=final_score)
