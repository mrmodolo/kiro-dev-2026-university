# Property-based testing (IDE only)

Property-based testing (PBT) increases confidence in AI-generated code by
moving away from example-based tests and instead enforcing a general rule that
must always hold. Kiro extracts properties from your spec requirements,
determining what can be logically tested, and creates and runs hundreds of
tests against the high-level intent of that requirement. Since PBTs are
optional by default, you can get your core implementation right, and then run
property checks to ensure the output matches your intent.


Kiro will generate PBTs by default during the design phase of your project. All
PBTs are optional, so you can apply correctness to requirements and behavior
that matter in your own project. PBT is only available in the Kiro IDE, so
start your project there or import your .kiro configuration from your preferred
Kiro tool.
