from game_context import GameContext
import tensorflow as tf


def live_training(predictor, iterations=1000):

    for i in range(iterations):
        ctx = run_game(predictor)
        ctx.dump()
        print(ctx.result())


def run_game(predictor):
    move_history = []

    ctx = GameContext.make()
    done = False
    while not done:
        move_history.append(ctx)
        next_ctx = run_round(ctx, predictor)
        if next_ctx is not None:
            ctx = next_ctx
        else:
            done = True

    if ctx.result().winner and len(move_history) % 2 == 0:
        predictor.train_win(map(to_tensor, move_history[::2]))
    else:
        predictor.train_loss(map(to_tensor, move_history[1::2]))

    return ctx


def run_round(ctx, predictor):
    moves = ctx.getMoves()
    best_score = 0
    best_move = None
    for move in moves:
        score = predictor.score(move.getTensor())
        if best_move is None or score > best_score:
            best_score = score
            best_move = move
    return best_move


def to_tensor(ctx):
    return ctx.getTensor()


# TODO: Machine learning
class Predictor(object):
    def score(self, tensor):
        return 0.1

    def train_win(self, tensors):
        pass

    def train_loss(self, tensors):
        pass


live_training(Predictor(), iterations=1)
