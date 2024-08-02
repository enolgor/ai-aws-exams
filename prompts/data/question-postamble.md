Please generate 2 question per each domain/task statement with one correct answer and three incorrect answers (total of four) in JSON with the following format (example for a question of task 1.1):

```
[
  {
    "domain": 1,
    "task": 1,
    "question": "<the question>",
    "answers": [
      {
        "answer": "<the answer>",
        "correct": true/false,
        "explanation": "<why is correct or incorrect>"
      }
    ]
  },
]
```