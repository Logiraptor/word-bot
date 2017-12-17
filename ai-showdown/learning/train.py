import random

from game_context import GameContext
from predictor import Predictor

Y = 0.85

def live_training(predictor, iterations=1000):
    """train the given predictor for {iterations} games"""
    with predictor.session() as session:
        train_data, train_labels = [], []
        for i in range(iterations):
            new_data, new_labels = run_game(session)
            train_data.extend(new_data)
            train_labels.extend(new_labels)
            train_data, train_labels = session.train(train_data, train_labels)
            print("Finished iteration %d / %d" % (i + 1, iterations))


def run_game(session):
    move_history = []

    ctx = GameContext.make()
    done = False
    while not done:
        move_history.append(ctx)
        next_ctx = run_round(ctx, session)
        if next_ctx is not None:
            ctx = next_ctx
        else:
            done = True

    move_history.reverse()
    if ctx.result().winner and len(move_history) % 2 == 0:
        label = 1.0
    else:
        label = -1.0

    train_data, train_labels = [], []
    for move in move_history:
        train_data.append(to_tensor(move))
        train_labels.append([label])
        label *= -Y

    shuffled = list(zip(train_data, train_labels))
    random.shuffle(shuffled)
    train_data[:], train_labels[:] = zip(*shuffled)

    return [train_data, train_labels]


def run_round(ctx, predictor):
    moves = ctx.get_moves()
    best_score = 0
    best_move = None
    for move in moves:
        score = predictor.score(move.get_tensor())
        if best_move is None or score > best_score:
            best_score = score
            best_move = move
    return best_move


def to_tensor(ctx):
    return ctx.get_tensor()


live_training(Predictor(), iterations=10000)
