# Step 4: storing temporary results in Memcache

In the previous step we ended up with an application that was wasting time
and network resources by sending a new request to the weather API for each event,
every time someone listed the events.

Now we will fix that by using App Engine Memcache. This is simply a fully
managed Memcache instance, you don't need to do anything to start using it!
Isn't that cool?

Once you're done with this exercise you will see that your application is much faster,
and it consumes less resources.

Congratulations, you're done! You can check your code comparing it to the one in
[step 5](../step5), or go back to the [instructions](../../section09/README.md#congratulations)
for the last instructions.
