# go_cache

- IsExpired:

I decided to move with a very simple idea. Nothing special, just added a certain duration into the current time. 

-  GetPricesFor

In order to make calls in concurrently, I decided to use a goroutine. I`m not sure if it`s well-designed since i`m not so familiar with coroutines patterns.
Due to concurrently approach for the GetPricesFor, I faced some issues writing cache into map. In order to avoid using a mutex, I decided to:
    - Duplicate the "get in the external system and persist to cache" approach.