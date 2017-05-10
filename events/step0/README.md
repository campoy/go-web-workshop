# Step 0: basic architecture

The goal for this step is to define the two entry points to your web application.

You can try running before you start coding to see the current behavior.

    $ dev_appserver.py .

Once you've implemented this, visiting `localhost:8080` should display a list of conferences,
and clicking "New Event" should not do anything.

In order to deploy you'll need to first install and configure [gcloud](https://cloud.google.com/sdk/downloads):

    $ gcloud init

_Note_: You do not need to set up any Google Compute Engine zone.

Then simply run

    $ gcloud app deploy --version=step0 app.yaml

And then visit https://step0.your-project-id.appspot.com or running

    $ gcloud app browse --version=step0

This will display a basic Events page but also an alert will be alerted. That's fine, for now.

Once you're done, let's go back to the [instructions](../../section05/README.md#congratulations).
