from contextlib import contextmanager

import tensorflow as tf
from tensorflow.contrib import slim
import numpy as np

BATCH_SIZE = 128
NUM_FEATURES = 15 * 15 * 3 + 27 * 2
NUM_LABELS = 1


class Predictor(object):
    saver: tf.train.Saver
    graph: tf.Graph

    def __init__(self):
        self.graph = tf.Graph()
        with self.graph.as_default():
            self.score_dataset = tf.placeholder(tf.float32, shape=(1, NUM_FEATURES))
            self.train_dataset = tf.placeholder(tf.float32, shape=(BATCH_SIZE, NUM_FEATURES))
            self.train_labels = tf.placeholder(tf.float32, shape=(BATCH_SIZE, NUM_LABELS))

            def layer(input_node, output):
                return slim.fully_connected(input_node, output,
                                            activation_fn=tf.nn.sigmoid)

            def model(input_node):
                layer1 = layer(input_node, 64)
                layer2 = layer(layer1, 48)
                layer3 = layer(layer2, 24)
                return (layer(layer3, 1) - 0.5) * 2

            self.score_model = model(self.score_dataset)
            self.model = model(self.train_dataset)
            self.loss = tf.reduce_sum(tf.square(self.model - self.train_labels))

            self.optimizer = tf.train.GradientDescentOptimizer(learning_rate=0.01).minimize(self.loss)
            self.saver = tf.train.Saver()

    @contextmanager
    def session(self):
        with tf.Session(graph=self.graph) as session:
            tf.global_variables_initializer().run()
            yield PredictorSession(self, session)


class PredictorSession(object):
    predictor: Predictor

    def __init__(self, predictor, session):
        self.predictor = predictor
        self.session = session

    def score(self, tensor):
        return self.session.run([self.predictor.score_model], feed_dict={
            self.predictor.score_dataset: np.array([tensor])
        })

    def train(self, positions, labels):
        offset = 0
        while (len(positions) - offset) > BATCH_SIZE:
            batch_data = positions[offset:offset + BATCH_SIZE]
            batch_labels = labels[offset:offset + BATCH_SIZE]

            _, predictions, loss = self.session.run(
                [self.predictor.optimizer, self.predictor.model, self.predictor.loss],
                feed_dict={
                    self.predictor.train_dataset: np.array(batch_data),
                    self.predictor.train_labels: np.array(batch_labels)
                })
            print('Predictions: ', list(zip(predictions, batch_labels)))
            print(f'Loss: {loss:f}')
            offset += BATCH_SIZE

        save_path = self.predictor.saver.save(self.session, './saved/maybe-better-y')
        print(f"Saved checkpoint to {save_path}")

        return positions[offset:], labels[:]
