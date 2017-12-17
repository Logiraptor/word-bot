# Q-Learning notes

## Definitions

- Policy
   - Learned rewards for each possible action and a way to chose the optimal one.
- Policy Gradient
   - Gradient descent trained policy through environment feedback
- Value Functions
   - Learn to predict how good a state is to be in.

- Action Selection Strategies
   - E-greedy - Mostly choose highest expected reward, but with _e_ probability, choose randomly

- Policy Loss equation `Loss = -log(π)*A`
    - A - Advantage - A tuning parameter to give a 'baseline' to rewards.
    - π - Policy (the action weight)

- Markov Decision Process
    - Consists of:
        - S - set of all states
        - s - a specific state that the agent experiences
        - A - Set of all actions
        - a - a specific action
        - T(s, a) - probability of transition into new state s' by taking action a in state s
        - R(s, a) - reward of taking action a in state s

- Experience traces
    - A set of state/action pairs that lead to a delayed reward

## Questions

The text mentions action-selection strategies in Pt 7. So far only e-greedy

If I train a classifier to recognize scrabble tiles, I might be able to use the yolo method
to find all tiles on a board image

## Interestings

- Tensorflow doesn't just do matrix math. You can index arrays and stuff.
    - This is useful when you want to update a single weight during training
