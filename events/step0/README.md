# Step 0: basic architecture

The goal for this step is to define the two entry points to your web application.

You can try running before you start coding to see the current behavior.

    $ goapp serve

Once you've implemented this, visiting `localhost:8080` should display a list of conferences,
and clicking "New Event" should not do anything.

You can also deploy it.

    $ goapp deploy --version=step0 --application=your-project-id

And then visit https://step0.your-project-id.appspot.com.

Once you're done, let's go back to the [instructions](../../section05/README.md#congratulations).
