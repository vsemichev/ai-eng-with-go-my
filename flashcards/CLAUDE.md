- after each task is completed, go build the project to verify your output

## Logging standards

- Use the standard library `log/slog` package for all logging.
- Log at the **start** of an operation and again on its **successful**
  completion, using the **info** level.
- Log every error using the **error** level, but only at the point where the
  error actually originates. Do not re-log the same error again as it
  propagates up the call stack — wrap and return it instead.