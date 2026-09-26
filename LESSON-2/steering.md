---
inclusion: always
---

# Temperature Input Validation

Create a steering file to enforce that **all temperature inputs are validated as Celsius values within the range of -5°C (minimum) to +200°C (maximum)**, so that **Kiro consistently rejects out-of-range or invalid temperature readings before they reach business logic, preventing sensor errors, physically impossible values, and downstream bugs in any feature that consumes temperature data**. For example, the following implementation:

```typescript
const TEMP_MIN_CELSIUS = -5;
const TEMP_MAX_CELSIUS = 200;

/**
 * Validates a temperature value in degrees Celsius.
 * @throws {TypeError} if the value is not a valid number.
 * @throws {RangeError} if the value is outside the allowed range.
 */
function validateTemperatureCelsius(value: number): number {
  if (typeof value !== "number" || Number.isNaN(value)) {
    throw new TypeError("Temperature must be a valid number in degrees Celsius.");
  }
  if (value < TEMP_MIN_CELSIUS || value > TEMP_MAX_CELSIUS) {
    throw new RangeError(
      `Temperature ${value}°C is outside the allowed range (${TEMP_MIN_CELSIUS}°C to ${TEMP_MAX_CELSIUS}°C).`
    );
  }
  return value;
}
```

ensures that **every temperature input is a valid Celsius value within the range of -5°C to +200°C** is followed.

## Rules Kiro must enforce

- Every temperature input MUST be treated as degrees Celsius.
- The minimum accepted value is **-5°C** (inclusive).
- The maximum accepted value is **+200°C** (inclusive).
- Reject non-numeric values (`NaN`, `null`, `undefined`, non-convertible strings) with a `TypeError`.
- Reject out-of-range values with a `RangeError` and a descriptive message.
- Apply validation **at the input boundary** (API, form, sensor reading) before any calculation or persistence.
