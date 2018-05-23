# An example of a recursive Fn function

## it's not a real tail recursion, but the way to keep function running

That's said, it's not a tail recursion, it's just a way to keep your function running as long as it necessary.

## "as long as it necessary" means it is limited

Indeed, there's no point in turning a function into a daemon, so the recursion or basically, a redirects must be limited.
Otherwise, you need a microservice, not a function.

## HTTP 3XX Redirect is not the case

Redirect requires to keep original connection up and you HTTP client will likely follow the redirects.

## recursive Fn function needs to be asynchronous

In case of sync execution, Fn will keep the connection between a caller and a function within the call timeout, but not in case of async execution.
Fn behaves differently in terms of async calls. A caller will get the call ID for each execution, so, there would be no blocking connections.

## how to limit recursion depth?

Personally, i found query parameters quite useful for bypassing recursion depth remaining.
However, no matter which road to take, here would be a risk of spinning up too much function containers 
and the number of calls may vary (be greater than max recursion depth). That's why there's always has to be "Plan B" - define default value for the recursion.

## why recursive serverless function matter?

Recursive serverless function matters in those cases when you need to trigger a function based on certain event, let's say, 46 files created in store.
Basically, it allows developers to define more custom "match criteria" to exit the recursion that is not possible to implement 
using any of built-it notifications API that 3rd-party service exposes.

