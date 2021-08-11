# go_cache

# IsExpired:

I decided to move with a very simple idea. Nothing special, just added a certain duration into the current time. 

# GetPricesFor

In order to make calls in concurrently, I decided to use a goroutine. Im not sure if its well-designed since im not so familiar with coroutines patterns. <br>

Due to concurrently approach for the GetPricesFor, I faced some issues writing cache into map. In order to avoid using a mutex, I decided to:

    - Duplicate the "get in the external system and persist to cache" approach.

# Pending items

I was planning to create chunks for the operation <b>GetPricesFor</b> in order to limit parallel calls to do not take the external system down.
I was not sure if it is a requirement. 
