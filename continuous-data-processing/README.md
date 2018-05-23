# Continuous data processing

## purpose

The idea here is to do a sync data processing. First you cut of the last|first item and make function call "itself" with data.

## cut of the last|first item

Depending on function's configuration, there are two options:

 - pop the first item
 - pop the last item

It helps to build a sequence of data processing, would that be direct (similar to FIFO) or revers (FILO).

## call "itself"?

Basically you make a function call itself through Fn with reduced data.
