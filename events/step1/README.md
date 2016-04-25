# Step 1: JSON and local storage

The goal of this step is to start encoding and encoding events.
For now they will be stored in a global variable protected by a mutex (to avoid data races).
This is not optimal because at any point the instance could disappear and we would lose our data.

Don't worry too much about that, the next step takes care of storage with datastore.

Once done, go back to the [JSON section](../../section06/README.md#congratulations).
